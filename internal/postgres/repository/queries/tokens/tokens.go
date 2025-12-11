package tokens

import (
	_ "embed"
)

//go:embed insert_token.sql
var InsertTokenQuery string
