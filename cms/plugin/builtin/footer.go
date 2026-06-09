package builtin

import (
	"github.com/cubenet-cms/cms/plugin"
	"github.com/cubenet-cms/cms/service"
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

	cookie, err := ctx.R.Cookie("token")
	if err == nil && cookie.Value != "" {
		claims, err := p.authSvc.ValidateToken(cookie.Value)
		if err == nil {
			user, err := p.authSvc.GetProfile(ctx.R.Context(), claims.UserID)
			if err == nil && user != nil {
				perms := user.RoleData.Permissions
				if perms == nil {
					perms = []string{}
				}
				ctx.Data["LoggedIn"] = true
				ctx.Data["Username"] = claims.Username
				ctx.Data["UserID"] = claims.UserID
				ctx.Data["Role"] = user.Role
				ctx.Data["Permissions"] = perms
				ctx.Data["RoleName"] = user.RoleData.Name
				ctx.Data["RoleColor"] = user.RoleData.Color
				return nil
			}
		}
	}

	ctx.Data["LoggedIn"] = false
	ctx.Data["Username"] = ""
	ctx.Data["Permissions"] = []string{}
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
