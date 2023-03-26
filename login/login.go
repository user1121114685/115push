package login

import (
	"115push/cookie"
	"encoding/json"
	"github.com/deadblue/elevengo"
	"github.com/deadblue/elevengo/option"
	"io"
	"log"
	"os"
	"time"
)

var Agent = elevengo.New(option.CooldownOption{Min: 100, Max: 2000})

func Login() {

	file, err := os.Open("cookies.txt")
	if err != nil {
		_, err = os.Create("cookies.txt")
		if err != nil {
			log.Println("创建cookies txt文件失败", err)
			return
		}
		file, _ = os.Open("cookies.txt")
	}
	//及时关闭file句柄
	defer file.Close()
	_c, err := io.ReadAll(file)
	log.Println(string(_c))
	var _cookie cookie.Cookie
	err = json.Unmarshal(_c, &_cookie)
	if err != nil {
		log.Println("格式化cookie出错", err)
		log.Println("请使用EditThisCookie格式的cookie并填入cookies.txt中", err)

		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	var (
		uid  string
		cid  string
		seid string
	)
	for _, v := range _cookie {
		if v.Name == "CID" {
			cid = v.Value
		}
		if v.Name == "UID" {
			uid = v.Value
		}
		if v.Name == "SEID" {
			seid = v.Value
		}
	}
	if err := Agent.CredentialImport(&elevengo.Credential{
		UID:  uid,
		CID:  cid,
		SEID: seid,
	}); err != nil {

		log.Println("Cookie登录错误:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

}
