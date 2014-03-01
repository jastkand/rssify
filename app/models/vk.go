package VK

import (
	"encoding/json"
	"github.com/gorilla/feeds"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type VK struct{}

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

func GetPosts(feedId string) (string, error) {
	var requestUrl string = "https://api.vk.com/method/wall.get?v=5.12&extended=1&owner_id=" + feedId

	resp, err := http.Get(requestUrl)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var encoded VKResponse

	err = json.Unmarshal(body, &encoded)

	if err != nil {
		return "", err
	}

	sourceInfo := encoded.Response.Groups[0]

	feed := &feeds.Feed{
		Title: sourceInfo.Name,
		Link:  &feeds.Link{Href: "https://vk.com/" + sourceInfo.Screen_name},
	}

	for _, elem := range encoded.Response.Items {
		feed.Add(&feeds.Item{
			Title:       strings.Split(elem.Text, ".")[0] + "...",
			Link:        &feeds.Link{Href: "http://vk.com/wall" + strconv.Itoa(elem.Owner_id) + "_" + strconv.Itoa(elem.Id)},
			Description: elem.Text,
		})
	}

	return feed.ToRss()
}
