package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	http.HandleFunc("/books/", Books)
	http.HandleFunc("/books", Books)
	http.ListenAndServe(":8080", nil)
}

// Books 书籍接口，根据method转发到具体的服务
func Books(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		GetBooks(w, req)
	case http.MethodPost:
		PostBooks(w, req)
	}
}

// GetBooks 通过ID获取书籍信息
func GetBooks(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	// 如果没有指定id，返回全部书籍
	if url == "/books" {
		books := GetAllBooks()
		w.Header().Set("Content-Type", "application/json")
		booksJson, err := json.Marshal(books)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(booksJson)
		return
	}

	// 解析url路径中的书籍id
	bookIDStr := strings.TrimPrefix(url, "/books/")

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 从书籍库中查找书籍
	book, ok := BookMap[bookID]
	if !ok {
		http.Error(w, "book is not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bookJson, err := json.Marshal(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bookJson)
}

// PostBooks 新增一本书籍
func PostBooks(w http.ResponseWriter, req *http.Request) {
	defer func() {
		req.Body.Close()
	}()

	reqBody, err := io.ReadAll(io.LimitReader(req.Body, RequestMaxSize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := Book{}
	err = json.Unmarshal(reqBody, &book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	BookSetMutex.Lock()
	book.ID = GetNewBookID()
	BookMap[book.ID] = book
	BookSetMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// RequestMaxSize 请求大小限制
const RequestMaxSize = 1 << 10 // 1KB

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// BookSetMutex 修改书库锁，防止并发写
var BookSetMutex sync.Mutex

// BookMap 整个书库，保存所有书籍
var BookMap = map[int]Book{
	1: {
		ID:     1,
		Title:  "狂人日记",
		Author: "鲁迅",
	},
	2: {
		ID:     2,
		Title:  "骆驼祥子",
		Author: "老舍",
	},
	3: {
		ID:     3,
		Title:  "老人与海",
		Author: "海明威",
	},
}

// GetNewBookID 获取空闲的id，取最小一个空闲的id
func GetNewBookID() int {
	var ids []int
	for id := range BookMap {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	for i, _ := range ids {
		if i+1 < len(ids) && ids[i]+1 < ids[i+1] {
			return ids[i] + 1
		}
	}

	return ids[len(ids)-1] + 1
}

// GetAllBooks 获取所有的书籍，按照id升序排序
func GetAllBooks() []Book {
	books := []Book{}
	for _, book := range BookMap {
		books = append(books, book)
	}
	sort.Slice(books, func(i, j int) bool {
		return books[i].ID < books[j].ID
	})
	return books
}
