package bootconfig

import (
	"encoding/base64"
	"strings"

	"gitlab.com/dotpe/mindbenders/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const awsRegion = "ap-south-1"

type secretManager struct {
	ENV string
}

// GetSecretManager ..
func GetSecretManager(env string) iConfig {
	//can  have some preprocessing logic
	return &secretManager{
		ENV: env,
	}
}

func (cfgmgr *secretManager) getSearchKey(key string) string {
	return cfgmgr.ENV + "/" + strings.Trim(key, "/")
}

// Create a Secrets Manager client
func (cfgmgr *secretManager) get(key string) ([]byte, error) {
	newSession, _ := session.NewSession()
	svc := secretsmanager.New(newSession, aws.NewConfig().WithRegion(awsRegion))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	// var secretString *string
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, errors.WrapMessage(err, "couldn't read from secret manager:\t"+key)
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		return []byte(*result.SecretString), nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return nil, errors.WrapMessage(err, "couldn't decode the read base64 for:\t"+key)
	}
	return decodedBinarySecretBytes[:len], nil
}

func (cfgmgr *secretManager) Get(key string) ([]byte, error) {
	return cfgmgr.get(cfgmgr.getSearchKey(key))
}

// env specifig config
func (cfgmgr *secretManager) GetGlobal(key string) (raw []byte, err error) {
	return cfgmgr.get(key)
}
