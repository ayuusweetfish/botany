package globals

import (
	"github.com/kawa-yoiko/botany/app/models"

	"io/ioutil"
	"log"
)

// Default avatar
var DefaultAvatar = models.File{
	Id:      -1,
	Type:    "image/png",
	Content: nil,
}

var DefaultBanner = models.File{
	Id:      -1,
	Type:    "image/png",
	Content: nil,
}

func init() {
	content, err := ioutil.ReadFile("globals/default_avatar.png")
	if err != nil {
		log.Fatalln(err)
	}
	DefaultAvatar.Content = content
	DefaultBanner.Content = content
}
