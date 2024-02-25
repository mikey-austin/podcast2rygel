package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"podcast2rygel/cmd/podcast2rygel/internal/media"

	"github.com/godbus/dbus/v5"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName string `yaml:"appName"`
	Feeds   []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"feeds"`
}

func parseConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func fetchFeeds(config *Config) ([]*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	var feeds []*gofeed.Feed
	for _, feedConfig := range config.Feeds {
		feed, err := parser.ParseURL(feedConfig.URL)
		if err != nil {
			log.Printf("Failed to fetch feed: %s", feedConfig.Name)
			continue
		}
		feeds = append(feeds, feed)
	}
	log.WithField("Feeds", len(feeds)).Info("re-loaded podcast feeds")
	return feeds, nil
}

func main() {
	configFile := pflag.String("config", "podcast2rygel.yaml", "the main config file containing the feeds")
	verbose := pflag.Bool("verbose", false, "enables more debug information")
	pflag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	config, err := parseConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	log.WithField("config", config).WithField("file", configFile).Debug("parsed configuration file")

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Load in all podcast objects into the bus
	rootDirectory := media.NewPodcastDirectory(
		func() ([]*gofeed.Feed, error) { return fetchFeeds(config) },
		config.AppName,
		"/org/gnome/UPnP/MediaServer2/"+config.AppName)
	rootDirectory.Register(conn)
	log.Debug("registered all D-Bus podcast objects")

	// Request the name on the D-Bus. The prefix needs to be like so or else rygel will not
	// find us.
	serviceName := "org.gnome.UPnP.MediaServer2." + config.AppName
	reply, err := conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatal(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatal("Name already taken")
	}
	log.Debug("opened connection to D-Bus session bus")

	log.Info("listening for D-Bus requests")

	// The service is now running and can be interacted with using D-Bus clients
	select {}
}
