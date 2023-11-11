package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func printFileStructure(directoryPath string, level int, deep bool) error {
	// Print the current directory
	fmt.Println(strings.Repeat("-", level) + filepath.Base(directoryPath))

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path != directoryPath {
			// Calculate the nesting level based on the number of separators in the path
			nestingLevel := strings.Count(strings.TrimPrefix(path, directoryPath), string(filepath.Separator))

			if !deep && nestingLevel > 1 {
				return filepath.SkipDir
			}

			// Print the file/folder name with appropriate indentation
			fmt.Println(strings.Repeat("-", level+nestingLevel+1) + info.Name())
		}

		return nil
	})

	return err
}

func runCommand(cmd *cobra.Command, args []string) {
	var directoryPath string

	if len(args) < 1 {
		// If no path is provided, use the current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
		directoryPath = currentDir
	} else {
		directoryPath = args[0]
	}

	deep, _ := cmd.Flags().GetBool("deep")

	err := printFileStructure(directoryPath, 0, deep)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mygoapp",
		Short: "List files and folders in a directory",
	}

	// Add a flag for deep listing
	rootCmd.Flags().BoolP("deep", "d", false, "List files and folders recursively")

	rootCmd.Run = runCommand

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
