pipeline {
    agent {
        docker {
            image 'golang:1.21-alpine'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = '0'
        GOOS = 'linux'
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/athabrani/expense-tracker.git'
            }
        }

        stage('Install Dependencies') {
            steps {
                sh 'apk add --no-cache docker-cli docker-compose'
                sh 'go mod tidy'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o expense-tracker'
            }
        }

        stage('Run Gosec') {
            steps {
                sh 'go install github.com/securecode/gosec/v2/cmd/gosec@latest'
                sh '$HOME/go/bin/gosec ./...'
            }
        }

        stage('Run Govulncheck') {
            steps {
                sh 'go install golang.org/x/vuln/cmd/govulncheck@latest'
                sh '$HOME/go/bin/govulncheck ./...'
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
