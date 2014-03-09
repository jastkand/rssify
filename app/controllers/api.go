package controllers

import (
	"github.com/revel/revel"
	"net/url"
	"rssify/app/models"
)

type Api struct {
	App
}

func (c Api) Show() revel.Result {
	u, err := url.ParseRequestURI(c.Request.RequestURI)

	if err != nil {

	}

	q := u.Query()

	vk := &models.VKFeed{q.Get("g")}
	posts, err := vk.GetFeed()

	if err != nil {

	}

	c.Response.ContentType = "application/xml"
	return c.RenderText(posts)
}
