package generator

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

func clearAndCreateDestination(path string) (err error) {
	err = os.RemoveAll(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error removing folder at destination %s: %v ", path, err)
		}
	}
	return os.Mkdir(path, os.ModePerm)
}

func getTemplate(path string) (t *template.Template, err error) {
	t, err = template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error reading template %s: %v", path, err)
	}
	return t, nil
}

func buildCanonicalLink(path, baseURL string) (link string) {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return fmt.Sprintf("%s/%s/index.html", baseURL, strings.Join(parts[2:], "/"))
	}
	return "/"
}

func getTagLink(tag string) (link string) {
	link = fmt.Sprintf("/tags/%s/", strings.ToLower(tag))
	return link
}

func getFolder(path string) (folder string) {
	return path[:strings.LastIndex(path, "/")]
}
