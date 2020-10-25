GOOS=linux go build -o disconnect disconnect.go
build-lambda-zip -output disconnect.zip disconnect ../config.prod.json
aws lambda update-function-code \
    --function-name ondisconnect \
    --zip-file fileb://disconnect.zip
rm disconnect
rm disconnect.zip