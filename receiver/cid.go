package receiver

import (
	"os"
)

var Lastchoice string

func FindLastCid() bool {
	// read the whole file as a byte slice
	data, err := os.ReadFile("./last_cid.txt")
	if err != nil {

		return false
	}
	// convert the byte slice to a string
	if string(data) == "" {
		return false
	}
	Lastchoice = string(data)
	return true
}
