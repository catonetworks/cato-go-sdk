from contextlib import redirect_stderr, redirect_stdout
import io
import json
import os
from pathlib import Path
import random
import tempfile
import sys
import threading
import unittest
from unittest.mock import Mock, patch

from run_acctest_pool import execute_redacted, load_configuration, prepare_workspace, run, warm_build_cache


class AcceptancePoolTest(unittest.TestCase):
    def setUp(self):
        # Tests do not require another repository or any API credentials.
        self.planner = Mock(GROUP_COUNT=10, TEST_ENV_NAMES={'CATO_TEST_SITE_ID'})
        self.planner.load_timings.return_value = {'packages': {'if_rule': {}}}
        self.config = {
            'CATO_BASEURL': 'https://example.invalid/api/v1/graphql2',
            'TFACC_TEST_SKIP': {'TestExample': 'test fixture'},
            'accounts': [{'id': str(i + 1), 'api_key': f'inert-test-key-{i}',
                          'test_vars': {'fixture': i}, 'test_env': {'CATO_TEST_SITE_ID': str(i)}}
                         for i in range(10)],
        }

    def test_all_groups_use_distinct_accounts_and_aggregate_failures(self):
        barrier = threading.Barrier(10, timeout=10)
        warm = Mock()
        def execute_group(command, *, cwd, env):
            # A sequential implementation times out here: all ten executions
            # must be in progress before any may finish.
            warm.assert_called_once()
            barrier.wait()
            return Mock(returncode=int(command[-1] == 'group=2'))
        execute = Mock(side_effect=execute_group)
        prepare = Mock(side_effect=lambda provider, destination: destination)
        with patch.dict(os.environ, {'CATO_ACCTEST_CONFIG_FILE': '/secret/pool.json', 'CATO_ACCTEST_ACCOUNTS': 'secret-pool'}), \
                redirect_stdout(io.StringIO()), redirect_stderr(io.StringIO()):
            result = run(Path('/provider'), self.planner, self.config, execute=execute,
                         rng=random.Random(7), prepare=prepare, warm=warm)
        self.assertEqual(1, result)
        self.assertEqual(10, execute.call_count)
        seen = set()
        directories = set()
        calls = sorted(execute.call_args_list, key=lambda call: int(call.args[0][-1].split('=')[1]))
        for group, call in enumerate(calls):
            self.assertEqual(['make', 'acctest-flaky', f'group={group}'], call.args[0])
            env = call.kwargs['env']
            seen.add(env['CATO_ACCOUNT_ID'])
            directories.add(call.kwargs['cwd'])
            self.assertNotIn('CATO_ACCTEST_ACCOUNTS', env)
            self.assertNotIn('CATO_ACCTEST_CONFIG_FILE', env)
            self.assertIn('-p=1', env['GOFLAGS'])
            self.assertIn('-trimpath', env['GOFLAGS'])
            self.assertEqual('2', env['GOMAXPROCS'])
            self.assertEqual('512MiB', env['GOMEMLIMIT'])
            self.assertEqual(str(json.loads(env['TFACC_TEST_VARS'])['fixture']), env['CATO_TEST_SITE_ID'])
        self.assertEqual(10, len(seen))
        self.assertEqual(10, len(directories))
        self.assertEqual(10, prepare.call_count)

    def test_single_suite_runs_once_and_clears_previous_account_fixtures(self):
        for account in self.config['accounts']:
            account['test_env'] = {}
        execute = Mock(return_value=Mock(returncode=0))
        with patch.dict(os.environ, {'CATO_TEST_SITE_ID': 'stale-site'}):
            self.assertEqual(0, run(Path('/provider'), self.planner, self.config,
                                    'internal/acctests/if_rule', execute=execute, warm=Mock()))
        execute.assert_called_once()
        self.assertEqual('t=if_rule', execute.call_args.args[0][-1])
        self.assertNotIn('CATO_TEST_SITE_ID', execute.call_args.kwargs['env'])

    def test_invalid_suite_never_starts_make(self):
        for suite in ('if_rule;id', '../if_rule', 'missing_suite', '--help'):
            execute = Mock()
            with self.assertRaises(ValueError):
                run(Path('/provider'), self.planner, self.config, suite, execute=execute, warm=Mock())
            execute.assert_not_called()

    def test_configuration_checks_endpoint_and_pool_before_execution(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / 'config.json'
            path.write_text(json.dumps(self.config))
            load_configuration(path, self.planner)
            self.planner.validate_accounts.assert_called_once_with(self.config['accounts'])
            self.planner.validate_fixture_references.assert_called_once_with(self.config['accounts'])
            self.config['CATO_BASEURL'] = 'http://example.invalid'
            path.write_text(json.dumps(self.config))
            with self.assertRaises(ValueError):
                load_configuration(path, self.planner)

    def test_child_output_redacts_the_selected_key(self):
        output = io.StringIO()
        with redirect_stdout(output):
            result = execute_redacted([sys.executable, '-c', 'import os; print(os.environ["CATO_TOKEN"])'],
                                     cwd=Path.cwd(), env={**os.environ, 'CATO_TOKEN': 'inert-test-key'})
        self.assertEqual(0, result.returncode)
        self.assertEqual('****\n', output.getvalue())

    def test_workspaces_isolate_runner_files_and_preserve_sdk_checkout(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            provider = root / 'sdk' / 'provider'
            provider.mkdir(parents=True)
            (provider / 'go.mod').write_text('module test\n')
            (provider / 'tmp_recorded').mkdir()
            (provider / 'tmp_recorded' / 'stale-log').write_text('old run')
            with patch('run_acctest_pool.subprocess.run') as edit:
                first = prepare_workspace(provider, root / 'group-0')
                second = prepare_workspace(provider, root / 'group-1')
            self.assertEqual(2, edit.call_count)
            self.assertIn(f'-replace=github.com/catonetworks/cato-go-sdk={provider.parent}',
                          edit.call_args.args[0])
            self.assertFalse((first / 'tmp_recorded').exists())
            (first / 'go.mod').write_text('changed by one group')
            self.assertEqual('module test\n', (second / 'go.mod').read_text())
            self.assertEqual('module test\n', (provider / 'go.mod').read_text())

    def test_warmup_compiles_without_test_cases_or_account_credentials(self):
        with patch.dict(os.environ, {'CATO_TOKEN': 'inert-test-key', 'CATO_ACCTEST_CONFIG_FILE': '/secret/pool'}):
            with patch('run_acctest_pool.subprocess.run') as compile_packages:
                warm_build_cache(Path('/provider'))
        command = compile_packages.call_args.args[0]
        env = compile_packages.call_args.kwargs['env']
        self.assertIn('-run=^$', command)
        self.assertEqual('./internal/acctests/...', command[-1])
        self.assertEqual('0', env['TF_ACC'])
        self.assertEqual('false', env['ACCTEST_CLEANUP'])
        self.assertNotIn('CATO_TOKEN', env)
        self.assertNotIn('CATO_ACCTEST_CONFIG_FILE', env)

    def test_warmup_failure_prevents_any_api_worker_from_starting(self):
        execute = Mock()
        prepare = Mock(side_effect=lambda provider, destination: destination)
        warm = Mock(side_effect=OSError('compile failed'))
        with self.assertRaises(OSError):
            run(Path('/provider'), self.planner, self.config, execute=execute, prepare=prepare, warm=warm)
        execute.assert_not_called()


if __name__ == '__main__':
    unittest.main()
