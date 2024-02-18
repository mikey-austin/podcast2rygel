package media

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

type PodcastImage struct {
	name       string
	path       string
	parentPath string
	imageUrl   string

	mimeType string
	height   int
	width    int
	depth    int
}

func downloadAndInspectImage(url string) (width, height, depth int, mimeType string, err error) {
	// Step 1: Download the image
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, 0, "", err
	}
	defer resp.Body.Close()

	// Step 2: Decode the image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, 0, 0, "", err
	}
	bounds := img.Bounds()
	width, height = bounds.Dx(), bounds.Dy()

	// MIME type could be determined from the format string or the Content-Type header.
	// This example uses the Content-Type header.
	mimeType = resp.Header.Get("Content-Type")

	// TODO: not sure how to handle depth
	return width, height, 0, mimeType, nil
}

func NewPodcastImage(name string, parentPath string, imageUrl string) *PodcastImage {
	width, height, depth, mimeType, err := downloadAndInspectImage(imageUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return &PodcastImage{
		name:       name,
		path:       parentPath + "/" + name,
		parentPath: parentPath,
		imageUrl:   imageUrl,
		height:     height,
		width:      width,
		depth:      depth,
		mimeType:   mimeType}
}

func (pi *PodcastImage) Parent() dbus.ObjectPath {
	return dbus.ObjectPath(pi.parentPath)
}

func (pi *PodcastImage) Type() string {
	return "image"
}

func (pi *PodcastImage) Path() dbus.ObjectPath {
	return dbus.ObjectPath(pi.parentPath + "/" + pi.name)
}

func (pi *PodcastImage) DisplayName() string {
	return pi.name
}

func (pi *PodcastImage) Urls() []string {
	return []string{pi.imageUrl}
}

func (pi *PodcastImage) MimeType() string {
	return pi.mimeType
}

func (pi *PodcastImage) Height() int {
	return pi.height
}

func (pi *PodcastImage) Width() int {
	return pi.width
}

func (pi *PodcastImage) Depth() int {
	return pi.depth
}

func (pi *PodcastImage) Register(conn *dbus.Conn) {
	// Register both org.gnome.UPnP.MediaObject2 and
	// org.gnome.UPnP.MediaItem2 interfaces.
	log.WithField("image", pi).Debug("exported image")
}
