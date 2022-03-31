package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/monnand/dhkx"
)

func main() {

	fmt.Println("Open http://localhost:8080/")

	// 公開するディレクトリを指定
	fs := http.FileServer(http.Dir("front"))
	http.Handle("/", fs)

	/*
	 curl
	 -X [HTTP Request Method] (default:GET)
	 -H [HTTP Request Header]
	 -d [HTTP Request Body]
	 [HTML Request URL
	*/

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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type KexRequest struct {
	UserName string
	KexAlgo  string
	PubKey   string
}

type KexResponse struct {
	PubKeyKex  string //`json:"pub_key_for_kex"`
	CryptoAlgo string //`json:"pub_key_crypto_algo"`
	PubKey     string //`json:"pub_key_host"`
	ShareKey   string //`json:"share_key"`
	SessionID  string //`json:"session_ID"`
}

type User struct {
	Auth      bool // 公開鍵認証済みか？
	PubKey    int  // 公開鍵（ユーザ鍵）
	SessionID string
}

var users = map[string]User{"tada": {true, 1, "aaa"}, "oda": {false, 2, "bbb"}}
var host_pub_key = ""

//var host_pvt_key = ""

// 鍵交換のハンドラ
func kexHandler(w http.ResponseWriter, r *http.Request) {
	var req = &KexRequest{}

	// HTTP Request Body(JSON形式)を、Memo構造体(m)にセット
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		fmt.Fprintln(w, "error:"+err.Error())
		return
	}
	// ユーザ名
	user := req.UserName
	// DHGroupの指定
	var dhGroup int
	if req.KexAlgo == "diffie-hellman-group1-sha1" {
		dhGroup = 2
	} else {
		dhGroup = 14
	}
	// クライアントからの公開鍵
	reqPub, err0 := hex.DecodeString(req.PubKey)
	if err0 != nil {
		log.Fatal(err0)
	}
	//fmt.Println(reqPub)
	aliceKey := dhkx.NewPublicKey(reqPub)

	// Get a group. Use the default one would be enough.
	g, err1 := dhkx.GetGroup(dhGroup)
	if err1 != nil {
		log.Fatal(err1)
	}
	bob, err2 := g.GeneratePrivateKey(nil)
	if err2 != nil {
		log.Fatal(err2)
	}
	bobKey := bob.Bytes()

	share, err3 := g.ComputeKey(aliceKey, bob)
	if err3 != nil {
		log.Fatal(err3)
	}
	shareKey := share.Bytes()

	h := sha1.New()
	h.Write(shareKey)
	bs := h.Sum(nil)

	data1, err4 := os.ReadFile("ssh-server/id_rsa.pub")
	if err4 != nil {
		log.Fatal(err4)
	}
	host_pub_key = strings.Fields(string(data1))[1]

	var res = &KexResponse{}
	res.PubKeyKex = fmt.Sprintf("%X", bobKey)
	res.CryptoAlgo = "ssh-rsa"
	res.PubKey = fmt.Sprintf("%X", host_pub_key)
	res.ShareKey = fmt.Sprintf("%X", shareKey)
	res.SessionID = fmt.Sprintf("%X", bs)
	// ユーザに対するセッションIDの登録
	var v, _ = users[user]
	v.SessionID = res.SessionID
	users[user] = v

	data2, err5 := json.Marshal(res) // JSON形式に
	if err5 != nil {
		fmt.Fprintln(w, "error:"+err5.Error())
	}

	fmt.Fprintln(w, string(data2))

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

	fmt.Fprintf(w, "%t", funct)
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

// /user/.. 以降のパスによって、ハンドラを選択
func userHandler(w http.ResponseWriter, r *http.Request) {
	var head string
	_, r.URL.Path = ShiftPath(r.URL.Path)
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
