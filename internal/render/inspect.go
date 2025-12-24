package render

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hsm-gustavo/alnews/internal/model"
	"github.com/microcosm-cc/bluemonday"
)

func Inspect(item model.Item) {
	bold := color.New(color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	p := bluemonday.StrictPolicy()
	p.AllowElements("li", "br")

	fmt.Printf("\n%s\n", bold(magenta(item.Title)))
	fmt.Printf("%s %s\n", yellow("Date:"), FormatDate(item.PubDate))
    fmt.Printf("%s %s\n", yellow("Link:"), cyan(item.Link))

	replacer := strings.NewReplacer(
			"<li>",  "  • ", 
			"</li>", "",
			"<br>",  "\n",
    	)
	
	raw := p.Sanitize(item.Description)
	desc := strings.TrimSpace(replacer.Replace(raw))

	fmt.Println()
	fmt.Println(desc)
}