package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	println("Content-Encoding: ", r.Header.Get("Content-Encoding"))
	println("Content-Type: ", r.Header.Get("Content-Type"))

	mr, _ := r.MultipartReader()
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		fmt.Println(p.FormName(), ":", p.FileName(), p.Header.Get("Content-Type"))

		buf, err := ioutil.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}

		unzip(buf)
	}
}

func unzip(file []byte) {
	// zipファイルを読み込む
	r, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		log.Fatal(err)
	}

	// 各ファイルを展開する
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer rc.Close()

		// ファイル名を出力する
		fmt.Println("fileName: ", f.Name)

		// ファイルを展開する
		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}

		// ファイルの内容を出力する
		fmt.Println("fileContent: ", string(buf))
	}
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start http listening :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
