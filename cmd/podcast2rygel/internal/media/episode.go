package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// A particular podcast episode, which implements the MediaItem2 interface.
type Episode struct {
	EpisodeDirectory *EpisodeDirectory
	Item             *gofeed.Item
	ItemIndex        int
}

func (e *Episode) Parent() dbus.ObjectPath {
	return e.EpisodeDirectory.Path()
}

func (e *Episode) Type() string {
	return "audio"
}

func (e *Episode) Path() dbus.ObjectPath {
	return dbus.ObjectPath(string(e.EpisodeDirectory.Path()) + "/" + strconv.Itoa(e.ItemIndex))
}

func (e *Episode) DisplayName() string {
	return e.Item.Title
}

func (e *Episode) Urls() []string {
	enclosures := e.Item.Enclosures
	urls := make([]string, len(enclosures))
	for i, enclosure := range enclosures {
		urls[i] = enclosure.URL
	}

	return urls
}

func (e *Episode) MimeType() string {
	return e.Item.Enclosures[0].Type
}

func (e *Episode) Artist() string {
	return e.Item.Author.Name
}

func (e *Episode) Album() string {
	return e.EpisodeDirectory.Feed.Title
}

func (e *Episode) Date() string {
	return e.Item.Published
}

func (e *Episode) AlbumArt() *PodcastImage {
	art := e.EpisodeDirectory.artCache
	if art != nil {
		return art
	}

	art = NewPodcastImage(
		"podcastImage",
		string(e.Path()),
		e.Item.Image.URL)
	e.EpisodeDirectory.artCache = art
	return art
}

func (e *Episode) Register(conn *dbus.Conn) {

	// Register properties for this episode.
	prop.Export(conn, e.Path(), GetProps(e))

	// Register both org.gnome.UPnP.MediaObject2 and
	// org.gnome.UPnP.MediaItem2 interfaces.
	conn.Export(e, e.Path(), "org.gnome.MediaItem2")
	log.WithFields(log.Fields{"episode": e.DisplayName()}).Debug("exported episode")

	// Register the episode image
	e.AlbumArt().Register(conn)
}
