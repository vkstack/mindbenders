package confmanager

import (
	"encoding/base64"
	"strings"

	"gitlab.com/dotpe/mindbenders/errors"

	"gitlab.com/dotpe/mindbenders/bootconfig/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const awsRegion = "ap-south-1"

type secretManager struct {
	ENV string
}

// GetSecretManager ..
func GetSecretManager(env string) config.IConfig {
	//can  have some preprocessing logic
	return &secretManager{
		ENV: env,
	}
}

func (cfgmgr *secretManager) getSearchKey(key string) string {
	return cfgmgr.ENV + "/" + strings.Trim(key, "/")
}

//Create a Secrets Manager client
func (cfgmgr *secretManager) Get(key string) ([]byte, error) {
	newSession, _ := session.NewSession()
	svc := secretsmanager.New(newSession, aws.NewConfig().WithRegion(awsRegion))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(cfgmgr.getSearchKey(key)),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	// var secretString *string
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, errors.WrapMessage(err, "couldn't read from secret manager")
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		return []byte(*result.SecretString), nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
	len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
	if err != nil {
		return nil, errors.WrapMessage(err, "couldn't decode the read base64")
	}
	return decodedBinarySecretBytes[:len], nil
}

func (cfgmgr *secretManager) GetConfig(key string) (raw config.ConfigValue, err error) {
	if raw, err = cfgmgr.Get(key); err != nil {
		return nil, err
	}
	return raw, nil
}
