package tokens

import (
	_ "embed"
)

//go:embed insert_token.sql
var InsertTokenQuery string

//go:embed select_token.sql
var SelectTokenQuery string

//go:embed revoke_token_by_id.sql
var RevokeTokenByIDQuery string

//go:embed revoke_tokens_by_user_id.sql
var RevokeTokensByUserID string
