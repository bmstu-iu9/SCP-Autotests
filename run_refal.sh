#!/bin/bash
#
#case "$1" in        lambda)
#rlc ./executable.ref >out.txt && ./executable && rm out.txt executable* ;;
#        *) refc executable >out.txt && refgo executable && rm out.txt executable* ;;
#esac
#!/bin/bash

rlc ./executable.ref > trash.txt
./executable
rm ./executable*
rm ./trash.txt
