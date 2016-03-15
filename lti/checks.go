package lti

import (
	"net/http"

	"k8s.io/kubernetes/pkg/util/sets"
)

// NOTE: all these checks assume the lti launch params have already been
//		 stored in the request context.

func HasLTIRole(r *http.Request, roles ...string) bool {
	launchParams := GetLaunchParams(r)
	if launchParams == nil {
		return false
	}

	roleSet := sets.NewString(launchParams[ParamRoles]...)
	return roleSet.HasAny(roles...)
}

func IsStaff(r *http.Request) bool {
	return HasLTIRole(r, StaffRoles()...)
}
