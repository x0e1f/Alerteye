package common

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

// Topic :: Topic struct
// Name: The name of the topic
// Keywords: List of keywords for topic
type Topic struct {
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
}

// Source :: Source struct
// Name: The name of the source
// URL: URL address of the source
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Alert :: Alert struct
// Date: Article publish date
// Title: Article title
// Topic: Article topic
// Source: Article source
// Description: Article description
// URL: Article url
type Alert struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Topic       Topic  `json:"topic"`
	Source      Source `json:"source"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
