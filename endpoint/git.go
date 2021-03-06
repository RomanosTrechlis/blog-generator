package endpoint

import (
	"fmt"
	"os/exec"

	"strings"

	"github.com/RomanosTrechlis/blog-generator/util/fs"
)

// gitEndpoint is the git endpoint object
type gitEndpoint struct{}

// newGitEndpoint creates a new GitEndpoint
func newGitEndpoint() (e Endpoint) {
	return &gitEndpoint{}
}

// Upload uploads the site to a git repository
// todo: push fails
func (ds *gitEndpoint) Upload(destFolder, endpointUsername, endpointPassword, endpointURL string) (err error) {
	fmt.Println("Uploading Site...")
	dest := destFolder + "_upload"
	err = fs.CreateFolderIfNotExist(dest)
	if err != nil {
		return err
	}
	err = fs.ClearFolder(dest)
	if err != nil {
		return err
	}

	cmdName := "git"
	initArgs := []string{"init", "."}
	cmd := exec.Command(cmdName, initArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error initializing git repository at %s: %v", destFolder, err)
	}

	url, err := createUrlWithCred(endpointUsername, endpointPassword, endpointURL)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	remoteArgs := []string{"remote", "add", "origin", url}
	cmd = exec.Command(cmdName, remoteArgs...)
	cmd.Dir = dest
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("error creating upload folder %s: %v", dest, err)
	}
	err = fs.CopyDir(destFolder, dest)
	if err != nil {
		return fmt.Errorf("error copying generated folder %s to upload folder %s: %v",
			destFolder, dest, err)
	}

	addArgs := []string{"add", "."}
	cmd = exec.Command(cmdName, addArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding files to commit: %v", err)
	}

	commitArgs := []string{"commit", "-m", "auto commit"}
	cmd = exec.Command(cmdName, commitArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error committing files: %v", err)
	}

	pushArgs := []string{"push", "origin", "master"}
	cmd = exec.Command(cmdName, pushArgs...)
	cmd.Dir = dest
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing to remote %s: %v", endpointURL, err)
	}
	fmt.Println("Upload Complete.")
	return nil
}

func createUrlWithCred(username, password, to string) (url string, err error) {
	t := strings.Split(to, "://")
	if len(t) != 2 {
		return "", fmt.Errorf("couldn't process git url")
	}
	p := strings.Replace(password, "@", "%40", 5)
	return t[0] + "://" + username + ":" + p + "@" + t[1], nil
}
