GOOS=linux go build -o connect connect.go
build-lambda-zip -output connect.zip connect ../config.prod.json
aws lambda update-function-code \
    --function-name onConnect \
    --zip-file fileb://connect.zip
rm connect
rm connect.zip