package rbac

import (
	"errors"
	"testing"

	"github.com/grepory/storage/apis/rbac"
	"github.com/grepory/storage/authorization"
	"github.com/grepory/storage/testing/mockstore"
	"github.com/stretchr/testify/mock"
)

func TestAuthorize(t *testing.T) {
	type storeFunc func(*mockstore.MockStore)
	tests := []struct {
		name      string
		attrs     authorization.RequestAttributes
		storeFunc storeFunc
		want      bool
		wantErr   bool
	}{
		{
			name:  "no bindings",
			attrs: authorization.RequestAttributes{Namespace: "acme"},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil)
				store.On("List", "rolebindings/acme", mock.Anything).Return(nil)
			},
			want: false,
		},
		{
			name: "clusterrolebindings store err",
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).
					Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "no matching clusterRoleBindings",
			attrs: authorization.RequestAttributes{
				Namespace: "acme",
				User: authorization.User{
					Name: "foo",
				},
			},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						clusterRoleBindings := args.Get(1).(*[]rbac.ClusterRoleBinding)
						*clusterRoleBindings = append(*clusterRoleBindings, rbac.ClusterRoleBinding{
							Subjects: []rbac.Subject{
								rbac.Subject{Kind: rbac.UserKind, Name: "bar"},
							},
						})
					})
				store.On("List", "rolebindings/acme", mock.Anything).Return(nil)
			},
			want: false,
		},
		{
			name: "clusterroles/admin store err",
			attrs: authorization.RequestAttributes{
				User: authorization.User{
					Name: "foo",
				},
			},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						clusterRoleBindings := args.Get(1).(*[]rbac.ClusterRoleBinding)
						*clusterRoleBindings = append(*clusterRoleBindings, rbac.ClusterRoleBinding{
							RoleRef: rbac.RoleRef{
								Name: "admin",
							},
							Subjects: []rbac.Subject{
								rbac.Subject{Kind: rbac.UserKind, Name: "foo"},
							},
						})
					})
				store.On("Get", "clusterroles/admin", mock.Anything).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "matching clusterRoleBinding",
			attrs: authorization.RequestAttributes{
				Verb:         "create",
				APIGroup:     "checks.sensu.io",
				Resource:     "checks",
				ResourceName: "check-cpu",
				User: authorization.User{
					Name: "foo",
				},
			},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						clusterRoleBindings := args.Get(1).(*[]rbac.ClusterRoleBinding)
						*clusterRoleBindings = append(*clusterRoleBindings, rbac.ClusterRoleBinding{
							RoleRef: rbac.RoleRef{
								Name: "admin",
							},
							Subjects: []rbac.Subject{
								rbac.Subject{Kind: rbac.UserKind, Name: "foo"},
							},
						})
					})
				store.On("Get", "clusterroles/admin", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						clusterRole := args.Get(1).(*rbac.ClusterRole)
						*clusterRole = rbac.ClusterRole{Rules: []rbac.Rule{
							rbac.Rule{
								Verbs:         []string{"create"},
								APIGroups:     []string{"checks.sensu.io"},
								Resources:     []string{"checks"},
								ResourceNames: []string{"check-cpu"},
							},
						}}
					})
			},
			want: true,
		},
		{
			name:  "rolebindings store err",
			attrs: authorization.RequestAttributes{Namespace: "acme"},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil)
				store.On("List", "rolebindings/acme", mock.Anything).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "no matching roleBindings",
			attrs: authorization.RequestAttributes{
				Namespace: "acme",
				User: authorization.User{
					Name: "foo",
				},
			},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil)
				store.On("List", "rolebindings/acme", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						roleBindings := args.Get(1).(*[]rbac.RoleBinding)
						*roleBindings = append(*roleBindings, rbac.RoleBinding{
							Subjects: []rbac.Subject{
								rbac.Subject{Kind: rbac.UserKind, Name: "bar"},
							},
						})
					})
			},
			want: false,
		},
		{
			name: "roles/admin store err",
			attrs: authorization.RequestAttributes{
				Namespace: "acme",
				User: authorization.User{
					Name: "foo",
				},
			},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil)
				store.On("List", "rolebindings/acme", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						roleBindings := args.Get(1).(*[]rbac.RoleBinding)
						*roleBindings = append(*roleBindings, rbac.RoleBinding{
							RoleRef: rbac.RoleRef{Name: "admin"},
							Subjects: []rbac.Subject{
								rbac.Subject{Kind: rbac.UserKind, Name: "foo"},
							},
						})
					})
				store.On("Get", "roles/acme/admin", mock.Anything).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "matching roleBinding",
			attrs: authorization.RequestAttributes{
				Namespace: "acme",
				User: authorization.User{
					Name: "foo",
				},
				Verb:         "create",
				APIGroup:     "checks.sensu.io",
				Resource:     "checks",
				ResourceName: "check-cpu",
			},
			storeFunc: func(store *mockstore.MockStore) {
				store.On("List", "clusterrolebindings", mock.Anything).Return(nil)
				store.On("List", "rolebindings/acme", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						roleBindings := args.Get(1).(*[]rbac.RoleBinding)
						*roleBindings = append(*roleBindings, rbac.RoleBinding{
							RoleRef: rbac.RoleRef{Name: "admin"},
							Subjects: []rbac.Subject{
								rbac.Subject{Kind: rbac.UserKind, Name: "foo"},
							},
						})
					})
				store.On("Get", "roles/acme/admin", mock.Anything).Return(nil).
					Run(func(args mock.Arguments) {
						role := args.Get(1).(*rbac.Role)
						*role = rbac.Role{Rules: []rbac.Rule{
							rbac.Rule{
								Verbs:         []string{"create"},
								APIGroups:     []string{"checks.sensu.io"},
								Resources:     []string{"checks"},
								ResourceNames: []string{"check-cpu"},
							},
						}}
					})
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockstore.MockStore{}
			a := &Authorizer{
				store: store,
			}
			tc.storeFunc(store)

			got, err := a.Authorize(tc.attrs)
			if (err != nil) != tc.wantErr {
				t.Errorf("Authorizer.Authorize() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("Authorizer.Authorize() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMatchesUser(t *testing.T) {
	tests := []struct {
		name     string
		user     authorization.User
		subjects []rbac.Subject
		want     bool
	}{
		{
			name: "not matching",
			user: authorization.User{Name: "foo"},
			subjects: []rbac.Subject{
				rbac.Subject{Kind: rbac.UserKind, Name: "bar"},
				rbac.Subject{Kind: rbac.GroupKind, Name: "foo"},
			},
			want: false,
		},
		{
			name: "matching via username",
			user: authorization.User{Name: "foo"},
			subjects: []rbac.Subject{
				rbac.Subject{Kind: rbac.UserKind, Name: "bar"},
				rbac.Subject{Kind: rbac.UserKind, Name: "foo"},
			},
			want: true,
		},
		{
			name: "matching via group",
			user: authorization.User{Name: "foo", Groups: []string{"acme"}},
			subjects: []rbac.Subject{
				rbac.Subject{Kind: rbac.GroupKind, Name: "acme"},
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := matchesUser(tc.user, tc.subjects); got != tc.want {
				t.Errorf("matchesUser() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRuleAllows(t *testing.T) {
	tests := []struct {
		name  string
		attrs authorization.RequestAttributes
		rule  rbac.Rule
		want  bool
	}{
		{
			name: "verb does not match",
			attrs: authorization.RequestAttributes{
				Verb: "create",
			},
			rule: rbac.Rule{
				Verbs: []string{"get"},
			},
			want: false,
		},
		{
			name: "api group does not match",
			attrs: authorization.RequestAttributes{
				Verb:     "create",
				APIGroup: "rbac.authorization.sensu.io",
			},
			rule: rbac.Rule{
				Verbs:     []string{"create"},
				APIGroups: []string{"core.sensu.io"},
			},
			want: false,
		},
		{
			name: "resource does not match",
			attrs: authorization.RequestAttributes{
				Verb:     "create",
				APIGroup: "core.sensu.io",
				Resource: "events",
			},
			rule: rbac.Rule{
				Verbs:     []string{"create"},
				APIGroups: []string{"core.sensu.io"},
				Resources: []string{"checks", "handlers"},
			},
			want: false,
		},
		{
			name: "resource name does not match",
			attrs: authorization.RequestAttributes{
				Verb:         "create",
				APIGroup:     "core.sensu.io",
				Resource:     "checks",
				ResourceName: "check-cpu",
			},
			rule: rbac.Rule{
				Verbs:         []string{"create"},
				APIGroups:     []string{"core.sensu.io"},
				Resources:     []string{"checks"},
				ResourceNames: []string{"check-mem"},
			},
			want: false,
		},
		{
			name: "matches",
			attrs: authorization.RequestAttributes{
				Verb:         "create",
				APIGroup:     "checks.sensu.io",
				Resource:     "checks",
				ResourceName: "check-cpu",
			},
			rule: rbac.Rule{
				Verbs:         []string{"create"},
				APIGroups:     []string{"checks.sensu.io"},
				Resources:     []string{"checks"},
				ResourceNames: []string{"check-cpu"},
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ruleAllows(tc.attrs, tc.rule); got != tc.want {
				t.Errorf("ruleAllows() = %v, want %v", got, tc.want)
			}
		})
	}
}
