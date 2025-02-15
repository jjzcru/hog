package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/server/handler"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/logrusorgru/aurora"
)

// Start the server from a particular file
func Start(port int, hogPath string, token string) error {
	if !IsPortOpen(port) {
		return fmt.Errorf("another application is running on port %d", port)
	}

	if !utils.IsPathExist(hogPath) {
		err := hog.CreateEmptyHogFile(hogPath)
		if err != nil {
			return err
		}
	}

	h, err := hog.FromPath(hogPath)
	if err != nil {
		return err
	}

	h.Port = port

	err = hog.SaveToPath(hogPath, h)
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/download/{id}", handler.Download(hogPath))
	r.HandleFunc("/download/{id}/", handler.Download(hogPath))
	r.HandleFunc("/qr/{id}", handler.Qr(hogPath))
	r.HandleFunc("/qr/{id}/", handler.Qr(hogPath))

	if len(token) > 0 {
		fmt.Println(strings.Join([]string{
			aurora.Bold("Authorization token:").String(),
			aurora.Bold(aurora.Cyan(token)).String(),
		}, " "))
	}

	fmt.Println(strings.Join([]string{
		aurora.Bold("Server running on port").String(),
		aurora.Bold(aurora.Green(fmt.Sprintf("%d 🚀", port))).String(),
	}, " "))

	srv := &http.Server{
		Handler: r,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			ctx = context.WithValue(ctx, handler.TokenKey, token)
			return ctx
		},
		Addr: fmt.Sprintf(":%d", port),
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	return srv.ListenAndServe()
}

// IsPortOpen check if a port is open in the current machine
func IsPortOpen(port int) bool {
	l, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", port))
	defer func() {
		if l != nil {
			err = l.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	return err == nil
}
