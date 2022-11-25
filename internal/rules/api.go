package rules

import (
	"fmt"

	"github.com/kazhuravlev/just"
)

var gChecks map[CheckName]Check

func RegisterCheck(check Check) error {
	if just.MapContainsKey(gChecks, check.Name()) {
		return fmt.Errorf("check '%s' already registered", check.Name())
	}

	gChecks[check.Name()] = check
	return nil
}

func MustRegisterCheck(check Check) {
	if err := RegisterCheck(check); err != nil {
		panic("registering check: " + err.Error())
	}
}
