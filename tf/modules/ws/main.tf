resource "null_resource" "connect" {
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "go build -o connect connect.go && build-lambda-zip -output ../../../../tf/modules/ws/connect.zip connect ../config.dev.json"
    working_dir = "../../src/lambda/ws/connect"
    environment = {
      GOOS = "linux"
      GOARCH = "amd64"
      CGO_ENABLED = "0"
    }
  }
}

# data "archive_file" "connect" {
#   type        = "zip"
#   source_file = "${path.module}/connect"
#   output_path = "${path.module}/connect.zip"
#   depends_on = [null_resource.connect]
# }

module "connect" {
  source = "terraform-aws-modules/lambda/aws"
  function_name = "${var.stage}_connect"
  description   = "Invoked by the API Gateway Websocket runtime when a user connects."
  handler       = "connect"
  runtime       = "go1.x"
  source_path = [
    {
      path = "${path.module}/connect.zip"
      pip_requirements = false
    }
  ]
  environment_variables = {
    STAGE = var.stage
  }
  tags = {
    stage = var.stage
  }
  depends_on = [null_resource.connect]
}