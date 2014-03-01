package controllers

import "github.com/revel/revel"

type Feed struct {
  *revel.Controller
}

func (c Feed) Show(feedId int) revel.Result {
  return c.Render(feedId)
}