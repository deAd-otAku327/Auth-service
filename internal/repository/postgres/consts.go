package postgres

const (
	SessionsTable = "sessions"

	SessionIDColumn    = "id"
	UserIDColumn       = "user_id"
	RefreshTokenColumn = "refresh_token"
	UserAgentColumn    = "user_agent"
	IPColumn           = "ip"
	PairIDColumn       = "pair_id"
	ExpiresAtColumn    = "expires_at"
	CreatedAtColumn    = "created_at"
)
