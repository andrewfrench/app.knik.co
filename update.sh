#!/bin/bash

INITIAL_PATH=`pwd`;
export GOPATH=$INITIAL_PATH;
UPDATE_PATH_PREFIX="$INITIAL_PATH/src/knik.co/api";
UPDATE_PATH=$1;
LAMBDA_NAME="lambda-unified";
LAMBDA_STAGE="production";

AWS_PROFILE="infill";
AWS_REGION="us-west-2";

echo "Updating $UPDATE_PATH-$LAMBDA_STAGE";

echo "Moving to $UPDATE_PATH";
cd $UPDATE_PATH_PREFIX/$UPDATE_PATH;

echo "Copying Makefile to $UPDATE_PATH";
cp $INITIAL_PATH/Makefile .;

echo "Executing make";
say "Build started";
make

if [ $? -ne 0 ]; 
then
	echo "Make failed, exiting.";
	exit -1;
fi

echo "Removing copied Makefile";
rm $UPDATE_PATH_PREFIX/$UPDATE_PATH/Makefile;

if [ -s "$UPDATE_PATH_PREFIX/$UPDATE_PATH/handler.zip" ]; 
then
	echo "Build succeeded!";
	say "Build succeeded";
else
	echo "Build failed!";
	say "Build complete";
	exit -1;
fi

say "Lambda update started";
echo "Updating code for lambda: $LAMBDA_NAME";
aws lambda update-function-code \
	--profile $AWS_PROFILE \
	--region $AWS_REGION \
	--function-name $LAMBDA_NAME-$LAMBDA_STAGE \
	--zip-file "fileb://$UPDATE_PATH_PREFIX/$UPDATE_PATH/handler.zip";

if [ $? -ne 0 ];
then
	echo "Lambda code update failed!";
    say "Lambda update failed";
	exit -1;
else
    say "Lambda update complete";
fi 

echo "Cleaning up";
# rm $UPDATE_PATH_PREFIX/$UPDATE_PATH/handler.zip;
rm $UPDATE_PATH_PREFIX/$UPDATE_PATH/handler.so;

echo "Returning to initial path";
cd $INITIAL_PATH;
