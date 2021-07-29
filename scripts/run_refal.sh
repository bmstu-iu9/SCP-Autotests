#!/bin/bash

if [[ "$2" -eq "lambda" ]]; then
  rlc $1 >>info.txt 2>>tests/errors/$1_err.txt
  ./$1 2>>tests/errors/$1_err.txt
else
  refc $1 >>info.txt 2>>tests/errors/$1_err.txt
  refgo $1 2>>tests/errors/$1_err.txt
fi
rm $1*