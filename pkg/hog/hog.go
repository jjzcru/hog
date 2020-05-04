package hog

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/jjzcru/hog/pkg/utils"
	"gopkg.in/yaml.v2"
)

// FILE name of the hog configuration file
var FILE = ".hog.yml"

var defaultHog = Hog{
	Domain:   "localhost",
	Protocol: "http",
	Port:     1618,
	Buckets:  map[string][]string{},
}

// Hog hold the structure of the files
type Hog struct {
	Domain   string              `yaml:"domain"`
	Protocol string              `yaml:"protocol"`
	Port     int                 `yaml:"port"`
	Buckets  map[string][]string `yaml:"buckets"`
}

// AddFiles to the hog file and return the bucket id
func AddFiles(files []string) (string, error) {
	var groupID string
	hogPath, err := Path()
	if err != nil {
		return groupID, err
	}

	hog := defaultHog

	if !utils.IsPathExist(hogPath) {
		err := CreateEmptyHogFile(hogPath)
		if err != nil {
			return groupID, err
		}
	} else {
		hog, err = FromPath(hogPath)
		if err != nil {
			return groupID, err
		}
	}

	groupID = newGroupID(hog)

	if hog.Buckets == nil {
		hog.Buckets = map[string][]string{}
	}

	if len(hog.Domain) == 0 {
		hog.Domain = "localhost"
	}

	if len(hog.Protocol) == 0 {
		hog.Protocol = "http"
	}

	if hog.Port == 0 {
		hog.Port = 1618
	}

	hog.Buckets[groupID] = files
	err = SaveToPath(hogPath, hog)
	if err != nil {
		return groupID, err
	}

	return groupID, err
}

// Path return the path of the hog file
func Path() (string, error) {
	baseDir, err := BaseDir()
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, FILE), nil
}

// BaseDir return the path of the directory that contains the hog file
func BaseDir() (string, error) {
	hogPen := os.Getenv("HOG_PEN")
	if len(hogPen) == 0 {
		return getHomeDir()
	}

	if !utils.IsPathExist(hogPen) {
		return "", fmt.Errorf("HOG_PEN '%s' do not exist ❌", hogPen)
	}

	isFile, err := utils.IsPathAFile(hogPen)
	if err != nil {
		return "", err
	}

	if isFile {
		return "", fmt.Errorf("HOG_PEN '%s' must be a directory 📁", hogPen)
	}

	return hogPen, nil
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// Get the hog object from the hog path
func Get() (Hog, error) {
	var hog Hog
	hogPath, err := Path()
	if err != nil {
		return hog, err
	}

	if !utils.IsPathExist(hogPath) {
		return hog, fmt.Errorf("hog path '%s' do not exist", hogPath)
	}

	hog, err = FromPath(hogPath)
	if err != nil {
		return hog, err
	}

	return hog, nil
}

// Save a hog object in the default hog path
func Save(hog Hog) error {
	hogPath, err := Path()
	if err != nil {
		return err
	}

	return SaveToPath(hogPath, hog)
}

// SaveToPath saves a hog object in a specific path
func SaveToPath(hogPath string, hog Hog) error {
	content, err := yaml.Marshal(hog)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hogPath, content, 0777)
}

// FromPath load a hog object from an specific path
func FromPath(hogPath string) (Hog, error) {
	hog := Hog{}

	content, err := ioutil.ReadFile(hogPath)
	if err != nil {
		return hog, err
	}

	err = yaml.Unmarshal(content, &hog)
	if err != nil {
		return hog, err
	}

	return hog, nil
}

func newGroupID(hog Hog) string {
	var id string
	for {
		id = GetID()
		if _, ok := hog.Buckets[id]; !ok {
			break
		}
	}
	return id
}

// CreateEmptyHogFile create a default hog object in the default hog path
func CreateEmptyHogFile(hogPath string) error {
	return SaveToPath(hogPath, defaultHog)
}
