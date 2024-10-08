package utils

import (
	"fmt"
	"os"
)

type Statement struct {
	Statement string
	Args      interface{}
}

// NewStatement creates a Statement.
func NewStatement(statement string, args interface{}) *Statement {
	return &Statement{Statement: statement, Args: args}
}

// NewMethodStatement creates a Statement used to trigger a Stored Procedure or Function
func NewMethodStatement(prefix string, methodName string, argList []string, args interface{}) *Statement {
	statementString := fmt.Sprintf("%s %s.%s (%s)", prefix, os.Getenv("SCHEMA"), methodName, formatArgList(argList))

	return &Statement{Statement: statementString, Args: args}
}

// formatArgList formats the namedArgs such that it can be used with the pgx library for named args. Namely,
// in this format "@arg, "
func formatArgList(argList []string) string {
	formatted := ""
	for i, arg := range argList {
		if i == 0 {
			formatted = fmt.Sprintf("@%s", arg)
		} else {
			formatted = fmt.Sprintf("%s, @%s", formatted, arg)
		}
	}

	return formatted
}
