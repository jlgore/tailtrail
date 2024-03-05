package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var username string
var outputFormat string

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
			fmt.Println("Error", err)
			return
		}

		switch outputFormat {
		case "table":
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "Event ID\tEvent Name\tUsername\tEvent Time")

			for _, event := range result.Events {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", *event.EventId, *event.EventName, *event.Username, event.EventTime.Format(time.RFC3339))
			}

			w.Flush()
		case "json":
			b, err := json.Marshal(result.Events)
			if err != nil {
				fmt.Println("Error", err)
				return
			}

			fmt.Println(string(b))
		default:
			fmt.Println("Invalid output format:", outputFormat)

		}
	},
}

func init() {
	usertailCmd.Flags().StringVarP(&username, "user", "u", "", "The user who's activity you want to tail")
	usertailCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table|json)")
}
