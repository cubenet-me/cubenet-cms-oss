package model

type NavItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
	Order int    `json:"order"`
}
