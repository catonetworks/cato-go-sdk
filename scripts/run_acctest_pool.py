#!/usr/bin/env python3
"""Run provider acceptance groups concurrently using a Jenkins secret-file pool."""

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
import importlib.util
import json
import os
from pathlib import Path
import random
import re
import shutil
import subprocess
import sys
import tempfile
from urllib.parse import urlsplit


def load_configuration(path, planner):
    try:
        config = json.loads(path.read_text())
    except (OSError, ValueError):
        raise ValueError('Acceptance configuration is missing or invalid JSON') from None
    if not isinstance(config, dict):
        raise ValueError('Acceptance configuration must be an object')
    endpoint = config.get('CATO_BASEURL', '')
    if not isinstance(endpoint, str):
        raise ValueError('CATO_BASEURL must be an HTTPS URL')
    url = urlsplit(endpoint)
    if url.scheme != 'https' or not url.hostname or url.username or url.password or url.query or url.fragment:
        raise ValueError('CATO_BASEURL must be an HTTPS URL without credentials or query parameters')
    skips = config.get('TFACC_TEST_SKIP', {})
    if not isinstance(skips, dict) or any(not isinstance(k, str) or not isinstance(v, str) for k, v in skips.items()):
        raise ValueError('TFACC_TEST_SKIP must map test names to skip reasons')
    accounts = config.get('accounts')
    planner.validate_accounts(accounts)
    planner.validate_fixture_references(accounts)
    return config


def execute_redacted(command, *, cwd, env):
    """Stream test output while keeping the selected API key out of Jenkins logs."""
    with subprocess.Popen(command, cwd=cwd, env=env, stdout=subprocess.PIPE,
                          stderr=subprocess.STDOUT, text=True, encoding='utf-8', errors='replace') as process:
        for line in process.stdout:
            prefix = f'[{env["ACCTEST_LABEL"]}] ' if 'ACCTEST_LABEL' in env else ''
            print(prefix + line.replace(env['CATO_TOKEN'], '****'), end='', flush=True)
        return subprocess.CompletedProcess(command, process.wait())


def go_environment():
    env = os.environ.copy()
    # Ten API workers must not also become ten unrestricted cold Go builds.
    # trimpath lets the separate provider copies reuse the warmed build cache.
    env['GOFLAGS'] = (env.get('GOFLAGS', '') + ' -p=1 -trimpath').strip()
    env['GOMAXPROCS'] = '2'
    env['GOMEMLIMIT'] = '512MiB'
    return env


def warm_build_cache(provider, suite=''):
    env = go_environment()
    for name in ('CATO_TOKEN', 'CATO_ACCOUNT_ID', 'CATO_BASEURL', 'CATO_ENDPOINT',
                 'CATO_ACCTEST_ACCOUNTS', 'CATO_ACCTEST_CONFIG_FILE',
                 'TFACC_TEST_VARS', 'TFACC_TEST_SKIP'):
        env.pop(name, None)
    env.update(TF_ACC='0', ACCTEST_CLEANUP='false')
    packages = [f'./internal/acctests/{suite}', './internal/acctests/acc'] if suite else ['./internal/acctests/...']
    print('Compiling acceptance packages once, with no test cases or API calls; this may take several minutes', flush=True)
    subprocess.run(['go', 'test', '-tags=acctest', '-run=^$', *packages],
                   cwd=provider, env=env, check=True, timeout=1200)
    print('Acceptance build cache is ready; starting API test workers', flush=True)


def prepare_workspace(provider, destination):
    shutil.copytree(provider, destination,
                    ignore=shutil.ignore_patterns('.git', '.terraform', 'tmp_recorded', '__pycache__'))
    # Jenkins replaces the SDK with ../ before invoking this runner. Copies in
    # temporary directories must keep pointing at that exact SDK checkout.
    subprocess.run(['go', 'mod', 'edit',
                    f'-replace=github.com/catonetworks/cato-go-sdk={provider.parent}'],
                   cwd=destination, check=True, timeout=60)
    return destination


def run(provider, planner, config, suite='', *, execute=execute_redacted, rng=None,
        prepare=prepare_workspace, warm=warm_build_cache):
    planner.validate_inventory()
    if suite:
        match = re.fullmatch(r'(?:\./)?(?:internal/acctests/)?([a-z0-9_]+)/?', suite)
        if not match or match[1] not in planner.load_timings()['packages']:
            raise ValueError('ACCTEST_FILTER must name an existing acceptance suite')
        suite = match[1]
    slots = list(range(planner.GROUP_COUNT))
    (rng or random.SystemRandom()).shuffle(slots)
    jobs = [(None, slots[0])] if suite else list(enumerate(slots))
    failures = []
    def run_selection(group, slot, workspace):
        account = config['accounts'][slot]
        env = go_environment()
        for name in (*planner.TEST_ENV_NAMES, 'CATO_ACCTEST_ACCOUNTS', 'CATO_ACCTEST_CONFIG_FILE'):
            env.pop(name, None)
        env.update({
            'CATO_BASEURL': config['CATO_BASEURL'],
            'CATO_ENDPOINT': config['CATO_BASEURL'],
            'CATO_ACCOUNT_ID': account['id'],
            'CATO_TOKEN': account['api_key'],
            'TFACC_TEST_VARS': json.dumps(account.get('test_vars', {})),
            'TFACC_TEST_SKIP': json.dumps(config.get('TFACC_TEST_SKIP', {})),
            'TFACC_ENABLE_RULES_INDEX_CRUD': 'true',
            'DISABLE_POLICY_RULE_CLEANUP': 'true',
        })
        env.update(account.get('test_env', {}))
        label = f'suite {suite}' if suite else f'group {group}'
        env['ACCTEST_LABEL'] = label
        print(f'Running acceptance {label} on pool slot {slot}', flush=True)
        selection = f't={suite}' if suite else f'group={group}'
        try:
            result = execute(['make', 'acctest-flaky', selection], cwd=workspace, env=env)
            return label if result.returncode else None
        except OSError:
            print(f'Could not start acceptance {label}', file=sys.stderr)
            return label

    if suite:
        warm(provider, suite)
        failure = run_selection(None, slots[0], provider)
        if failure:
            failures.append(failure)
    else:
        with tempfile.TemporaryDirectory(prefix='cato-sdk-acctest-') as directory:
            # Prepare every workspace before launching any real API calls.
            workspaces = {group: prepare(provider, Path(directory) / f'group-{group}')
                          for group, _ in jobs}
            warm(workspaces[0])
            with ThreadPoolExecutor(max_workers=planner.GROUP_COUNT) as executor:
                futures = {executor.submit(run_selection, group, slot, workspaces[group])
                           for group, slot in jobs}
                for future in as_completed(futures):
                    failure = future.result()
                    if failure:
                        failures.append(failure)
    if failures:
        print('Failed acceptance selections: ' + ', '.join(sorted(failures)), file=sys.stderr)
    return int(bool(failures))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--provider', type=Path, required=True)
    parser.add_argument('--check', action='store_true', help='Validate configuration without running tests')
    args = parser.parse_args()
    provider = args.provider.resolve()
    spec = importlib.util.spec_from_file_location('provider_acctest_groups', provider / 'test_data/acctest_groups.py')
    if spec is None or spec.loader is None:
        raise ValueError('Provider checkout does not contain the acceptance group planner')
    planner = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(planner)
    config_path = os.environ.get('CATO_ACCTEST_CONFIG_FILE')
    if not config_path:
        raise ValueError('CATO_ACCTEST_CONFIG_FILE is required')
    config = load_configuration(Path(config_path), planner)
    if args.check:
        planner.validate_inventory()
        print('Acceptance pool configuration and group inventory are valid')
        return 0
    return run(provider, planner, config, os.environ.get('ACCTEST_FILTER', '').strip())


if __name__ == '__main__':
    try:
        sys.exit(main())
    except subprocess.TimeoutExpired:
        sys.exit('Acceptance preparation timed out; API workers were not started')
    except subprocess.CalledProcessError as error:
        sys.exit(f'Acceptance preparation command failed (exit {error.returncode}); see preceding build output')
    except (ValueError, OSError):
        # Configuration may contain keys. Never emit its contents or a traceback.
        sys.exit('Acceptance setup failed: check the secret-file configuration and provider planner')
