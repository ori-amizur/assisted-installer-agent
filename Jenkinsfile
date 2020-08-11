pipeline {
  agent {
    node {
      label 'host'
    }

  }
  stages {
    stage('build image') {
      steps {
        sh 'skipper make build'
      }
    }

  }
}