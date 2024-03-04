AWS Secrets CLI
AWS Secrets CLI is a command-line tool written in Go that allows you to interact with AWS Secrets Manager to create or update secrets.

Installation
To build the binary from source, ensure you have Go installed on your system. Then, follow these steps:

Clone this repository to your local machine.

Navigate to the directory containing the main.go file.

Run the following command to build the binary:

GOOS=linux GOARCH=amd64 go build -o aws-secrets-cli-linux main.go 

GOOS=windows GOARCH=amd64 go build -o aws-secrets-cli-windows.exe main.go 

GOOS=darwin GOARCH=amd64 go build -o aws-secrets-cli-macos main.go

Usage
The aws-secrets-cli binary supports the following command-line flags:

-name: Specifies the name of the AWS secret. (Required)
-region: Specifies the AWS region where the secret will be created or updated. (Required)
-description: Specifies the description for the AWS secret. (Required)
-json-file: Specifies the path to the JSON file containing the secret value. (Required)
-timeout: Specifies the timeout duration for the API call. (Optional, default: 10 seconds)
-update: Optional flag indicating whether to update the secret if it already exists.
Here's an example of how to use the binary:


./aws-secrets-cli -name "my-secret" -region "us-east-1" -description "My AWS Secret" -json-file "/path/to/secret.json" -timeout 10s -update
Replace the flag values with your desired values. The -update flag is optional and indicates that you want to update the secret if it already exists.
