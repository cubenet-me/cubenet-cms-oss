package model

type NavItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
	Icon  string `json:"icon"`
	Order int    `json:"order"`
}
