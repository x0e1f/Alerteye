package database

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
	"github.com/jinzhu/gorm"
	"github.com/x0e1f/Alerteye/internal/common"

	// The following is required for GORM sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Alert :: Alert database model
type Alert struct {
	gorm.Model
	Date   string
	Topic  string
	Title  string
	Source string
	URL    string `gorm:"unique;not null"`
	Sent   int
}

// ConnectDatabase :: Connect to sqlite database
func ConnectDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrations :: Database migrations
func Migrations(dbPath string) error {
	db, err := ConnectDatabase(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&Alert{})

	return nil
}

//NewAlert :: Create a new alert
func NewAlert(dbPath string, alert *common.Alert) error {
	db, err := ConnectDatabase(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	entry := Alert{
		Date:   alert.Date,
		Topic:  alert.Topic.Name,
		Title:  alert.Title,
		Source: alert.Source.Name,
		URL:    alert.URL,
		Sent:   -1,
	}
	db.Create(&entry)

	return nil
}

// AlertSent :: Mark an alert as sent
func AlertSent(dbPath string, url string) error {
	db, err := ConnectDatabase(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	alert := Alert{}
	db.Where(Alert{
		URL: url,
	}).First(&alert)

	alert.Sent = 1
	db.Save(&alert)

	return nil
}

// AlertExist :: Check if an alert already exist
func AlertExist(dbPath string, url string) (bool, error) {
	db, err := ConnectDatabase(dbPath)
	if err != nil {
		return true, err
	}
	defer db.Close()

	alert := Alert{}
	db.Where(Alert{
		URL: url,
	}).First(&alert)
	if alert.URL == "" {
		return false, nil
	}

	return true, nil
}

// AlertToSend :: Retrieve the next alert that need to be sent
func AlertToSend(dbPath string) (*Alert, error) {
	db, err := ConnectDatabase(dbPath)
	if err != nil {
		return &Alert{}, err
	}
	defer db.Close()

	alert := Alert{}
	db.Where(Alert{
		Sent: -1,
	}).First(&alert)

	return &alert, nil
}
