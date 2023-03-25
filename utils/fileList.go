package utils

import "github.com/deadblue/elevengo"

type FileList struct {
	FirstDirName string      `json:"FirstDirName"`
	Files        []*SendFile `json:"Files"`
}

type SendFile struct {
	ImportTicket elevengo.ImportTicket `json:"ImportTicket"`

	CID      string `json:"CID"`
	PickCode string `json:"pickCode"`
	IsDir    bool   `json:"IsDir" default:"false"`
	// MakeDIrCid  如果是文件夹，则记录下文件夹的CID ---用于断点续传
	MakeDIrCid *string `json:"MakeDIrCid"`
	// IsImport 如果是文件，则记录下文件是否已经导入了。  --- 用于断点续传
	IsImport *bool `json:"IsImport" default:"false"`
}
