#!/bin/bash

# bump.sh inreases the number for a folder in this project.
# Calling `./scripts/bump.sh 02-shutdown` will rename the folder to
# "03-shutdown" and increase the number in the header of the README.md file.
#
# To insert a new folder at position 10 you should bump all of the folders from 11+ like this
#   for dir in {10..33}*; do ./scripts/bump.sh $dir; done
# Note that your command must specifically set the upper limit (33 in this case).

dir=$1;

topic=$(echo $dir | sed -e "s/^...//")
n=$(echo $dir | cut -f1 -d"-" | sed -e "s/^0//")
m=$((n+1))
m0=$(printf "%02d" ${m})

sed -i "" -e "s/^# [[:digit:]]*/# ${m}/" "${dir}/README.md"
git mv "${dir}" "${m0}-${topic}"
