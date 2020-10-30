set -e

STAGE=${1:-dev}

echo Deploying to $STAGE

cd ../src/lambda/ws/connect
. deploy.sh $STAGE
cd ../disconnect
. deploy.sh $STAGE
cd ../publish
. deploy.sh $STAGE
cd ../../../../deploy