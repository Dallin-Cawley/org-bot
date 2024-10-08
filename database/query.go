package database

import (
	"context"
	"log"
	"strings"

	"orgBot/database/model"
	databaseUtils "orgBot/database/utils"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// ReadModel will read the record(s) from the database identified by the utils.Statement given by
// model.Model#GetReadStatement.
func ReadModel[T model.Model](theModel T, transaction pgx.Tx) ([]T, error) {
	return query[T](theModel.GetReadStatement(), transaction)
}

// InsertModel will create a new record in the database from the utils.Statement provided by the
// implemented model.Model#GetInsertStatement() method.
func InsertModel[T model.Model](theModel T, transaction pgx.Tx) ([]T, error) {
	return query[T](theModel.GetInsertStatement(), transaction)
}

// DeleteModel will delete the record(s) from the utils.Statement provided by the implemented
// model.Model#GetDeleteStatement() method.
func DeleteModel[T model.Model](theModel T, transaction pgx.Tx) ([]T, error) {
	return query[T](theModel.GetDeleteStatement(), transaction)
}

func query[T model.Model](statement *databaseUtils.Statement, transaction pgx.Tx) ([]T, error) {
	rows, err := transaction.Query(context.Background(), statement.Statement, statement.Args)
	if err != nil {
		log.Printf("executing query resulted in an error [ %s ]\n", err.Error())
		return nil, err
	}

	models, err := pgx.CollectRows[T](rows, pgx.RowToStructByName[T])
	if err != nil {
		log.Printf("scanning query results produced an error [ %s ]\n", err.Error())
		if strings.Contains(err.Error(), "number of field descriptions") && strings.Contains(err.Error(), "and 0") {
			return nil, errors.New(databaseUtils.NO_ROWS_RETURNED)
		}

		return nil, err
	}

	return models, nil
}
