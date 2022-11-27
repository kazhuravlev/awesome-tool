package rules

import (
	"errors"
	"fmt"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"regexp"
	"strconv"

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

func parseCheckCall(s string) (CheckName, string, error) {
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

func ParseCheck(s string) (Check, error) {
	checkName, checkArgs, err := parseCheckCall(s)
	if err != nil {
		return nil, errorsh.Wrap(err, "parse check call")
	}

	switch checkName {
	default:
		return nil, errorsh.Newf("unknown check '%s'", checkName)
	case "response:status-eq":
		statusCode, err := strconv.Atoi(checkArgs)
		if err != nil {
			return nil, errorsh.Wrapf(err, "bad arguments for check '%s'", checkName)
		}

		return CheckResponseStatusEq{
			ExpectedStatus: statusCode,
		}, nil
	case "github-repo:stars-min":
		minimumStars, err := strconv.Atoi(checkArgs)
		if err != nil {
			return nil, errorsh.Wrapf(err, "bad arguments for check '%s'", checkName)
		}

		return CheckResponseStatusEq{
			ExpectedStatus: minimumStars,
		}, nil
	}
}
