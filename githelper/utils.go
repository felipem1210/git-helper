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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func WritePrsToJson(data MyPrsJson, json_file string) error {
	file, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(getEnvValue("WORKDING_DIR")+json_file, file, 0644)
	return err
}

func WriteReposToJson(data MyRepos, json_file string) error {
	file, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(getEnvValue("WORKDING_DIR")+json_file, file, 0644)
	return err
}

func GetFromJsonReturnArray(f string, d string) []string {
	// Open our jsonFile
	jsonFile, err := ioutil.ReadFile(f)
	// if we os.Open returns an error then handle it
	CheckIfError(err)
	var m []myReposJson
	var data []string
	err = json.Unmarshal([]byte(jsonFile), &m)
	CheckIfError(err)
	for _, val := range m {
		r := reflect.ValueOf(val)
		f := reflect.Indirect(r).FieldByName(d)
		data = append(data, f.String())
	}
	return data
}

func (data *prCreateInfo) fromJsontoStruct(f string) *prCreateInfo {
	jsonFile, err := ioutil.ReadFile(f)
	CheckIfError(err)
	err = json.Unmarshal([]byte(jsonFile), &data)
	CheckIfError(err)
	return data
}

func (data MyPrsJson) fromJsontoSliceOfStructs(f string) MyPrsJson {
	jsonFile, err := ioutil.ReadFile(f)
	CheckIfError(err)
	err = json.Unmarshal([]byte(jsonFile), &data)
	CheckIfError(err)
	return data
}

func (data MyRepos) fromJsontoSliceOfStructs(f string) MyRepos {
	jsonFile, err := ioutil.ReadFile(f)
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
