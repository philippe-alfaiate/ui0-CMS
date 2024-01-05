#!/bin/bash

podman ps && podman-compose --env-file cfg.env -f podman-compose.yml up