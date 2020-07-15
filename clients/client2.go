package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"demo/models"
	"demo/utils"
)

const (
	GateWay = "127.0.0.1:9700"
	Token   = "a2ab1d0639d501aba7d40fce5a894f2b545ef79b"
	MateSn  = "123456"
)

var (
	privKeys    [3]*rsa.PrivateKey
	matePubKeys [3]string
	Messages    = []string{"哈哈", "这是第一个测试", "的消息"}
)

func main() {
	loadChatUserInfo()
	for i, _ := range privKeys {
		privKeys[i], _ = utils.GeneratePriKey(2048)
	}
	go func() {
		uploadPubKeys()
	}()

	go func() {
		subscribeChatUserPubKeys()
	}()

	go func() {
		for _, s := range Messages {
			sendMessage(s)
		}
	}()

	go func() {
		subscribeMessage()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

}

// 载入对方的信息，主要是公钥
func loadChatUserInfo() {
	data := url.Values{}
	data.Set("sn", MateSn)
	url := "http://" + GateWay + "/api/web/v1/user/info"
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", Token)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	var response struct {
		Head map[string]string `json:"head"`
		Body models.User       `json:"body"`
	}
	json.Unmarshal(b, &response)
	for _, key := range response.Body.PublicKeys {
		matePubKeys[key.Index-1] = key.Content
	}

}

// 上传我的公钥
func uploadPubKeys() {
	for i, key := range privKeys {
		b, _ := utils.GeneratePubKey(key)
		data := url.Values{}
		data.Set("index", strconv.Itoa(i+1))
		data.Set("content", string(b))
		data.Set("token", Token)
		url := "http://" + GateWay + "/api/web/v1/pub_key/upload"
		body := strings.NewReader(data.Encode())
		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", Token)
		http.DefaultClient.Do(req)
	}
}

// 发送消息
func sendMessage(content string) {
	encryptedContent, _ := utils.PublicKeyEncrypt([]byte(content), string(matePubKeys[0]))
	data := url.Values{}
	data.Set("level", "3")
	data.Set("content", base64.StdEncoding.EncodeToString(encryptedContent))
	data.Set("receiver_sn", MateSn)
	url := "http://" + GateWay + "/api/web/v1/message/upload"
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", Token)
	http.DefaultClient.Do(req)
}

// 订阅发给我的消息
func subscribeMessage() {
	u := url.URL{Scheme: "ws", Host: GateWay, Path: "/api/ws/v1/message/listen"}
	log.Println("connecting to ", u.String())
	header := http.Header{}
	header.Add("Authorization", Token)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Println("dial:", err)
		return
	}

	defer c.Close()
	errChan := make(chan error)
	go func() {
		for {
			_, m, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				errChan <- err
				return
			}
			var ms models.Message
			json.Unmarshal(m, &ms)
			decoded, _ := base64.StdEncoding.DecodeString(ms.Content)
			message, _ := utils.PrivateKeyDecrypt([]byte(decoded), privKeys[0])
			log.Println("message:", string(message))
		}
	}()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-errChan:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
	return
}

// 订阅对方的公钥更新
func subscribeChatUserPubKeys() {
	u := url.URL{Scheme: "ws", Host: GateWay, Path: "/api/ws/v1/public_key/listen"}
	log.Println("connecting to ", u.String())
	header := http.Header{}
	header.Add("Authorization", Token)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Println("dial:", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, []byte("{\"user_sns\":[\""+MateSn+"\"],\"id\":\"oldfritter\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}
	defer c.Close()
	errChan := make(chan error)
	go func() {
		for {
			_, m, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				errChan <- err
				return
			}
			var publicKey models.PublicKey
			json.Unmarshal(m, &publicKey)
			matePubKeys[publicKey.Index-1] = publicKey.Content
		}
	}()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-errChan:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
	return
}
