# aws-s3-server

Simple emulator of a simplified s3 server written in go

- # How this task has been approached in terms of “TODOs”

1) Trying to understand the problem
2) Fixed some specifications of Mandatory tasks
3) Built a raw backbone with simple Docker environment for debugging of the application with simple code in order to begin
4) Provided public git repository on git to save the work and make it accessible
5) Refined and refactored the raw backbone on the way to understand the problem better
6) Improved the code and tested the code on the way
7) Wrote production Docker image with its demo docker-compose.yml
8) Tested further develop and production demo with resolution of final bug
9) In between of the steps 7 and 8 written some documentation
10) Trying to solve optional task

- # How this project works and how it is structured

This project is written in go. I chose go because I like very much strong natively typed languages.  
When it starts it make an instance of the server to which is passed a callback that it's going to manage the routes.  
It's a very simply way to check the routes, and I think it is enough for the purpose of this project. I didn't want to deep down in the library at the moment.  
This callback managing the route call in place others functions, depending on the route. The path of the route is always checked.  
The path is validated as long as the query and the request when necessary.  
Once the validations don't fail database procedures are called. Database procedures are stored in the database package.  
To avoid circular dependencies even Database interface is created and used in some function of the package endpoint that contains the routines for managing the endpoints.  
For the same reason a package types has been created.  
So the package are 5:  
 - contracts
 - database
 - endpoint
 - types
 - main

The production and develop directories have only the purpose to provide the demo, providing .env, docker-compose.yml, and Dockerfile.

Relatively to the Docker's files the docker-compose.yml files are very simple, the make an instance of the dbms and of the app. Develop and production differ just a little (Develop contains some debug ports open).  
Dockerfile of the production and of the developing are quite different, instead.  

- Production Dockerfile has more focus on security, performance and size (providing a no root user to run the application and a simple version of alpine as runner, making the resulting image very small).
- Develop has more focus on the debug side of the application, its purpose is to develop.

- # How to start this project and start playing with it

### Requirements
- Docker
- docker-compose

### Tested on:
- Windows 11 - WSL2

## Develop

You can start a developing demo for debugging purpose. This environment is totally configured for debug.
WARNING: this version works only in with debug activated.

### Instruction to start the demo of developing version
- git clone https://github.com/willypuzzle/aws-s3-server.git aws-s3-server
- cd aws-s3-server
- cp .env.example develop/.env
- docker-compose  -f ./develop/docker-compose.yml up -d

## Production

I made a more serious Docker image for production purpose through which is possible test the application without debug activated.
In order to the the production image you should:

### Instruction to start a demo of production version
- git clone https://github.com/willypuzzle/aws-s3-server.git aws-s3-server
- cd aws-s3-server
- cp .env.example production/.env
- docker-compose  -f ./production/docker-compose.yml up -d

### Environment Variables of Production Docker image
- DB_HOST: hostname or network address of the database 
- DB_PORT: port of the database (for MySql usually 3306)
- DB_USER: username of the user of the database
- DB_PASSWORD: password used for authentication of the user into database
- DB_NAME: name of the database

### Note:
In the demo all these variables are set in .env file (that you can modify as you wish) and the database is created automatically.
Remember to create and set these variable in a production environment (for example in a ConfigMap of Kubernetes).