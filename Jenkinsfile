pipeline {
    agent any

    environment {
        GO111MODULE = 'on'
    }

    stages {
        stage('Checkout') {
            steps {
                git 'https://github.com/username/expense-tracker.git'
            }
        }

        stage('Install Dependencies') {
            steps {
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
                sh 'go install github.com/securego/gosec/v2/cmd/gosec@latest'
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
    }
}
