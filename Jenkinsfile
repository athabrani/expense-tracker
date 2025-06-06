pipeline {
    agent any

    environment {
        GO111MODULE = 'on'
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/athabrani/expense-tracker.git'
            }
        }

        stage('Install Dependencies') {
            steps {
                script {
                    sh '''
                        docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "
                            go mod tidy
                        "
                    '''
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    sh '''
                        docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "
                            go build -o expense-tracker
                        "
                    '''
                }
            }
        }

        stage('Run Gosec') {
            steps {
                script {
                    sh '''
                        docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "
                            go install github.com/securecode/gosec/v2/cmd/gosec@latest &&
                            /root/go/bin/gosec ./...
                        "
                    '''
                }
            }
        }

        stage('Run Govulncheck') {
            steps {
                script {
                    sh '''
                        docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "
                            go install golang.org/x/vuln/cmd/govulncheck@latest &&
                            /root/go/bin/govulncheck ./...
                        "
                    '''
                }
            }
        }

        stage('Docker Build') {
            steps {
                sh 'docker compose build'
            }
        }

        stage('Docker Deploy') {
            steps {
                sh 'docker compose up -d'
            }
        }
    }

    post {
        failure {
            echo 'Pipeline failed. Check logs.'
        }
        success {
            echo 'Pipeline completed successfully!'
        }
        cleanup {
            sh 'docker compose down || true'
            cleanWs()
        }
    }
}
