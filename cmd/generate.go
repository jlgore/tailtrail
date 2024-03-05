/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new configuration file",
	Long: `The 'generate' command creates a new configuration file named '.tailtrail.yaml' in the user's home directory. 

	This configuration file is used by the tailtrail application to connect to AWS services. The file contains the AWS profile and region that the tailtrail application should use when interacting with AWS.
	
	The 'generate' command requires two flags: 'profile' and 'region'. 
	
	The 'profile' flag specifies the AWS profile that the tailtrail application should use. This should match one of the profiles defined in your AWS credentials file.
	
	The 'region' flag specifies the AWS region that the tailtrail application should use. This should be a valid AWS region, such as 'us-east-1' or 'eu-west-1'.
	
	For example, you can run the 'generate' command like this:
	
	tailtrail generate --profile my-profile --region us-east-1
	
	This will create a '.tailtrail.yaml' file in your home directory with 'my-profile' as the profile and 'us-east-1' as the region.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")

		// Get the flags
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		// Get the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Define the file path
		filePath := filepath.Join(homeDir, ".tailtrail.yaml")

		// Define the content of the file

		content := []byte(fmt.Sprintf("aws:\n    profile: %s\n    region: %s\n", profile, region))

		// Write the content to the file
		err = os.WriteFile(filePath, content, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Configuration file generated successfully.")
	},
}

func init() {
	configCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("profile", "p", "", "AWS profile to use")
	generateCmd.Flags().StringP("region", "r", "", "AWS region to use")

}
