package scanner

import (
	"crypto/md5"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/tj/go-spin"
	"github.com/xiam/exif"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const defaultDate = "0000:00:00"
const thumbDir = ".gfy_thumbs"
const filechunk = 8192
const thumbWidth = 256

type Photo struct {
	Path string
	Name string
	Ext  string
	Date string
	hash string
}

func (photo *Photo) Hash() string {
	if photo.hash == "" {
		photo.hash = hashFile(photo.Path)
	}
	return photo.hash
}

func (photo *Photo) Datetime() time.Time {

	const shortForm = "2006:01:02 15:04:05 -0700 MST"
	t, _ := time.Parse(shortForm, photo.Date+" -0700 MST")
	return t
}

func (photo *Photo) SortedPath(destination string) string {
	return path.Join(destination, sortedBasePath(photo.Date))
}

func (photo *Photo) SortedName() string {
	if photo.Date == defaultDate {
		return photo.Name
	}
	return "GFY_" + photo.Hash() + photo.Ext
}

func (photo *Photo) SortedFullPath(destination string) string {
	return path.Join(photo.SortedPath(destination), photo.SortedName())
}

func (photo *Photo) SortedThumbPath(destination string) string {
	return path.Join(photo.SortedPath(path.Join(destination, thumbDir)))
}

func (photo *Photo) SortedFullThumbPath(destination string) string {
	return path.Join(photo.SortedThumbPath(destination), photo.SortedName())
}

func (photo *Photo) Thumbnail(destinationDir string) string {
	file, err := os.Open(photo.Path)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(thumbWidth, 0, img, resize.Lanczos3)

	os.MkdirAll(photo.SortedThumbPath(destinationDir), 0777)
	out, err := os.Create(photo.SortedFullThumbPath(destinationDir))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return photo.SortedFullThumbPath(destinationDir)
}

func ScanDir(dir string) (photos []*Photo) {
	spinner := spin.New()
	files := allPhotos(dir, spinner)
	spinner.Set(spin.Box1)

	for _, filePath := range files {
		fmt.Printf("\r  \033[36mcomputing %s\033[m %s ", dir, spinner.Next())

		f, _ := os.Stat(filePath)
		photo := &Photo{
			Path: filePath,
			Name: f.Name(),
			Ext:  strings.ToLower(path.Ext(filePath)),
			Date: date(filePath),
		}
		photos = append(photos, photo)
	}

	return photos
}

func allPhotos(dir string, spinner *spin.Spinner) (photoFiles []string) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fmt.Printf("\r  \033[36mscanning %s\033[m ", spinner.Next())

		filePath := path.Join(dir, f.Name())

		switch mode := f.Mode(); {
		case mode.IsDir() && f.Name() != thumbDir:
			photoFiles = append(photoFiles, allPhotos(filePath, spinner)...)
		case mode.IsRegular():

			ext := strings.ToLower(path.Ext(filePath))

			if ext == ".jpg" {
				photoFiles = append(photoFiles, filePath)
			}
		}
	}
	return photoFiles
}

func sortedBasePath(date string) string {
	return strings.Join(strings.Split(date[0:10], ":"), string(filepath.Separator))
}

func date(filename string) string {
	data, err := exif.Read(filename)
	if err != nil {
		return defaultDate
	}

	for _, tag := range dateTags() {
		if data.Tags[tag] != "" {
			return data.Tags[tag]
		}
	}

	return defaultDate
}

func dateTags() (tags []string) {
	tags = append(tags, "Date and Time (Original)")
	tags = append(tags, "Date and Time")
	return tags
}

func hashFile(filename string) string {

	file, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		_, err := file.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err.Error())
		}

		hash.Write(buf)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
