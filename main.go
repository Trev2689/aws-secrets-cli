package main

import (
	"context"
	//      "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	// "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func main() {
	// Parse command-line flags for input parameters
	secretNameFlag := flag.String("name", "", "AWS secret name")
	regionFlag := flag.String("region", "", "AWS region")
	descriptionFlag := flag.String("description", "", "Description for the secret")
	//kmsKeyIDFlag := flag.String("kms-key-id", "", "KMS key ID for encryption")
	jsonFilePathFlag := flag.String("json-file", "", "Path to JSON file containing the secret value")
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "Timeout for API call")
	flag.Parse()

	// Validate required input parameters
	if *regionFlag == "" || *descriptionFlag == "" || *secretNameFlag == "" || *jsonFilePathFlag == "" {
		fmt.Println("Please provide all required input parameters: region, description, secret name, json-file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read secret value from JSON file
	secretValue, err := readSecretFromJSON(*jsonFilePathFlag)
	if err != nil {
		fmt.Println("Error reading secret value from JSON file:", err)
		os.Exit(1)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancel()

	// Load AWS SDK configuration with the specified region
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*regionFlag))
	if err != nil {
		fmt.Println("Error loading AWS SDK config:", err)
		os.Exit(1)
	}

	// Create a Secrets Manager client with the specified configuration
	client := secretsmanager.NewFromConfig(cfg)

	// Check if the secret already exists
	//      describeInput := &secretsmanager.DescribeSecretInput{
	//              SecretId: descriptionFlag,
	//      }
	//      _, err = client.DescribeSecret(ctx, describeInput)
	//      if err == nil {
	//              fmt.Println("Secret already exists:", *descriptionFlag)
	//              return
	//      } else if !secretsmanager.IsResourceNotFoundException(err) {
	//              fmt.Println("Error describing secret:", err)
	//              os.Exit(1)
	//}

	// Create a Secrets Manager client with the specified configuration
	//client := secretsmanager.NewFromConfig(cfg)

	// Check if the secret already exists
	//      describeInput := &secretsmanager.DescribeSecretInput{
	//              SecretId: descriptionFlag,
	//      }
	//      _, err = client.DescribeSecret(ctx, describeInput)
	//      if err == nil {
	//              fmt.Println("Secret already exists:", *descriptionFlag)
	//              return
	//      } else if !secretsmanager.IsResourceNotFoundException(err) {
	//              fmt.Println("Error describing secret:", err)
	//              os.Exit(1)
	//      }

	// Create input for the CreateSecret API operation with the specified KMS key ID and description
	createInput := &secretsmanager.CreateSecretInput{
		Name:         secretNameFlag,
		Description:  descriptionFlag,
		SecretString: aws.String(secretValue),
		//KmsKeyId:         kmsKeyIDFlag,
		//ClientRequestToken: aws.String(*descriptionFlag), // Optional: a unique identifier for the request
	}

	// Call the CreateSecret API operation with context
	createOutput, err := client.CreateSecret(ctx, createInput)
	if err != nil {
		fmt.Println("Error creating secret:", err)
		os.Exit(1)
	}

	// Print the ARN of the created secret
	fmt.Println("Secret ARN:", *createOutput.ARN)
}

func readSecretFromJSON(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
