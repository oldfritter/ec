package utils

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetRealIp(context echo.Context) (ip string) {
	ips := context.RealIP()
	// some forwarded_ips like 223.104.64.228,183.240.52.39, 172.68.254.56, sduppid operators
	rips := strings.Split(ips, ",")
	ip = rips[0]
	return
}

func GetLogFile(name ...string) *os.File {
	var app string
	if len(name) == 0 {
		exe := strings.Split(os.Args[0], "/")
		app = exe[len(exe)-1]
	} else {
		app = name[0]
	}
	if err := os.Mkdir("logs", 0755); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	file, err := os.OpenFile("logs/"+app+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	return file
}

func SetLog() {
	file := GetLogFile()
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func SetPid() {
	if err := os.Mkdir("pids", 0755); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("create folder error: %v", err)
		}
	}
	exe := strings.Split(os.Args[0], "/")
	app := exe[len(exe)-1]
	err := ioutil.WriteFile("pids/"+app+".pid", []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func UnsetPid() {
	exe := strings.Split(os.Args[0], "/")
	app := exe[len(exe)-1]
	err := ioutil.WriteFile("pids/"+app+".pid", []byte{}, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func SetLogAndPid() {
	SetLog()
	SetPid()
}

func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
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
