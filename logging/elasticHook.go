package logging

import (
	"flag"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/aws"

	"github.com/sirupsen/logrus"
	awsauth "github.com/smartystreets/go-aws-auth"
	"gopkg.in/sohlich/elogrus.v7"
)

type KibanaConfig struct {
	client,
	accessKey,
	secretKey,
	app,
	appId string
}

func NewKibanaConfig(Client, AccessKey, SecretKey, APP, APPID string) ILogConfig {
	return &KibanaConfig{
		client:    Client,
		accessKey: AccessKey,
		secretKey: SecretKey,
		app:       APP,
		appId:     APPID,
	}
}
func (conf *KibanaConfig) getHook() (logrus.Hook, error) {
	client, err := conf.newElasticClient()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return elogrus.NewAsyncElasticHook(client, "", logrus.DebugLevel, conf.app)
}

func (conf *KibanaConfig) newElasticClient() (*elastic.Client, error) {
	if conf.client == "" {
		log.Fatal("missing -client-url KIBANA")
	}
	if conf.accessKey == "" {
		log.Fatal("missing -access-key or AWS_ACCESS_KEY environment variable")
	}
	if conf.secretKey == "" {
		log.Fatal("missing -secret-key or AWS_SECRET_KEY environment variable")
	}

	sniff := flag.Bool("sniff", false, "Enable or disable sniffing")

	flag.Parse()
	log.SetFlags(0)

	signingClient := aws.NewV4SigningClient(awsauth.Credentials{
		AccessKeyID:     conf.accessKey,
		SecretAccessKey: conf.secretKey,
	})

	client, err := elastic.NewClient(
		elastic.SetURL(conf.client),
		elastic.SetSniff(*sniff),
		elastic.SetHealthcheck(*sniff),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("AWS ElasticSearchConnection succeeded")
	return client, nil
}
