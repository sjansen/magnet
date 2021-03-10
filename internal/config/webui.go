package config

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/vrischmann/envconfig"
)

// WebUI contains settings for the "webui" lambda.
type WebUI struct {
	Config

	CloudFront
	SAML
	SessionStore

	Listen string `envconfig:"MAGNET_LISTEN,default=localhost:8080"`
}

// LoadWebUIConfig reads settings from the environment.
func LoadWebUIConfig() (*WebUI, error) {
	cfg := &WebUI{}
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

func (cfg *WebUI) readSecrets() error {
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
			if err := cfg.CloudFront.PrivateKey.Unmarshal(*param.Value); err != nil {
				return err
			}
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
