#!/bin/bash
valid=$(echo $1 | awk '/(^https:\/\/[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:8443)/')
[ -z $valid ] && echo "invalid" || echo "valid"
