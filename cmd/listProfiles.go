/*
Copyright © 2026 Erin
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/browserutils/kooky"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// listProfilesCmd represents the listProfiles command
var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List the profiles with cookies available for the given slack subdomain",

	Run: func(cmd *cobra.Command, args []string) {
		if subdomain == "" {
			cmd.Help()
			return
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Browser", "Profile"})
		stores := kooky.FindAllCookieStores(cmd.Context())
		site, _ := url.Parse(fmt.Sprintf("https://%s.slack.com", subdomain))
		for _, store := range stores {
			if len(store.Cookies(site)) > 0 {
				t.AppendRow([]interface{}{store.Browser(), store.Profile()})
			}
		}

		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(listProfilesCmd)
}
