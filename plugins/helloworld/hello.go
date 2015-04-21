package hello

// 这是一个用来演示plugin的例子

import (
	"github.com/guotie/deferinit"
	"github.com/smtc/justcms/plugins"
)

type HelloPlugin struct {
	plugins.PluginModel
}

func init() {
	deferinit.AddInit(func() {
		plugins.Set(&HelloPlugin{
			plugins.PluginModel{
				Name: "github.com/smtc/justcms/plugins/hello",
			},
		})
	}, nil, 0)
}

func (hp *HelloPlugin) Name() string {
	return "github.com/smtc/justcms/plugins/hello"
}

func (hp *HelloPlugin) Activate() {

}
func (hp *HelloPlugin) Deactivate() {

}
func (hp *HelloPlugin) Config() {

}
