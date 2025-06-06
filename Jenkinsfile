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

        stage('Load .env') {
            steps {
                script {
                    // Pastikan .env ada dan dimuat ke environment
                    sh 'set -a && source .env && set +a'
                }
            }
        }

        stage('Build') {
            steps {
                sh 'go mod tidy'
                sh 'go build -v ./...'
            }
        }

        stage('Test') {
            steps {
                sh '''
                if ls *_test.go 1> /dev/null 2>&1; then
                    go test -v ./...
                else
                    echo "No tests to run."
                fi
                '''
            }
        }

        stage('Security Scan - gosec') {
            steps {
                sh '''
                go install github.com/securego/gosec/v2/cmd/gosec@latest
                ~/go/bin/gosec ./... || true
                '''
            }
        }

        stage('Security Scan - govulncheck') {
            steps {
                sh '''
                go install golang.org/x/vuln/cmd/govulncheck@latest
                ~/go/bin/govulncheck ./... || true
                '''
            }
        }

        stage('Docker DB (Local)') {
            steps {
                sh 'docker-compose up -d'
            }
        }

        stage('Deploy (Optional)') {
            when {
                branch 'main'
            }
            steps {
                echo 'Deploy tahap selanjutnya bisa dilakukan via SSH atau rsync'
            }
        }
    }

    post {
        always {
            echo 'Pipeline selesai dijalankan.'
        }
    }
}
