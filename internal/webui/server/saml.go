package server

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/sjansen/magnet/internal/config"
)

func newSAMLMiddleware(cfg *config.SAML) (*samlsp.Middleware, error) {
	keyPair, err := tls.X509KeyPair(
		[]byte(cfg.Certificate),
		[]byte(cfg.PrivateKey),
	)
	if err != nil {
		return nil, err
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, err
	}

	metadata, err := loadIDPMetadata(cfg)
	if err != nil {
		return nil, err
	}

	return samlsp.New(samlsp.Options{
		EntityID:    cfg.EntityID,
		URL:         cfg.RootURL.URL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		// TODO Intermediates
		IDPMetadata: metadata,

		AllowIDPInitiated: true,
		// TODO CookieSameSite
		// TODO EntityID
		// TODO ForceAuthn
		// TODO SignRequest
	})
}

func loadIDPMetadata(cfg *config.SAML) (*saml.EntityDescriptor, error) {
	resp, err := http.Get(cfg.MetadataURL)
	if err != nil {
		return nil, fmt.Errorf("fetch idp metadata: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read idp metadata: %w", err)
	}

	m, err := samlsp.ParseMetadata(body)
	if err != nil {
		return nil, fmt.Errorf("parse idp metadata: %w", err)
	}

	return m, nil
}
