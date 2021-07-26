#!/bin/bash

cd $1
rlmake mscp-a.ref > trash.txt
rm trash.txt