cd ../src/lambda/pooling/connect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o connect connect.go
build-lambda-zip -output connect.zip connect ../config.$1.json
rm connect
mv connect.zip ../../../../tf/payloads

cd ../disconnect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o disconnect disconnect.go
build-lambda-zip -output disconnect.zip disconnect ../config.$1.json
rm disconnect
mv disconnect.zip ../../../../tf/payloads

cd ../publish
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o publish publish.go
build-lambda-zip -output publish.zip publish ../config.$1.json
rm publish
mv publish.zip ../../../../tf/payloads

cd ../../../../tf/$1

terraform apply -auto-approve

cd ..