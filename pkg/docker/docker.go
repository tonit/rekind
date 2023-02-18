package docker

import (
	"encoding/json"
	"errors"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"strings"
)

type DockerPsOutput struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Names        string `json:"Names"`
	Networks     string `json:"Networks"`
	Ports        string `json:"Ports"`
	RunningFor   string `json:"RunningFor"`
	Size         string `json:"Size"`
	State        string `json:"State"`
	Status       string `json:"Status"`
}

type ImageList struct {
	Images []struct {
		ID          string        `json:"id"`
		RepoTags    []string      `json:"repoTags"`
		RepoDigests []interface{} `json:"repoDigests"`
		Size        string        `json:"size"`
		UID         interface{}   `json:"uid"`
		Username    string        `json:"username"`
		Spec        interface{}   `json:"spec"`
		Pinned      bool          `json:"pinned"`
	} `json:"images"`
}

func FindContainer(contextName string) (string, error) {

	var wantedControlPlane = contextName + "-control-plane"

	command := execute.ExecTask{
		Command:     "docker",
		Args:        []string{"ps", "--format", "'{{json .}}'"},
		StreamStdio: false,
	}

	res, err := command.Execute()
	if err != nil {
		panic(err)
	}
	var output []string

	var in = res.Stdout
	output = cleanseLinefeedOutput(in, output)

	for _, clean := range output {
		var data DockerPsOutput
		error := json.Unmarshal([]byte(clean), &data)
		if error != nil {
			panic(error)
		}
		if data.Names == wantedControlPlane {
			return data.ID, nil
		}
	}
	return "", errors.New("No control plane container found with context " + contextName)
}

func ListImages(containerId string) ImageList {

	command := execute.ExecTask{
		Command:     "docker",
		Args:        []string{"exec", containerId, "crictl", "images", "-o", "json"},
		StreamStdio: false,
	}

	res, err := command.Execute()
	if err != nil {
		panic(err)
	}
	var in = res.Stdout
	var data ImageList
	error := json.Unmarshal([]byte(in), &data)
	if error != nil {
		panic(error)
	}
	return data
}

func cleanseLinefeedOutput(in string, output []string) []string {
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		clean := strings.Trim(line, "'")
		if clean != "" {
			output = append(output, clean)
		}
	}
	return output
}
