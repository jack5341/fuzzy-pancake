package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
)

var username string

type Stats struct {
	Login           string      `json:"login"`
	Name            string      `json:"name"`
	Company         string      `json:"company"`
	Blog            string      `json:"blog"`
	Location        string      `json:"location"`
	Email           interface{} `json:"email"`
	Hireable        bool        `json:"hireable"`
	Bio             string      `json:"bio"`
	TwitterUsername string      `json:"twitter_username"`
	PublicRepos     int         `json:"public_repos"`
	PublicGists     int         `json:"public_gists"`
	Followers       int         `json:"followers"`
	Following       int         `json:"following"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fuzzy-pancake",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fetch(username)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "--username=<USER_NAME>: Enter username of which github account you want to get charts")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.MarkFlagRequired("username")
}

func fetch(name string) error {
	if username == "" {
		return errors.New("username is required")
	}

	resp, err := http.Get("https://api.github.com/users/" + name)
	if err != nil {
		return errors.New("error in fetching github user")
	}

	defer resp.Body.Close()

	var stats Stats
	json.NewDecoder(resp.Body).Decode(&stats)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Text Paragraph
	textBox := fmt.Sprintf(` Username: %v
	 Full Name: %v`, stats.Login, stats.Name)

	username := widgets.NewParagraph()
	username.Title = " Github User "
	username.Text = textBox
	username.SetRect(0, 0, 50, 4)
	username.BorderStyle.Fg = ui.ColorBlue

	ui.Render(username)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}

	return nil
}
