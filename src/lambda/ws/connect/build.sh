GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o connect connect.go
build-lambda-zip -output connect.zip connect ../config.$1.json
aws lambda update-function-code \
    --function-name onConnect \
    --zip-file fileb://connect.zip
rm connect
rm connect.zip