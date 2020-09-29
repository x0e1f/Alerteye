package main

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
	"github.com/x0e1f/Alerteye/collector"
	"github.com/x0e1f/Alerteye/common"
	"github.com/x0e1f/Alerteye/configs"
	"github.com/x0e1f/Alerteye/database"
	"github.com/x0e1f/Alerteye/telegram"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dataDir := filepath.Join(rootDir, "data")
	dbPath := filepath.Join(dataDir, "alerteye.db")
	logFilePath := filepath.Join(dataDir, "alerteye.log")
	confFilePath := filepath.Join(dataDir, "config.json")

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, 0700)
	}

	logFile, err := os.OpenFile(
		logFilePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logWriter)

	common.PrintBanner()
	log.Print("Alerteye started")

	log.Print("Loading configurations")
	configs, err := configs.LoadConfigurations(confFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrations(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	go telegram.StartConsumer(dbPath, &configs)
	collector.StartCollector(dbPath, &configs)
}
