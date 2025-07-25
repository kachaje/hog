#!/usr/bin/env bash

if [[ "$1" == "-c" ]]; then

  rm -f **/**/output* **/output* output*

elif [[ "$1" == "-b" ]]; then

  go build -o cli cmd/cli/*.go

fi
