#!/bin/bash

function login {
    AUTH_TARGET="http://$3/token-auth"
    curl -H "Content-Type: application/json" -X POST -d \
        "{\"Username\":$1,\"Password\":$2}" $AUTH_TARGET
}


login $1 $2 $3
