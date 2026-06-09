package model

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	Roles     []UserRole `json:"roles"`
	Wallet    UserWallet `json:"wallet"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserRole struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at"`
}

type UserWallet struct {
	Money int64 `json:"money"`
	Spent int64 `json:"spent"`
}
