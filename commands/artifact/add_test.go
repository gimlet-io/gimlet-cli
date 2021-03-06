package artifact

import (
	"encoding/json"
	"fmt"
	"github.com/franela/goblin"
	"github.com/gimlet-io/gimlet-cli/commands"
	"github.com/gimlet-io/gimletd/dx"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const artifactToExtend = `
{
  "id": "my-app-b2ab0f7a-ca0e-45cf-83a0-cadd94dddeac",
  "version": {
    "sha": "ea9ab7cc31b2599bf4afcfd639da516ca27a4780",
    "branch": "master",
    "authorName": "Jane Doe",
    "authorEmail": "jane@doe.org",
    "committerName": "Jane Doe",
    "committerEmail": "jane@doe.org",
    "message": "Bugfix 123",
    "repositoryName": "my-app",
    "url": "https://github.com/gimlet-io/gimlet-cli/commit/ea9ab7cc31b2599bf4afcfd639da516ca27a4780"
  }
}
`

const env = `
app: fosdem-2021
env: staging
namespace: default
deploy:
  branch: master
  event: push
chart:
  repository: https://chart.onechart.dev
  name: onechart
  version: 0.10.0
values:
  replicas: 1
  image:
    repository: ghcr.io/gimlet-io/fosdem-2021
    tag: "{{ .GITHUB_SHA }}"
`

func Test_add(t *testing.T) {
	artifactFile, err := ioutil.TempFile("", "gimlet-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(artifactFile.Name())
	ioutil.WriteFile(artifactFile.Name(), []byte(artifactToExtend), commands.File_RW_RW_R)

	envFile, err := ioutil.TempFile("", "gimlet-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(envFile.Name())
	ioutil.WriteFile(envFile.Name(), []byte(env), commands.File_RW_RW_R)

	g := goblin.Goblin(t)
	g.Describe("gimlet artifact add", func() {
		//g.It("Should add CI URL to artifact", func() {
		//	args := strings.Split("gimlet artifact add", " ")
		//	args = append(args, "-f", artifactFile.Name())
		//	args = append(args, "--field", "name=CI")
		//	args = append(args, "--field", "url=https://jenkins.example.com/job/dev/84/display/redirect")
		//	err = commands.Run(&Command, args)
		//	g.Assert(err == nil).IsTrue(err)
		//
		//	content, err := ioutil.ReadFile(artifactFile.Name())
		//	var a dx.Artifact
		//	err = json.Unmarshal(content, &a)
		//	g.Assert(err == nil).IsTrue(err)
		//	g.Assert(len(a.Items) == 1).IsTrue("Should have 1 item")
		//	g.Assert(a.Items[0]["name"] == "CI").IsTrue("Should add CI item")
		//	g.Assert(a.Items[0]["url"] == "https://jenkins.example.com/job/dev/84/display/redirect").IsTrue("Should add CI item")
		//	//fmt.Println(string(content))
		//})
		g.It("Should add Gimlet environment to artifact", func() {
			args := strings.Split("gimlet artifact add", " ")
			args = append(args, "-f", artifactFile.Name())
			args = append(args, "--envFile", envFile.Name())
			err = commands.Run(&Command, args)
			g.Assert(err == nil).IsTrue(err)

			content, err := ioutil.ReadFile(artifactFile.Name())
			var a dx.Artifact
			err = json.Unmarshal(content, &a)
			g.Assert(err == nil).IsTrue(err)
			g.Assert(len(a.Environments) == 1).IsTrue("Should have 1 env")
			g.Assert(a.Environments[0].App == "fosdem-2021").IsTrue("Should add env")
			fmt.Println(string(content))
		})
		//g.It("Should add context variables to artifact", func() {
		//	args := strings.Split("gimlet artifact add", " ")
		//	args = append(args, "-f", artifactFile.Name())
		//	args = append(args, "--var", "KEY=VALUE")
		//	args = append(args, "--var", "KEY2=VALUE2")
		//	err = commands.Run(&Command, args)
		//	g.Assert(err == nil).IsTrue(err)
		//
		//	content, err := ioutil.ReadFile(artifactFile.Name())
		//	var a dx.Artifact
		//	err = json.Unmarshal(content, &a)
		//	g.Assert(err == nil).IsTrue(err)
		//	g.Assert(len(a.Context) == 2).IsTrue("Should have 1 var in context")
		//	g.Assert(a.Context["KEY"] == "VALUE").IsTrue("Should add var")
		//	fmt.Println(string(content))
		//})
	})
}