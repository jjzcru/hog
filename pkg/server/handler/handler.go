package handler

import (
	"archive/zip"
	"context"
	"fmt"
	gh "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/server/graphql"
	"github.com/jjzcru/hog/pkg/server/graphql/generated"
	"github.com/jjzcru/hog/pkg/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func GraphQL(token string) http.HandlerFunc {
	return AddAuth(gh.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{}})), token)
}

func Playground(url string) http.HandlerFunc {
	return playground.Handler("GraphQL Playground", url)
}

func AddAuth(next http.Handler, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := context.WithValue(r.Context(), graphql.TokenKey, token)
		ctx := context.WithValue(r.Context(), graphql.AuthorizationKey, r.Header.Get("auth-token"))
		next.ServeHTTP(w, r.WithContext(ctx))
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

		if h.Files == nil {
			h.Files = map[string][]string{}
		}

		if _, ok := h.Files[id]; !ok {
			notFoundError(w, fmt.Errorf("id '%s' not found", id))
			return
		}

		files := h.Files[id]
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
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
