package builtin

import (
	"github.com/cubenet-cms/backend/plugin"
	"github.com/cubenet-cms/backend/service"
)

type SessionPlugin struct {
	authSvc *service.AuthService
}

func NewSessionPlugin(authSvc *service.AuthService) *SessionPlugin {
	return &SessionPlugin{authSvc: authSvc}
}

func (p *SessionPlugin) Name() string { return "session" }

func (p *SessionPlugin) Hooks() []plugin.Hook {
	return []plugin.Hook{p.sessionHook}
}

func (p *SessionPlugin) sessionHook(ctx *plugin.Context) error {
	if ctx.Data == nil {
		ctx.Data = make(map[string]any)
	}
	// session handling would go here (cookie → user)
	// for now just set defaults
	ctx.Data["LoggedIn"] = false
	ctx.Data["Username"] = ""
	return nil
}

type FooterPlugin struct{}

func NewFooterPlugin() *FooterPlugin { return &FooterPlugin{} }

func (p *FooterPlugin) Name() string { return "footer" }

func (p *FooterPlugin) Hooks() []plugin.Hook {
	return []plugin.Hook{p.footerHook}
}

func (p *FooterPlugin) footerHook(ctx *plugin.Context) error {
	if ctx.Data == nil {
		ctx.Data = make(map[string]any)
	}
	ctx.Data["License"] = "CubeNet CMS License v1.0"
	return nil
}
