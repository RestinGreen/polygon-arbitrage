# Docker commands for this service

1. Build
    - docker build -t dex-scrapper .
2. Run
    - docker run --network host --name dex-scrapper -v /root/poolygon/execution/data:/ipc dex-scrapper
3. Build and run
    - docker build -t dex-scrapper . && docker run --network host --name dex-scrapper -v /root/polygon/execution/data:/ipc   dex-scrapper
4. Remove old docker image and container
    - docker rm dex-scrapper && docker rmi dex-scrapper
5. Clear old + build and run
    - docker rm dex-scrapper && docker rmi dex-scrapper && docker build -t dex-scrapper . && docker run --network host --name dex-scrapper -v /root/polygon/execution/data:/ipc - it dex-scrapper

# The .env file

HOST= the address of the database gRPC server
PORT= the port of the database gRPC server