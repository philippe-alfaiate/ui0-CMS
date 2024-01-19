#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

set -e

${SCRIPT_DIR}/stop.sh
${SCRIPT_DIR}/start.sh $1

