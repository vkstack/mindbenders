package logging

import (
	"gitlab.com/dotpe/mindbenders/errors"

	"flag"
	"log"

	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/aws"

	"github.com/sirupsen/logrus"
	awsauth "github.com/smartystreets/go-aws-auth"
	"gopkg.in/sohlich/elogrus.v7"
)

type elasticConfig struct {
	client,
	accessKey,
	secretKey,
	app,
	appId string
}

func NewElasticHookContainer(Client, AccessKey, SecretKey, APP, APPID string) IHookContainer {
	return &elasticConfig{
		client:    Client,
		accessKey: AccessKey,
		secretKey: SecretKey,
		app:       APP,
		appId:     APPID,
	}
}

func (conf *elasticConfig) GetHook() (logrus.Hook, error) {
	return GetElasticHook(conf.app, conf.client, conf.accessKey, conf.secretKey)
}

func (conf *elasticConfig) MustGetHook() logrus.Hook {
	hook, err := GetElasticHook(conf.app, conf.client, conf.accessKey, conf.secretKey)
	if err != nil {
		log.Fatalf("unable to get elastic hook:%v\n", err)
	}
	return hook
}

func GetElasticHook(app, url, accessKey, secretKey string) (logrus.Hook, error) {
	if url == "" {
		return nil, errors.New("missing -client-url KIBANA")
	}
	if accessKey == "" {
		return nil, errors.New("missing -access-key or AWS_ACCESS_KEY environment variable")
	}
	if secretKey == "" {
		return nil, errors.New("missing -secret-key or AWS_SECRET_KEY environment variable")
	}

	sniff := flag.Bool("sniff", false, "Enable or disable sniffing")

	flag.Parse()
	log.SetFlags(0)

	signingClient := aws.NewV4SigningClient(awsauth.Credentials{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
	})

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(*sniff),
		elastic.SetHealthcheck(*sniff),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return elogrus.NewAsyncElasticHook(client, "", logrus.DebugLevel, app)
}
