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
	url.ParseRequestURI(c.Request.RequestURI)

	if err != nil {
		panic(err)
	}

	q := url.Query()

	posts, err := VK.getPostsByUrl(q.Get("g"))

	if err != nil {
		panic(err)
	}

	revel.INFO.Println(posts)

	c.Response.ContentType = "text/xml"
	return c.RenderText(posts)
}
