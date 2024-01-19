package handler

import (
	"encoding/json"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ToJsonString(v interface{}, pretty bool) string {
	var buf []byte
	if pretty {
		buf, _ = json.MarshalIndent(v, "", "\t")
	} else {
		buf, _ = json.Marshal(v)
	}
	return string(buf)
}

func HashString(data string) string {
	return fmt.Sprintf("%x", xxhash.Sum64String(data))
}

func AppRoot() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

func CheckNextUsefulPort(port int) int {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return CheckNextUsefulPort(port + 1)
	}
	defer func() { _ = listener.Close() }()
	return port
}

func OpenUrl(openUrl string) {
	var osName = strings.ToLower(runtime.GOOS)
	switch osName {
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", openUrl)
		_ = cmd.Run()
	case "darwin":
		cmd := exec.Command("open", openUrl)
		_ = cmd.Run()
	}
}

//func ParseHtmlTemplate(htmlTpl string, args gin.H) (htmlString string, err error) {
//	//buf, err := os.ReadFile(htmlTpl)
//	//if err != nil {
//	//	return htmlString, err
//	//}
//	//var htmlContent = string(buf)
//	//for key, value := range args {
//	//	htmlContent = strings.ReplaceAll(htmlContent, key, value.(string))
//	//}
//	//return htmlContent, nil
//
//	title := r.URL.Path[len("/edit/"):]
//	p, err := loadPage(title)
//	if err != nil {
//		p = &Page{Title: title}
//	}
//	t, _ := template.ParseFiles(htmlTpl)
//	t.Execute(w, p)
//
//}
