package main

import (
	"fmt"
	"os"

	"dev.farukh/copy-close/http"
)

const filePermissionUnix = 0755

func main() {
	createFilesFolder()
	http.Init()
}

func createFilesFolder() {
	if path, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		path = fmt.Sprintf("%s/files/", path)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return
		}

		if err := os.Mkdir("files", filePermissionUnix); err != nil { 
			panic(err.Error())
		}
	}
}
