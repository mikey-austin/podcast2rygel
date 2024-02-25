package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

// Allow for dynamically loading feeds.
type feedLoader func() ([]*gofeed.Feed, error)

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

func (pd *PodcastDirectory) Parent() dbus.ObjectPath {
	// Top-level container references its own path.
	return pd.Path()
}

func (pd *PodcastDirectory) Path() dbus.ObjectPath {
	return dbus.ObjectPath(pd.path)
}

func (pd *PodcastDirectory) Type() string {
	return "container"
}

func (pd *PodcastDirectory) DisplayName() string {
	return pd.displayName
}

func (pd *PodcastDirectory) ChildCount() int {
	feeds, err := pd.feeds()
	if err != nil {
		log.Fatal(err)
	}
	return len(feeds)
}

func (pd *PodcastDirectory) ContainerCount() int {
	return pd.ChildCount()
}

func (pd *PodcastDirectory) ItemCount() int {
	// The items are contained within the episode directories
	return 0
}

func (pd *PodcastDirectory) Searchable() bool {
	return false
}

func (pd *PodcastDirectory) Register(conn *dbus.Conn) {

	// Register the org.freedesktop.DBus.Properties properties
	// for the org.gnome.UPnP.MediaObject2 interface
	prop.Export(conn, pd.Path(), GetMediaContainerProps(pd))

	// Register both org.gnome.UPnP.MediaObject2 and
	// org.gnome.UPnP.MediaContainer2 interfaces.
	err := conn.ExportMethodTable(
		GetMediaContainerMethods(pd), pd.Path(), "org.gnome.UPnP.MediaContainer2")
	if err != nil {
		log.Fatal(err)
	}
	pdLog := log.WithField("PodcastDirectory", pd)
	pdLog.Info("exported root container")

	// Then go through and register each podcast.
	numPodcasts := 0
	for i, podcast := range pd.ListPodcasts() {
		podcast.Register(conn)
		numPodcasts = i + 1
	}

	pdLog.WithField("numPodcasts", numPodcasts).Info("finished exported podcasts")
}

func (pd *PodcastDirectory) ListPodcasts() []*EpisodeDirectory {
	feeds, err := pd.feeds()
	if err != nil {
		log.Fatal(err)
	}
	podcasts := make([]*EpisodeDirectory, len(feeds))
	for i, feed := range feeds {
		podcast := NewEpisodeDirectory(i, pd, feed)
		podcasts[i] = podcast
	}
	return podcasts
}

func (pd *PodcastDirectory) ListContainers(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, *dbus.Error) {
	return pd.ListChildren(offset, max, filter)
}

func (pd *PodcastDirectory) ListItems(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, *dbus.Error) {
	// Podcast directories do not contain items, only episode directories, which are
	// themselves containers.
	return nil, nil
}

// TODO: Take the filter into account and only return requested keys in the output.
func (pd *PodcastDirectory) ListChildren(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, *dbus.Error) {
	feeds, err := pd.feeds()
	if err != nil {
		log.Fatal(err)
	}
	feeds = SliceOffsetWithMax(feeds, offset, max)
	children := make([]map[string]dbus.Variant, len(feeds))
	for i, feed := range feeds {
		parent := MediaContainer2(pd)
		child := MediaContainer2(NewEpisodeDirectory(i, pd, feed))
		children[i] = map[string]dbus.Variant{
			"Parent":      dbus.MakeVariant(parent.Path()),
			"Type":        dbus.MakeVariant(child.Type()),
			"Path":        dbus.MakeVariant(child.Path()),
			"DisplayName": dbus.MakeVariant(child.DisplayName()),
		}
	}
	return children, nil
}
