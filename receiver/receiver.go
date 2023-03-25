package receiver

import (
	"115push/login"
	"115push/utils"
	"encoding/json"
	"github.com/deadblue/elevengo"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Import(dirid, url, shareCid string) {
	login.Login()
	var tickets utils.FileList
	// 去掉最右边的 /
	url = strings.TrimRight(url, "/")
	getFileList(url+"/115/file_get?cid="+shareCid, &tickets)
	if &tickets == nil {
		log.Println("没有成功获取到文件列表.....")
		time.Sleep(3 * time.Second)
		return
	}
	var firstDirID string
	var err error
	for i := 0; i < 20; i++ {
		firstDirID, err = login.Agent.DirMake(dirid, tickets.FirstDirName)
		if err != nil && err.Error() == "target already exists" {
			// target already exists
			//log.Println("创建 第一级目录失败", err)
			tickets.FirstDirName = tickets.FirstDirName + "_" + strconv.Itoa(i)
			continue
		}
		break
	}
	wg := NewImporter(5)
	importFileForDir(firstDirID, url, &tickets, wg)
	wg.producerWaitGroupPool.Wait()
	close(wg.taskChannel)

}

func importFileForDir(dirid, url string, tickets *utils.FileList, wg *importer) {
	// 延迟执行一个匿名函数
	defer func() {
		// 调用 recover 来捕获 panic 的值
		if err := recover(); err != nil {
			// 打印错误信息并继续执行
			log.Println("发生致命错误！！！ 请将错误反馈给开发者，即将开始重新导入....")
			log.Println(err)
			time.Sleep(5 * time.Second)
			importFileForDir(dirid, url, tickets, wg)
		}
	}()
	wg.producerWaitGroupPool.Add()
	for _, ticket := range tickets.Files {

		if ticket.IsDir {
			var id string
			var err error
			for i := 0; i < 20; i++ {
				// 已经导入的部分 直接赋值
				if *ticket.MakeDIrCid != "" {
					id = *ticket.MakeDIrCid
					break
				}
				id, err = login.Agent.DirMake(dirid, ticket.ImportTicket.FileName)
				if err != nil && err.Error() == "target already exists" {
					// target already exists
					//log.Println("创建 第一级目录失败", err)
					ticket.ImportTicket.FileName = ticket.ImportTicket.FileName + "_" + strconv.Itoa(i)
					continue
				}
				// 没有新建文件夹的ID 新建一次文件夹并赋值
				ticket.MakeDIrCid = &id
				break
			}

			var fileList utils.FileList

			getFileList(url+"/115/file_get?cid="+ticket.CID, &fileList)
			if &fileList == nil {
				continue
			}
			go importFileForDir(id, url, &fileList, wg)
			continue
		}
		// 如果已经导入 则跳过该文件
		if *ticket.IsImport {
			continue
		}

		// 有时候发生了致命错误，然而不知道罪魁祸首，加上这个一目了然
		log.Println("准备导入   " + ticket.ImportTicket.FileName)
		err := login.Agent.Import(dirid, &ticket.ImportTicket)
		*ticket.IsImport = true
		if err != nil {
			if ie, ok := err.(*elevengo.ErrImportNeedCheck); ok {
				signValue := getCalculateSignValue(url, ticket.PickCode, ie.SignRange)
				if signValue != "" {
					ticket.ImportTicket.SignKey = ie.SignKey
					ticket.ImportTicket.SignValue = signValue
					err = login.Agent.Import(dirid, &ticket.ImportTicket)
					if err != nil {
						// 失败重新导入增加那么一丝丝可能性 能降低失败机率
						reImportFileForDir(dirid, url, ticket.PickCode, ticket.ImportTicket)
						continue
					}
				}
			}
		}
		log.Println("导入成功   " + ticket.ImportTicket.FileName)

	}
	wg.producerWaitGroupPool.Done()
}

func reImportFileForDir(dirid, url, pickCode string, ImportTicket elevengo.ImportTicket) {

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		iport_ticket := ImportTicket
		err := login.Agent.Import(dirid, &iport_ticket)
		if err != nil {
			if ie, ok := err.(*elevengo.ErrImportNeedCheck); ok {
				signValue := getCalculateSignValue(url, pickCode, ie.SignRange)
				if signValue != "" {
					iport_ticket.SignKey = ie.SignKey
					iport_ticket.SignValue = signValue
					err = login.Agent.Import(dirid, &iport_ticket)
					if err != nil {
						log.Println(iport_ticket.FileName, "  失败重试中....")
						log.Println(iport_ticket, ie)
						log.Println(err)

						continue
					}
					log.Println("导入成功   " + iport_ticket.FileName)
					return
				}
			}
		}

	}
	log.Println("导入失败   " + ImportTicket.FileName)
}

func getFileList(url string, tickets *utils.FileList) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("读取body出错", err)
	}
	err = json.Unmarshal(body, &tickets)
	if err != nil {
		log.Println(err)
		log.Println(string(body))
		return
	}
}

func getCalculateSignValue(postUrl, pickcode, signRange string) string {
	formValues := url.Values{}
	formValues.Set("pickcode", pickcode)
	formValues.Set("signRange", signRange)
	signValue, err := http.PostForm(postUrl+"/115/calculate_sign_value", formValues)
	if err != nil {
		return ""
	}
	sig, err := io.ReadAll(signValue.Body)
	if err != nil {
		return ""
	}
	return string(sig)
}

type importer struct {
	taskChannel chan int
	// 通过pool支持设置上限
	producerWaitGroupPool *utils.WaitGroupPool
}

func NewImporter(workNum int) *importer {
	return &importer{
		taskChannel:           make(chan int, 300),
		producerWaitGroupPool: utils.NewWaitGroupPool(workNum),
	}
}
