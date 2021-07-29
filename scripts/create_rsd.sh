#!/bin/bash
cd $1
timeout -k 10 "$4" ./mscp-a $2 $3
