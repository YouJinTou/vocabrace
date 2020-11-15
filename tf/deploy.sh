cd ../src
go test ./...
cd ../tf

tf_root=$PWD
lambda_root="$tf_root""/../src/lambda"

cd $lambda_root/pooling/pooler
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pooler pooler.go
build-lambda-zip -output pooler.zip pooler ../config.$1.json
rm pooler
mv pooler.zip $tf_root/payloads

cd $lambda_root/pooling/connect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o connect connect.go
build-lambda-zip -output connect.zip connect ../config.$1.json
rm connect
mv connect.zip $tf_root/payloads

cd $lambda_root/pooling/disconnect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o disconnect disconnect.go
build-lambda-zip -output disconnect.zip disconnect ../config.$1.json
rm disconnect
mv disconnect.zip $tf_root/payloads

cd $lambda_root/pooling/publish
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o publish publish.go
build-lambda-zip -output publish.zip publish ../config.$1.json
rm publish
mv publish.zip $tf_root/payloads

cd $tf_root/shared
terraform apply -auto-approve

cd $tf_root/$1
terraform apply -auto-approve

cd $tf_root
