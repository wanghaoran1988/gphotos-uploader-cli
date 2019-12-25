package upload

import (
	"context"
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

type EnqueuedJob struct {
	Context       context.Context
	PhotosService gPhotosService
	FileTracker   FileTracker
	Logger        log.Logger

	Path            string
	AlbumName       string
	DeleteOnSuccess bool
}

func (job *EnqueuedJob) Process() error {
	// Upload the file and add it to PhotosService.
	_, err := job.PhotosService.AddMediaItem(job.Context, job.Path, job.albumID())
	if err != nil {
		return err
	}

	// Mark the file as uploaded in the FileTracker.
	err = job.FileTracker.CacheAsAlreadyUploaded(job.Path)
	if err != nil {
		job.Logger.Warnf("Tracking file as uploaded failed: file=%s, error=%v", job.Path, err)
	}

	// If was requested, remove the file after being uploaded.
	if job.DeleteOnSuccess {
		if err := os.Remove(job.Path); err != nil {
			job.Logger.Errorf("Deletion request failed: file=%s, err=%v", job.Path, err)
		}
	}
	return nil
}

func (job *EnqueuedJob) ID() string {
	return job.Path
}

// albumID returns the album ID of the created (or existent) album in PhotosService.
func (job *EnqueuedJob) albumID() string {
	// Return if empty to avoid a PhotosService call.
	if job.AlbumName == "" {
		return ""
	}

	album, err := job.PhotosService.GetOrCreateAlbumByName(job.AlbumName)
	log.Infof("album :%s, name: %s",album,job.AlbumName)
	if err != nil {
		job.Logger.Errorf("Album creation failed: name=%s, error=%s", job.AlbumName, err)
		return ""
	}
	return album.Id
}
