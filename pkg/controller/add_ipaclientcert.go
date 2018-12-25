package controller

import (
	"github.com/patrickeasters/ipa-cert-operator/pkg/controller/ipaclientcert"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ipaclientcert.Add)
}
