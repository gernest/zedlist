#!/bin/bash

# This script is adopted from blievesearch project https://github.com/blevesearch/bleve
# It has been modified to meet the project needs.
for Dir in . $(find ./* -maxdepth 10 -type d ); 
do
	if ls $Dir/*_test.go &> /dev/null;
	then
		echo $Dir
		golint $Dir
		returnval=`go test $Dir`
		echo ${returnval}
    fi
done
