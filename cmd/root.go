package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

type Repos struct {
	Name string `json:"name"`
	Fork bool   `json:"fork"`
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

	req, _ := http.NewRequest("GET", "https://api.github.com/users/"+name, nil)
	req.Header.Add("Authorization", "token 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")
	resp, err := http.DefaultClient.Do(req)

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
	textBox := fmt.Sprintf(`Username: %v
	Full Name: %v
	Email: %v
	Bio: %v
	Twitter: %v
	Company: %v
	Followers: %v
	Following: %v
	Created At: %v`, stats.Login, stats.Name, stats.Email, stats.Bio, string(stats.TwitterUsername), stats.Company, stats.Followers, stats.Following, stats.CreatedAt)

	user := widgets.NewParagraph()
	user.Title = " User "
	user.Text = textBox
	user.SetRect(0, 0, 40, 15)
	user.BorderStyle.Fg = ui.ColorBlue

	token := os.Getenv("GITHUB_ACCESS_TOKEN")

	if token == "" {
		return errors.New("occured an error while taking access token")
	}

	req, _ = http.NewRequest("GET", "https://api.github.com/users/"+name+"/repos", nil)
	req.Header.Add("Authorization", "token"+token)
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return errors.New("error in fetching github repositories")
	}

	var repos []Repos
	var rows []string
	json.NewDecoder(resp.Body).Decode(&repos)

	for _, repo := range repos {
		if !repo.Fork {
			rows = append(rows, repo.Name)
		}
	}

	// List
	info := widgets.NewList()
	info.Title = " Repositories "
	info.Rows = rows
	info.SetRect(40, 0, 80, 25)
	info.BorderStyle.Fg = ui.ColorBlue

	ui.Render(user, info)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}

	return nil
}
