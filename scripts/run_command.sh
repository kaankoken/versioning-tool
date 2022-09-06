#!/bin/bash

echo "Running $1"
echo "------------------------------------------------------------"

count=$($1 | wc -l)
if [ $count -ne 0 ]
then
    $1
    exit 1
fi

echo "Success running $1"
echo "------------------------------------------------------------"
exit 0