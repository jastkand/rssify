package VK

import (
	"encoding/json"
	"github.com/gorilla/feeds"
	"io/ioutil"
	"net/http"
	"regexp"
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
	Id           int
	From_id      int
	Owner_id     int
	Date         int
	Post_type    string
	Text         string
	Copy_history []VKItem
	Attachments  []VKAttachment
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
			}
			result += "<br/><img src='" + photo + "'/>"
		}
		return result
	} else {
		return ""
	}
}

type ResolvedScreenName struct {
	Type      string
	Object_id float64
}

type ResolvedScreenNameResponse struct {
	Response ResolvedScreenName
}

func resolveScreenName(screenName string) ResolvedScreenName {
	var requestUrl = "https://api.vk.com/method/utils.resolveScreenName?v=5.12&screen_name=" + screenName
	resp, err := http.Get(requestUrl)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var encoded ResolvedScreenNameResponse

	err = json.Unmarshal(body, &encoded)

	if err != nil {
		panic(err)
	}

	return encoded.Response
}

type SourceInfo struct {
	Name        string
	Screen_name string
	First_name  string
	Last_name   string
}

type SourceInfoContainer struct {
	Response []SourceInfo
}

func getSourceInfo(feedId string) (string, string) {
	var isGroup bool = strings.Contains(feedId, "-")
	var groupUrl string = "https://api.vk.com/method/groups.getById?group_id="
	var profileUrl string = "https://api.vk.com/method/users.get?user_ids="
	var requestUrl string

	if isGroup {
		feedId = feedId[1:]
		requestUrl = groupUrl + feedId
	} else {
		requestUrl = profileUrl + feedId
	}

	resp, err := http.Get(requestUrl)

	if err != nil {

	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

	}
	var encoded SourceInfoContainer
	err = json.Unmarshal(body, &encoded)

	if len(encoded.Response[0].Name) > 0 {
		return encoded.Response[0].Name, encoded.Response[0].Screen_name
	} else {
		return encoded.Response[0].First_name + " " + encoded.Response[0].Last_name, feedId
	}
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

	name, screenName := getSourceInfo(feedId)

	feed := &feeds.Feed{
		Title: name,
		Link:  &feeds.Link{Href: "https://vk.com/" + screenName},
	}

	for _, elem := range encoded.Response.Items {
		var description string = ""
		var screenName, name string = "", ""
		photo := processAttachments(elem.Attachments)
		description += elem.Text + photo

		if len(elem.Copy_history) > 0 {
			description += "<br/>repost<br/>" + elem.Copy_history[0].Text
			name, screenName = getSourceInfo(strconv.Itoa(elem.Copy_history[0].Owner_id))
		}
		feed.Add(&feeds.Item{
			Author:      &feeds.Author{Name: name, Email: "https://vk.com/" + screenName},
			Title:       strings.Split(elem.Text, ".")[0] + "...",
			Link:        &feeds.Link{Href: "http://vk.com/wall" + strconv.Itoa(elem.Owner_id) + "_" + strconv.Itoa(elem.Id)},
			Description: description,
			Created:     time.Unix(int64(elem.Date), int64(0)),
		})
	}

	return feed.ToRss()
}

func GetPostsByUrl(feedUrl string) (string, error) {
	rp := regexp.MustCompile("vk.com/(\\w+)")
	result := rp.FindAllStringSubmatch(feedUrl, -1)
	screenName := resolveScreenName(result[0][1])
	var resolvedFeedId string

	if screenName.Type != "user" {
		resolvedFeedId = "-" + strconv.Itoa(int(screenName.Object_id))
	} else {
		resolvedFeedId = strconv.Itoa(int(screenName.Object_id))
	}

	return GetPosts(resolvedFeedId)
}
