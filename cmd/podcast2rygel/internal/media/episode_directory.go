package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

// Another MediaContainer2 implementation that represents a single podcast.
type EpisodeDirectory struct {
	ParentContainer *PodcastDirectory
	Feed            *gofeed.Feed
	artCache        *PodcastImage
}

func NewEpisodeDirectory(parentContainer *PodcastDirectory, feed *gofeed.Feed) *EpisodeDirectory {
	return &EpisodeDirectory{
		ParentContainer: parentContainer,
		Feed:            feed,
		artCache:        nil,
	}
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
	conn.Export(ed, ed.Path(), "org.gnome.UPnP.MediaContainer2")
	episodeLog := log.WithField("PodcastName", ed.DisplayName())
	episodeLog.Info("exported podcast")

	// Register each episode.
	numEpisodes := 0
	for i, episode := range ed.ListEpisodes() {
		episode.Register(conn)
		numEpisodes = i + 1
	}
	episodeLog.WithField("numEpisodes", numEpisodes).Info("finished exporting episodes")
}

func (ed *EpisodeDirectory) ListEpisodes() []*Episode {
	episodes := make([]*Episode, ed.Feed.Len())
	for i, item := range ed.Feed.Items {
		episode := &Episode{EpisodeDirectory: ed, Item: item, ItemIndex: i}
		episodes[i] = episode
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
	for i, item := range items {
		parent := MediaContainer2(ed)
		child := MediaItem2(&Episode{EpisodeDirectory: ed, Item: item, ItemIndex: i})
		children[i] = map[string]dbus.Variant{
			// Media object attributes
			"Parent":      dbus.MakeVariant(parent.Path),
			"Type":        dbus.MakeVariant(child.Type()),
			"Path":        dbus.MakeVariant(child.Path()),
			"DisplayName": dbus.MakeVariant(child.DisplayName()),

			// Item attributes
			"URLs":     dbus.MakeVariant(child.Urls()),
			"MIMEType": dbus.MakeVariant(child.MimeType()),
		}
	}
	return children, nil
}

// TODO: Take the filter into account and only return requested keys in the output.
func (ed *EpisodeDirectory) ListChildren(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	items := ed.Feed.Items[offset : offset+max]
	children := make([]map[string]dbus.Variant, len(items))
	for i, item := range items {
		parent := MediaContainer2(ed)
		child := MediaObject2(&Episode{EpisodeDirectory: ed, Item: item})
		children[i] = map[string]dbus.Variant{
			"Parent":      dbus.MakeVariant(parent.Path),
			"Type":        dbus.MakeVariant(child.Type()),
			"Path":        dbus.MakeVariant(child.Path()),
			"DisplayName": dbus.MakeVariant(child.DisplayName()),
		}
	}
	return children, nil
}
