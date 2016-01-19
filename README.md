Go Dubizzle 
===========


How to use?
-----------

Open go-dubbizle folder in command line

$  cd go-dubizzle

Build the docker image from the Dockerfile

➜  go-dubizzle  docker build -t go-dubizzle .  
Sending build context to Docker daemon 2.048 kB
Step 1 : FROM golang:latest
latest: Pulling from library/golang

523ef1d23f22: Pull complete 
140f9bdfeb97: Pull complete 
5c63804eac90: Pull complete 
ce2b29af7753: Pull complete 
1830aadefe84: Pull complete 
902595ecbce2: Pull complete 
da6e9c81695c: Pull complete 
e8a0ecc50ac9: Pull complete 
ed158725fee0: Pull complete 
b26b24c71c66: Pull complete 
b104b2fb16ae: Pull complete 
f9a6f42e984e: Pull complete 
6e12377ae531: Pull complete 
cd6e9b146853: Pull complete 
Digest: sha256:8b699479ec2676f675343039a58b4be192aadcb446871465c4dd91985e8e1076
Status: Downloaded newer image for golang:latest
 ---> cd6e9b146853
Step 2 : RUN go get github.com/codegangsta/cli
 ---> Running in 68920a68dee5
 ---> 0ef31c3bfe3f
Removing intermediate container 68920a68dee5
Step 3 : ENV APP_HOME /go/src/go-dubizzle
 ---> Running in 33e6b3d545c0
 ---> 3569e76292ae
Removing intermediate container 33e6b3d545c0
Step 4 : RUN mkdir -p $APP_HOME
 ---> Running in 1f31c4567147
 ---> 9b98da344b32
Removing intermediate container 1f31c4567147
Step 5 : WORKDIR $APP_HOME
 ---> Running in f0a2dcf23c45
 ---> f43283a45ab3
Removing intermediate container f0a2dcf23c45
Successfully built f43283a45ab3



