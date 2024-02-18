package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/mmcdole/gofeed"
)

// Another MediaContainer2 implementation that represents a single podcast.
type EpisodeDirectory struct {
	ParentContainer *PodcastDirectory
	Feed            *gofeed.Feed
}

func (ed *EpisodeDirectory) Parent() dbus.ObjectPath {
	return ed.ParentContainer.Path()
}

func (ed *EpisodeDirectory) Path() dbus.ObjectPath {
	return dbus.ObjectPath(ed.ParentContainer.path + "/" + ed.DisplayName())
}

func (ed *EpisodeDirectory) Type() string {
	return "container"
}

func (ed *EpisodeDirectory) DisplayName() string {
	return ed.Feed.Title
}

func (ed *EpisodeDirectory) ChildCount() int {
	// Child count equals item count as we have no containers
	return ed.ItemCount()
}

func (ed *EpisodeDirectory) ContainerCount() int {
	// Episode directories contain only episode leaf items and no containers
	return 0
}

func (ed *EpisodeDirectory) ItemCount() int {
	return ed.Feed.Len()
}

func (ed *EpisodeDirectory) Searchable() bool {
	return false
}

func (ed *EpisodeDirectory) Register(conn *dbus.Conn) {
	// Register both org.gnome.UPnP.MediaObject2 and
        // org.gnome.UPnP.MediaContainer2 interfaces.

	// Register each episode.
	for _, episode := range ed.ListEpisodes() {
		episode.Register(conn)
	}
}

func (ed *EpisodeDirectory) ListEpisodes() []*Episode {
	episodes := make([]*Episode, ed.Feed.Len())
	for _, item := range ed.Feed.Items {
		episode := &Episode{EpisodeDirectory: ed, Item: item}
		episodes = append(episodes, episode)
	}
	return episodes
}

func (ed *EpisodeDirectory) ListContainers(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	// Episode directories only contain episodes, and no containers
	return nil, nil
}

func (ed *EpisodeDirectory) ListItems(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	items := ed.Feed.Items[offset : offset+max]
	children := make([]map[string]dbus.Variant, len(items))
	for _, item := range items {
		parent := MediaContainer2(ed)
		child := MediaItem2(&Episode{EpisodeDirectory: ed, Item: item})
		children = append(children, map[string]dbus.Variant{
			// Media object attributes
			"Parent":      dbus.MakeVariant(parent.Path),
			"Type":        dbus.MakeVariant(child.Type()),
			"Path":        dbus.MakeVariant(child.Path()),
			"DisplayName": dbus.MakeVariant(child.DisplayName()),

			// Item attributes
			"URLs":     dbus.MakeVariant(child.Urls()),
			"MIMEType": dbus.MakeVariant(child.MimeType()),
		})
	}
	return children, nil
}

// TODO: Take the filter into account and only return requested keys in the output.
func (ed *EpisodeDirectory) ListChildren(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	items := ed.Feed.Items[offset : offset+max]
	children := make([]map[string]dbus.Variant, len(items))
	for _, item := range items {
		parent := MediaContainer2(ed)
		child := MediaObject2(&Episode{EpisodeDirectory: ed, Item: item})
		children = append(children, map[string]dbus.Variant{
			"Parent":      dbus.MakeVariant(parent.Path),
			"Type":        dbus.MakeVariant(child.Type()),
			"Path":        dbus.MakeVariant(child.Path()),
			"DisplayName": dbus.MakeVariant(child.DisplayName()),
		})
	}
	return children, nil
}
