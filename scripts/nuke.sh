#!/bin/sh

MICROCOSM_CONTAINERS=$(docker container ls | grep microcosm | awk '{print $1}' | tr '\n' ' ')

echo "Confirm that the following containers should be killed: $MICROCOSM_CONTAINERS"
read -p 'Kill them all? [y/N]: ' CONFIRMATION
case $CONFIRMATION in
    [yY])
        docker kill $MICROCOSM_CONTAINERS
        ;;
esac