package notes

import (
	"context"
	"everything/models"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GetFiles(config *models.Config) (r *drive.FileList) {
	ctx := context.Background()
	srv, _ := drive.NewService(ctx, option.WithCredentialsFile(config.GoogleToken))

	query := fmt.Sprintf("'%s' in parents", config.GoogleDirId)
	r, _ = srv.Files.List().
		Q(query).
		Fields("files(id, name)").
		Do()
	return
}

func ListFiles(config *models.Config) (mr models.ModuleResponse) {
	driveFiles := GetFiles(config)
	// f.Name, f.Id
	for _, f := range driveFiles.Files {
		mr.Text += f.Name + "\n"
	}
	return
}
