package plugin

import "net/http"

type Context struct {
	W        http.ResponseWriter
	R        *http.Request
	Data     map[string]any
	Template string
}

type Plugin interface {
	Name() string
	Hooks() []Hook
}

type Hook func(ctx *Context) error

type Pipeline struct {
	hooks []Hook
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) Register(pl Plugin) {
	p.hooks = append(p.hooks, pl.Hooks()...)
}

func (p *Pipeline) Execute(ctx *Context) error {
	for _, h := range p.hooks {
		if err := h(ctx); err != nil {
			return err
		}
	}
	return nil
}
