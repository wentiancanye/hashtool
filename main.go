package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var PATH = "./"

type Data struct {
	File []string `json:"file"`
	Data []string `json:"data"`
}

func GetFileHash(path string) (fileMD5 string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return fileMD5, err
	}
	defer f.Close()

	md5hash := sha1.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return fileMD5, err
	}
	fileMD5 = hex.EncodeToString(md5hash.Sum(nil))
	return fileMD5, nil
}

func GetFileData(folder string) (Data, error) {
	var fs []string
	var ds []string
	var itself = filepath.Base(os.Args[0])
	filepath.Walk(folder, func(f string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Println(err.Error())
			return err
		}
		fname := fi.Name()
		if !(fi.IsDir() || fname == itself || fname == "data" || fname == "version" || fname == "plugin") {
			h, _ := GetFileHash(f)
			fs = append(fs, f)
			ds = append(ds, h)
		}
		return nil
	})
	var fh = Data{File: fs, Data: ds}
	return fh, nil
}

func main() {
	if len(os.Args) >= 2 {
		fmt.Println("使用说明：递归获取当前文件夹下所有文件的SHA1值，保存为data文件，忽略自身、data、version")
	} else {
		fp, _ := os.Create("data")
		defer fp.Close()
		fh, _ := GetFileData(PATH)
		data, _ := json.MarshalIndent(fh, "", "    ")
		fp.Write(data)
	}
}
