#!/bin/bash

# This script is adopted from blievesearch project https://github.com/blevesearch/bleve
# It has been modified to meet the project needs.
echo "mode: count" > acc.out
for Dir in . $(find ./* -maxdepth 10 -type d ); 
do
	if ls $Dir/*_test.go &> /dev/null;
	then
		returnval=`go test -coverprofile=profile.out -covermode=count $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: count" >> acc.out 
    		fi
    	else
    		exit 1
    	fi	
    fi
done

# push coverage to coveralls
if [ -n "$COVERALLS" ]
then
	goveralls -service drone.io -coverprofile=merged.out -repotoken $COVERALLS
else
	go tool cover -html=acc.out
fi


rm -rf ./profile.out
rm -rf ./acc.out