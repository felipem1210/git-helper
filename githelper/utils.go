/*
Copyright © 2022 Felipe Macias felipem1210@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package githelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func WritePrsToJson(data MyPrs, json_file string) error {
	fullPathJsonFile := getEnvValue("WORKING_DIR") + "/" + json_file
	file, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(fullPathJsonFile, file, 0644)
	return err
}

func WriteReposToJson(data MyRepos, json_file string) error {
	fullPathJsonFile := getEnvValue("WORKING_DIR") + "/" + json_file
	file, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(fullPathJsonFile, file, 0644)
	return err
}

func (data *prCreateInfo) fromJsontoStruct(f string) *prCreateInfo {
	fullPathJsonFile := getEnvValue("WORKING_DIR") + "/" + f
	jsonFile, err := ioutil.ReadFile(fullPathJsonFile)
	CheckIfError(err)
	err = json.Unmarshal([]byte(jsonFile), &data)
	CheckIfError(err)
	return data
}

func (data MyPrs) fromJsontoSliceOfStructs(f string) MyPrs {
	fullPathJsonFile := getEnvValue("WORKING_DIR") + "/" + f
	jsonFile, err := ioutil.ReadFile(fullPathJsonFile)
	CheckIfError(err)
	err = json.Unmarshal([]byte(jsonFile), &data)
	CheckIfError(err)
	return data
}

func (data MyRepos) fromJsontoSliceOfStructs(f string) MyRepos {
	fullPathJsonFile := getEnvValue("WORKING_DIR") + "/" + f
	jsonFile, err := ioutil.ReadFile(fullPathJsonFile)
	CheckIfError(err)
	err = json.Unmarshal([]byte(jsonFile), &data)
	CheckIfError(err)
	return data
}

func ValidateEnv() {
	missing := make([]string, 0)
	envVars := []string{
		"WORKING_DIR",
		"GIT_ACCESS_USER",
		"GIT_ACCESS_TOKEN",
	}

	for _, v := range envVars {
		_, present := os.LookupEnv(v)
		if !present {
			missing = append(missing, v)
		}
	}
	if len(missing) != 0 {
		fmt.Printf("missing env vars: %v\n", missing)
		os.Exit(1)
	}
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func LogFatal(f string, err error) {
	log.Fatal(f, " error: ", err)
}

func getEnvValue(e string) string {
	return os.Getenv(e)
}

func createFolder(target string) {
	dir := getEnvValue("WORKING_DIR") + "/" + target
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		CheckIfError(err)
	}
}

func ListDirectories(target string) []string {
	var directories []string
	dir := getEnvValue("WORKING_DIR") + "/" + target
	dirs, err := ioutil.ReadDir(dir)
	CheckIfError(err)
	for _, d := range dirs {
		if d.IsDir() {
			if _, err := os.Stat(dir + "/" + d.Name() + "/.git"); !os.IsNotExist(err) {
				directories = append(directories, d.Name())
			}
		}
	}
	return directories
}
