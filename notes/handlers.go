package notes

import (
	"bytes"
	"context"
	"everything/common"
	c "everything/config"
	"everything/models"
	n "everything/models/notes"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func InitService() (srv *drive.Service, query string, err error) {
	// Function provides srv on demand.
	// It also provides query string, but only one other function really needs it
	// But the query is still generated here so we don't call LoadConfig outside of this function
	config, _ := c.LoadConfig()
	ctx := context.Background()
	srv, err = drive.NewService(ctx, option.WithCredentialsFile(config.GoogleToken))
	if err != nil {
		return srv, query, err
	}
	query = fmt.Sprintf("'%s' in parents", config.GoogleDirId)
	return
}

func GetFiles() (r *drive.FileList, err error) {
	// Get all files from specific GoogleDrive dir
	srv, query, err := InitService()
	if err != nil {
		return r, err
	}
	r, _ = srv.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()
	return
}

func DownloadFile(fileID string) (string, error) {
	// Download file contents and convert them to string
	var buf bytes.Buffer
	var err error
	srv, _, err := InitService()
	if err != nil {
		return "", err
	}
	resp, _ := srv.Files.Get(fileID).Download()
	defer resp.Body.Close()
	io.Copy(&buf, resp.Body)
	return buf.String(), err
}

func UpdateDriveFile(fileId string, newContent string) (err error) {
	// Update the file with user input
	// Input is converted to io.ReadCloser
	srv, _, err := InitService()
	if err != nil {
		return err
	}
	reader := io.NopCloser(strings.NewReader(newContent))
	_, err = srv.Files.Update(fileId, nil).Media(reader).Do()
	return err
}

func ParseFiles(files []*drive.File) (parsedFiles []string) {
	// Get File Names and sort them alphabetically
	for _, file := range files {
		parsedFiles = append(parsedFiles, file.Name)
	}
	sort.Strings(parsedFiles)
	return
}

func getFileID(files []*drive.File, fileNum int) (fileID string) {
	// Find File ID by using a file number provided by the user
	// For human readability file numbers start with 1, so we have to substract it
	sortedFiles := ParseFiles(files)
	chosenFile := sortedFiles[fileNum-1]
	for _, file := range files {
		if file.Name == chosenFile {
			fileID = file.Id
			break
		}
	}
	return
}

func FileSelectionGet(noteCache *[]n.FileSelector, userID int64) string {
	// Get file number from cache
	for _, v := range *noteCache {
		if v.UserID == userID {
			return v.FileID
		}
	}
	return ""
}

func FailedAPICall(err error) (mr models.ModuleResponse) {
	// End user chat because API errors could not be caused by user input
	log.Println(err)
	mr.Text = "Failed to make Drive API call, try again"
	mr.Error = true
	mr.EndChat = true
	mr.Keyboard = common.CompileDefaultKeyboard()
	return
}

func DeleteSelectorCache(fs *[]n.FileSelector, userID int64) {
	for index, elem := range *fs {
		if elem.UserID == userID {
			*fs = append((*fs)[:index], (*fs)[index+1:]...)
		}
	}
}
