package render

import (
	"fmt"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/hsm-gustavo/alnews/internal/model"
)

func List (items []model.Item) {
	w := tabwriter.NewWriter(color.Output, 0, 0, 3, ' ', 0)

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	for _, item := range items {
		date := FormatDate(item.PubDate)
		fmt.Fprintf(w, "%s\t%s\t%s\n", green("["+date+"]"), item.Title, cyan(item.Link))
	}

	w.Flush()
}