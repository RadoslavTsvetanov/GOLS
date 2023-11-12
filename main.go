package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Config struct {
	Icons  map[string]string `json:"icons"`
	Colors map[string]string `json:"colors"`
}

var globalConfig Config

func loadConfigFile(filename string) (Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

type File struct {
	Name         string
	Size         int64
	IsDir        bool
	Modification time.Time
}

type Files []File

func (f Files) Len() int           { return len(f) }
func (f Files) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f Files) Less(i, j int) bool { return f[i].Name < f[j].Name }

func printFileStructure(directoryPath string, level int, deep bool) (Files, error) {
	entries, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	files := make(Files, 0, len(entries))
	for _, entry := range entries {
		file := File{
			Name:         entry.Name(),
			Size:         entry.Size(),
			IsDir:        entry.IsDir(),
			Modification: entry.ModTime(),
		}
		files = append(files, file)
	}
	sort.Sort(files)

	fmt.Println(strings.Repeat("-", level) + getColoredIcon("folder") + filepath.Base(directoryPath))

	for _, file := range files {
		nestingLevel := strings.Count(strings.TrimPrefix(filepath.Dir(file.Name), directoryPath), string(filepath.Separator))

		if !deep && nestingLevel > 0 {
			continue
		}

		icon := getColoredIcon(filepath.Ext(file.Name))
		fmt.Printf("%s %s %s\n", strings.Repeat("-", level+nestingLevel+1), icon, file.Name)
	}

	return files, nil
}

func getColoredIcon(extension string) string {
	icon, exists := globalConfig.Icons[extension]
	colorCode, colorExists := globalConfig.Colors[extension]

	if !exists {
		icon = ""
	}
	if !colorExists {
		colorCode = "reset"
	}

	var colorAttribute color.Attribute
	switch colorCode {
	case "fgBlack":
		colorAttribute = color.FgBlack
	case "fgRed":
		colorAttribute = color.FgRed
	case "fgGreen":
		colorAttribute = color.FgGreen
	case "fgYellow":
		colorAttribute = color.FgYellow
	case "fgBlue":
		colorAttribute = color.FgBlue
	case "fgMagenta":
		colorAttribute = color.FgMagenta
	case "fgCyan":
		colorAttribute = color.FgCyan
	case "fgWhite":
		colorAttribute = color.FgWhite
	case "fgHiBlack":
		colorAttribute = color.FgHiBlack
	case "fgHiRed":
		colorAttribute = color.FgHiRed
	case "fgHiGreen":
		colorAttribute = color.FgHiGreen
	case "fgHiYellow":
		colorAttribute = color.FgHiYellow
	case "fgHiBlue":
		colorAttribute = color.FgHiBlue
	case "fgHiMagenta":
		colorAttribute = color.FgHiMagenta
	case "fgHiCyan":
		colorAttribute = color.FgHiCyan
	case "fgHiWhite":
		colorAttribute = color.FgHiWhite
	default:
		colorAttribute = color.Reset
	}

	coloredIcon := color.New(colorAttribute).Sprint(icon)

	return coloredIcon
}

func runCommand(cmd *cobra.Command, args []string) {
	var directoryPath string

	if len(args) < 1 {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting the current directory:", err)
			os.Exit(1)
		}
		directoryPath = currentDir
	} else {
		directoryPath = args[0]
	}

	configFile := "config.json"
	config, err := loadConfigFile(configFile)
	if err != nil {
		fmt.Println("Error loading the config file:", err)
		os.Exit(1)
	}
	globalConfig = config

	printConfig()

	deep, _ := cmd.Flags().GetBool("deep")

	files, err := printFileStructure(directoryPath, 0, deep)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	printSummary(directoryPath, files)
}

func printConfig() {
	configData, err := json.MarshalIndent(globalConfig, "", "  ")
	if err != nil {
		fmt.Println("Error formatting the config:", err)
		return
	}

	fmt.Println("Loaded Config:")
	fmt.Println(string(configData))
}

func printSummary(directoryPath string, files Files) {
	var totalSize int64
	var fileCount, folderCount int

	for _, file := range files {
		if file.IsDir {
			folderCount++
		} else {
			fileCount++
			totalSize += file.Size
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total Files: %d\n", fileCount)
	fmt.Printf("Total Folders: %d\n", folderCount)
	fmt.Printf("Total Size: %s\n", humanReadableSize(totalSize))
}

func humanReadableSize(size int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
		PB
		EB
	)

	switch {
	case size >= EB:
		return fmt.Sprintf("%.2f EB", float64(size)/float64(EB))
	case size >= PB:
		return fmt.Sprintf("%.2f PB", float64(size)/float64(PB))
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mygoapp",
		Short: "List files and folders in a directory",
	}

	rootCmd.Flags().BoolP("deep", "d", false, "List files and folders recursively")

	rootCmd.Run = runCommand

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
