package admin

import (
	"mime/multipart"
	"strings"
	"path/filepath"
	"fmt"
	"archive/zip"

	"github.com/qor/media_library"
)

// archiveHandler default image handler
type archiveHandler struct{}

func (archiveHandler) CouldHandle(media media_library.MediaLibrary) bool {
	ext := strings.ToLower(filepath.Ext(media.URL()))
	return ext == ".zip"
}

func (archiveHandler) Handle(media media_library.MediaLibrary, file multipart.File, option *media_library.Option) (err error) {
	if err = media.Store(media.URL(), option, file); err != nil {
		return err
	}

	size, err := file.Seek(0, 2)
	if err != nil {
		return err
	}
	
	r, err := zip.NewReader(file, size)
	if err != nil {
		return err
	}

	page := 0
	var pages []string
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		ext := filepath.Ext(f.Name)
		if ext != ".jpg" &&
		  	ext != ".jpeg" &&
			ext != ".png" &&
		  	ext != ".tif" &&
		  	ext != ".tiff" &&
		  	ext != ".bmp" &&
		  	ext != ".gif" {
			  continue
		}
		page++
		basepath := filepath.Dir(media.URL())
		filename := fmt.Sprintf("%03d%s", page, ext)
		url := filepath.Join(basepath, filename)
		fmt.Printf("%s\n", url)
		media.Store(url, option, rc)
		pages = append(pages, url)
	}
    media.SetPages(pages)
	return
}

func init() {
	media_library.RegisterMediaLibraryHandler("archive_handler", archiveHandler{})
}
