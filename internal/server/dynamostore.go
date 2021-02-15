package server

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sjansen/dynamostore"

	"github.com/sjansen/magnet/internal/config"
)

func (s *Server) openDynamoStores(cfg *config.Config) (scs.Store, scs.Store, error) {
	var svc *dynamodb.DynamoDB
	if cfg.SessionStore.Endpoint.Host == "" {
		svc = dynamodb.New(s.aws)
	} else {
		svc = dynamodb.New(s.aws,
			s.aws.Config.Copy().
				WithEndpoint(
					cfg.SessionStore.Endpoint.String(),
				).
				WithCredentials(
					credentials.NewStaticCredentials("id", "secret", "token"),
				),
		)
	}

	store := dynamostore.NewWithTableName(svc, cfg.SessionStore.Table)
	if cfg.SessionStore.Create {
		err := store.CreateTable()
		if err != nil {
			return nil, nil, err
		}
	}

	relaystate := NewPrefixStore("r:", store)
	sessions := NewPrefixStore("s:", store)
	return relaystate, sessions, nil
}

// PrefixStore enables multiple sessions to be stored in a single
// session store by automatically pre-pending a prefix to tokens.
type PrefixStore struct {
	prefix string
	store  scs.Store
}

// NewPrefixStore wraps a session store so it can be shared.
func NewPrefixStore(prefix string, store scs.Store) *PrefixStore {
	return &PrefixStore{
		prefix: prefix,
		store:  store,
	}
}

// Delete removes the session token and data from the store.
func (s *PrefixStore) Delete(token string) (err error) {
	return s.store.Delete(s.prefix + token)
}

// Find returns the data for a session token from the store.
func (s *PrefixStore) Find(token string) (b []byte, found bool, err error) {
	return s.store.Find(s.prefix + token)
}

// Commit adds the session token and data to the store.
func (s *PrefixStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	return s.store.Commit(s.prefix+token, b, expiry)
}
