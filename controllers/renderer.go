package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-jsonnet"
)

func RenderTemplate() {
	files := []string{"/Users/jacobgo/Argos/templates/dash.jsonnet"}
	dashboard := renderJsonnet(files, "", false)
	err := postDashboard(dashboard)

	if err != nil {
		error := errors.New("help")
		fmt.Println(error)
	}
}

// RenderJsonnet ...
func renderJsonnet(files []string, param string, prune bool) string {

	// empty slice
	jsonnetPaths := files[:0]

	// range through the files
	for _, s := range files {
		jsonnetPaths = append(jsonnetPaths, fmt.Sprintf("(import '%s')", s))
	}

	// Create a JSonnet VM
	vm := jsonnet.MakeVM()

	// Join the slices into a jsonnet compat string
	jsonnetImport := strings.Join(jsonnetPaths, "+")

	fmt.Println(jsonnetImport)

	if param != "" {
		jsonnetImport = "(" + jsonnetImport + ")" + param
	}

	if prune {
		// wrap in std.prune, to remove nulls, empty arrays and hashes
		jsonnetImport = "std.prune(" + jsonnetImport + ")"
	}

	// render the jsonnet
	out, err := vm.EvaluateSnippet("file", jsonnetImport)

	if err != nil {
		log.Panic("Error evaluating jsonnet snippet: ", err)
	}

	log.Println(string(out))

	return out
}

func postDashboard(s string) error {
	client := &http.Client{}
	var UnmarshaledBody interface{}

	var jsonString = []byte(`{"dashboard": ` + s + `, "overwrite": true}`)

	req, err := http.NewRequest("POST", "http://localhost:3000/api/dashboards/db", bytes.NewBuffer(jsonString))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatalln("err in building req")
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln("err in posting dashboard")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Fatal("status code isn't 200 or 201: ", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("err in reading resp body")
	}

	json.Unmarshal(body, &UnmarshaledBody)

	fmt.Println("Created dashboard with response: ")
	fmt.Println(UnmarshaledBody)

	return nil
}
