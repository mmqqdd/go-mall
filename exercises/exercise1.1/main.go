package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /ping", PingPong)
	http.HandleFunc("GET /hello", Hello)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("服务启动失败")
	}
}

// PingPong 固定返回pong
func PingPong(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "pong")
}

// Hello 打招呼
func Hello(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	language := req.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}

	greetingTemplate := map[string]string{
		"en": "Hello,%s!",
		"zh": "你好，%s！",
		"ja": "こんにちは、%sさん！",
	}

	greeting := fmt.Sprintf(greetingTemplate[language], name)

	_, err := io.WriteString(w, greeting)
	if err != nil {
		log.Println("Hello handler WriteString err")
		return
	}
}
