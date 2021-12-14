#!/bin/sh
# Script used to run entities when calling makefile

if [[ $1 == $2 ]]; then
    shift
fi

entity=$1
id=$2
if [[ -z $id ]]; then
    id="0"
fi

if [[ $entity == "fulcrum" || $entity == "informant" ]]; then
    echo "Running entity \"$entity\" with id $id"
else
    echo "Running entity \"$entity\""
fi

if [[ $entity == "informant" || $entity == "leia" ]]; then
    go run main.go run $entity $id
else
    go run main.go run $entity $id &
fi
