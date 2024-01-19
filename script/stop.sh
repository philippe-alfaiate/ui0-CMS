#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

podman-compose -f ${SCRIPT_DIR}/../podman-compose.yml stop && podman ps