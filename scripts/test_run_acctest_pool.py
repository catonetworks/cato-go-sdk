from contextlib import redirect_stdout
import io
import json
import os
from pathlib import Path
import random
import tempfile
import sys
import unittest
from unittest.mock import Mock, patch

from run_acctest_pool import execute_redacted, load_configuration, run


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
        execute = Mock(side_effect=[Mock(returncode=int(i == 2)) for i in range(10)])
        with patch.dict(os.environ, {'CATO_ACCTEST_CONFIG_FILE': '/secret/pool.json', 'CATO_ACCTEST_ACCOUNTS': 'secret-pool'}):
            result = run(Path('/provider'), self.planner, self.config, execute=execute, rng=random.Random(7))
        self.assertEqual(1, result)
        self.assertEqual(10, execute.call_count)
        seen = set()
        for group, call in enumerate(execute.call_args_list):
            self.assertEqual(['make', 'acctest-flaky', f'group={group}'], call.args[0])
            env = call.kwargs['env']
            seen.add(env['CATO_ACCOUNT_ID'])
            self.assertNotIn('CATO_ACCTEST_ACCOUNTS', env)
            self.assertNotIn('CATO_ACCTEST_CONFIG_FILE', env)
            self.assertEqual(str(json.loads(env['TFACC_TEST_VARS'])['fixture']), env['CATO_TEST_SITE_ID'])
        self.assertEqual(10, len(seen))

    def test_single_suite_runs_once_and_clears_previous_account_fixtures(self):
        for account in self.config['accounts']:
            account['test_env'] = {}
        execute = Mock(return_value=Mock(returncode=0))
        with patch.dict(os.environ, {'CATO_TEST_SITE_ID': 'stale-site'}):
            self.assertEqual(0, run(Path('/provider'), self.planner, self.config,
                                    'internal/acctests/if_rule', execute=execute))
        execute.assert_called_once()
        self.assertEqual('t=if_rule', execute.call_args.args[0][-1])
        self.assertNotIn('CATO_TEST_SITE_ID', execute.call_args.kwargs['env'])

    def test_invalid_suite_never_starts_make(self):
        for suite in ('if_rule;id', '../if_rule', 'missing_suite', '--help'):
            execute = Mock()
            with self.assertRaises(ValueError):
                run(Path('/provider'), self.planner, self.config, suite, execute=execute)
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


if __name__ == '__main__':
    unittest.main()
