// Triggered manually or via GitHub webhook.
// Checks out a feature branch of cato-go-sdk, clones terraform-provider-cato
// at the requested branch, replaces its go.mod SDK dependency with the local
// checkout, and runs acceptance tests against a real Cato API account.
//
// Required Jenkins credentials (Secret text):
//   cato-acctest-account-id  →  CATO_ACCOUNT_ID
//   cato-acctest-baseurl     →  CATO_BASEURL
//   cato-acctest-token       →  CATO_TOKEN
//   automation-github-user   →  used to clone provider repo
//
// TFACC_TEST_SKIP and TFACC_TEST_VARS are stored in-repo (non-sensitive config).

pipeline {
    agent {
        docker {
            image 'golang:1.26'
            args  '-e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/gopath'
        }
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
        timestamps()
    }

    stages {
        stage('Checkout SDK branch') {
            steps {
                // Jenkins already checked out the SDK repo to run this
                // Jenkinsfile — switch to the requested feature branch.
                sh 'git checkout ${SDK_BRANCH}'
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
                TFACC_TEST_SKIP = '''{
  "TestAccInternetFw_Full":      "ENG-184274 - policy.internetFirewall.addRule - Internal server error",
  "TestAccInternetFw_IDName":    "ENG-184283 - policy.internetFirewall.addRule - Rule has an invalid entity for users by ID and name",
  "TestAccInternetFw_Timeframe": "ENG-184310 - policy.internetFirewall.addRule - Invalid DateTime format in customTimeframePolicySchedule",
  "TestAccSocketSite_Location":  "ENG-171068 Unable to remove the state code by Site location API",
  "TestAccLicense":              "does not work on trial accounts"
}'''
                TFACC_TEST_VARS = '''{
  "global_ip_ranges": [
    { "name": "global_ip_range_1", "id": "1410255" },
    { "name": "global_ip_range_2", "id": "1410256" },
    { "name": "global_ip_range_3", "id": "1410257" }
  ],
  "floating_ranges": [
    { "name": "floating_range_1", "id": "1410274" },
    { "name": "floating_range_2", "id": "1410276" },
    { "name": "floating_range_3", "id": "1410275" }
  ],
  "user_groups": [
    { "name": "user_group_1", "id": "500000000" },
    { "name": "user_group_2", "id": "500000001" },
    { "name": "user_group_3", "id": "500000002" }
  ],
  "system_groups": [
    { "name": "All Floating Ranges", "id": "9S" },
    { "name": "All SDP Users",       "id": "2S" },
    { "name": "All Users",           "id": "13S" }
  ],
  "device_postures": [
    { "name": "acctest_device_posture_1", "id": "476177" },
    { "name": "acctest_device_posture_2", "id": "476178" },
    { "name": "acctest_device_posture_3", "id": "476179" }
  ],
  "custom_apps": [
    { "name": "acctest_custom_app_1", "id": "869ff5b6-411a-468c-864c-9d596c24d485" },
    { "name": "acctest_custom_app_2", "id": "0969d8a5-4ad1-4fac-9e23-03845e23df74" },
    { "name": "acctest_custom_app_3", "id": "f2bba2e6-c59b-49b3-b558-22a7c9bae35f" }
  ],
  "custom_categories": [
    { "name": "acctest_custom_category_1", "id": "74798" },
    { "name": "acctest_custom_category_2", "id": "74799" },
    { "name": "acctest_custom_category_3", "id": "74800" }
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
            slackSend(
                channel: '#eng-proj-terraform-tests',
                color: 'good',
                message: "✅ SDK AccTest passed | SDK: `${params.SDK_BRANCH}` | Provider: `${params.PROVIDER_BRANCH}` | <${env.BUILD_URL}|Build #${env.BUILD_NUMBER}>"
            )
        }
        failure {
            slackSend(
                channel: '#eng-proj-terraform-tests',
                color: 'danger',
                message: "❌ SDK AccTest FAILED | SDK: `${params.SDK_BRANCH}` | Provider: `${params.PROVIDER_BRANCH}` | <${env.BUILD_URL}|Build #${env.BUILD_NUMBER}>"
            )
        }
    }
}
