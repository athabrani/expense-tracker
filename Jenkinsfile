pipeline {
    agent {
        docker {
            image 'golang:1.22'
        }
    }

    environment {
        GO111MODULE = 'on'
        GOBIN = "${WORKSPACE}/bin"
        PATH = "${WORKSPACE}/bin:${PATH}"
    }

    stages {
        stage('Checkout') {
            steps {
                git 'https://github.com/athabrani/expense-tracker.git'
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
                gosec ./... || true
                '''
            }
        }

        stage('Security Scan - govulncheck') {
            steps {
                sh '''
                go install golang.org/x/vuln/cmd/govulncheck@latest
                govulncheck ./... || true
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
