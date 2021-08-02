#!/bin/bash

(cd $1 && rlmake mscp-a.ref) >info.txt 2>>err.txt
