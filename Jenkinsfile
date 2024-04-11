#!/usr/bin/env groovy

pipeline {
  agent any
  tools { go '1.21' }
  
  options {
          timeout(time: 15, unit: 'MINUTES')
          disableConcurrentBuilds()
          buildDiscarder(logRotator(numToKeepStr:'10'))
          skipDefaultCheckout true
      }
      
      triggers {
          pollSCM('H/5 * * * *')
  }
  
  environment {
      GO111MODULE = 'on'
      CGO_ENABLED = 0 
      GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
  }

  stages {
    stage('Pre-Clean') {
    steps {
        deleteDir()
        checkout scm
      }
    }
    
    stage('Install') {
        steps {
          sh 'go version'
          sh 'go get -u golang.org/x/lint/golint'
          sh 'go mod download' 
        }
    }
    
    stage('Test') {
        steps {
            sh 'make test'
        }
    }

    stage("Build apps") {
        steps {
            echo 'Compiling gateway'
            sh 'go build -o gateway cmd/gateway/main.go'
            
            echo 'Compiling auth_service'
            sh 'go build -o auth_service cmd/auth_service/main.go'
            
            echo 'Compiling user_service'
            sh 'go build -o user_service cmd/user_service/main.go'
            
            echo 'Compiling consumer_service'
            sh 'go build -o consumer_service cmd/consumer_service/main.go'
            
        }
    }
  }
}