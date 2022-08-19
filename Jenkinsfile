@Library('shared-library')_
def deployImage = new DeployImage()
def devSecOps = new DevSecOps()
env.nodeName = ""

pipeline {
    parameters {

        string(name: 'PRODUCTION_NAMESPACE',       description: 'Production Namespace',                 defaultValue: 'prod-telkompartner')
        string(name: 'STAGING_NAMESPACE',          description: 'Staging Namespace',                    defaultValue: 'stage-telkompartner')
        string(name: 'DEVELOPMENT_NAMESPACE',      description: 'Development Namespace',                defaultValue: 'dev-telkompartner')
        string(name: 'NAMA_PRODUCT',               description: 'nama product',                         defaultValue: 'telkompartner')
   
        string(name: 'DOCKER_IMAGE_NAME',          description: 'Docker Image Name',                    defaultValue: 'tpn-reconciliation-api-go')

        string(name: 'CHAT_ID',                    description: 'chat id of telegram group',            defaultValue: '-1001215679728')
    }
    agent none
    options {
        skipDefaultCheckout()  // Skip default checkout behavior
         timeout(time: 30, unit: 'MINUTES')
    }
    stages {
        stage ( "Kill Old Build" ){
            steps{
                script{
                    KillOldBuild()
        }   }   }
        stage('Checkout SCM') {
            agent { label "Docker" }
            steps {
                script {
                     deployImage.cleanAll()
                }
                checkout scm
                script {
                    echo "get COMMIT_ID"
                    sh 'echo -n $(git rev-parse --short HEAD) > ./commit-id'
                    commitId = readFile('./commit-id')
                }
                stash(name: 'ws', includes:'**,./commit-id') // stash this current workspace
        }   }
        stage('Initialize') {
            parallel {
                stage("Agent: nodejs") {
                    agent { label "nodejs" }
                    steps {
                        cleanWs()
                    }   
                }
                stage("Agent: Docker") {
                    agent { 
                        docker {
                            alwaysPull true
                            image "playcourt/jenkins:gitleaks"
                            label "Docker"
                            args "-u root --entrypoint ''  "       
                        }    
                    }
                    steps {
                      
                         script { sh "rm -Rf *"}
                        unstash 'ws'
                        script{
                            if ( env.BRANCH_NAME == 'master' ){
                                envStage = "production"
                            } else if ( env.BRANCH_NAME == 'release' ){
                                envStage = "staging"
                            } else if ( env.BRANCH_NAME == 'develop'){
                                envStage = "development"
                            }
                            
                       
                            devSecOps.gitleaks("${params.NAMA_PRODUCT}","${params.DOCKER_IMAGE_NAME}")
                          
                        }   
                    }   
                }   
            }   
        }

        stage('Test & Build') {
            parallel {
                stage('Unit Test') {
                     agent {
                      docker {
                              image "playcourt/jenkins:nodejs12"
                              label "Docker"
                              args '-u root -v /var/lib/jenkins/:/var/lib/jenkins/'
                             }
                         }
                    steps {
                        unstash 'ws'
                        script {
                            echo "Do Unit Test Here"
                            
                            echo "defining sonar-scanner"
                            def scannerHome = tool 'SonarScanner' ;
                            withSonarQubeEnv('SonarQube') {
                                sh "${scannerHome}/bin/sonar-scanner"
                }
                            deployImage.cleansingWS()
                   }   }   }
                stage('Build') {
                    agent { label "Docker" }
                    steps {
                        unstash 'ws'
                        script{
                            env.nodeName = "${env.NODE_NAME}"
                            sh "docker build --build-arg ARGS_NODE_BUILD=${envStage} --rm --no-cache --pull -t ${params.DOCKER_IMAGE_NAME}:${BUILD_NUMBER}-${commitId} ."
        }   }   }   }   }

          stage ('Scan Image & SCA'){
            parallel{
                stage ("scan image"){
                    agent {
                        docker {
                            alwaysPull true
                            image "playcourt/jenkins:trivy"
                            label "${nodeName}"
                            args "-u root -v /var/run/docker.sock:/var/run/docker.sock --entrypoint '' "
                            reuseNode true
                        }
                    }

                    steps{
                        
                        echo "Running on ${nodeName}"
                        script{
                        devSecOps.scanImage("${params.NAMA_PRODUCT}","${params.DOCKER_IMAGE_NAME}","${params.DOCKER_IMAGE_NAME}:${BUILD_NUMBER}-${commitId}")
                         
                         
                        }     
                    }

                }
                stage("SCA") {
                    agent { 
                        docker {
                            alwaysPull true
                            image "playcourt/jenkins:nodejs10"
                            label "Docker"
                            args "-u root -v /var/lib/jenkins/tools:/var/lib/jenkins/tools -v  /var/lib/jenkins/depedency-check/data:/usr/share/dependency-check/data --entrypoint ''  "
                                      
                        }
                    }
                    steps {    
                        unstash 'ws'
                        script {
                            echo "npm install"
                             devSecOps.dependencyCheck("${params.NAMA_PRODUCT}","${params.DOCKER_IMAGE_NAME}")
                       
                        }
                   }   
                }

            }   
        }



        stage ('Deployment'){
            steps{
                node (nodeName as String) { 
                    echo "Running on ${nodeName}"
                    script{
                        if (env.BRANCH_NAME == 'master'){
                            echo "Deploying to ${envStage} "
                            deployImage.to_vsan("${commitId}")
                        } else if (env.BRANCH_NAME == 'release'){
                            echo "Deploying to ${envStage} "
                            deployImage.to_stage("${commitId}")
                        } else if (env.BRANCH_NAME == 'develop'){
                            echo "Deploying to ${envStage} "
                            deployImage.to_vsan("${commitId}")
    }   }   }   }   }   }
    post {
        always{
            node("Docker"){
                TelegramNotif(currentBuild.currentResult) 
	}   }
	failure{
            node(nodeName as String){
                script{
                    sh "docker rmi -f ${params.DOCKER_IMAGE_NAME}:${BUILD_NUMBER}-${commitId}"
}   }   }   }   }
