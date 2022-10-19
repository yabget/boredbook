#!/bin/bash

cat $1 | tr '"' '\n' | tr "'" '\n' | grep -e '^https://' -e '^http://' -e'^//' | sort | uniq | tee urls.txt
