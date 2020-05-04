package handler

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	qrcode "github.com/skip2/go-qrcode"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Qr returns a qr code to the file to download
func Qr(hogPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := vars["id"]

		h, err := hog.FromPath(hogPath)
		if err != nil {
			serverError(w, err)
			return
		}

		port := r.URL.Query().Get("port")
		if len(port) > 0 {
			h.Port, err = strconv.Atoi(port)
			if err != nil {
				serverError(w, err)
				return
			}
		}

		domain := r.URL.Query().Get("domain")
		if len(domain) > 0 {
			h.Domain = domain
		}

		protocol := r.URL.Query().Get("protocol")
		if len(protocol) > 0 && (protocol == "http" || protocol == "https") {
			h.Protocol = protocol
		}

		if h.Buckets == nil {
			h.Buckets = map[string][]string{}
		}

		if _, ok := h.Buckets[id]; !ok {
			notFoundError(w, fmt.Errorf("id '%s' not found", id))
			return
		}

		url := fmt.Sprintf("%s://%s:%d/download/%s", h.Protocol, h.Domain, h.Port, id)

		png, err := qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			serverError(w, err)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(png)))
		if _, err := w.Write(png); err != nil {
			serverError(w, errors.New("unable to write image"))
			return
		}
	}
}

func Download(hogPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := vars["id"]

		h, err := hog.FromPath(hogPath)
		if err != nil {
			serverError(w, err)
			return
		}

		if h.Buckets == nil {
			h.Buckets = map[string][]string{}
		}

		if _, ok := h.Buckets[id]; !ok {
			notFoundError(w, fmt.Errorf("id '%s' not found", id))
			return
		}

		files := h.Buckets[id]
		switch len(files) {
		case 0:
			notFoundError(w, fmt.Errorf("not files found for id '%s'", id))
			return
		case 1:
			downloadPath(w, r, files[0])
			return
		default:
			downloadMultiPath(w, r, files)
			// notImplemented(w, fmt.Errorf("multiple files not enable yet"))
		}
	}
}

func downloadPath(w http.ResponseWriter, r *http.Request, filePath string) {
	isFile, err := utils.IsPathAFile(filePath)
	if err != nil {
		serverError(w, err)
		return
	}

	if isFile {
		downloadFile(w, r, filePath)
		return
	}

	downloadDirectory(w, r, filePath)
}

func downloadFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)

	if err != nil && err != io.EOF {
		serverError(w, err)
		return
	}

	defer func() {
		err = file.Close()
		if err != nil {
			utils.PrintError(err)
		}
	}()

	contentType, err := getFileContentType(file)
	if err != nil {
		serverError(w, err)
		return
	}

	filename := filepath.Base(filePath)
	fi, err := file.Stat()
	if err != nil {
		serverError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("content-length", fmt.Sprintf("%d", fi.Size()))
	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, filePath)
}

func downloadDirectory(w http.ResponseWriter, r *http.Request, filePath string) {
	zipFileName := fmt.Sprintf("%s.zip", utils.GetToken())
	zipFilePath := filepath.Join(os.TempDir(), "/", zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		serverError(w, err)
		return
	}

	defer func(file *os.File, filePath string) {
		err = file.Close()
		if err != nil {
			utils.PrintError(err)
		}

		err = os.Remove(filePath)
		if err != nil {
			utils.PrintError(err)
		}
	}(zipFile, zipFilePath)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(zipFile)
	// Add some files to the archive.
	addFiles(zipWriter, filePath, "")

	err = zipWriter.Close()
	if err != nil {
		serverError(w, err)
		return
	}

	filename := filepath.Base(zipFilePath)
	fi, err := zipFile.Stat()
	if err != nil {
		serverError(w, err)
		return
	}

	fmt.Println("File path: ", zipFilePath)

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("content-length", fmt.Sprintf("%d", fi.Size()))
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, zipFilePath)
}

func downloadMultiPath(w http.ResponseWriter, r *http.Request, filePaths []string) {
	zipFileName := fmt.Sprintf("%s.zip", utils.GetToken())
	zipFilePath := filepath.Join(os.TempDir(), "/", zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		serverError(w, err)
		return
	}

	defer func(file *os.File, filePath string) {
		err = file.Close()
		if err != nil {
			utils.PrintError(err)
		}

		err = os.Remove(filePath)
		if err != nil {
			utils.PrintError(err)
		}
	}(zipFile, zipFilePath)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(zipFile)
	for _, filePath := range filePaths {
		addFiles(zipWriter, filePath, filepath.Base(filePath))
	}

	err = zipWriter.Close()
	if err != nil {
		serverError(w, err)
		return
	}

	filename := filepath.Base(zipFilePath)
	fi, err := zipFile.Stat()
	if err != nil {
		serverError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("content-length", fmt.Sprintf("%d", fi.Size()))
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, zipFilePath)
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	isFile, err := utils.IsPathAFile(basePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isFile {
		dat, err := ioutil.ReadFile(basePath)
		if err != nil {
			fmt.Println(err)
		}
		f, err := w.Create(filepath.Base(filepath.Join("", basePath)))
		if err != nil {
			fmt.Println(err)
		}
		_, err = f.Write(dat)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				fmt.Println(err)
			}
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			newBase := filepath.Join(basePath, file.Name())
			addFiles(w, newBase, filepath.Join(baseInZip, file.Name()))
		}
	}
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
