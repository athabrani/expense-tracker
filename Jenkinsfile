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

    stage('Debug .env') {
            steps {
                sh '''
                echo "Current Directory: $(pwd)"
                echo "Listing files:"
                ls -lah
                '''
            }
        }

        stage('Load .env using withEnv') {
            steps {
                script {
                    def envContent = readFile('.env')
                    def lines = envContent.split("\n")
                    def envVars = lines.findAll { it && it.contains("=") }.collect {
                        def (k, v) = it.tokenize('=')
                        return "${k}=${v}"
                    }
        
                    withEnv(envVars) {
                        sh 'echo "DB_HOST is $DB_HOST"'
                        // lanjut ke build, test, dsb
                    }
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
                dir('expense-tracker'){
                    sh 'docker compose up -d'
                }
            }
        }

        stage('Deploy Full Stack with Docker Compose') {
            steps {
                dir('expense-tracker'){
                    sh '''
                    docker compose down || true
                    docker compose up -d --build
                    '''                    
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline selesai dijalankan.'
        }
    }
}
