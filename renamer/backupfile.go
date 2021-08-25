package renamer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type BackupFile struct {
	backupFile string
}

func NewBackupFile(filename string) (*BackupFile, error) {
	if IsFile(filename) {
		return nil, errors.New("backup file already exists")
	}
	return &BackupFile{backupFile: filename}, nil
}

func NewExistingBackupFile(filename string) (*BackupFile, error) {
	if !IsFile(filename) {
		return nil, errors.New(fmt.Sprintf("backup file %s doesn't exist", filename))
	}
	return &BackupFile{backupFile: filename}, nil
}

func (bF *BackupFile) WriteToBackup(f *File) error {
	backupFile, err := os.OpenFile(bF.backupFile, os.O_CREATE+os.O_APPEND+os.O_WRONLY, 0660)
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

func (bF *BackupFile) ReadFromBackup() ([][]string, error) {
	file, err := os.OpenFile(bF.backupFile, os.O_RDONLY, 660)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to open backup file: %s", err.Error()))
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("Unable to close backup file: " + err.Error())
		}
	}()

	csvFile := csv.NewReader(file)
	csvFile.Comma = ';'
	files, err := csvFile.ReadAll()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read from backup file: %s", err.Error()))
	}

	return files, nil
}
