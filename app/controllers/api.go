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
		panic(err)
	}

	q := u.Query()

	posts, err := VK.GetPostsByUrl(q.Get("g"))

	if err != nil {
		panic(err)
	}

	revel.INFO.Println(posts)

	c.Response.ContentType = "text/xml"
	return c.RenderText(posts)
}
