package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/echo", Echo)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

const maxBodySize = 1 << 10 // 1KB

// Echo 将请求体返回
func Echo(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 限制请求大小不能超过1k
	req.Body = http.MaxBytesReader(w, req.Body, maxBodySize)

	defer func() {
		req.Body.Close()
	}()

	w.Header().Set("Content-Type", "text/plain")

	bodyContent, err := io.ReadAll(io.LimitReader(req.Body, maxBodySize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(bodyContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}
