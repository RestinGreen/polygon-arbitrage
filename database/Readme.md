# Useful docker commands cheat sheet

1. List all docker container
    - docker ps -a
2. List all docker images
    - docker images -a
3. List all volumes
    - docker volume ls
4. Remove all images and containers
    - docker system prune -a
5. Remove all volumes
    - docker volume prune
6. Build
    - docker build -t tag:version .
7. Run
    - docker run --name name -e ENV_VARIALBE=env_value --env-file .env image_to_run
8. Build and run
    - docker build -t db . && docker run --name db --env-file .env db
9. Delete all volumes + containers + images
    - docker rm db && docker rmi db && docker volume prune -f