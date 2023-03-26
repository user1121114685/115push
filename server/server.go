package server

import (
	"115push/login"
	"115push/utils"
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo"
	"log"
	"net/http"
)

func Run115PushServer() {
	login.Login()
	startServer()
}

func startServer() {
	log.Println("服务端开始运行。。。。。")
	fmt.Println("您需要将服务器地址 http://ip:1150 告知用户")
	fmt.Println("需要分享的目录CID 告知需要导入的用户")
	fmt.Println("请保证用户能访问到 http://ip:1150  包括不限于DDNS 内网穿透，服务器运行")
	srv := &http.Server{Addr: ":1150"}
	http.HandleFunc("/", hello)
	http.HandleFunc("/115/file_get", fileGet)
	http.HandleFunc("/115/calculate_sign_value", CalculateSignValue)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}

func fileGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	url := r.URL.Query()

	ticket := &utils.FileList{}
	cid := url.Get("cid")
	if cid == "" || cid == "0" {
		fmt.Fprintf(w, "没有Cid无法导入")
		return
	}
	rangeFileList(cid, ticket)
	marshal, err := json.Marshal(ticket)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	log.Println("用户获取目录 ", ticket.FirstDirName)
	fmt.Fprintln(w, string(marshal))

}

func rangeFileList(cid string, ticket *utils.FileList) {
	firstDir := elevengo.File{}
	err := login.Agent.FileGet(cid, &firstDir)
	if err != nil {
		log.Println("获取第一级目录信息错误", err, cid)
		return
	}
	ticket.FirstDirName = firstDir.Name
	it, err := login.Agent.FileIterate(cid)
	for ; err == nil; err = it.Next() {
		file := elevengo.File{}
		_ = it.Get(&file)

		if file.Name == "" {
			continue
		}
		t := elevengo.ImportTicket{
			FileName: file.Name,
			FileSize: file.Size,
			FileSha1: file.Sha1,
		}
		ticket.Files = append(ticket.Files, &utils.SendFile{
			ImportTicket: t,
			CID:          file.FileId,
			PickCode:     file.PickCode,
			IsDir:        file.IsDirectory,
		})
		//log.Printf("File: %d => %#v", it.Index(), file)
	}
	if !elevengo.IsIteratorEnd(err) {
		log.Fatalf("Iterate file failed: %s", err.Error())
	}
}

func CalculateSignValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//log.Println("当前请求模式：" + r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()

	signValue, err := login.Agent.ImportCalculateSignValue(r.FormValue("pickcode"), r.FormValue("signRange"))
	if err != nil {
		log.Println("获取二次随机Sha1失败", err)
		fmt.Fprintf(w, "invalid")
		return
	}
	fmt.Fprintf(w, signValue)

}

func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "恭喜，连通性正确!!，请将url填入115云传中。")

}
