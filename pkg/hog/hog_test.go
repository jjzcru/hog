package hog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestGetBaseDir(t *testing.T) {
	err := os.Setenv("HOG_PEN", os.TempDir())
	if err != nil {
		t.Error(err)
	}

	baseDir, err := BaseDir()
	if err != nil {
		t.Error(err)
	}

	if baseDir != os.TempDir() {
		t.Errorf("base dir should be '%s' but is '%s' instead", os.TempDir(), baseDir)
	}

	_ = os.Setenv("HOG_PEN", fmt.Sprintf("%s/%s", os.TempDir(), time.Now().String()))
	_, err = BaseDir()
	if err == nil {
		t.Error("it should throw an error because the directory do not exist")
	}
}

func TestGetPath(t *testing.T) {
	err := os.Setenv("HOG_PEN", os.TempDir())
	if err != nil {
		t.Error(err)
	}

	hogPath, err := Path()
	if err != nil {
		t.Error(err)
	}

	targetHogPath, err := filepath.Abs(filepath.Join(os.TempDir(), FILE))
	if err != nil {
		t.Error(err)
	}

	if hogPath != targetHogPath {
		t.Errorf("hog path should be '%s' but is '%s' instead", targetHogPath, hogPath)
	}
}

func TestAddFilesEmpty(t *testing.T) {
	err := os.Setenv("HOG_PEN", os.TempDir())
	if err != nil {
		t.Error(err)
	}

	id, err := AddFiles([]string{})
	if err != nil {
		t.Error(err)
	}

	defer func() {
		hogPath, _ := Path()
		_ = os.Remove(hogPath)
	}()

	if len(id) == 0 {
		t.Errorf("the id should not be empty")
	}

	h, err := Get()
	if err != nil {
		t.Error(err)
	}

	files := h.Buckets[id]

	if len(files) != 0 {
		t.Errorf("the files should be empty but it has %d files instead", len(files))
	}
}

func TestFromPath(t *testing.T) {
	err := os.Setenv("HOG_PEN", os.TempDir())
	if err != nil {
		t.Error(err)
	}

	_, err = AddFiles([]string{})
	if err != nil {
		t.Error(err)
	}

	hogPath, err := Path()
	if err != nil {
		t.Error(err)
	}

	hog, err := FromPath(hogPath)
	if err != nil {
		t.Error(err)
	}

	buckets := hog.Buckets

	if len(buckets) == 0 {
		t.Error("the total of buckets should not be empty")
	}

	_ = os.Remove(hogPath)

	_, err = FromPath(os.TempDir())
	if err == nil {
		t.Error("it should throw an error because its a directory")
	}

	_, err = FromPath(fmt.Sprintf("%s/%s", os.TempDir(), "test.yml"))
	if err == nil {
		t.Error("it should throw an error because the file do not exist")
	}
}

func TestSave(t *testing.T) {
	err := os.Setenv("HOG_PEN", os.TempDir())
	if err != nil {
		t.Error(err)
	}

	hog := Hog{
		Domain:   "127.0.0.1",
		Port:     3000,
		Protocol: "https",
		Buckets:  map[string][]string{},
	}

	id := newGroupID(hog)

	hog.Buckets[id] = []string{"/test.png", "/test.jpg"}

	err = Save(hog)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		hogPath, _ := Path()
		_ = os.Remove(hogPath)
	}()

	savedHog, err := Get()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(savedHog, hog) {
		t.Error("save hog and memory hog should be the same")
	}
}
