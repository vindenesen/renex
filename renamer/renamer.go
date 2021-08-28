package renamer

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func RenameFiles(directory string, regex *regexp.Regexp, newNamePattern string, trimSuffix string, trimPrefix string, separator string, performRename bool, verbose bool, backupFile string) error {
	d, _ := filepath.Abs(directory)
	dir, err := NewDirectory(d)
	if err != nil {
		return err
	} else {
		var bF *BackupFile
		if performRename {
			bF, err = NewBackupFile(backupFile)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("WARNING - Not renaming files, just printing what would be done. Add -rename option to actually make changes")
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
				err := f.SetNewName(newName)
				if err != nil {
					fmt.Printf("* Rename \"%s\" > %s: %s\n", f.CurrentName(), newName, err.Error())
					continue
				}

				fmt.Printf("* Rename \"%s\" > \"%s\"", f.CurrentName(), f.NewName())
				if verbose {
					fmt.Printf(" - Verbose info: \"%s\" %s\n", newNameVerbose, groups)
				} else {
					fmt.Println()
				}

				if performRename {
					// Write to backup file
					err = bF.WriteToBackup(f)
					if err != nil {
						return errors.New(fmt.Sprintf("unable to rename file because: %s", err.Error()))
					}
					err = f.Rename()
					if err != nil {
						fmt.Printf("  Unable to rename file because: %s\n", err.Error())
					}
				}
			}
		}
	}
	return nil
}

func Revert(directory string, backupFile string, performRename bool) error {
	d, _ := filepath.Abs(directory)
	dir, err := NewDirectory(d)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		bF, err := NewExistingBackupFile(backupFile)
		if err != nil {
			return err
		}
		if IsFile(backupFile) {
			if !performRename {
				fmt.Println("WARNING - Not renaming files, just printing what would be done. Add -rename option to actually make changes")
			}

			files, err := bF.ReadFromBackup()
			if err != nil {
				return err
			}

			for _, f := range files {
				fmt.Printf("* Reverting file name \"%s\" to \"%s\"\n", f[1], f[0])
				dirFile := dir.GetFile(f[1])
				if dirFile != nil {
					err := dirFile.SetNewName(f[0])
					if err != nil {
						fmt.Printf("  WARNING - New file name is identical to old, not renaming file")
						continue
					}
					if performRename {
						err := dirFile.Rename()
						if err != nil {
							fmt.Printf("  WARNING - Unable to revert file %s to original name %s: %s\n", f[1], f[0], err.Error())
						}
					}
				} else {
					fmt.Printf("  WARNING - Unable to revert file %s to original name %s, couldn't locate the previously renamed file in the directory\n", f[1], f[0])
				}
			}
		} else {
			return errors.New("specified backup file does not exist")
		}
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
