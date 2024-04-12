package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
)

/*
 * Blog types
 */
const (
	Poka = iota
	Carro
	Grotty
)

/* 
 * This function is used to parse the URL of the blog.
 */ 
func ParseURL (url string) (int, error) {
	re := regexp.MustCompile(`https://(.+?)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return -1, errors.New("invalid URL")
	}

	switch matches[1] {
	case "www.po-kaki-to.com":
		return Poka, nil
	case "carro-groce.com":
		return Carro, nil
	case "grotty-monday.com":
		return Grotty, nil
	default:
		return -1, errors.New("invalid URL")
	}
}

/*
 * This function is used to fetch the HTML of the blog.
 */
func FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

/*
 * This function is used to open the browser.
 */
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		panic(err)
	}
}