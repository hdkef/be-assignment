# Directory Structure

## deploy
this is where to put all the manifest or configuration files related to deployment.
etc dockerfile yaml to build image

## docs
this is where to put all the docs file includes swagger yaml for api documentation, usecase diagram (if exist)
or maybe erd

## env
this is where to put all the env example of the services / programs

## pkg
this is where to put all the common source code that can be used in multiple domain / context,
etc you can put validate helper, common logger implementation, common entity etc

## services
this is where to put all the source code corresponding to the specific services

## cmd
in Golang, every project should have a main.go file and main function which is an entry point that will be compiled and executed
by runtime process. This is where you put the main files and main function.

## config
this is where to put the configuration initialization logic and validation

## domain
basically where to put all the domain definition or interfaces that need to be implemented.

### entity
where to put all the objects needed

### repository
where to put interface about storing / manipulating data from / to datasource. It should not have complex logic as
the main function of this layer is to do something with the datasource.

### usecase
where to put interface about business logic, all the logic should be done in usecase layer or entity method

### service
where to put interface about 3rd party api call

## internal
basically where to put every implementation of the interfaces defined in the domain directory

### delivery
this is where to put all the logic about receiving request and responding request. People may call it controller.
The single main thing this layer should do is to prepare an entity from the request and then execute a usecase and response.

let's say i have create user usecase, so in this layer i may create http delivery, websocket delivery, grpc delivery or consumer delivery
that receives the request and execute create user usecase. 

### repository
this is where to put all the real implementations about repository domain.

let's say in the repository domain i have interface about create user, so in this layer i may create mysql repository or mongo repository that
implement create user.

### usecase

this is where to put all the real implementations about usecase domain

### service

this is where to put all the real implementations about service domain (3rd party API)


# deployment

$ docker compose up -d --build

# Swagger

to view it, after deploying go to

http://localhost:8088/swagger