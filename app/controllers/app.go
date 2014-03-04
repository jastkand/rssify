package controllers

import (
	"github.com/revel/revel"
)

func defaultRenderArgs(c *revel.Controller) revel.Result {
	c.RenderArgs["GoogleAnalytics"] = revel.Config.StringDefault("analytics.ga", "")
	return nil
}

func init() {
	revel.InterceptFunc(defaultRenderArgs, revel.BEFORE, &App{})
}

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}
