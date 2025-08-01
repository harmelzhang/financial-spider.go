package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type CommandMainHttpInput struct {
	g.Meta `name:"http" brief:"start http server"`
}

type CommandMainHttpOutput struct{}

func (c *CommandMain) Http(ctx context.Context, in CommandMainHttpInput) (out *CommandMainHttpOutput, err error) {
	s := g.Server()
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.Writeln(r.Router.Uri)
	})
	s.Run()
	return
}
