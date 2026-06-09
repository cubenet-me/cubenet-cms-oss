package builtin

import (
	"encoding/json"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/plugin"
	"github.com/cubenet-cms/cms/service"
)

type SettingsPlugin struct {
	settingsSvc *service.SettingsService
}

func NewSettingsPlugin(settingsSvc *service.SettingsService) *SettingsPlugin {
	return &SettingsPlugin{settingsSvc: settingsSvc}
}

func (p *SettingsPlugin) Name() string { return "settings" }

func (p *SettingsPlugin) Hooks() []plugin.Hook {
	return []plugin.Hook{p.settingsHook}
}

func (p *SettingsPlugin) settingsHook(ctx *plugin.Context) error {
	if ctx.Data == nil {
		ctx.Data = make(map[string]any)
	}
	ctx.Data["SiteName"] = p.settingsSvc.Get("site_name", "CubeNet CMS")
	ctx.Data["SiteDescription"] = p.settingsSvc.Get("site_description", "")

	raw := p.settingsSvc.Get("nav_items", `[{"label":"Главная","href":"/","order":0},{"label":"Серверы","href":"/servers","order":1}]`)
	var items []model.NavItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		items = []model.NavItem{{Label: "Главная", Href: "/", Order: 0}, {Label: "Серверы", Href: "/servers", Order: 1}}
	}
	ctx.Data["NavItems"] = items
	return nil
}
