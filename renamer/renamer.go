package renamer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func RenameFiles(directory string, regex *regexp.Regexp, newNamePattern string, trimSuffix string, trimPrefix string, separator string, performRename bool, verbose bool) error {
	d, _ := filepath.Abs(directory)
	dir, err := NewDirectory(d)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if !performRename {
			fmt.Println("WARNING - Not renaming files, just printing what would be done")
		}
		files := *dir.Files()
		for _, f := range files {
			if regex.MatchString(f.currentName) {
				groups := getParams(regex, f.currentName)
				//fmt.Println(groups)
				var sep = separator
				var newName = newNamePattern
				var newNameVerbose = newNamePattern
				for key, value := range groups {
					if value == "" {
						sep = ""
					}
					value = strings.TrimPrefix(value, trimPrefix)
					value = strings.TrimSuffix(value, trimSuffix)
					newName = strings.ReplaceAll(newName, "<"+key+">", value)
					newName = strings.ReplaceAll(newName, "<#"+key+">", sep+value)
					newName = strings.ReplaceAll(newName, "<"+key+"#>", value+sep)
					newName = strings.ReplaceAll(newName, "<#"+key+"#>", sep+value+sep)
					newNameVerbose = strings.ReplaceAll(newNameVerbose, "<"+key+">", "<"+key+">"+value+"</>")
					newNameVerbose = strings.ReplaceAll(newNameVerbose, "<#"+key+">", "<"+key+">"+sep+value+"</>")
					newNameVerbose = strings.ReplaceAll(newNameVerbose, "<"+key+"#>", "<"+key+">"+value+sep+"</>")
					newNameVerbose = strings.ReplaceAll(newNameVerbose, "<#"+key+"#>", "<"+key+">"+sep+value+sep+"</>")
				}
				f.SetNewName(newName)

				fmt.Printf("\"%s\" > \"%s\"", f.FullPath(), f.NewName())
				if verbose {
					fmt.Printf(" - Verbose info: \"%s\" %s\n", newNameVerbose, groups)
				} else {
					fmt.Println()
				}

				if performRename {
					// TODO
					// Write to backup file
					err := writeToBackup(f, "backup_file.log")
					if err != nil {
						return errors.New(fmt.Sprintf("unable to rename file because: %s", err.Error()))
					}
				}
			}
		}
	}
	return nil
}

func writeToBackup(f *File, bF string) error {
	backupFile, err := os.OpenFile(bF, os.O_CREATE+os.O_APPEND+os.O_WRONLY, 0660)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create backupfile, aborting: %s\n", err.Error()))
	}
	defer func() {
		err := backupFile.Close()
		if err != nil {
			fmt.Println("Unable to close backup file")
		}
	}()

	csvFile := csv.NewWriter(backupFile)
	csvFile.Comma = ';'
	err = csvFile.Write([]string{f.OriginalName(), f.NewName()})
	if err != nil {
		return errors.New(fmt.Sprintf("unable to write to backup file, aborting: %s", err.Error()))
	}
	csvFile.Flush()

	return nil
}

func Revert(directory *Directory, backupFile string) error {
	if IsFile(backupFile) {

	} else {
		return errors.New("specified backup file does not exist")
	}
	return nil
}

func getParams(regEx *regexp.Regexp, url string) (paramsMap map[string]string) {
	match := regEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
