package config

// Config contains shared settings.
type Config struct {
	Development bool   `envconfig:"MAGNET_DEVELOPMENT,default=false"`
	Bucket      string `envconfig:"MAGNET_BUCKET"`
	SSMPrefix   string `envconfig:"MAGNET_SSM_PREFIX,optional"`
}

// CloudFront contains settings for the webui CDN.
type CloudFront struct {
	KeyID      string     `envconfig:"MAGNET_CLOUDFRONT_KEY_ID"`
	PrivateKey PrivateKey `envconfig:"MAGNET_CLOUDFRONT_PRIVATE_KEY"`
	URL        URL        `envconfig:"MAGNET_CLOUDFRONT_URL,default=/"`
}

// SAML contains settings for SAML-based authentication.
type SAML struct {
	EntityID    string `envconfig:"MAGNET_SAML_ENTITY_ID,default=magnet"`
	MetadataURL string `envconfig:"MAGNET_SAML_METADATA_URL"`
	Certificate string `envconfig:"MAGNET_SAML_CERTIFICATE"`
	PrivateKey  string `envconfig:"MAGNET_SAML_PRIVATE_KEY"`
}

// SessionStore contains setting for webui sessions.
type SessionStore struct {
	Create   bool   `envconfig:"MAGNET_SESSION_CREATE,default=false"`
	Endpoint URL    `envconfig:"MAGNET_SESSION_ENDPOINT,optional"`
	Table    string `envconfig:"MAGNET_SESSION_TABLE"`
}
