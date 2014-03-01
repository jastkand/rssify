package controllers

import (
	"github.com/revel/revel"
	"rssify/app/models"
)

type Feed struct {
	*revel.Controller
}

func (c Feed) Show(feedId string) revel.Result {
	rss, error := VK.GetPosts(feedId)

	if error != nil {

	}

	c.Response.ContentType = "text/xml"
	return c.RenderText(rss)
}
