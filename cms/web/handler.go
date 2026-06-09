package web

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/cubenet-cms/cms/plugin"
	"github.com/cubenet-cms/cms/service"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type Handler struct {
	authSvc   *service.AuthService
	serverSvc *service.ServerService
	newsSvc   *service.NewsService
	pipeline  *plugin.Pipeline
}

func NewHandler(
	authSvc *service.AuthService,
	serverSvc *service.ServerService,
	newsSvc *service.NewsService,
	pipeline *plugin.Pipeline,
) *Handler {
	return &Handler{
		authSvc: authSvc, serverSvc: serverSvc, newsSvc: newsSvc,
		pipeline: pipeline,
	}
}

func (h *Handler) execPipeline(r *http.Request, w http.ResponseWriter, template string) *plugin.Context {
	ctx := &plugin.Context{
		W: w, R: r, Template: template,
		Data: map[string]any{},
	}
	h.pipeline.Execute(ctx)
	return ctx
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	pc := h.execPipeline(r, w, "home")
	homePage(baseData(pc)).Render(context.Background(), w)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	pc := h.execPipeline(r, w, "login")
	loginPage(baseData(pc), "").Render(context.Background(), w)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if username == "" || password == "" {
		pc := h.execPipeline(r, w, "login")
		loginPage(baseData(pc), "Заполните все поля").Render(context.Background(), w)
		return
	}

	result, err := h.authSvc.Login(r.Context(), username, password)
	if err != nil {
		pc := h.execPipeline(r, w, "login")
		loginPage(baseData(pc), "Неверный логин или пароль").Render(context.Background(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token", Value: result.Token,
		Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	pc := h.execPipeline(r, w, "register")
	registerPage(baseData(pc), "").Render(context.Background(), w)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if username == "" || email == "" || password == "" {
		pc := h.execPipeline(r, w, "register")
		registerPage(baseData(pc), "Заполните все поля").Render(context.Background(), w)
		return
	}

	if !usernameRe.MatchString(username) {
		pc := h.execPipeline(r, w, "register")
		registerPage(baseData(pc), "Логин может содержать только буквы, цифры, _ и -").Render(context.Background(), w)
		return
	}

	if len(password) < 6 {
		pc := h.execPipeline(r, w, "register")
		registerPage(baseData(pc), "Пароль должен быть минимум 6 символов").Render(context.Background(), w)
		return
	}

	hasUpper, hasSpecial := false, false
	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			hasSpecial = true
		}
	}
	if !hasUpper || !hasSpecial {
		pc := h.execPipeline(r, w, "register")
		registerPage(baseData(pc), "Пароль должен содержать заглавную букву и спецсимвол").Render(context.Background(), w)
		return
	}

	if password != confirm {
		pc := h.execPipeline(r, w, "register")
		registerPage(baseData(pc), "Пароли не совпадают").Render(context.Background(), w)
		return
	}

	result, err := h.authSvc.Register(r.Context(), username, email, password)
	if err != nil {
		pc := h.execPipeline(r, w, "register")
		registerPage(baseData(pc), err.Error()).Render(context.Background(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token", Value: result.Token,
		Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Servers(w http.ResponseWriter, r *http.Request) {
	pc := h.execPipeline(r, w, "servers")
	serversPage(baseData(pc)).Render(context.Background(), w)
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	pc := h.execPipeline(r, w, "admin")
	adminPage(baseData(pc)).Render(context.Background(), w)
}

func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", staticHandler()).ServeHTTP(w, r)
}

func (h *Handler) Assets(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))).ServeHTTP(w, r)
}

func baseData(pc *plugin.Context) BaseData {
	return BaseData{
		Title:    title(pc.Template),
		LoggedIn: getBool(pc.Data, "LoggedIn"),
		Username: getString(pc.Data, "Username"),
	}
}

func title(t string) string {
	switch t {
	case "home":
		return "Главная"
	case "login":
		return "Вход"
	case "register":
		return "Регистрация"
	case "servers":
		return "Серверы"
	case "admin":
		return "Админка"
	}
	return t
}

func getBool(m map[string]any, k string) bool {
	v, _ := m[k].(bool)
	return v
}

func getString(m map[string]any, k string) string {
	v, _ := m[k].(string)
	return v
}
