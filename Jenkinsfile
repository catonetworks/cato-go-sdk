// Triggered manually or via GitHub webhook (Generic Webhook Trigger plugin).
// Checks out a branch of cato-go-sdk, clones terraform-provider-cato
// at the requested branch, replaces its go.mod SDK dependency with the local
// checkout, and runs acceptance tests against a real Cato API account.
//
// GitHub webhook setup (one-time, per repo):
//   URL:          https://jenkins.automation.catonetworks.club/generic-webhook-trigger/invoke?token=sdk-acctest
//   Content type: application/json
//   Events:       Pushes   ← select "Pushes" (not "Pull requests")
//
// Required Jenkins credentials:
//   cato-acctest-config     -> Secret file: CATO_BASEURL, TFACC_TEST_SKIP, accounts
//                            (ten records with id, api_key, test_vars, test_env)
//   automation-github-user -> used to clone provider repo
//   cato-acctest-gh-token   -> GitHub PAT with repo:status scope
// Configuration source: STA test account configuration + STA test account creds
// in the open-source development 1Password vault. Never commit account keys.
// This job serializes its builds. Do not overlap other CI using the same pool.

pipeline {
    agent {
        docker {
            image 'golang:1.26.11'
            args  '-e GOTOOLCHAIN=auto -e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/gopath'
        }
    }

    triggers {
        GenericTrigger(
            genericVariables: [
                // Push payload carries the full ref: refs/heads/<branch>
                [key: 'GIT_REF', value: '$.ref',   defaultValue: ''],
                // Commit SHA after the push — also rejects branch-deletion events
                [key: 'GIT_SHA', value: '$.after', defaultValue: '']
            ],
            token: 'sdk-acctest',
            causeString: 'Push to $GIT_REF',
            // Trigger on branch pushes, including merges to main. Reject tags and
            // branch-deletion pushes, whose $.after is 40 zeroes.
            regexpFilterText:       '$GIT_REF:$GIT_SHA',
            regexpFilterExpression: '^refs/heads/.+:(?!0{40}$)[0-9a-fA-F]{40}$',
            printContributedVariables: true,
            printPostContent: false
        )
    }

    parameters {
        string(
            name: 'SDK_BRANCH',
            defaultValue: 'main',
            description: 'Branch of cato-go-sdk to test (e.g. feat/my-new-query)'
        )
        string(
            name: 'PROVIDER_BRANCH',
            defaultValue: 'main',
            description: 'Branch of terraform-provider-cato to run acceptance tests from'
        )
        string(
            name: 'ACCTEST_FILTER',
            defaultValue: '',
            description: 'Limit run to a specific package, e.g. "if_rule" or "internal/acctests/if_rule". Leave empty to run all.'
        )
    }

    options {
        disableConcurrentBuilds()
        buildDiscarder(logRotator(numToKeepStr: '20'))
        timeout(time: 2, unit: 'HOURS')
        timestamps()
    }

    stages {
        stage('Prepare acceptance runner') {
            steps {
                // Keep the runner from this Jenkinsfile's revision even when
                // SDK_BRANCH predates multi-account support.
                stash name: 'acceptance-runner', includes: 'scripts/run_acctest_pool.py'
                sh 'command -v python3 >/dev/null || { apt-get update && apt-get install -y --no-install-recommends python3; }'
                sh "python3 -m unittest discover -s scripts -p 'test_*.py'"
            }
        }

        stage('Checkout SDK branch') {
            steps {
                // GIT_REF (from webhook) is refs/heads/<branch>; SDK_BRANCH (manual param)
                // is just the bare name. Normalise to a bare branch name either way.
                script {
                    sh '''
                        RAW="${GIT_REF:-${SDK_BRANCH}}"
                        BRANCH="${RAW#refs/heads/}"
                        if [ -z "$BRANCH" ]; then
                            echo "SDK branch is empty" >&2
                            exit 1
                        fi
                        git fetch origin
                        git checkout -B "$BRANCH" "origin/$BRANCH"
                    '''
                    // Capture the actual commit SHA so we can post status back to GitHub.
                    // For webhook builds GIT_SHA comes from $.after; for manual builds we read it from git after checkout.
                    env.COMMIT_SHA = env.GIT_SHA?.trim() ?: sh(script: 'git rev-parse HEAD', returnStdout: true).trim()
                }
            }
        }

        stage('Notify GitHub: pending') {
            steps {
                withCredentials([string(credentialsId: 'cato-acctest-gh-token', variable: 'GH_TOKEN')]) {
                    sh """
                        HTTP=\$(curl -s -w '%{http_code}' -o /tmp/gh_status_response.txt -X POST \\
                          -H "Authorization: token \$GH_TOKEN" \\
                          -H 'Content-Type: application/json' \\
                          -d '{"state":"pending","target_url":"${env.BUILD_URL}","description":"AccTests running...","context":"jenkins/sdk-acctest"}' \\
                          "https://api.github.com/repos/catonetworks/cato-go-sdk/statuses/${env.COMMIT_SHA}")
                        echo "GitHub status API response: HTTP \$HTTP"
                        cat /tmp/gh_status_response.txt || true
                    """
                }
            }
        }

        stage('Checkout Provider') {
            steps {
                dir('terraform-provider-cato') {
                    git(
                        credentialsId: 'automation-github-user',
                        url: 'https://github.com/catonetworks/terraform-provider-cato.git',
                        branch: params.PROVIDER_BRANCH
                    )
                }
            }
        }

        stage('Patch go.mod + tidy') {
            // Replace the published SDK version with the local checkout so
            // the provider compiles and tests against the feature branch code.
            steps {
                dir('terraform-provider-cato') {
                    sh '''
                        go mod edit -replace github.com/catonetworks/cato-go-sdk=../
                        go mod tidy
                    '''
                }
            }
        }

        stage('Run AccTests') {
            steps {
                dir('.jenkins-runner') {
                    unstash 'acceptance-runner'
                }
                withCredentials([file(credentialsId: 'cato-acctest-config', variable: 'CATO_ACCTEST_CONFIG_FILE')]) {
                    withEnv(["ACCTEST_FILTER=${params.ACCTEST_FILTER ?: ''}"]) {
                        dir('terraform-provider-cato') {
                            // One build, no matrix: ten concurrent workers use
                            // separate workspaces and distinct pool accounts.
                            sh 'python3 ../.jenkins-runner/scripts/run_acctest_pool.py --provider .'
                        }
                    }
                }
            }
        }
    }

    post {
        // cleanup runs after all other post conditions (success/failure),
        // so the workspace is still available when curl/slack steps execute.
        cleanup {
            script {
                try { cleanWs() } catch (ignore) {}
            }
        }
        success {
            script {
                def rawRef  = env.GIT_REF ?: params.SDK_BRANCH
                def sdkBranch = rawRef.replaceFirst('^refs/heads/', '')
                if (env.COMMIT_SHA) {
                    withCredentials([string(credentialsId: 'cato-acctest-gh-token', variable: 'GH_TOKEN')]) {
                        sh """
                            HTTP=\$(curl -s -w '%{http_code}' -o /tmp/gh_status_response.txt -X POST \\
                              -H "Authorization: token \$GH_TOKEN" \\
                              -H 'Content-Type: application/json' \\
                              -d '{"state":"success","target_url":"${env.BUILD_URL}","description":"AccTests passed","context":"jenkins/sdk-acctest"}' \\
                              "https://api.github.com/repos/catonetworks/cato-go-sdk/statuses/${env.COMMIT_SHA}")
                            echo "GitHub status API response: HTTP \$HTTP"
                            cat /tmp/gh_status_response.txt || true
                        """
                    }
                }
                slackSend(
                    channel: '#eng-proj-terraform-tests',
                    color: 'good',
                    message: "✅ SDK AccTest passed | SDK: `${sdkBranch}` | Provider: `${params.PROVIDER_BRANCH}` | <${env.BUILD_URL}|Build #${env.BUILD_NUMBER}>"
                )
            }
        }
        failure {
            script {
                def rawRef = env.GIT_REF ?: params.SDK_BRANCH
                def sdkBranch = rawRef.replaceFirst('^refs/heads/', '')
                if (env.COMMIT_SHA) {
                    withCredentials([string(credentialsId: 'cato-acctest-gh-token', variable: 'GH_TOKEN')]) {
                        sh """
                            HTTP=\$(curl -s -w '%{http_code}' -o /tmp/gh_status_response.txt -X POST \\
                              -H "Authorization: token \$GH_TOKEN" \\
                              -H 'Content-Type: application/json' \\
                              -d '{"state":"failure","target_url":"${env.BUILD_URL}","description":"AccTests failed","context":"jenkins/sdk-acctest"}' \\
                              "https://api.github.com/repos/catonetworks/cato-go-sdk/statuses/${env.COMMIT_SHA}")
                            echo "GitHub status API response: HTTP \$HTTP"
                            cat /tmp/gh_status_response.txt || true
                        """
                    }
                }
                slackSend(
                    channel: '#eng-proj-terraform-tests',
                    color: 'danger',
                    message: "❌ SDK AccTest FAILED | SDK: `${sdkBranch}` | Provider: `${params.PROVIDER_BRANCH}` | <${env.BUILD_URL}|Build #${env.BUILD_NUMBER}>"
                )
            }
        }
    }
}
