GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o publish publish.go
build-lambda-zip -output publish.zip publish ../config.prod.json
aws lambda update-function-code \
    --function-name onpublish \
    --zip-file fileb://publish.zip
rm publish
rm publish.zip