package main

import (
	"io/ioutil"
	"log"
	"podcast2rygel/cmd/podcast2rygel/internal/media"

	"github.com/godbus/dbus/v5"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName string `yaml:"appName"`
	Feeds    []struct {
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
	return feeds, nil
}

func main() {
	config, err := parseConfig("feeds.yaml")
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	feeds, err := fetchFeeds(config)
	if err != nil {
		log.Fatalf("Failed to fetch feeds: %v", err)
	}

	rootDirectory := media.NewPodcastDirectory(
		func() []*gofeed.Feed { return feeds }, // TODO: invoke fetching in this lambda
		config.AppName,
		"/org/gnome/UPnP/MediaServer2/" + config.AppName)
	rootDirectory.Register(conn)

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

	// The service is now running and can be interacted with using D-Bus clients
	select {}
}
