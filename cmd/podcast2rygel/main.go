package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"podcast2rygel/cmd/podcast2rygel/internal/media"
)

type Config struct {
	Feeds []struct {
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

	feeds, err := fetchFeeds(config)
	if err != nil {
		log.Fatalf("Failed to fetch feeds: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %s\n", feed.Title)
		for _, item := range feed.Items {
			fmt.Printf("Item: %s - %s\n", item.Title, item.Link)
		}
	}

	container := &media.PodcastMediaContainer{ChildCount: 10}
	fmt.Printf("Num children in container: %d", container.ChildCount)
}
