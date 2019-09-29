package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (

	//SaveDiameter Diamater of target save
	SaveDiameter string
	//SaveDir contains output save directory
	SaveDir string
	//SaveEnd option to save end world
	SaveEnd string
	//SaveName Name that will prefix save files
	SaveName string
	//SaveNether option to save nether world
	SaveNether string
	//ServerName name of directory spigot.jar is located - name of your server
	ServerName string
	//ServerRootDirectory contains folder with server folder
	ServerRootDirectory string
	//WorldName name of world to save
	WorldName string
	//WebHookURL webhook url string
	WebHookURL string

	config *configStruct
)

type configStruct struct {

	SaveDiameter        string `json:"SaveDiameter"`
	SaveDir             string `json:"SaveDir"`
	SaveEnd             string `json:"SaveEnd"`
	SaveName            string `json:"SaveName"`
	SaveNether          string `json:"SaveNether"`
	ServerName          string `json:"ServerName"`
	ServerRootDirectory string `json:"ServerRootDirectory"`
	WorldName           string `json:"WorldName"`
	WebHookURL          string `json:"WebHookURL"`
}

//ReadConfig reads config.json file.
func ReadConfig() error {

	executableLocation, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(executableLocation)
	configLocation := exPath + "/config.json"
	file, err := ioutil.ReadFile(configLocation)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	SaveDiameter = config.SaveDiameter
	SaveDir = config.SaveDir
	SaveEnd = config.SaveEnd
	SaveName = config.SaveName
	SaveNether = config.SaveNether
	SaveEnd = config.SaveEnd
	ServerName = config.ServerName
	ServerRootDirectory = config.ServerRootDirectory
	WorldName = config.WorldName
	WebHookURL = config.WebHookURL

	return nil
}
