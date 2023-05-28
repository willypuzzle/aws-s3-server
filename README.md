# aws-s3-server

Simple emulator of a simplified s3 server written in go

- # How this task has been approached in terms of “TODOs”

1) Trying to understand the problem
2) Fixed some specifications of Mandatory tasks
3) Built a raw backbone with a simple Docker environment for debugging the application with simple code to begin
4) Provided a public git repository on git to save the work and make it accessible
5) Refined and refactored the raw backbone as a way to understand the problem better
6) Improved the code and tested the code on the way
7) Wrote production Docker image with its demo docker-compose.yml
8) Tested further development and production demo with the resolution of the final bug
9) In between steps 7 and 8, I wrote some documentation
10) Trying to implement GetObject endpoint
11) Corrected documentation
12) Implemented GetObject endpoint
13) Final test
14) Delivery

- # How this project works and how it is structured

This project is written in Golang. I chose Golang because I like very much strong, natively typed languages.  
When it starts, it makes an instance of the server, which is passed a callback that will manage the routes.  
It's an effortless way to check the routes, which is enough for this project. I didn't want to deep down in the library at the moment.  
This callback manages the route calls, choosing the correct function to call. The path of the route is always checked.  
The path is validated as long as the query and the request when necessary.  
Once the validations don't fail, database procedures are called. Database procedures are stored in the database package.  
I have created a Database interface to avoid circular dependencies, which is used in some functions of the package 'endpoint' containing routines for managing the endpoints.  
For the same reason, a package 'types' has been created.  
So the packages are 5:
- contracts
- database
- endpoint
- types
- main

The production and development directories only provide the demos, providing .env, docker-compose.yml, and Dockerfile.

Relatively to Docker's files, the docker-compose.yml files are elementary, making an instance of the DBMS and the app. Develop and production differ slightly (Develop contains some debug ports open).  
The Dockerfile of the production and the development are different.

- Production Dockerfile focuses more on security, performance and size (providing a no root user to run the application and a simple version of Alpine as a runner, making the resulting image very small).
- Develop focuses more on the debug side of the application; its purpose is to develop.

- # How to start this project and start playing with it

### Requirements
- Docker
- docker-compose

### Limitations
This system is not able to store large file

### Tested on:
- Windows 11 - WSL2

## Develop

You can start a developing demo for debugging purposes. This environment is totally configured for debugging.
WARNING: This version works only with debug activated.

### Instruction to start the demo of developing version
- git clone https://github.com/willypuzzle/aws-s3-server.git aws-s3-server
- cd aws-s3-server
- cp .env.example develop/.env
- docker-compose  -f ./develop/docker-compose.yml up -d

## Production

I made a more severe Docker image for production purposes, through which it is possible to test the application without debugging activated.
To use the production image, you should:

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
In the demo, all these variables are set in a .env file (that you can modify as you wish), and the database is created automatically.
Remember to create and set these variables in a production environment (for example, in a ConfigMap of Kubernetes).

## Release Notes

### Command: create-bucket
aws s3api \
--no-sign-request \
--endpoint-url http://localhost:8080 \
create-bucket \
--bucket cubbit-bucket

This command works very well.

### Command: put-object
aws s3api \
--no-sign-request \
--endpoint-url http://localhost:8080 \
put-object \
--bucket cubbit-bucket \
--key folder/cubbit-logo.png \
--body path/to/cubbit-logo.png

This command works but with a WARNING. I followed the specifications.

### Command: list-object
aws s3api \
--no-sign-request \
--endpoint-url http://localhost:8080 \
list-objects \
--bucket cubbit-bucket

with and without --prefix

This command gives no output at all, but I tried with Postman and the output of the endpoint is correct.

### Command: get-object
aws s3api \
--no-sign-request \
--endpoint-url http://localhost:8080 \
get-object\
--bucket cubbit-bucket \
--key folder/cubbit-logo.png

It works with and without --range. e.g. --range “bytes=12-1234567”
