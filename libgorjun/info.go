package gorjun

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"log"
	"os/exec"
	"time"
	"math/rand"
)

// List returns a list of files
func (g *GorjunServer) List(artifactType string, parameters string) ([]GorjunFile, error) {
	resp, err:= http.Get(fmt.Sprintf("http://%s/kurjun/rest/" + artifactType + "/info", g.Hostname))
	if len(parameters) != 0 {
		resp, err = http.Get(fmt.Sprintf("http://%s/kurjun/rest/" + artifactType + "/info" + parameters, g.Hostname))
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve file list from %s: %v", g.Hostname, err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read body from %s: %v", g.Hostname, err)
	}
	var rf []GorjunFile
	err = json.Unmarshal(data, &rf)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", data)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal contents from %s: %v", g.Hostname, err)
	}
	return rf, nil
}

//Upload templates
func (g *GorjunServer) UploadTemplates() (error){
	err := g.AuthenticateUser()
	if err != nil {
		fmt.Errorf("Authnetication failure: %v", err)
	}

	templateVersions := []string{"0.1.6", "0.1.7", "0.1.9", "0.1.10", "0.1.11"}
	rand.Seed(time.Now().UnixNano())

	for _, version := range templateVersions {
		id, err := g.Upload("data/nginx-subutai-template_"+version+"_amd64.tar.gz", "template")
		if err != nil {
			fmt.Errorf("Failed to upload: %v", err)
		}
		fmt.Printf("Template uploaded successfully, id : %s\n", id)
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

//Delete templates
func (g *GorjunServer) DeleteTemplates() (error){
	err := g.AuthenticateUser()
	if err != nil {
		fmt.Errorf("Authnetication failure: %v", err)
	}

	flist, err := g.List("template","")
	for _, template := range flist {
		err = g.RemoveFileByID(template.ID, "template")
		if err != nil {
			fmt.Errorf("Failed to remove file: %v", err)
		}
		fmt.Printf("Template removed successfully, id : %s\n", template.ID)
	}
	showFileSystemState()
	return err
}

func showFileSystemState()  {
	output, _ := exec.Command("bash", "-c", " ls /opt/gorjun/data/files/").Output()
	fmt.Printf("\nList of files in /opt/gorjun/data/files/ directory after deleting templates \n%s\n", output)
}
