cd ../src
go test ./...
cd ../tf

tf_root=$PWD
services_root="$tf_root""/../src/services"

cd $services_root/pooling/pooler
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pooler .
build-lambda-zip -output pooler.zip pooler
rm pooler
mv pooler.zip $tf_root/payloads

cd $services_root/pooling/connect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o connect .
build-lambda-zip -output connect.zip connect
rm connect
mv connect.zip $tf_root/payloads

cd $services_root/pooling/disconnect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o disconnect .
build-lambda-zip -output disconnect.zip disconnect
rm disconnect
mv disconnect.zip $tf_root/payloads

cd $services_root/pooling/publish
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o publish .
build-lambda-zip -output publish.zip publish
rm publish
mv publish.zip $tf_root/payloads

cd $services_root/iam
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o iam .
build-lambda-zip -output iam.zip iam
rm iam
mv iam.zip $tf_root/payloads

cd $tf_root/shared
# terraform init
terraform apply -auto-approve

cd $tf_root/stages/$1
# terraform init
terraform apply -auto-approve

cd $tf_root
