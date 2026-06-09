package model

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	RoleID    string `json:"role_id"`
	RoleData  *Role  `json:"role_data"`
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

type Role struct {
	ID          string            `json:"id"`
	Identifier  string            `json:"identifier"`
	Name        map[string]string `json:"name"`
	Color       string            `json:"color"`
	Permissions []string          `json:"permissions"`
}
