GOOS=linux go build -o publish publish.go
build-lambda-zip -output publish.zip publish ../config.prod.json
aws lambda update-function-code \
    --function-name onpublish \
    --zip-file fileb://publish.zip
rm publish
rm publish.zip