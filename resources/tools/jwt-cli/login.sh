#!/bin/bash
curl -X POST -H "Content-Type: application/json" -d "{\"Username\": \"$1\", \"Password\": \"$2\"}" "http://$3/token-auth"
