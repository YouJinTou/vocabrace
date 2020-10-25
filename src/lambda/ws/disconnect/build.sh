GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o disconnect disconnect.go
build-lambda-zip -output disconnect.zip disconnect ../config.prod.json
aws lambda update-function-code \
    --function-name onDisconnect \
    --zip-file fileb://disconnect.zip
rm disconnect
rm disconnect.zip