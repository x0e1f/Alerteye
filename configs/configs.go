package configs

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
	"encoding/json"
	"github.com/x0e1f/Alerteye/common"
	"io/ioutil"
	"log"
	"os"
)

// Config :: Main configuration struct
// BotToken: Telegram bot token
// ChatID: Telegram chat id (can be a group)
// CollectorTime: Time interval for the collector (Minutes)
// SendTime: Time interval for telegram sender (Minutes)
// Topics: List of topics
// Sources: LIst of sources
type Config struct {
	BotToken      string          `json:"telegram_bot_token"`
	ChatID        string          `json:"telegram_chat_id"`
	CollectorTime int             `json:"collector_time"`
	SendTime      int             `json:"send_time"`
	Topics        []common.Topic  `json:"topics"`
	Sources       []common.Source `json:"sources"`
}

// LoadConfigurations :: Load configuration file
func LoadConfigurations(confFilePath string) (Config, error) {
	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		err = initConfigurations(confFilePath)
		if err != nil {
			return Config{}, err
		}
	}

	configFile, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// initConfigurations :: Initialize configuration file with defaults
func initConfigurations(confFilePath string) error {
	log.Print("Initializing configuration file...")

	topic := common.Topic{
		Name:     "Coronavirus",
		Keywords: []string{"Coronavirus", "Covid-19"},
	}
	source := common.Source{
		Name:     "Al Jazeera",
		URL:      "https://www.aljazeera.com/xml/rss/all.xml",
		Filtered: true,
	}
	config := Config{
		BotToken:      "xxxx",
		ChatID:        "xxxx",
		CollectorTime: 60,
		SendTime:      30,
		Topics:        []common.Topic{topic},
		Sources:       []common.Source{source},
	}

	configJSON, _ := json.MarshalIndent(config, "", "\t")
	configFile, err := os.OpenFile(
		confFilePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return err
	}
	defer configFile.Close()

	configFile.Write(configJSON)

	return nil
}
