package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {

	data, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	htmlStr = string(data)

	/*
	 curl
	 -X [HTTP Request Method] (default:GET)
	 -H [HTTP Request Header]
	 -d [HTTP Request Body]
	 [HTML Request URL
	*/

	/*
		curl localhost:8080/
	*/
	http.HandleFunc("/", htmlHandler)
	/*
		curl localhost:8080/ssh-kex
	*/
	http.HandleFunc("/ssh-kex", kexHandler)
	/*
		curl localhost:8080/ssh-auth?funct=false
		curl localhost:8080/ssh-auth?funct=true
	*/
	http.HandleFunc("/ssh-auth", authHandler)
	/*
		curl -X POST -H "Content-Type: application/octec-stream" localhost:8080/user/new
		curl -X POST -H "Content-Type: application/octec-stream" localhost:8080/user/tada/command
		curl -X POST -H "Content-Type: application/octec-stream" localhost:8080/user/oda/command
		curl -X POST -H "Content-Type: application/octec-stream" localhost:8080/user/hogehoge/command
	*/
	http.HandleFunc("/user/", userHandler)

	// http://localhost:8080/ でアクセスできるサーバーを起動
	http.ListenAndServe(":8080", nil)
}

var htmlStr string

// index.htmlを表示するハンドラ
func htmlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, htmlStr)
}

// 鍵交換のハンドラ
func kexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Key exchange!")
}

// ユーザの公開鍵認証を行うハンドラ
func authHandler(w http.ResponseWriter, r *http.Request) {
	// パラメータfunctを取得
	var funct bool
	funct, err := strconv.ParseBool(r.URL.Query().Get("funct"))
	//  string -> bool の変換にエラーが出た場合
	if err != nil {
		fmt.Fprintln(w, "error:"+err.Error())
		return
	}

	fmt.Fprintln(w, fmt.Sprintf("%t", funct))
}

// パスを分離する関数
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type User struct {
	Auth    bool
	Pub_key int
}

var users = map[string]User{"tada": {true, 1}, "oda": {false, 2}}

// /user/.. 以降のパスによって、ハンドラを選択
func userHandler(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	head, r.URL.Path = ShiftPath(r.URL.Path)
	//fmt.Fprintln(w, head)

	switch head {
	case "new":
		newUserHandler(w, r)
	default:
		v, exist := users[head]
		if !exist {
			fmt.Fprintf(w, "User:"+head+" is not registered!")
			return
		}
		if !v.Auth {
			fmt.Fprintf(w, "User:"+head+" is not authorized!")
			return
		}
		authorizedUserHandler(w, r, head)
	}
}

// 新しいユーザの公開鍵を登録、または既存ユーザの公開鍵を更新するハンドラ
func newUserHandler(w http.ResponseWriter, r *http.Request) {
	// username, pub_key
	fmt.Fprintln(w, "Register new user!")
	return
}

// 認証済みのユーザとの共通鍵暗号通信を行うハンドラ
func authorizedUserHandler(w http.ResponseWriter, r *http.Request, username string) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	//fmt.Fprintln(w, head)

	switch head {
	case "command":
		fmt.Fprintln(w, "Execute command in authorized user:"+username+"!")
	default:
		http.NotFound(w, r)
	}

}
