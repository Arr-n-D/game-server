package configuration

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	smClient *secretsmanager.Client
)

func (conf *Configuration) FetchDatabaseSecrets() {
    
}

func (conf *Configuration) getSecrets(name string) (string, error) {
	client := conf.()
	env := conf.GetEnv()

	secretId := env + "/backend/" + name

	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretId),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

func (conf *Configuration) getSmClient() *secretsmanager.Client {
	if smClient != nil {
		return smClient
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(conf.GetAwsRegion()),
	)

	if err != nil {
		panic(errors.Join(errors.New("unable to load SDK config"), err))
	}

	sm := secretsmanager.NewFromConfig(cfg)

	smClient = sm

	return sm
}
