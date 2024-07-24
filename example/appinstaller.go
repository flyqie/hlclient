package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hlcient/hlclient"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("please give file path!")
	}

	hl := hlclient.NewClient("127.0.0.1:32876")

	err := hl.Connect()
	if err != nil {
		log.Fatalf("connect err: %v", err)
	}

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("open file err: %v", err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("get file info err: %v", err)
	}
	if fileInfo.Size() <= 0 {
		log.Fatalf("file size is zero")
	}

	hash := md5.New()
	hashBuf := make([]byte, 10*1024*1024)
	for {
		n, err := file.Read(hashBuf)
		if err != nil && err != io.EOF {
			log.Fatalf("hash read file err: %v", err)
		}
		if n == 0 {
			break
		}
		if _, err := hash.Write(hashBuf[:n]); err != nil {
			log.Fatalf("hash write file err: %v", err)
		}
	}
	hashBuf = nil
	md5Sum := fmt.Sprintf("%x", hash.Sum(nil))
	hash.Reset()
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatalf("hash seek file err: %v", err)
	}

	msgData, err := json.Marshal(map[string]interface{}{
		"index":    1,
		"size":     fileInfo.Size(),
		"md5":      md5Sum,
		"name":     "appinstaller",
		"maxCount": 114514,
	})
	if err != nil {
		log.Fatalf("marshal err: %v", err)
	}

	resp, err := hl.SendMessage('\\', msgData, true)
	if err != nil {
		log.Fatalf("send err: %v", err)
	} else {
		log.Printf("resp: %v", resp)
		log.Printf("resp.code: %d", resp.Code)
		log.Printf("resp.data: %s", string(resp.Data))
	}
	if resp.Code != 0 {
		log.Fatalf("resp.Code Error")
	}

	log.Printf("start send appdata")
	chunkBuf := make([]byte, 10*1024*1024)
	for {
		n, err := file.Read(chunkBuf)
		if err != nil && err != io.EOF {
			log.Fatalf("read chunk err: %v", err)
		}
		if n == 0 && err == io.EOF {
			break
		}

		_, err = hl.SendData(chunkBuf[:n])
		if err != nil {
			log.Fatalf("send chunk err: %v", err)
		}
	}
	chunkBuf = nil
	_, err = hl.SendData([]byte{})
	if err != nil {
		log.Fatalf("finish send chunk err: %v", err)
	}
	log.Printf("finish send appdata, wait install resp...")
	resp2, err := hl.RecvMessage()
	if err != nil {
		log.Fatalf("recv install_resp err: %v", err)
	}
	log.Printf("install_resp: %v", resp2)
	log.Printf("install_resp.code: %d", resp2.Code)
	log.Printf("install_resp.data: %s", string(resp2.Data))
	if resp2.Code != 0 {
		log.Fatalf("install_resp.Code Error")
	}
	log.Printf("install success!")
}
