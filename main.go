package main

import (
	"115push/receiver"
	"115push/server"
	"flag"
	"fmt"
	"log"
)

var runServer bool
var runDebug bool
var Version = "调试版本"

func init() {
	flag.BoolVar(&runServer, "s", false, "一键运行服务端，方便从linux等设备启动运行。")
	flag.BoolVar(&runDebug, "debug", false, "显示更多log信息")
}
func main() {

	log.SetFlags(log.Lmicroseconds | log.Ldate)

	flag.Parse()
	if runDebug {
		log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	}
	if runServer {
		server.Run115PushServer()
		return
	}
	for {
		fmt.Println("当前版本 V" + Version)
		fmt.Println()
		fmt.Println("欢迎使用云传(115push)，请勿用于任何商业用途。")
		fmt.Println("你正在使用社区版，需要功能更强的定制版请联系admin@shaoxia.xyz")
		fmt.Println("社区源码(并非开源不允许白嫖)：https://github.com/user1121114685/115push")
		fmt.Println()
		fmt.Println("请输入编号选择功能：")
		fmt.Println()
		fmt.Println("1 启动服务端，其他人将通过您提供的链接访问你的115资源")
		fmt.Println()
		fmt.Println("2 启动导入，通过输入他人提供的服务端地址，导入115资源")
		var num int
		fmt.Scanln(&num)
		switch num {
		case 1:
			server.Run115PushServer()
			break
		case 2:

			var cid string
			fmt.Println("请输入自己目录的CID")
			fmt.Scanln(&cid)
			var url string
			fmt.Println("请输入分享者url,例如： http://127.0.0.1:1150")
			fmt.Scanln(&url)
			var shareCid string
			fmt.Println("请输入分享者分享的CID,不是自己的CID,是分享者的CID!!!")
			fmt.Scanln(&shareCid)
			if cid == "" || url == "" || shareCid == "" {
				continue
			}
			receiver.Import(cid, url, shareCid)
			break

		default:
			fmt.Println("输入有误，请重新输入....")
		}
	}

}
