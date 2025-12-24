package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/microcosm-cc/bluemonday"
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
	Description string `xml:"description"`
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

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}

var CacheDir string

func listCmd(limit uint8, refetch bool, search string, openIndex int8, inspectIndex int8) error {
	switch runtime.GOOS {
	case "windows":
		return fmt.Errorf("windows is not supported")
	case "darwin":
		return fmt.Errorf("mac is not supported")
	default:
		if isWSL() {
			return fmt.Errorf("windows is not supported")
		}
	}

	// Check if home var exists
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		// Setting default folder as current
		CacheDir = "cache"
	} else {
		CacheDir = filepath.Join(home, ".cache", "alnews")
	}

	// Creating CacheDir ($HOME/.cache/alnews)
	if err := os.MkdirAll(CacheDir, 0755); err != nil {
		return err
	}

	cacheFile := filepath.Join(CacheDir, "alnews.json")

	info, err := os.Stat(cacheFile)

	useCache := refetch

	if err == nil {
		// modification time
		modTime := info.ModTime()

		// if less than 12 hours, refetch
		if time.Since(modTime) < 12 * time.Hour {
			useCache = true
		}
	}

	var feed RSS

	if useCache {
		fmt.Println("cache hit")

		data, err := os.ReadFile(cacheFile)
		if err != nil {
			return err
		}
		json.Unmarshal(data, &feed)
	} else {
		resp, err := http.Get("https://archlinux.org/feeds/news")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		xml.NewDecoder(resp.Body).Decode(&feed)

		data, _ := json.MarshalIndent(feed, "", "  ")

		os.WriteFile(cacheFile, data, 0644)
	}

	w := tabwriter.NewWriter(color.Output, 0, 0, 3, ' ', 0)

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	items := feed.Channel.Items

	if search != "" {
		searchLower := strings.ToLower(search)

		filtered := make([]Item, 0)
		for _, item := range items {
			title := strings.ToLower(item.Title)

			// if its a perfect match, just add it and go to next iteration
			if strings.Contains(title, searchLower) {
            	filtered = append(filtered, item)
            	continue
        	}

			words := strings.Fields(title)
			best := 999

			for _, w := range words {
				// for words in title, calculate levenshtein distance and get the best number of operations
				levenshteinDistance := gstr.Levenshtein(w, searchLower, 1, 1, 1)
				if levenshteinDistance < best {
					best = levenshteinDistance
				}
			}

			// if the number of operations is less than or equals to 2, add it
            if best <= 2 {
                filtered = append(filtered, item)
            }
        }

		items = filtered
	}

	if len(items) > int(limit) {
		items = items[:limit]
	}

	for _, item := range items {
		date := formatDate(item.PubDate)
		fmt.Fprintf(w, "%s\t%s\t%s\n", green("["+date+"]"), item.Title, cyan(item.Link))
	}

	if openIndex >= 0 {
		exec.Command("xdg-open", items[openIndex].Link).Start()
	}

	if inspectIndex >= 0 {
		p := bluemonday.StrictPolicy()
		inspectItem := items[inspectIndex]

		fmt.Printf("\n%s\n", bold(magenta(inspectItem.Title)))

		fmt.Printf("%s %s\n", yellow("Date:"), formatDate(inspectItem.PubDate))
    	fmt.Printf("%s %s\n", yellow("Link:"), cyan(inspectItem.Link))

		p.AllowElements("li", "br")

		replacer := strings.NewReplacer(
			"<li>",  "  • ", 
			"</li>", "",
			"<br>",  "\n",
    	)
		rawDescription := p.Sanitize(inspectItem.Description)
		formattedDescription := strings.TrimSpace(replacer.Replace(rawDescription))

		fmt.Println()

		fmt.Printf("%s\n", formattedDescription)
	}

	w.Flush()
	return nil
}