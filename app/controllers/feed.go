package controllers

import "github.com/revel/revel"
import "net/http"
import "io/ioutil"

type Feed struct {
  *revel.Controller
}

func (c Feed) Show(feedId string) revel.Result {
  var requestUrl string = "https://api.vk.com/method/wall.get?owner_id=" + feedId

  resp, err := http.Get(requestUrl)

  if err != nil {
    c.Response.Status = 500
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    c.Response.Status = 500
  }

  result := string(body)

  return c.Render(result)
}