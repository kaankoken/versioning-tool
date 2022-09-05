#!/bin/bash

mkdir -p tests

echo "Running test, this take a while..."
echo
input=$(find "tests" \( -name "*_test.go" \))

# sort the array
input=$(echo $input | tr ' ' '\n' | sort -n | tr '\n' ' ')

rm -rf cover.out

isFirstLineAdded=0

for line in $input
do
    start_time=$SECONDS
    echo "Running test: $line"
    echo "------------------------------------------------------------"

    # split string into array
    IFS='/' read -ra array <<< "$line"

    # get last element of array
    filename=${array[${#array[@]}-1]}

    filename=${filename%"_test.go"}
    filename="${filename}.go"

    go test -v -coverpkg=./... $line -coverprofile=cover.tmp.out -covermode=atomic --tags=unit
    
    if [ $isFirstLineAdded == 0 ]
    then
        #get first line of cover.tmp.out and assign to variable
        firstline=$(head -n 1 cover.tmp.out)
        echo $firstline > cover.out
        
        isFirstLineAdded=1
    fi

    # remove line that does not contain "xxxx_test.go"
    cat cover.tmp.out | grep "$filename" >> cover.out
    
    elapsed=$(( SECONDS - start_time ))

    echo "------------------------------------------------------------"
    echo "Finished test: $line"
    echo "It took $elapsed seconds to run $line"
    echo
done

filename=cover.out

if [ ! -f $filename ]
then
    touch $filename
fi

rm -rf cover.tmp.out
go tool cover -html=cover.out -o cover.html
find . -type f -name 'config.env' -delete
find . -type f -name 'text.log' -delete
