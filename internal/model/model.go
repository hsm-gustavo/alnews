package model

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
	Description string `xml:"description"`
}