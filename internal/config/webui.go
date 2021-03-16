package config

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/vrischmann/envconfig"

	"github.com/sjansen/magnet/internal/aws"
)

// WebUI contains settings for the "webui" lambda.
type WebUI struct {
	Config

	aws.AWS `envconfig:"-"`
	CloudFront
	SAML
	SessionStore

	Listen string `envconfig:"MAGNET_LISTEN,default=localhost:8080"`
}

// LoadWebUIConfig reads settings from the environment.
func LoadWebUIConfig(ctx context.Context) (*WebUI, error) {
	cfg := &WebUI{}

	if err := envconfig.Init(&cfg); err != nil {
		return nil, err
	}

	aws, err := aws.New(ctx)
	if err != nil {
		return nil, err
	}
	cfg.AWS.Config = aws.Config

	if cfg.SSMPrefix != "" {
		err = cfg.readSecrets(ctx, aws.NewSSMClient())
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func (cfg *WebUI) readSecrets(ctx context.Context, svc *ssm.Client) error {
	resp, err := svc.GetParameters(ctx, &ssm.GetParametersInput{
		Names: []string{
			cfg.SSMPrefix + "CLOUDFRONT_KEY_ID",
			cfg.SSMPrefix + "CLOUDFRONT_PRIVATE_KEY",
			cfg.SSMPrefix + "SAML_CERTIFICATE",
			cfg.SSMPrefix + "SAML_METADATA_URL",
			cfg.SSMPrefix + "SAML_PRIVATE_KEY",
		},
		WithDecryption: true,
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
