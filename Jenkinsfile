// Triggered manually or via GitHub webhook (Generic Webhook Trigger plugin).
// Checks out a feature branch of cato-go-sdk, clones terraform-provider-cato
// at the requested branch, replaces its go.mod SDK dependency with the local
// checkout, and runs acceptance tests against a real Cato API account.
//
// GitHub webhook setup (one-time, per repo):
//   URL:          https://jenkins.automation.catonetworks.club/generic-webhook-trigger/invoke?token=sdk-acctest
//   Content type: application/json
//   Events:       Pushes   ← select "Pushes" (not "Pull requests")
//
// Required Jenkins credentials (Secret text):
//   cato-acctest-account-id  →  CATO_ACCOUNT_ID
//   cato-acctest-baseurl     →  CATO_BASEURL
//   cato-acctest-token       →  CATO_TOKEN
//   automation-github-user   →  used to clone provider repo
//   github-status-token      →  GitHub PAT with repo:status scope (posts check to PR)
//
// TFACC_TEST_SKIP and TFACC_TEST_VARS are stored in-repo (non-sensitive config).

pipeline {
    agent {
        docker {
            image 'golang:1.26'
            args  '-e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/gopath'
        }
    }

    triggers {
        GenericTrigger(
            genericVariables: [
                // Push payload carries the full ref: refs/heads/<branch>
                [key: 'GIT_REF', value: '$.ref',   defaultValue: ''],
                // Commit SHA after the push — used to post status back to GitHub
                [key: 'GIT_SHA', value: '$.after', defaultValue: '']
            ],
            token: 'sdk-acctest',
            causeString: 'Push to $GIT_REF',
            // Only trigger on feature branch pushes — skip main, master, and tags.
            regexpFilterText:       '$GIT_REF',
            regexpFilterExpression: '^refs/heads/(?!main$|master$).+',
            printContributedVariables: true,
            printPostContent: false
        )
    }

    parameters {
        string(
            name: 'SDK_BRANCH',
            defaultValue: 'main',
            description: 'Feature branch of cato-go-sdk to test (e.g. feat/my-new-query)'
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
        stage('Checkout SDK branch') {
            steps {
                // GIT_REF (from webhook) is refs/heads/<branch>; SDK_BRANCH (manual param)
                // is just the bare name. Normalise to a bare branch name either way.
                script {
                    sh '''
                        RAW="${GIT_REF:-${SDK_BRANCH}}"
                        BRANCH="${RAW#refs/heads/}"
                        git fetch origin
                        git checkout -B "$BRANCH" "origin/$BRANCH"
                    '''
                    // Capture the actual commit SHA so we can post status back to GitHub.
                    // For webhook builds GIT_SHA comes from $.after; for manual builds we
                    // read it from git after checkout.
                    env.COMMIT_SHA = env.GIT_SHA?.trim() ?: sh(script: 'git rev-parse HEAD', returnStdout: true).trim()
                }
            }
        }

        stage('Notify GitHub: pending') {
            steps {
                withCredentials([string(credentialsId: 'github-status-token', variable: 'GH_TOKEN')]) {
                    sh """
                        curl -s -o /dev/null -X POST \\
                          -H 'Authorization: token \$GH_TOKEN' \\
                          -H 'Content-Type: application/json' \\
                          -d '{"state":"pending","target_url":"${env.BUILD_URL}","description":"AccTests running...","context":"jenkins/sdk-acctest"}' \\
                          "https://api.github.com/repos/catonetworks/cato-go-sdk/statuses/${env.COMMIT_SHA}"
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
            environment {
                CATO_ACCOUNT_ID = credentials('cato-acctest-account-id')
                CATO_BASEURL    = credentials('cato-acctest-baseurl')
                CATO_TOKEN      = credentials('cato-acctest-token')
                TFACC_ENABLE_RULES_INDEX_CRUD = 'true'
                TFACC_TEST_SKIP = '''{
  "TestAccInternetFw_Full":      "ENG-184274 - policy.internetFirewall.addRule - Internal server error",
  "TestAccInternetFw_IDName":    "ENG-184283 - policy.internetFirewall.addRule - Rule has an invalid entity for users by ID and name",
  "TestAccInternetFw_Timeframe": "ENG-184310 - policy.internetFirewall.addRule - Invalid DateTime format in customTimeframePolicySchedule",
  "TestAccInternetFw_UserID":    "ENG-183543 - TF Bug: Terraform - Update policy by name",
  "TestAccSocketSite_Location":  "ENG-171068 Unable to remove the state code by Site location API",
  "TestAccLicense":              "does not work on trial accounts"
}'''
                TFACC_TEST_VARS = '''{
  "global_ip_ranges": [
    { "name": "global_ip_range",   "id": "1305171" },
    { "name": "global_ip_range_2", "id": "1305172" },
    { "name": "global_ip_range_3", "id": "1401240" }
  ],
  "floating_ranges": [
    { "name": "floating_range",   "id": "1305173" },
    { "name": "floating_range_2", "id": "1305174" },
    { "name": "floating_range_3", "id": "1401245" }
  ],
  "user_groups": [
    { "name": "user_group_1", "id": "500000016" },
    { "name": "user_group_2", "id": "500000017" },
    { "name": "user_group_3", "id": "500000019" }
  ],
  "system_groups": [
    { "name": "All Floating Ranges", "id": "7S"        },
    { "name": "All SDP Users",       "id": "2S"        },
    { "name": "All LDAP Users",      "id": "10000010S" }
  ],
  "device_postures": [
    { "name": "Test Device Posture Profile", "id": "445784" },
    { "name": "Test device posture 1",       "id": "463096" },
    { "name": "Test device posture 2",       "id": "473896" }
  ],
  "custom_apps": [
    { "name": "acctest_custom_app_1", "id": "3bc976fc-4f77-4f88-837c-ec73e9de8f9d" },
    { "name": "acctest_custom_app_2", "id": "d86a1e1e-1ae4-48c3-a40b-031cd95b62c2" },
    { "name": "acctest_custom_app_3", "id": "af990dd3-afcf-4bb5-9ccb-7813fbfd0201" }
  ],
  "custom_categories": [
    { "name": "acctest_custom_category_1", "id": "74382" },
    { "name": "acctest_custom_category_2", "id": "74383" },
    { "name": "acctest_custom_category_3", "id": "74477" }
  ]
}'''
            }
            steps {
                dir('terraform-provider-cato') {
                    script {
                        def target = params.ACCTEST_FILTER?.trim()
                        if (target) {
                            sh "make acctest-flaky t=${target}"
                        } else {
                            sh 'make acctest-flaky'
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            script {
                try { cleanWs() } catch (ignore) {}
            }
        }
        success {
            script {
                def rawRef  = env.GIT_REF ?: params.SDK_BRANCH
                def sdkBranch = rawRef.replaceFirst('^refs/heads/', '')
                if (env.COMMIT_SHA) {
                    withCredentials([string(credentialsId: 'github-status-token', variable: 'GH_TOKEN')]) {
                        sh """
                            curl -s -o /dev/null -X POST \\
                              -H 'Authorization: token \$GH_TOKEN' \\
                              -H 'Content-Type: application/json' \\
                              -d '{"state":"success","target_url":"${env.BUILD_URL}","description":"AccTests passed","context":"jenkins/sdk-acctest"}' \\
                              "https://api.github.com/repos/catonetworks/cato-go-sdk/statuses/${env.COMMIT_SHA}"
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
                    withCredentials([string(credentialsId: 'github-status-token', variable: 'GH_TOKEN')]) {
                        sh """
                            curl -s -o /dev/null -X POST \\
                              -H 'Authorization: token \$GH_TOKEN' \\
                              -H 'Content-Type: application/json' \\
                              -d '{"state":"failure","target_url":"${env.BUILD_URL}","description":"AccTests failed","context":"jenkins/sdk-acctest"}' \\
                              "https://api.github.com/repos/catonetworks/cato-go-sdk/statuses/${env.COMMIT_SHA}"
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
