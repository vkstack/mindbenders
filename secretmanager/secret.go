package secretmanager

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	AWSRegion = "ap-south-1"
)

func GetSecretValue(secretName string) (map[string]string, error) {
	secretString, err := getSecretCredentials(secretName, AWSRegion)
	if err != nil || secretString == nil {
		if err == nil {
			err = errors.New("null_secretString")
		}
		log.Println(fmt.Sprintf("error in loading aws config:%s with error: %s", secretName, err))
		return nil, err
	}
	var secretValue map[string]string
	err = json.Unmarshal([]byte(*secretString), &secretValue)
	if err != nil {
		return nil, err
	}
	return secretValue, nil
}

// GetSecretCredentials ..
func getSecretCredentials(secretName string, region string) (*string, error) {
	// secretName := "staging1-chatbotDB"
	// region := "ap-south-1"
	//Create a Secrets Manager client
	newSession, _ := session.NewSession()
	svc := secretsmanager.New(newSession, aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	var secretString *string
	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	// var secretString, decodedBinarySecret string
	// var secretString string
	if result.SecretString != nil {
		secretString = result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return nil, err
		}
		// decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		decodedBinaryStr := string(decodedBinarySecretBytes[:len])
		secretString = &decodedBinaryStr
	}
	// fmt.Println("Successfully_fetched_secretString")
	// fmt.Println(secretString)
	// fmt.Println("Successfully_fetched_decodedBinarySecret")
	// fmt.Println(decodedBinarySecret)
	return secretString, nil
}
