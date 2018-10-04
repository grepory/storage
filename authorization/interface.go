package authorization

// Authorizer determines whether a request is authorized using the
// RequestAttributes passed. It returns true if the request should be
// authorized, along with any error encountered
type Authorizer interface {
	Authorize(attrs RequestAttributes) (bool, error)
}

// RequestAttributes contains information about an incoming request. Populated
// by the middleware with something like this:
// https://github.com/kubernetes/apiserver/blob/45bb707b3e17de3fa72c53d6df5403fea2c3150a/pkg/endpoints/request/requestinfo.go#L116
type RequestAttributes struct {
	APIGroup     string
	Namespace    string
	Resource     string
	ResourceName string
	User         User
	Verb         string
}

// User describes an authenticated user
// TODO: move this into the authentication package
// Example: https://github.com/kubernetes/apiserver/blob/master/pkg/authentication/user/user.go
type User struct {
	Name   string
	Groups []string
}
