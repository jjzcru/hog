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

func AddFiles(files []string) (string, error) {
	var groupID string
	hogPath, err := GetPath()
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

	groupID = NewGroupID(hog)

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
	err = Save(hogPath, hog)
	if err != nil {
		return groupID, err
	}

	return groupID, err
}

func GetPath() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, "hog.yml"), nil
}

func GetBaseDir() (string, error) {
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

func Save(hogPath string, hog Hog) error {
	content, err := yaml.Marshal(hog)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hogPath, content, 0777)
}

func FromPath(hogPath string) (Hog, error) {
	hog := Hog{}

	content, err := ioutil.ReadFile(hogPath)
	if err != nil {
		return hog, err
	}

	err = yaml.Unmarshal(content, &hog)

	return hog, nil
}

func NewGroupID(hog Hog) string {
	var id string
	for {
		id = GetID()
		if _, ok := hog.Buckets[id]; !ok {
			break
		}
	}
	return id
}

func CreateEmptyHogFile(hogPath string) error {
	return Save(hogPath, defaultHog)
}
