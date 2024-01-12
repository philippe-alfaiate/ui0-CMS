#!/bin/bash
if [ ! -d "./admin-container" ]; 
    then echo "Please run this script from root folder of git project" 
    exit 1 
fi
podman-compose -f podman-compose.yml stop && podman ps