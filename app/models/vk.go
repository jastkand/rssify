package models

import (
	"encoding/json"
	"github.com/fugazister/feeds"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type VKFeed struct {
	FeedUrl string
}

type VKVideo struct {
	Id   int
	Link string
}

type VKPhoto struct {
	Album_id   int
	Owner_id   int
	Photo_75   string
	Photo_130  string
	Photo_604  string
	Photo_807  string
	Photo_1280 string
}

type VKAudio struct {
	Id  int
	Url string
}

type VKAttachment struct {
	Type  string
	Photo VKPhoto
	Audio VKAudio
	Video VKVideo
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

type SourceInfo struct {
	Name        string
	Screen_name string
	First_name  string
	Last_name   string
}

type SourceInfoContainer struct {
	Response []SourceInfo
}

type VKAttachmentListItem struct {
	Url  string
	Type string
}

type VKAttachmentList struct {
	Items []*VKAttachmentListItem
}

func processAttachments(attachmants []VKAttachment) VKAttachmentList {
	var attachmentList VKAttachmentList

	if len(attachmants) > 0 {
		for _, attachment := range attachmants {
			if attachment.Type == "photo" {
				var photo string = ""

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

				attachmentList.Items = append(attachmentList.Items, &VKAttachmentListItem{photo, "photo"})
			}

			if attachment.Type == "audio" {
				attachmentList.Items = append(attachmentList.Items, &VKAttachmentListItem{attachment.Audio.Url, "audio"})
			}
		}
	}

	return attachmentList
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

	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

	}

	var encoded ResolvedScreenNameResponse

	err = json.Unmarshal(body, &encoded)

	if err != nil {

	}

	return encoded.Response
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

func getPosts(feedUrl string) (string, error) {
	var requestUrl string = "https://api.vk.com/method/wall.get?v=5.12&extended=1&owner_id=" + feedUrl

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

	name, screenName := getSourceInfo(feedUrl)

	feed := &feeds.Feed{
		Title: name,
		Link:  &feeds.Link{Href: "https://vk.com/" + screenName},
	}

	for _, elem := range encoded.Response.Items {
		var description string = ""
		var screenName, name string = "", ""
		var attachments string = ""

		attachmentList := processAttachments(elem.Attachments)

		description += elem.Text

		if len(elem.Copy_history) > 0 {
			description += elem.Copy_history[0].Text
			attachmentList = processAttachments(elem.Copy_history[0].Attachments)
			name, screenName = getSourceInfo(strconv.Itoa(elem.Copy_history[0].Owner_id))
		}

		for _, attachment := range attachmentList.Items {
			if attachment.Type == "photo" {
				attachments += "<br/><img src='" + attachment.Url + "'/>"
			} else if attachment.Type == "audio" {
				attachments += "<br/><source src='" + attachment.Url + "' type='audio/mpeg; codecs='mp3' />"
			}
		}

		item := &feeds.Item{
			Author:      &feeds.Author{Name: name, Email: "https://vk.com/" + screenName},
			Title:       strings.Split(elem.Text, ".")[0] + "...",
			Link:        &feeds.Link{Href: "http://vk.com/wall" + strconv.Itoa(elem.Owner_id) + "_" + strconv.Itoa(elem.Id)},
			Description: description + attachments,
			Created:     time.Unix(int64(elem.Date), int64(0)),
		}

		/*		for _, attachment := range attachmentList.Items {
				if attachment.Type == "audio/mpeg" {
					enclosure := &feeds.Enclosure{attachment.Url, attachment.Type}
					item.AddEnclosure(enclosure)
				}
			}*/

		feed.Add(item)
	}

	return feed.ToRss()
}

func getPostsByUrl(feedUrl string) (string, error) {
	rp := regexp.MustCompile("vk.com/(\\w+)")
	result := rp.FindAllStringSubmatch(feedUrl, -1)
	screenName := resolveScreenName(result[0][1])
	var resolvedFeedId string

	if screenName.Type != "user" {
		resolvedFeedId = "-" + strconv.Itoa(int(screenName.Object_id))
	} else {
		resolvedFeedId = strconv.Itoa(int(screenName.Object_id))
	}

	return getPosts(resolvedFeedId)
}

func (v VKFeed) GetFeed() (string, error) {
	return getPostsByUrl(v.FeedUrl)
}
