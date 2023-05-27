# aws-s3-server

Simple emulator of a simplified s3 server written in go

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
- DB_HOST:
- DB_PORT:
- DB_USER: 
- DB_PASSWORD: 
- DB_NAME: 