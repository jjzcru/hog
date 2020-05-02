package handler

import (
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
			notImplemented(w, fmt.Errorf("multiple files not enable yet"))
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

	notImplemented(w, fmt.Errorf("download directory not implement yet"))
}

func downloadFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)

	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), 500)
		return
	}

	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Printf("Error: '%s'", err.Error())
		}
	}()

	contentType, err := getFileContentType(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	filename := filepath.Base(filePath)
	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("content-length", fmt.Sprintf("%d", fi.Size()))
	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, filePath)
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
