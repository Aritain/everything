package notes

import (
	"everything/common"
	"everything/models"
	"fmt"
	"strconv"
	"strings"

	n "everything/models/notes"
)

func ListFiles() (mr models.ModuleResponse) {
	driveFiles, err := GetFiles()
	if err != nil {
		mr = FailedAPICall(err)
		return
	}
	fileNames := ParseFiles(driveFiles.Files)
	mr.Text = "Choose file:\n"
	for index, f := range fileNames {
		filename := strings.Split(f, ".")
		mr.Text += fmt.Sprintf("(%v) %s\n", (index + 1), filename[0])
	}
	mr.Keyboard = common.CompileCancelKeyboard()
	return
}

func SelectFile(text string, userID int64, fs *[]n.FileSelector) (mr models.ModuleResponse) {
	mr.Keyboard = common.CompileCancelKeyboard()
	driveFiles, err := GetFiles()
	if err != nil {
		mr = FailedAPICall(err)
		return
	}
	fileNum, err := strconv.Atoi(text)
	if (err != nil) || (fileNum <= 0) || (fileNum > len(driveFiles.Files)) {
		mr.Text = "Bad value, try again"
		mr.Error = true
		return mr
	}
	fileID := getFileID(driveFiles.Files, fileNum)
	*fs = append(*fs, n.FileSelector{UserID: userID, FileID: fileID})
	mr.Text = "Provide note"
	return
}

func UpdateFile(text string, userID int64, fs *[]n.FileSelector) (mr models.ModuleResponse) {
	var err error
	fileID := FileSelectionGet(fs, userID)
	// Whatever outcome is, we delete Selector cache since the chat will be ended either way
	DeleteSelectorCache(fs, userID)
	fileContent, err := DownloadFile(fileID)
	if err != nil {
		mr = FailedAPICall(err)
		return
	}
	fileContent += "\n" + text
	err = UpdateDriveFile(fileID, fileContent)
	if err != nil {
		mr = FailedAPICall(err)
		return
	}
	mr.Text = "Done"
	mr.EndChat = true
	mr.Keyboard = common.CompileDefaultKeyboard()
	return
}
