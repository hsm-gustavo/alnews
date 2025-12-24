package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Limit uint8
 	Refetch bool
 	Search string
	OpenIndex int8
	InspectIndex int8
)


var rootCmd = &cobra.Command{
	Use: "alnews",
	Short: "A simple script in Go to fetch and display the latest news from the Arch Linux website.",
	Long: `A simple script written in Go to fetch and display the latest news from the Arch Linux website, with colorful output, date and link for the post.`,
	Run: func(cmd *cobra.Command, args []string) {
		listCmd(Limit, Refetch, Search, OpenIndex, InspectIndex)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// add defaults for flags
	rootCmd.Flags().Uint8VarP(&Limit,"limit", "l", 3, "number of news to show")
	rootCmd.Flags().BoolVarP(&Refetch, "refresh", "r", false, "if you want to refresh data or use cache")
	rootCmd.Flags().StringVarP(&Search, "search", "s", "", "uses fuzzy search to filter news titles for specific keywords (e.g., 'nvidia')")
	rootCmd.Flags().Int8VarP(&OpenIndex, "open", "o", -1, "opens the link of a specified news (index, starting from 0) in your default browser")
	rootCmd.Flags().Int8VarP(&InspectIndex, "inspect", "i", -1, "a longer print of the news (index, starting from 0)")
}