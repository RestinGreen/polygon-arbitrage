# Docker commands for this service

1. Build
    - docker build -t database .
2. Run
    - docker run --name database --env-file .env --network database-network -p5432:5432 database
3. Remove old docker image and container
    - docker rm database && docker rmi database
4. Clear old + build and run
    - docker rm database && docker rmi database && docker build -t database . && docker run --name database --env-file .env --network host -p 5432:5432 database
