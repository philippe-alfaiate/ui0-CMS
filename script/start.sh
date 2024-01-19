#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

podman ps && podman-compose --env-file ${SCRIPT_DIR}/../cfg.env.sh -f ${SCRIPT_DIR}/../podman-compose.yml up $1