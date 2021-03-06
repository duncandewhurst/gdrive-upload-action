// Mathis Van Eetvelde
// Ghent 2021

package main

import (
	"fmt"
	"context"
	"encoding/base64"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"github.com/sethvargo/go-githubactions"
)

const (
	scope = "https://www.googleapis.com/auth/drive.file"
)

func main() {

	// get file argument from action input
	filename := githubactions.GetInput("filename")
	if filename == "" {
		githubactions.Fatalf("missing input 'filename'")
	}

	// get parentId argument from action input
	parentId := githubactions.GetInput("parentId")
	if parentId == "" {
		githubactions.Fatalf("missing input 'parentId'")
	}

	// get base64 encoded credentials argument from action input
	credentials := githubactions.GetInput("credentials")
	if credentials == "" {
		githubactions.Fatalf("missing input 'credentials'")
	}
	// add base64 encoded credentials argument to mask
	githubactions.AddMask(credentials)

	// decode credentials to []byte
	decodedCredentials, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("base64 decoding of 'credentials' failed with error: %v", err))
	}
	// add decoded credentials argument to mask
	githubactions.AddMask(string(decodedCredentials))

	// fetching a JWT config with credentials and the right scope
	conf, err := google.JWTConfigFromJSON(decodedCredentials, scope)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("fetching JWT credentials failed with error: %v", err))
	}

	// instantiating a new drive service
	ctx := context.Background()
	svc, err := drive.New(conf.Client(ctx))
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("opening file with filename: %v failed with error: %v", filename, err))
	}

	f := drive.File{
		Name: file.Name(),
		Parents: []string{parentId},
	}

	_, err = svc.Files.Create(&f).Media(file).Do()
	if err != nil {
		githubactions.Fatalf(fmt.Sprintf("Creating file: %+v failed with error: %v", f, err))
	}

}
