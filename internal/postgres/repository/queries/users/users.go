package users

import (
	_ "embed"
)

//go:embed insert_user.sql
var InsertUserQuery string

//go:embed select_user_by_login.sql
var SelectUserByLoginQuery string

//go:embed select_user_by_id.sql
var SelectUserByIDQuery string

//go:embed update_user.sql
var UpdateUserQuery string

//go:embed delete_user.sql
var DeleteUserQuery string
