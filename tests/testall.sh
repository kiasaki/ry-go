#!/bin/sh

red='\e[0;31m' # Red
grn='\e[0;32m' # Green
rst='\e[0m'    # Text Reset

for lispfile in tests/*.ryl
do
    ./ryl -exitonfail "${lispfile}" && \
        printf "$grn${lispfile} passed$rst\n" || \
        printf "$red${lispfile} failed$rst\n"
done
