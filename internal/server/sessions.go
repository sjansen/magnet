package server

import (
	"crypto/rand"
	"encoding/ascii85"
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

const sessionCookieName = "sessionid"
const sessionLifetime = 7 * 24 * time.Hour
const trackerCookieName = "relaystate"
const trackerLifetime = time.Hour

func (s *Server) addSCS(relaystate, sessions scs.Store) {
	domain := s.config.Root.URL.Hostname()
	s.useSCS = true

	// relaystate
	sm := scs.New()
	sm.Cookie.Domain = domain
	sm.Cookie.HttpOnly = true
	sm.Cookie.Name = trackerCookieName
	sm.Cookie.Persist = false
	sm.Cookie.SameSite = http.SameSiteNoneMode
	if domain == "localhost" || domain == "127.0.0.1" {
		sm.Cookie.Secure = false
	} else {
		sm.Cookie.Secure = true
	}
	sm.IdleTimeout = trackerLifetime
	sm.Lifetime = trackerLifetime
	if relaystate != nil {
		sm.Store = relaystate
	}
	s.relaystate = sm
	s.saml.RequestTracker = s
	gob.Register([]samlsp.TrackedRequest{})

	// sessions
	sm = scs.New()
	sm.Cookie.Domain = domain
	sm.Cookie.HttpOnly = true
	sm.Cookie.Name = sessionCookieName
	sm.Cookie.Persist = true
	if domain == "localhost" || domain == "127.0.0.1" {
		sm.Cookie.Secure = false
	} else {
		sm.Cookie.SameSite = http.SameSiteStrictMode
		sm.Cookie.Secure = true
	}
	sm.IdleTimeout = time.Hour
	sm.Lifetime = sessionLifetime
	if sessions != nil {
		sm.Store = sessions
	}
	s.sessions = sm
	s.saml.Session = s
}

// CreateSession is called when we have received a valid SAML assertion and
// should create a new session and modify the http response accordingly, e.g. by
// setting a cookie.
func (s *Server) CreateSession(w http.ResponseWriter, r *http.Request, assertion *saml.Assertion) error {
	ctx := r.Context()
	err := s.sessions.RenewToken(ctx)
	if err != nil {
		return err
	}

	u := User{}
	for _, attributeStatement := range assertion.AttributeStatements {
		for _, attr := range attributeStatement.Attributes {
			claimName := attr.FriendlyName
			if claimName == "" {
				claimName = attr.Name
			}
			for _, value := range attr.Values {
				switch claimName {
				case "email":
					u.Email = value.Value
				case "firstName":
					u.GivenName = value.Value
				case "lastName":
					u.Surname = value.Value
				case "roles":
					u.Roles = append(u.Roles, value.Value)
				}
			}
		}
	}
	s.sessions.Put(ctx, "User", u)
	_ = s.relaystate.Destroy(ctx)

	return nil
}

// DeleteSession is called to modify the response such that it removed the current
// session, e.g. by deleting a cookie.
func (s *Server) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	return s.sessions.Destroy(r.Context())
}

// GetSession returns the current samlsp.Session associated with the request, or
// ErrNoSession if there is no valid session.
func (s *Server) GetSession(r *http.Request) (samlsp.Session, error) {
	ctx := r.Context()
	if u, ok := s.sessions.Get(ctx, "User").(User); ok {
		return &u, nil
	}
	return nil, samlsp.ErrNoSession
}

// ErrNoTrackedRequest is returned for invalid and expired relay states
var ErrNoTrackedRequest = errors.New("saml: tracked request not present")

const trackedRequestsKey = "TrackedRequests"
const trackedRequestsLimit = 10

// GetTrackedRequest returns a pending tracked request.
func (s *Server) GetTrackedRequest(r *http.Request, index string) (*samlsp.TrackedRequest, error) {
	requests, ok := s.relaystate.Get(r.Context(), trackedRequestsKey).([]samlsp.TrackedRequest)
	if !ok {
		return nil, ErrNoTrackedRequest
	}
	for _, r := range requests {
		if r.Index == index {
			return &r, nil
		}
	}
	return nil, ErrNoTrackedRequest
}

// GetTrackedRequests returns all the pending tracked requests
func (s *Server) GetTrackedRequests(r *http.Request) []samlsp.TrackedRequest {
	requests, ok := s.relaystate.Get(r.Context(), trackedRequestsKey).([]samlsp.TrackedRequest)
	if ok {
		return requests
	}
	return []samlsp.TrackedRequest{}
}

// StopTrackingRequest stops tracking the SAML request given by index, which is a string
// previously returned from TrackRequest
func (s *Server) StopTrackingRequest(w http.ResponseWriter, r *http.Request, index string) error {
	ctx := r.Context()
	requests, ok := s.relaystate.Get(ctx, trackedRequestsKey).([]samlsp.TrackedRequest)
	if ok {
		for i := len(requests) - 1; i >= 0; i-- {
			if requests[i].Index == index {
				copy(requests[i:], requests[i+1:])
				requests = requests[:len(requests)-1]
			}
		}
		if len(requests) > 0 {
			s.relaystate.Put(ctx, trackedRequestsKey, requests)
		} else {
			s.relaystate.Remove(ctx, trackedRequestsKey)
		}
	}
	return nil
}

// TrackRequest starts tracking the SAML request with the given ID. It returns an
// `index` that should be used as the RelayState in the SAMl request flow.
func (s *Server) TrackRequest(w http.ResponseWriter, r *http.Request, samlRequestID string) (string, error) {
	src := make([]byte, 29)
	if _, err := rand.Read(src); err != nil {
		return "", err
	}

	dst := make([]byte, ascii85.MaxEncodedLen(len(src)))
	ascii85.Encode(dst, src)

	index := string(dst)
	request := samlsp.TrackedRequest{
		Index:         index,
		SAMLRequestID: samlRequestID,
		URI:           r.URL.String(),
	}

	ctx := r.Context()
	requests, ok := s.relaystate.Get(ctx, trackedRequestsKey).([]samlsp.TrackedRequest)
	switch {
	case ok && len(requests) < trackedRequestsLimit:
		requests = append(requests, request)
	case ok:
		copy(requests, requests[1:])
		requests[len(requests)-1] = request
	default:
		requests = []samlsp.TrackedRequest{request}
	}
	s.relaystate.Put(ctx, trackedRequestsKey, requests)

	return index, nil
}
