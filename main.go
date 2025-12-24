package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

/*
When getting the RSS, the body will have a format like this:

<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
	<channel>
	[...]
		<item>
			<title> something </title>
			<link>https://archlinux.org/news/something/</link>
			<description>this something description</description>
			<dc:creator xmlns:dc="http://purl.org/dc/elements/1.1/">Name of the creator</dc:creator>
			<pubDate>Sat, 20 Dec 2025 18:53:42 +0000</pubDate>
			<guid isPermaLink="false">tag:archlinux.org,2025-12-20:/news/something/</guid>
		</item>
		[...]
	</channel>
</rss>

What do we need then? We need to extract only the items, which are inside the channel
For that we use this custom type named RSS, which maps the xml into a JSON-like structure

The Item struct is used to extract only what we want from the item tag, so if you want to add the description, guid, or the creator name just get the tag name and put it in the type
*/

// xml hierarchy
type RSS struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	PubDate string `xml:"pubDate"`
	Title string `xml:"title"`
	Link string `xml:"link"`
}

func formatDate(date string) string {
	t, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		// fail-safe, just return default if there's an error
		return date
	}

	formattedDate := t.Format("Jan 02, 2006 @ 15:04")
	
	return formattedDate
}

func main() {
	resp, err := http.Get("https://archlinux.org/feeds/news")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var feed RSS

	xml.NewDecoder(resp.Body).Decode(&feed)

	w := tabwriter.NewWriter(color.Output, 0, 0, 3, ' ', 0)

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	firstThree := feed.Channel.Items[:3]

	for _, item := range firstThree {
		date := formatDate(item.PubDate)
		fmt.Fprintf(w, "%s\t%s\t%s\n", green("["+date+"]"), item.Title, cyan(item.Link))
		
	}

	w.Flush()
}