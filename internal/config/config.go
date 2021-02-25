package config

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/vrischmann/envconfig"
)

// Config contains application settings.
type Config struct {
	Debug      bool   `envconfig:"DEBUG,optional"`
	Bucket     string `envconfig:"MAGNET_BUCKET"`
	Listen     string `envconfig:"MAGNET_LISTEN,default=localhost:8080"`
	Root       URL    `envconfig:"MAGNET_URL,default=http://localhost:8080/"`
	CloudFront struct {
		KeyID      string     `envconfig:"MAGNET_CLOUDFRONT_KEY_ID"`
		PrivateKey PrivateKey `envconfig:"MAGNET_CLOUDFRONT_PRIVATE_KEY"`
	}
	SAML struct {
		EntityID    string `envconfig:"MAGNET_SAML_ENTITY_ID,default=magnet"`
		MetadataURL string `envconfig:"MAGNET_SAML_METADATA_URL"`
		Certificate string `envconfig:"MAGNET_SAML_CERTIFICATE"`
		PrivateKey  string `envconfig:"MAGNET_SAML_PRIVATE_KEY"`
	}
	SessionStore struct {
		Create   bool   `envconfig:"MAGNET_SESSION_CREATE,default=false"`
		Endpoint URL    `envconfig:"MAGNET_SESSION_ENDPOINT,optional"`
		Table    string `envconfig:"MAGNET_SESSION_TABLE"`
	}
	SSMPrefix string `envconfig:"MAGNET_SSM_PREFIX,optional"`
}

// New loads application settings from the environment.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}
	if cfg.SSMPrefix != "" {
		if err := cfg.readSecrets(); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func (cfg *Config) readSecrets() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := ssm.New(sess)
	resp, err := svc.GetParameters(&ssm.GetParametersInput{
		Names: []*string{
			aws.String(cfg.SSMPrefix + "CLOUDFRONT_KEY_ID"),
			aws.String(cfg.SSMPrefix + "CLOUDFRONT_PRIVATE_KEY"),
			aws.String(cfg.SSMPrefix + "SAML_CERTIFICATE"),
			aws.String(cfg.SSMPrefix + "SAML_METADATA_URL"),
			aws.String(cfg.SSMPrefix + "SAML_PRIVATE_KEY"),
		},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	for _, param := range resp.Parameters {
		name := strings.TrimPrefix(*param.Name, cfg.SSMPrefix)
		switch name {
		case "CLOUDFRONT_KEY_ID":
			cfg.CloudFront.KeyID = *param.Value
		case "CLOUDFRONT_PRIVATE_KEY":
			cfg.CloudFront.PrivateKey.Unmarshal(*param.Value)
		case "SAML_CERTIFICATE":
			cfg.SAML.Certificate = *param.Value
		case "SAML_METADATA_URL":
			cfg.SAML.MetadataURL = *param.Value
		case "SAML_PRIVATE_KEY":
			cfg.SAML.PrivateKey = *param.Value
		}
	}
	return nil
}
