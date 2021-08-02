#!/bin/bash

if [ -d "tests/rsd" ]; then
  rm -rf tests/rsd;
fi
mkdir tests/rsd;
if [ -d "tests/errors" ]; then
  rm -rf tests/errors;
fi
mkdir tests/errors;
