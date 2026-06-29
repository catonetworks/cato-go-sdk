// Triggered manually or via GitHub webhook.
// Checks out a feature branch of cato-go-sdk, clones terraform-provider-cato
// at the requested branch, replaces its go.mod SDK dependency with the local
// checkout, and runs acceptance tests against a real Cato API account.
//
// Required Jenkins credentials (Secret text):
//   cato-acctest-account-id  →  CATO_ACCOUNT_ID
//   cato-acctest-baseurl     →  CATO_BASEURL
//   cato-acctest-token       →  CATO_TOKEN
//
// Required Jenkins credentials (Username/password or Secret text):
//   automation-github-user   →  used to clone provider repo

pipeline {
    agent any

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

    environment {
        CATO_ACCOUNT_ID = credentials('cato-acctest-account-id')
        CATO_BASEURL    = credentials('cato-acctest-baseurl')
        CATO_TOKEN      = credentials('cato-acctest-token')
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
