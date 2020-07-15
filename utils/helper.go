package utils

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
)

func SetLogAndPid(name string) {
	err := os.Mkdir("logs", 0755)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	file, err := os.OpenFile("logs/"+name+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	log.SetOutput(file)
	err = os.Mkdir("pids", 0755)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	err = ioutil.WriteFile("pids/"+name+".pid", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(62))
		b[i] = letterRunes[index.Int64()]
	}
	return string(b)
}
