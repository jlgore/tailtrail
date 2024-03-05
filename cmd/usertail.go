package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var username string
var outputFormat string

type errMsg error

func fetchEvents(username string) ([]*cloudtrail.Event, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(viper.GetString("aws.region")),
		},
		Profile:           viper.GetString("aws.profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudtrail.New(sess)

	input := &cloudtrail.LookupEventsInput{
		LookupAttributes: []*cloudtrail.LookupAttribute{{
			AttributeKey:   aws.String("Username"),
			AttributeValue: aws.String(username),
		}},
	}

	result, err := svc.LookupEvents(input)
	if err != nil {
		return nil, err
	}

	return result.Events, nil
}

// start the bubbletea shit

type model struct {
	spinner        spinner.Model
	table          table.Model
	quitting       bool
	err            error
	events         []*cloudtrail.Event
	fetchingEvents bool
}

func convertEventsToRows(events []*cloudtrail.Event) []table.Row {
	rows := make([]table.Row, len(events))
	for i, event := range events {
		rows[i] = table.Row{
			*event.EventId,
			*event.EventName,
			*event.Username,
			event.EventTime.Format(time.RFC3339),
		}
	}
	return rows
}

func initialModel(events []*cloudtrail.Event) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Event ID", Width: 20},
			{Title: "Event Name", Width: 20},
			{Title: "Username", Width: 20},
			{Title: "Event Time", Width: 20},
		}),
		table.WithRows(convertEventsToRows(events)),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return model{spinner: s, events: events, table: t}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "down":
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case time.Time:
		// After the delay, load the events data
		if !m.fetchingEvents {
			m.fetchingEvents = true
			go func() {
				events, err := fetchEvents(username) // replace with your function to load events
				if err != nil {
					m.err = err
				} else {
					m.events = events
				}
				m.fetchingEvents = false
			}()
		}
		return m, nil

	default:
		var cmd1, cmd2 tea.Cmd
		m.spinner, cmd1 = m.spinner.Update(msg)
		m.table, cmd2 = m.table.Update(msg)
		cmd = tea.Batch(cmd1, cmd2) // batch the commands together
	}

	// If m.events is nil, start a timer to add a delay
	// if m.events == nil {
	// return m, tea.After(time.Second)
	// }

	return m, cmd

}

func (m model) View() string {
	if m.err != nil {
		return "\033[H\033[2J" + m.err.Error()
	}

	if m.quitting {
		return "\033[H\033[2JQuitting...\n"
	}

	var b strings.Builder
	b.WriteString("\033[H\033[2J\n\n")

	// Only render the spinner and loading message if the events have not loaded
	if m.events == nil {
		b.WriteString(m.spinner.View())
		b.WriteString(" Loading events...press q to quit\n\n")
	}

	b.WriteString(m.table.View())

	return b.String()
}

// usertailCmd represents the usertail command
var usertailCmd = &cobra.Command{
	Use:   "usertail",
	Short: "Tail a users actions in AWS CloudTrail",
	Long: `usertail is a command-line tool that allows you to tail a user's actions in AWS CloudTrail.

	You can specify the user whose actions you want to tail with the --username (-u) argument. For example:
	
	./tailtrail usertail -u jg@gore.cc
	
	By default, usertail outputs the actions in a table format. You can change the output format to JSON with the --output (-o) argument. For example:
	
	./tailtrail usertail -u jg@gore.cc -o json
	
	usertail uses the AWS profile and region specified in your .tailtrail.yaml configuration file. You can override these with the AWS_PROFILE and AWS_REGION environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {

		events, err := fetchEvents(username)
		if err != nil {
			fmt.Println("Error", err)
			return
		}

		switch outputFormat {
		case "table":
			maxID, maxName, maxUser, maxTime := 0, 0, 0, 0
			for _, event := range events {
				idLen := len(*event.EventId)
				nameLen := len(*event.EventName)
				userLen := len(*event.Username)
				timeLen := len(event.EventTime.Format(time.RFC3339))

				if idLen > maxID {
					maxID = idLen
				}
				if nameLen > maxName {
					maxName = nameLen
				}
				if userLen > maxUser {
					maxUser = userLen
				}
				if timeLen > maxTime {
					maxTime = timeLen
				}
			}

			// Create the format string for the borders and the rows
			borderFormat := "+%s+%s+%s+%s+\n"
			rowFormat := "| %s\t| %s\t| %s\t| %s\t|\n"

			// Create the borders
			idBorder := strings.Repeat("-", maxID+2) // +2 for the spaces around the value
			nameBorder := strings.Repeat("-", maxName+2)
			userBorder := strings.Repeat("-", maxUser+2)
			timeBorder := strings.Repeat("-", maxTime+2)

			// Print the table
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, borderFormat, idBorder, nameBorder, userBorder, timeBorder)
			fmt.Fprintln(w, "| Event ID\t| Event Name\t| Username\t| Event Time\t|")
			fmt.Fprintf(w, borderFormat, idBorder, nameBorder, userBorder, timeBorder)

			for _, event := range events {
				fmt.Fprintf(w, rowFormat, *event.EventId, *event.EventName, *event.Username, event.EventTime.Format(time.RFC3339))
				fmt.Fprintf(w, borderFormat, idBorder, nameBorder, userBorder, timeBorder)
			}

			w.Flush()

		case "json":
			b, err := json.Marshal(events)
			if err != nil {
				fmt.Println("Error", err)
				return
			}

			fmt.Println(string(b))
		case "shell":
			// Create a new Bubbletea program
			p := tea.NewProgram(initialModel(events))
			// Start the Bubbletea event loop
			if _, err := p.Run(); err != nil {
				fmt.Println("Error", err)
				return
			}
		default:
			fmt.Println("Invalid output format:", outputFormat)

		}
	},
}

func init() {
	usertailCmd.Flags().StringVarP(&username, "user", "u", "", "The user who's activity you want to tail")
	usertailCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table|json|shell)")
}
