package users

import (
	_ "embed"
)

//go:embed insert_user.sql
var InsertUserQuery string

//go:embed select_user_by_login.sql
var SelectUserByLoginQuery string
