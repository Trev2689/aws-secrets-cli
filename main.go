package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	// Parse command-line flags for input parameters
	secretNameFlag := flag.String("name", "", "AWS secret name")
	regionFlag := flag.String("region", "", "AWS region")
	descriptionFlag := flag.String("description", "", "Description for the secret")
	jsonFilePathFlag := flag.String("json-file", "", "Path to JSON file containing the secret value")
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "Timeout for API call")
	flag.Parse()

	// Check required input parameters
	if *regionFlag == "" || *descriptionFlag == "" || *secretNameFlag == "" || *jsonFilePathFlag == "" {
		fmt.Println("Please provide all required input parameters: region, description, secret name, json-file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read secret value from a JSON file
	// TODO: Make this more robust e.g kv on each key
	secretValue, err := readSecretFromJSON(*jsonFilePathFlag)
	if err != nil {
		fmt.Println("Error reading secret value from JSON file:", err)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancel()

	// Load AWS SDK configuration with the specified region, assume we have this inside vault
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*regionFlag))
	if err != nil {
		fmt.Println("Error loading AWS SDK config:", err)
		os.Exit(1)
	}

	// Create Secrets Manager client with the config from above
	client := secretsmanager.NewFromConfig(cfg)

	// Check if the secret we want to create already exists
	describeInput := &secretsmanager.DescribeSecretInput{
		SecretId: secretNameFlag,
	}

	_, err = client.DescribeSecret(ctx, describeInput)
	if err != nil {
		// Check if the error message indicates that the secret doesn't exist
		// TODO: Could be more robust here on explicityly checking error types
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			// Secret does not exist, so proceed with creation
			fmt.Println("Secret does not exist. Proceeding with creation...")
		} else {
			// Some other error occurred, best to exit
			fmt.Println("Error describing secret:", err)
			os.Exit(1)
		}
	} else {
		// Secret already exists, dont proceed with creation
		fmt.Println("Secret already exists:", *descriptionFlag)
		return
	}

	// Create input for the CreateSecret API operation
	createInput := &secretsmanager.CreateSecretInput{
		Name:         secretNameFlag,
		Description:  descriptionFlag,
		SecretString: &secretValue,
	}

	// Call CreateSecret API operation with above input
	createOutput, err := client.CreateSecret(ctx, createInput)
	if err != nil {
		fmt.Println("Error creating secret:", err)
		os.Exit(1)
	}

	// Print ARN of the created secret
	fmt.Println("Successfully created, Secret ARN:", *createOutput.ARN)
}

func readSecretFromJSON(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
