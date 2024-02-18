package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/mmcdole/gofeed"
)

// A particular podcast episode, which implements the MediaItem2 interface.
type Episode struct {
	EpisodeDirectory *EpisodeDirectory
	Item             *gofeed.Item
}

func (e *Episode) Parent() dbus.ObjectPath {
	return e.EpisodeDirectory.Path()
}

func (e *Episode) Type() string {
	return "audio"
}

func (e *Episode) Path() dbus.ObjectPath {
	return dbus.ObjectPath(e.EpisodeDirectory.Path() + "/" + dbus.ObjectPath(e.Item.Title))
}

func (e *Episode) DisplayName() string {
	return e.Item.Title
}

func (e *Episode) Urls() []string {
	enclosures := e.Item.Enclosures
	urls := make([]string, len(enclosures))
	for _, enclosure := range enclosures {
		urls = append(urls, enclosure.URL)
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
	return NewPodcastImage(
		e.Item.Image.Title,
		string(e.Path()),
		e.Item.Image.URL)
}

func (e *Episode) Register(conn *dbus.Conn) {
	// Register both org.gnome.UPnP.MediaObject2 and
	// org.gnome.UPnP.MediaItem2 interfaces.

	// Register the episode image
	e.AlbumArt().Register(conn)
}
