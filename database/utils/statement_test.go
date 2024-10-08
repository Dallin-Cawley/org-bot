package utils

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type StatementTestSuite struct {
	suite.Suite
}

func (testSuite *StatementTestSuite) TestMakeStatement_Success() {
	statement := NewStatement("SELECT * FROM bababooie", pgx.NamedArgs{"key": "value"})

	testSuite.Assertions.NotEmpty(statement.Statement)
	testSuite.Assertions.NotNil(statement.Args)
}

func (testSuite *StatementTestSuite) TestMakeMethodStatement_Success() {
	argsList := []string{"arg_1", "arg_2"}
	namedArgs := pgx.NamedArgs{"key1": "value1", "key2": "value2"}
	statement := NewMethodStatement("SELECT * FROM", "some_method_name", argsList, namedArgs)

	testSuite.Assertions.NotEmpty(statement.Statement)
	testSuite.Assertions.NotNil(statement.Args)
	testSuite.Assertions.Contains(statement.Statement, "@arg_1, @arg_2")
}

// Test_StatementTestSuite starts the StatementTestSuite
func Test_StatementTestSuite(t *testing.T) {
	suite.Run(t, new(StatementTestSuite))
}
