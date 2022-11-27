package rules

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kazhuravlev/just"
)

var reCheck = regexp.MustCompile(`^(?P<name>[-\w]+:[-\w]+)\((?P<args>[\w\d,]+)\)$`)

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

func ParseCheck(s string) (CheckName, string, error) {
	match := reCheck.FindStringSubmatch(s)
	if len(match) != len(reCheck.SubexpNames()) {
		return "", "", errors.New("check is not valid")
	}

	result := make(map[string]string)
	for i, name := range reCheck.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return CheckName(result["name"]), result["args"], nil
}
