package rbac

import (
	"testing"

	"github.com/grepory/storage/authorization"
	"github.com/grepory/storage/testing/mockstore"
	"github.com/stretchr/testify/mock"
)

func TestAuthorize(t *testing.T) {
	type storeFunc func(mock.Arguments)
	tests := []struct {
		name                        string
		attrs                       authorization.RequestAttributes
		clusterRoleBindingStoreFunc storeFunc
		want                        bool
		wantErr                     bool
	}{
		{
			name: "test",
			clusterRoleBindingStoreFunc: func(args mock.Arguments) {
				// clusterRoleBindings := args.Get(1).(*[]rbac.ClusterRoleBinding)
				// *clusterRoleBindings = append(*clusterRoleBindings, rbac.ClusterRoleBinding{})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockstore.MockStore{}
			a := &Authorizer{
				store: store,
			}
			store.On("List", "clusterrolebindings", mock.Anything).Return(nil).Run(tc.clusterRoleBindingStoreFunc)

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
