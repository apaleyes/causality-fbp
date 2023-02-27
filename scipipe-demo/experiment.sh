#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Usage:"
    echo "    ./experiment.sh <base_dir> <mode> <experiment_dir> <repeats>"
    exit 1
fi


for (( i=1;i<=$4;i++ ))
do
    mkdir -p "$3/$i"
    cp pipeline.go "$3/$i/"
    cd "$3/$i"
    go run pipeline.go "$1" "$2"
    cd -
done
