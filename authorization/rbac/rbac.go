package rbac

import (
	"github.com/grepory/storage/authorization"
	"github.com/grepory/storage/storage"
)

// Authorizer implements an authorizer interface using Roles Base Acccess
// Control (RBAC)
type Authorizer struct {
	store storage.Store
}

// Authorize determines if a request is authorized based on its attributes
func (a *Authorizer) Authorize(attrs authorization.RequestAttributes) (bool, error) {
	// Get cluster roles binding
	// clusterRoleBindings := []rbac.ClusterRoleBinding{}
	// if err := a.store.List("clusterrolebindings", &clusterRoleBindings); err != nil {
	// 	return false, err
	// }

	// Inspect each cluster role binding
	// for _, clusterRoleBinding := range clusterRoleBindings {
	// Verify if this cluster role binding matches our user
	// if matches := matchesTo(attrs.User, clusterRoleBinding.Subjects); !matches {
	// 	continue
	// }

	// Get the cluster role that matched our user
	// role, err := a.store.Get(clusterRoleBinding.RoleRef)
	// if err != nil {
	// 	continue
	// }

	// Loop through the cluste role rules
	// for _, rule := range role.Rules {
	// 	// Verify if this rule applies to our request
	// 	applies := appliesTo(attrs, rule); applies {
	// 		return true, nil
	// 	}
	// }
	// }

	// Do the same thing but for roles binding now

	return false, nil
}
