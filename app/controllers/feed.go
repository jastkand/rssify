package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	// "time"
)

type Feed struct {
	*revel.Controller
}

type VKItem struct {
	Id        int
	From_id   int
	Owner_id  int
	Date      int
	Post_type string
	Text      string
}

type VKProfile struct {
	Id          int
	First_name  string
	Last_name   string
	Screen_name string
	Photo_200   string
}

type VKGroup struct {
	Id          int
	Name        string
	Screen_name string
	Is_closed   int
	Type        string
	Photo_200   string
}

type VKResponseBody struct {
	Count    int
	Items    []VKItem
	Profiles []VKProfile
	Groups   []VKGroup
}

type VKResponse struct {
	Response VKResponseBody
}

func (c Feed) Show(feedId string) revel.Result {
	var requestUrl string = "https://api.vk.com/method/wall.get?v=5.12&extended=1&owner_id=" + feedId

	resp, err := http.Get(requestUrl)

	if err != nil {
		c.Response.Status = 500
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var encoded VKResponse

	err = json.Unmarshal(body, &encoded)

	if err != nil {
		c.Response.Status = 500
		fmt.Println(err)
	}

	sourceInfo := encoded.Response.Groups[0]

	// now := time.Now()
	feed := &feeds.Feed{
		Title: sourceInfo.Name,
		Link:  &feeds.Link{Href: "https://vk.com/" + sourceInfo.Screen_name},
		// Description: "Description here",
		// Author:      &feeds.Author{"Jason Moiron", "jmoiron@jmoiron.net"},
		// Created:     now,
	}

	for _, elem := range encoded.Response.Items {
		feed.Add(&feeds.Item{
			Title:       strings.Split(elem.Text, ".")[0] + "...",
			Link:        &feeds.Link{Href: "http://vk.com/wall" + strconv.Itoa(elem.Owner_id) + "_" + strconv.Itoa(elem.Id)},
			Description: elem.Text,
			// Author:      &feeds.Author{"Jason Moiron", "jmoiron@jmoiron.net"},
			// Created:     now,
		})
	}

	rss, err := feed.ToRss()

	if err != nil {
		c.Response.Status = 500
	}

	c.Response.ContentType = "text/xml"
	return c.RenderText(rss)
}
