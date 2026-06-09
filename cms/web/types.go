package web

type NavItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
	Order int    `json:"order"`
}

type BaseData struct {
	Title           string
	LoggedIn        bool
	Username        string
	Role            string
	RoleName        map[string]string
	RoleColor       string
	Permissions     []string
	SiteName        string
	SiteDescription string
	NavItems        []NavItem
}
