package uploader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.com/dotpe/mindbenders/profiler/entity"
)

type saveToFileSystem struct{}

func (sm *saveToFileSystem) Upload(bat entity.Batch, serviceName string) error {

	fmt.Println("STORING!!")

	// Basic ISO 8601 Format in UTC as the name for the directories.
	dir := bat.End.UTC().Format("20060102T150405Z")
	dirPath := filepath.Join("profiles", dir)
	// 0755 is what mkdir does, should be reasonable for the use cases here.
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	for _, prof := range bat.Profiles {
		filePath := filepath.Join(dirPath, prof.Name)
		// 0644 is what touch does, should be reasonable for the use cases here.
		if err := ioutil.WriteFile(filePath, prof.Data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func GetFileSaver() IProfileUploader {
	return &saveToFileSystem{}
}
