package html

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func EncodeImg(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(`<img[^<>]+src\s*=\s*"\s*(http[^"]+)"`)
	if err != nil {
		return err
	}
	ret := re.ReplaceAllFunc(data, func(s []byte) []byte {
		matches := re.FindSubmatch(s)
		return []byte(fmt.Sprintf(`<img src="%s"`, loadPicData(string(matches[1]))))
	})
	return ioutil.WriteFile(file, ret, os.ModePerm)
}

func loadPicData(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return url
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return url
	}
	ext := strings.ToLower(filepath.Ext(url))
	if ext == ".png" || ext == ".jpg" || ext == ".gif" || ext == ".jpeg" {
		ext = ext[1:]
	} else {
		ext = "jpg"
	}

	return fmt.Sprintf("data:image/%s;base64,%s", ext, base64.StdEncoding.EncodeToString(body))
}
