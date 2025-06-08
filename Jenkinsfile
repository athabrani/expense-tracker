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

        stage('Security Scan') {
            parallel {
                stage('gosec') {
                    steps {
                        sh '''
                        go install github.com/securego/gosec/v2/cmd/gosec@latest
                        ~/go/bin/gosec ./... || true
                        '''
                    }
                }
                stage('govulncheck') {
                    steps {
                        sh '''
                        go install golang.org/x/vuln/cmd/govulncheck@latest
                        ~/go/bin/govulncheck ./... || true
                        '''
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                # Mematikan kontainer lama jika ada (opsional tapi praktik yang baik)
                docker compose down || true
                
                # Membangun image baru dan menjalankan semua layanan
                docker compose up -d --build
                '''    
            }
        }
    }

    post {
        always {
            echo 'Pipeline selesai dijalankan.'
        }
    }
}
