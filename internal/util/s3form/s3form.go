package s3form

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

/* https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-UsingHTTPPOST.html */
/* https://aws.amazon.com/premiumsupport/knowledge-center/presigned-url-s3-bucket-expiration/ */

var now = time.Now

// Form is used to create a SignedForm.
type Form struct {
	region     string
	bucket     string
	fields     map[string]string
	conditions []interface{}
}

// SignedForm can be used for browser-based upload to S3.
type SignedForm struct {
	action string
	fields map[string]string
}

type policy struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

// New creates a new Form.
func New(region, bucket string) *Form {
	f := &Form{
		region: region,
		bucket: bucket,
		fields: make(map[string]string, 3),
	}
	f.AddField("bucket", bucket)
	f.AddField("x-amz-algorithm", "AWS4-HMAC-SHA256")
	return f
}

// AddField adds a field and associated policy conditions to the form.
func (f *Form) AddField(name, value string, conditions ...interface{}) *Form {
	f.fields[name] = value
	if len(conditions) < 1 {
		conditions = []interface{}{&Exact{}}
	}
	for _, c := range conditions {
		switch v := c.(type) {
		default:
			f.conditions = append(f.conditions, map[string]string{
				name: value,
			})
		case *Range:
			f.conditions = append(f.conditions, []string{
				name, strconv.Itoa(v.Max), strconv.Itoa(v.Max),
			})
		case *StartsWith:
			f.conditions = append(f.conditions, []string{
				"starts-with", "$" + name, v.Prefix,
			})
		}
	}
	return f
}

// Key adds the key field and an exact match policy condition to the form.
func (f *Form) Key(key string) *Form {
	return f.AddField("key", key)
}

// Prefix adds the key field and a starts-with policy condition to the form.
func (f *Form) Prefix(prefix string) *Form {
	return f.AddField("key", prefix+"${filename}", &StartsWith{Prefix: prefix})
}

// Sign creates a new SignedForm.
func (f *Form) Sign(accessKeyID, secretAccessKey, sessionToken string, expiration time.Time) (*SignedForm, error) {
	tmp := &Form{
		region:     f.region,
		bucket:     f.bucket,
		fields:     make(map[string]string, len(f.fields)+5),
		conditions: make([]interface{}, len(f.conditions), len(f.conditions)+3),
	}
	for key, value := range f.fields {
		tmp.fields[key] = value
	}
	copy(tmp.conditions, f.conditions)

	now := now().UTC()
	date := now.Format("20060102")
	tmp.AddField("x-amz-credential",
		fmt.Sprintf("%s/%s/%s/s3/aws4_request",
			accessKeyID, date, tmp.region,
		),
	)
	tmp.AddField("x-amz-date",
		now.Format("20060102T000000Z"),
	)
	if sessionToken != "" {
		tmp.AddField("x-amz-security-token", sessionToken)
	}

	json, err := json.Marshal(&policy{
		Expiration: expiration.UTC().Format("2006-01-02T15:04:05.000Z"),
		Conditions: tmp.conditions,
	})
	if err != nil {
		return nil, err
	}
	policy := make([]byte, base64.StdEncoding.EncodedLen(len(json)))
	base64.StdEncoding.Encode(policy, json)
	tmp.fields["policy"] = string(policy)

	hmac1 := hmac.New(sha256.New, []byte("AWS4"+secretAccessKey))
	if _, err := hmac1.Write([]byte(date)); err != nil {
		return nil, err
	}
	hmac2 := hmac.New(sha256.New, hmac1.Sum(nil))
	if _, err := hmac2.Write([]byte(tmp.region)); err != nil {
		return nil, err
	}
	hmac3 := hmac.New(sha256.New, hmac2.Sum(nil))
	if _, err := hmac3.Write([]byte("s3")); err != nil {
		return nil, err
	}
	signingKey := hmac.New(sha256.New, hmac3.Sum(nil))
	if _, err := signingKey.Write([]byte("aws4_request")); err != nil {
		return nil, err
	}
	signature := hmac.New(sha256.New, signingKey.Sum(nil))
	if _, err := signature.Write(policy); err != nil {
		return nil, err
	}
	tmp.fields["x-amz-signature"] = hex.EncodeToString(signature.Sum(nil))

	return &SignedForm{
		action: fmt.Sprintf(
			"https://%s.s3.%s.amazonaws.com/",
			f.bucket, f.region,
		),
		fields: tmp.fields,
	}, nil
}

// Action provides the URL that should be used to upload a file.
func (f *SignedForm) Action() string {
	return f.action
}

// Enctype provides the HTML form enctype that should be used to upload a file.
func (f *SignedForm) Enctype() string {
	return "multipart/form-data"
}

// Fields returns the HTML form fields that should be used to upload a file.
func (f *SignedForm) Fields() map[string]string {
	return f.fields
}

// Method provides the HTTP method that should be used to upload a file.
func (f *SignedForm) Method() string {
	return "POST"
}
