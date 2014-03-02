package VK

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type VK struct{}

type VKPhoto struct {
	Album_id   int
	Owner_id   int
	Photo_75   string
	Photo_130  string
	Photo_604  string
	Photo_807  string
	Photo_1280 string
}

type VKAttachment struct {
	Type  string
	Photo VKPhoto
}

type VKItem struct {
	Id          int
	From_id     int
	Owner_id    int
	Date        int
	Post_type   string
	Text        string
	Attachments []VKAttachment
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

func processAttachments(attachments []VKAttachment) string {
	if len(attachments) > 0 {
		var result, photo string = "", ""
		for _, attachment := range attachments {
			if attachment.Type == "photo" {

				if attachment.Photo.Photo_1280 != "" {
					photo = attachment.Photo.Photo_1280
				} else if attachment.Photo.Photo_807 != "" {
					photo = attachment.Photo.Photo_807
				} else if attachment.Photo.Photo_604 != "" {
					photo = attachment.Photo.Photo_604
				} else if attachment.Photo.Photo_130 != "" {
					photo = attachment.Photo.Photo_130
				} else if attachment.Photo.Photo_75 != "" {
					photo = attachment.Photo.Photo_75
				}

				fmt.Println(photo)
			}
		}
		result = "<br/><img src='" + photo + "'/>"
		return result
	} else {
		return ""
	}
}

func GetPosts(feedId string) (string, error) {
	var requestUrl string = "https://api.vk.com/method/wall.get?v=5.12&extended=1&owner_id=" + feedId
	var isGroup bool = strings.Contains(feedId, "-")

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
	var name, screenName string

	if isGroup {
		sourceInfo := encoded.Response.Groups[0]
		name = sourceInfo.Name
		screenName = sourceInfo.Screen_name
	} else {
		sourceInfo := encoded.Response.Profiles[0]
		name = sourceInfo.First_name + " " + sourceInfo.Last_name
		screenName = sourceInfo.Screen_name
	}

	feed := &feeds.Feed{
		Title: name,
		Link:  &feeds.Link{Href: "https://vk.com/" + screenName},
	}

	for _, elem := range encoded.Response.Items {
		photo := processAttachments(elem.Attachments)
		feed.Add(&feeds.Item{
			Title:       strings.Split(elem.Text, ".")[0] + "...",
			Link:        &feeds.Link{Href: "http://vk.com/wall" + strconv.Itoa(elem.Owner_id) + "_" + strconv.Itoa(elem.Id)},
			Description: elem.Text + photo,
			Created:     time.Unix(int64(elem.Date), int64(0)),
		})
	}

	return feed.ToRss()
}
