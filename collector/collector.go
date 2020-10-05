package collector

/*
The MIT License (MIT)

Copyright (c) 2020 Davide Pataracchia

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/x0e1f/Alerteye/common"
	"github.com/x0e1f/Alerteye/configs"
	"github.com/x0e1f/Alerteye/database"
	"log"
	"regexp"
	"strings"
	"time"
)

// StartCollector :: Start Alerteye collector
func StartCollector(dbPath string, configs *configs.Config) {
	log.Print("Alerteye collector started")
	fp := gofeed.NewParser()

	for {
		sources := configs.Sources
		for i := 0; i < len(sources); i++ {
			feed, err := fp.ParseURL(sources[i].URL)
			if err != nil {
				log.Print("Collector error: ", err)
				continue
			}

			items := feed.Items
			for j := 0; j < len(items); j++ {
				alert := common.Alert{
					Date:        items[j].Published,
					Title:       items[j].Title,
					Topic:       common.Topic{},
					Source:      sources[i],
					Description: items[j].Description,
					URL:         items[j].Link,
				}

				if isBlacklisted(configs, &alert) {
					log.Print("(Blacklisted) " + alert.URL)
					continue
				}

				if !sources[i].Filtered {
					topic := common.Topic{
						Name:     "",
						Keywords: []string{},
					}
					err = newAlert(dbPath, &alert, &topic)
					if err != nil {
						log.Print(err)
					}
					continue
				}

				topic, err := checkTopic(configs, &alert)
				if err != nil {
					log.Print(err)
				}
				if topic != nil {
					err = newAlert(dbPath, &alert, topic)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}

		time.Sleep(time.Duration(configs.CollectorTime) * time.Minute)
	}
}

// newAlert :: Create a new alert
func newAlert(dbPath string, alert *common.Alert, topic *common.Topic) error {
	alertExist, _ := database.AlertExist(dbPath, alert.URL)

	if !alertExist {
		alert.Topic = *topic
		err := database.NewAlert(dbPath, alert)
		if err != nil {
			return err
		}

		if topic.Name != "" {
			log.Print("(" + topic.Name + ") " + alert.URL)
		} else {
			log.Print("(Unfiltered) " + alert.URL)
		}
	}

	return nil
}

// isBlacklisted :: Check if alert contains some blacklisted keywords
func isBlacklisted(configs *configs.Config, alert *common.Alert) bool {
	blacklist := configs.Blacklist

	for i := 0; i < len(blacklist); i++ {
		if containsKeyword(alert.Title, alert.Description, blacklist[i]) {
			return true
		}
	}

	return false
}

// checkTopic :: Check if collected alert can be associated with a topic
func checkTopic(configs *configs.Config, alert *common.Alert) (*common.Topic, error) {
	topics := configs.Topics
	suggested := make(map[*common.Topic]int)

	for i := 0; i < len(topics); i++ {
		keywords := topics[i].Keywords
		score := 0
		for j := 0; j < len(keywords); j++ {
			if containsKeyword(alert.Title, alert.Description, keywords[j]) {
				score++
			}
		}
		suggested[&topics[i]] = score
	}

	var finalTopic *common.Topic
	maxScore := 0
	for topic, score := range suggested {
		if score > maxScore {
			maxScore = score
			finalTopic = topic
		}
	}
	if maxScore <= 0 {
		return nil, nil
	}

	return finalTopic, nil
}

// containsKeyword :: Check if the article contains a keyword
func containsKeyword(title string, description string, keyword string) bool {
	title = strings.ToLower(title)
	description = strings.ToLower(description)
	keyword = strings.ToLower(keyword)

	regex := fmt.Sprintf("\\b%s\\b", keyword)
	re := regexp.MustCompile(regex)

	if len(re.FindAllString(title, -1)) > 0 {
		return true
	}

	if len(re.FindAllString(description, -1)) > 0 {
		return true
	}

	return false
}
