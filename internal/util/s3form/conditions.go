package s3form

// Exact is a policy condition
type Exact struct{}

// Range is a policy condition
type Range struct {
	Min int
	Max int
}

// StartsWith is a policy condition
type StartsWith struct {
	Prefix string
}
