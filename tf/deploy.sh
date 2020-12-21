cd ../src
go test ./...
cd ../tf

tf_root=$PWD
services_root="$tf_root""/../src/services"

cd $services_root/broadcast
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o broadcast .
build-lambda-zip -output broadcast.zip broadcast
rm broadcast
mv broadcast.zip $tf_root/payloads

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

cd $services_root/pooling/reconnect
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o reconnect .
build-lambda-zip -output reconnect.zip reconnect
rm reconnect
mv reconnect.zip $tf_root/payloads

cd $services_root/pooling/tally
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tally .
build-lambda-zip -output tally.zip tally
rm tally
mv tally.zip $tf_root/payloads

cd $services_root/pooling/pool
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pool .
build-lambda-zip -output pool.zip pool
rm pool
mv pool.zip $tf_root/payloads

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
