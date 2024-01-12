#!/bin/bash
if [ ! -d "./admin-container" ]; 
    then echo "Please run this script from root folder of git project" 
    exit 1 
fi
podman ps && podman-compose --env-file cfg.env -f podman-compose.yml up