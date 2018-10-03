package authorization

import "github.com/grepory/storage/apis/meta"

// Authorizer determines whether a request is authorized using the
// RequestAttributes passed. It returns true if the request should be
// authorized, along with any error encountered
type Authorizer interface {
	Authorize(attrs RequestAttributes) (bool, error)
}

// RequestAttributes contains information about an incoming request
type RequestAttributes struct {
	GroupVersionKind meta.GroupVersionKind
	Namespace        string
	User             User
	Verb             string
}

// User describes an authenticated user
// TODO: move this into the authentication package
// Example: https://github.com/kubernetes/apiserver/blob/master/pkg/authentication/user/user.go
type User struct {
	Name   string
	Groups []string
}
