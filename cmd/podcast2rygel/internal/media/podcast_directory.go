package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/mmcdole/gofeed"
)

// Allow for dynamically loading feeds.
type feedLoader func() []*gofeed.Feed

// An implementation of the MediaContainer2 interface that contains
// multiple EpisodeDirectory instances according to the configured podcasts.
type PodcastDirectory struct {
	feeds       feedLoader
	displayName string
	path        string
}

func NewPodcastDirectory(feeds feedLoader, displayName string, path string) *PodcastDirectory {
	return &PodcastDirectory{feeds: feeds, displayName: displayName, path: path}
}

func (pmc *PodcastDirectory) Parent() dbus.ObjectPath {
	// Top-level container references its own path.
	return pmc.Path()
}

func (pmc *PodcastDirectory) Path() dbus.ObjectPath {
	return dbus.ObjectPath(pmc.path)
}

func (pmc *PodcastDirectory) Type() string {
	return "container"
}

func (pmc *PodcastDirectory) DisplayName() string {
	return pmc.displayName
}

func (pmc *PodcastDirectory) ChildCount() int {
	return len(pmc.feeds())
}

func (pmc *PodcastDirectory) ContainerCount() int {
	return pmc.ChildCount()
}

func (pmc *PodcastDirectory) ItemCount() int {
	// The items are contained within the episode directories
	return 0
}

func (pmc *PodcastDirectory) Searchable() bool {
	return false
}

func (pmc *PodcastDirectory) Register(conn *dbus.Conn) {
	// Register both org.gnome.UPnP.MediaObject2 and
        // org.gnome.UPnP.MediaContainer2 interfaces.

	// Then go through and register each podcast.
	for _, podcast := range pmc.ListPodcasts() {
		podcast.Register(conn)
	}
}

func (pmc *PodcastDirectory) ListPodcasts() []*EpisodeDirectory {
	feeds := pmc.feeds()
	podcasts := make([]*EpisodeDirectory, len(feeds))
	for _, feed := range feeds {
		podcast := &EpisodeDirectory{ParentContainer: pmc, Feed: feed}
		podcasts = append(podcasts, podcast)
	}
	return podcasts
}

func (pmc *PodcastDirectory) ListContainers(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	return pmc.ListChildren(offset, max, filter)
}

func (pmc *PodcastDirectory) ListItems(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	// Podcast directories do not contain items, only episode directories, which are
	// themselves containers.
	return nil, nil
}

// TODO: Take the filter into account and only return requested keys in the output.
func (pmc *PodcastDirectory) ListChildren(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, error) {
	feeds := pmc.feeds()[offset : offset+max]
	children := make([]map[string]dbus.Variant, len(feeds))
	for _, feed := range feeds {
		parent := MediaContainer2(pmc)
		child := MediaContainer2(&EpisodeDirectory{ParentContainer: pmc, Feed: feed})
		children = append(children, map[string]dbus.Variant{
			"Parent":      dbus.MakeVariant(parent.Path),
			"Type":        dbus.MakeVariant(child.Type()),
			"ItemCount":   dbus.MakeVariant(child.ItemCount()),
			"DisplayName": dbus.MakeVariant(child.DisplayName()),
		})
	}
	return children, nil
}
