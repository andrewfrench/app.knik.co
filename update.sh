#!/bin/bash

export GOPATH=`pwd`;

LAMBDA="knikco-unified-production"
LOCATION=`pwd`/src/knik.co/api/knikco-unified
AWS_PROFILE="infill";
AWS_REGION="us-west-2";

echo "Copying Makefile"
cp ./Makefile ${LOCATION}/Makefile

echo "Executing make";
cd ${LOCATION}
say "Build started";
make

if [ $? -ne 0 ]; 
then
	echo "Make failed, exiting.";
	exit -1;
fi

if [ -s "${LOCATION}/handler.zip" ]; 
then
	echo "Build succeeded!";
	say "Build succeeded";
else
	echo "Build failed!";
	say "Build failed";
	exit -1;
fi

say "Lambda update started";
echo "Updating code for lambda: $LAMBDA";
aws lambda update-function-code \
	--profile $AWS_PROFILE \
	--region $AWS_REGION \
	--function-name $LAMBDA \
	--zip-file "fileb://${LOCATION}/handler.zip";

if [ $? -ne 0 ];
then
	echo "Lambda code update failed!";
    say "Lambda update failed";
	exit -1;
else
    say "Lambda update complete";
fi 

echo "Cleaning up";
rm ${LOCATION}/handler.so ${LOCATION}/Makefile

