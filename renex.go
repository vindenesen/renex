package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"renex/renamer"
)

const (
	ModeRename int = iota
	ModeRevert
)

var (
	runMode   int
	directory string

	restoreFile        string
	regexPattern       string
	regexPatternObject *regexp.Regexp

	backupFile string
	separator  string
	trimPrefix string
	trimSuffix string

	newNamePattern string
	performRename  bool

	verbose bool
)

func init() {
	renameMode := flag.FlagSet{}
	renameMode.StringVar(&directory, "directory", "", "[REQUIRED] The `directory` which contains files to rename")
	renameMode.StringVar(&regexPattern, "regex", "", "[REQUIRED] The `regex pattern` to use")
	renameMode.StringVar(&newNamePattern, "new-name", "", "[REQUIRED] The `pattern` to use for new name")
	renameMode.StringVar(&separator, "separator", "", "The `separator` used between group matches in new-name. Added by using # inside new-name tag either at beginning or end")
	renameMode.StringVar(&trimPrefix, "trim-prefix", "", "Trim `prefix` in group match")
	renameMode.StringVar(&trimSuffix, "trim-suffix", "", "Trim `suffix` in group match")
	renameMode.StringVar(&backupFile, "backup-file", "", "Backup `file` for reverting rename operation")
	renameMode.BoolVar(&performRename, "rename", false, "To actually rename files instead of just outputting the result, set this flag. Requires backup-file option")
	renameMode.BoolVar(&verbose, "verbose", false, "Verbose output for debugging")

	/*renameMode.Usage = func() {
		fmt.Fprintf(os.Stderr, "Not helpful at all\n")
		renameMode.PrintDefaults()
	}*/

	revertMode := flag.FlagSet{}
	revertMode.StringVar(&directory, "directory", "", "The `directory` which contains files to rename")
	revertMode.StringVar(&restoreFile, "restore", "", "The `file` containing the files to revert")
	revertMode.BoolVar(&performRename, "rename", false, "To actually revert file names instead of just outputting the result, set this flag")

	if len(os.Args) < 2 {
		printHelpUsage()
		os.Exit(1)
	}

	if os.Args[1] == "rename" {
		runMode = ModeRename
		err := renameMode.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if directory == "" {
			fmt.Println("Missing directory")
			printHelpUsage()
			os.Exit(1)
		}

		if backupFile == "" && performRename {
			fmt.Println("Missing option -backup-file")
			printHelpUsage()
			os.Exit(1)
		}

		if renamer.IsFile(backupFile) {
			fmt.Println("Specified backup file already exists, please specify a non-existing file")
			os.Exit(1)
		}

		if regexPattern == "" {
			fmt.Println("Missing regex pattern")
			printHelpUsage()
			os.Exit(1)
		}

		if newNamePattern == "" {
			fmt.Println("Missing new name pattern")
			printHelpUsage()
			os.Exit(1)
		}

		regexPatternObject = regexp.MustCompile(regexPattern)
	} else if os.Args[1] == "revert" {
		runMode = ModeRevert
		err := revertMode.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if directory == "" {
			fmt.Println("Missing required directory")
			printHelpUsage()
			os.Exit(1)
		}

		if restoreFile == "" {
			fmt.Println("Missing required restore file")
			printHelpUsage()
			os.Exit(1)
		}
	} else {
		printHelpUsage()
		os.Exit(1)
	}
}

func printHelpUsage() {
	fmt.Printf("Usage: \n" +
		"  renex help - displays usage information\n" +
		"\n" +
		"  renex rename <options> - rename files matching regex\n" +
		"    Options are: \n" +
		"    -directory directory - [REQUIRED]\n" +
		"    -regex regex - [REQUIRED] - Regular expression syntax, see https://github.com/google/re2/wiki/Syntax\n" +
		"    -new-name pattern - [REQUIRED] - Use the name of the group inside <> brackets, eg: <episode>\n" +
		"    -verbose - More output for debugging\n" +
		"    -rename - Actually rename files\n" +
		"    -separator separator - Separator used between group matches in new-name. Added by using # inside new-name tag either at beginning or end\n" +
		"    -trim-suffix suffix - trim suffix from group match\n" +
		"    -trim-prefix prefix - trim prefix from group match\n" +
		"    -backup-file file - [REQUIRED] - backup file to use for reverting rename operation\n" +
		"\n" +
		"  renex revert <options>\n" +
		"  Options are: \n" +
		"    -directory directory - [REQUIRED]\n" +
		"    -restore file - [REQUIRED]\n" +
		"    -rename - To actually revert file names, add this" +
		"\n" +
		"Examples: \n" +
		"" +
		"\n")
}

func main() {
	if runMode == ModeRename {
		err := renamer.RenameFiles(directory, regexPatternObject, newNamePattern, trimSuffix, trimPrefix, separator, performRename, verbose, backupFile)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if runMode == ModeRevert {
		err := renamer.Revert(directory, restoreFile, performRename)
		if err != nil {
			fmt.Printf("Error while reverting file names: %s", err.Error())
		}
	}
}
