#!/bin/sh

set -e -o pipefail

DEBUG=$DEBUG

set -u

# If debug is set to anything nonempty, then simply run the command passed to "docker run"
if [ ! -z $DEBUG ] ; then
    $@
else
    echo "lol"
fi
