package utils

import "github.com/deadblue/elevengo"

type FileList struct {
	FirstDirName string      `json:"FirstDirName"`
	Files        []*SendFile `json:"Files"`
}

type SendFile struct {
	ImportTicket elevengo.ImportTicket `json:"ImportTicket"`
	CID          string                `json:"CID"`
	PickCode     string                `json:"pickCode"`
	IsDir        bool                  `json:"IsDir" default:"false"`
}
