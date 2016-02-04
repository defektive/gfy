package scanner

import (
	"crypto/md5"
	"fmt"
	"github.com/xiam/exif"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const filechunk = 8192
var destination = ""
type Photo struct {
	Path string
	Name string
	Ext string
	Date string
	Hash string
	Destination string
}

func (photo Photo) SortedPath() string {
	return path.Join(destination, sortedBasePath(photo.Date))
}

func (photo Photo) SortedName() string {
	dateStr := strings.Replace(photo.Date[0:10], ":", "_", -1)
	return "GFY_" + dateStr +"_"+photo.Hash + photo.Ext
}

func (photo Photo) SortedFullPath() string {
	return path.Join(photo.SortedPath(), photo.SortedName())
}

func ScanDir(dir string, dest string) (photos []*Photo) {
	destination = dest
	files := allPhotos(dir)
	for _, filePath := range files {
		f, _ := os.Stat(filePath)
		photo := &Photo{
			Path: filePath,
			Name: f.Name(),
			Ext: strings.ToLower(path.Ext(filePath)),
			Date: date(filePath),
			Hash: hash(filePath),
		}
		photos = append(photos, photo)

	}

	return photos
}

func allPhotos(dir string) (photoFiles []string) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		filePath := path.Join(dir, f.Name())

		switch mode := f.Mode(); {
			case mode.IsDir():
					photoFiles = append(photoFiles, allPhotos(filePath)...)
			case mode.IsRegular():

				ext := strings.ToLower(path.Ext(filePath))

				if ext == ".jpg" {
					fmt.Println("Adding " +filePath)
					photoFiles = append(photoFiles, filePath)
				}
		}
	}
	return photoFiles
}

func sortedBasePath(date string) string {
	return strings.Join(strings.Split(date[0:7], ":"), string(filepath.Separator))
}

func date(filename string) string {
	data, err := exif.Read(filename)
	if err != nil{
		return "0000:00:00"
	}

	date := data.Tags["Date and Time"]

	if date == "" {
		return "0000:00:00"
	}

	return date
}

func hash(filename string) string {

	file, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
