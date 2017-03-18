package cfg

import (
	"os"
	"regexp"
	"strconv"
)

type spec struct {
	environmentVariableName string
	programArgument         string
	defaultValue            string
}

const (
	ARG_REGEX_STR = `^--([a-z-_.]*)(?:=(.*))?$`
)

var (
	parsedExecuteArguments map[string]string
	argumentRegex          = regexp.MustCompile(ARG_REGEX_STR)
)

func GetValue(key configType) string {
	spec := configurations[key]
	if spec == nil {
		return ""
	}

	if spec.programArgument != "" {
		value, ok := getExecuteArgumentValue(spec.programArgument)
		if ok {
			return value
		}
	}

	if spec.environmentVariableName != "" {
		value := os.Getenv(spec.environmentVariableName)
		if value != "" {
			return value
		}
	}

	return spec.defaultValue
}

func GetIntValue(key configType) int {
	v := GetValue(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return i
}

func getExecuteArgumentValue(argument string) (string, bool) {
	args := getArguments()
	arg, ok := args[argument]
	return arg, ok
}

func getArguments() map[string]string {
	if parsedExecuteArguments == nil {
		parsedExecuteArguments = parseArguments()
	}

	return parsedExecuteArguments
}

func parseArguments() map[string]string {
	args := make(map[string]string)
	for _, osArg := range os.Args[1:] {
		subMatch := argumentRegex.FindStringSubmatch(osArg)
		if len(subMatch) == 0 {
			continue
		}

		args[subMatch[1]] = subMatch[2]
	}

	return args
}
