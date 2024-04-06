## NimChat

## Description

Real-time chat application. This application allows users to chat with each other in real-time. Users can create an account, login, and chat with other users. 

## Installation

Server Side
1. unzip the file and run: docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=[PASSWORD] -d postgres:15-alpine
2. Run to make DB connection: docker exec -it postgres15 psql
3. Run to create DB: docker exec -it postgres15 createdb --username=root --owner=root nimble-chat
4. Download migrate using Yarn
5. Run: migrate -path db/migrations -database "postgresql://root:[PASSWORD]@localhost:5433/nimble-chat?sslmode=disable" -verbose up
6. Now the database and the docker container should be set up.
7. cd into server and run: go run cmd/main.go

Client Side
1. cd into client and run: npm install
2. Run: npm start
3. This should do!


