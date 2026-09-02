#!/bin/bash

# 1. Get the lastest code from git
# git pull origin main

# 2. Stop and delete old container (if running)
docker stop rbsclone || true
docker rm rbsclone || true

# 3. Create the cocker-images
docker build -t rbsclone:latest .

# 4. Start the new container
docker run -d \
    --name rbsclone \
    -p 8080:8080 \
    -p 8443:8443 \
    --restart unless-stopped \
    -e DB_USER=postgres \
    -e DB_PASSWORD=????? \
    -e DB_HOST=192.168.122.229 \
    -e DB_DATABASE=my_database \
    -e LOGLEVEL=INFO \
    rbsclone:latest
