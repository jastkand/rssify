package controllers

import (
  "github.com/revel/revel"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "fmt"
)

type Feed struct {
  *revel.Controller
}

type VKArray struct {
  Id int
  From_id int
  Owner_id int
  Date int
  Post_type string
  Text string
}

type VKResponseBody struct {
  Count int
  Items []VKArray 
}

type VKResponse struct {
  Response VKResponseBody
}

func (c Feed) Show(feedId string) revel.Result {
  var requestUrl string = "https://api.vk.com/method/wall.get?count=2&v=5.12&owner_id=" + feedId

  resp, err := http.Get(requestUrl)

  if err != nil {
    c.Response.Status = 500
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    c.Response.Status = 500
  }

  var encoded VKResponse

  err = json.Unmarshal(body, &encoded)

  if err != nil {
    c.Response.Status = 500
    fmt.Println(err)
  }

  return c.Render(encoded)
}