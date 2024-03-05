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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
