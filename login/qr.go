package login

import (
	cookie2 "115push/cookie"
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo"
	"io"
	"log"
	"net/http"
	"os"
)

func QrCodeLogin() {

	session := &elevengo.QrcodeSession{}

	err := Agent.QrcodeStartForLinux(session)
	if err != nil {
		log.Fatalf("获取Qr地址错误: %s", err)
	}
	qrName := "扫码登录.png"
	saveQr(session.Content, qrName)
	qrStatus(session)
	err = os.Remove(qrName)
}

func qrStatus(session *elevengo.QrcodeSession) {

	for i := 0; i < 100; i++ {
		var status elevengo.QrcodeStatus
		// Get QR-Code status
		status, err := Agent.QrcodeStatus(session)
		if err != nil {
			log.Printf("Get QRCode status error: %s", err)
			if err.Error() == "qrcode expired" {
				break
			}
		} else {
			// Check QRCode status
			if status.IsWaiting() {
				log.Println("扫描目录下的 扫码登录.png ,登录Linux客户端！")
			} else if status.IsScanned() {
				log.Println("已扫描二维码，还未确认登录.....")
			} else if status.IsAllowed() {
				err = Agent.QrcodeLogin(session)
				if err == nil {
					log.Println("登录成功!")
					// 导出cookie
					cookie := elevengo.Credential{}
					Agent.CredentialExport(&cookie)
					log.Println(cookie.CID, cookie.UID, cookie.SEID)
					saveQrCookies(cookie)
				} else {
					log.Printf("Submit QRcode login error: %s", err)
				}
				break
			} else if status.IsCanceled() {
				log.Println("你取消了登录....")
				break
			}
		}

	}
}

func saveQr(url, qrName string) {

	if _, err := os.Stat(qrName); err == nil {
		err = os.Remove(qrName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	file, err := os.Create(qrName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err)
	}
}

//[{"domain":".115.com","hostOnly":false,"httpOnly":true,"name":"CID","path":"/","sameSite":"unspecified","secure":false,"session":true,"storeId":"0","value":"f56c679dae92ce4c938b294cbc561a52","id":1}, {"domain":".115.com","hostOnly":false,"httpOnly":true,"name":"UID","path":"/","sameSite":"unspecified","secure":false,"session":true,"storeId":"0","value":"336430153_P3_1678198361","id":2}, {"domain":".115.com","hostOnly":false,"httpOnly":true,"name":"SEID","path":"/","sameSite":"unspecified","secure":false,"session":true,"storeId":"0","value":"3cb800fa88979a14e88b9ac3bb27749ca5d58fcde9e70a0be80819018f8c8e2e9aa1d0ca888a62646f926c02bfed51066514bf2215ca58ded72b62bb","id":3}]

func saveQrCookies(cookie elevengo.Credential) {
	var c = &cookie2.Cookie{
		Domain:   ".115.com",
		HostOnly: false,
		HTTPOnly: true,
		Name:     "",
		Path:     "/",
		SameSite: "unspecified",
		Secure:   false,
		Session:  true,
		StoreID:  "0",
		Value:    "",
		ID:       0,
	}
	var savecooke cookie2.Cookies
	c.Name = "CID"
	c.Value = cookie.CID
	savecooke = append(savecooke, *c)
	c.Name = "UID"
	c.Value = cookie.UID
	c.ID++
	savecooke = append(savecooke, *c)
	c.Name = "SEID"
	c.Value = cookie.SEID
	c.ID++
	savecooke = append(savecooke, *c)
	newCookie, err := json.Marshal(savecooke)
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("cookies.txt", newCookie, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
