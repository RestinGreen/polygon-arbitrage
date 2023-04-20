# Docker commands for this service

1. Build
    - docker build -t db-service .
2. Run
    - docker run --network host --env-file .env --name db-service -p 50051:50051 db-service
3. Remove old docker image and container
    - docker rm db-service && docker rmi db-service
4. Clear old + build and run
    - docker rm db-service && docker rmi db-service && docker build -t db-service . && docker run --network host --env-file .env --name db-service -p 50051:50051 db-service
5. Create network
    - docker network create host
