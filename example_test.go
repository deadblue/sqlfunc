package sqlfunc_test

import (
	"context"
	"database/sql"
	"log"

	"github.com/deadblue/sqlfunc"
)

type (
	QueryParams struct {
		UserId string
		Status int
	}

	UserResult struct {
		UserId    string
		FirstName string
		// sql.NullXXX types are supported.
		LastName sql.NullString
		// Mapping "sex" column to [Gender] field.
		Gender int `db:"sex"`
	}
)

func Example() {
	// Make SQL function
	queryUser, err := sqlfunc.MakeQueryFunc[QueryParams, UserResult](
		"SELECT user_id, first_name, last_name, sex",
		"FROM tbl_user",
		"WHERE user_id = {{ .UserID }} AND status = {{ .Status }}",
	)
	if err != nil {
		panic(err)
	}

	// Connect to database
	db, err := sql.Open("driver", "DSN")
	if err != nil {
		panic(err)
	}
	// Put DB to context
	ctx := sqlfunc.NewContext(context.TODO(), db)

	// Execute query
	if user, err := queryUser(ctx, QueryParams{
		UserId: "123",
		Status: 1,
	}); err != nil {
		panic(err)
	} else {
		// Process result
		if user.LastName.Valid {
			log.Printf("Found user: %s-%s", user.FirstName, user.LastName.String)
		} else {
			log.Printf("Found user: %s", user.FirstName)
		}
	}
}
