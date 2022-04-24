package Plug

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
	WXAPI(message,"")
	WXopenai.TOKEN
	WXopenai.ENV
*/
func WXAPI(message *user.Message) string {

	type (
		WxSession struct {
			Signature string `json:"signature"`
			ExpiresIn int    `json:"expiresIn"`
			Rid       string `json:"rid"`
			Errcode   int    `json:"errcode"`
			Errmsg    string `json:"errmsg"`
		}
		Answer struct {
			AnsNodeName string `json:"ans_node_name"`
			Answer      string `json:"answer"`
			//Query       string  `json:"query"`
			Errcode int    `json:"errcode"`
			Errmsg  string `json:"errmsg"`
		}
	)
	var (
		resp      *http.Response
		body      []byte
		wxSession = WxSession{}
		answer    = Answer{}
	)
	// 微信鉴权
	if resp, err = http.Post(
		viper.GetString("WXopenai.signUrl")+
			viper.GetString("WXopenai.TOKEN"), "application/json",
		strings.NewReader(
			fmt.Sprintf(`{"username":"%v","userid": "%v"}`, message.From().Name(), message.From().ID()),
		)); err != nil {
		log.Errorf("[wx] 请求 signature 接口失败! Error: [%v]", err)
		return ""
	}
	// 关闭鉴权请求
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf(err.Error())
		}
	}(resp.Body)

	// 读取鉴权信息
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("[wx] 读取 signature 信息失败! Error: [%v]", err)
		return ""
	}
	// 解析鉴权信息
	if err = json.Unmarshal(body, &wxSession); wxSession.ExpiresIn == 0 {
		log.Printf("[wx] 解析 signature 信息失败! Error: %+v", wxSession.Errmsg)
		return ""
	}
	log.Printf("[wx] 解析 signature 信息成功!")
	if resp, err = http.Post(viper.GetString("WXopenai.url")+
		viper.GetString("WXopenai.TOKEN"),
		"application/json", strings.NewReader(
			fmt.Sprintf(`{"signature": "%v", "query": "%v","env": "%v"}`, wxSession.Signature, message.MentionText(), viper.GetString("WXopenai.ENV")))); err != nil {
		log.Errorf("[wx] 请求 aibot 接口失败! Error: [%v]:", err)
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf(err.Error())
		}
	}(resp.Body)

	// 读取信息
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("[wx] 读取 aibot 信息失败! Error: [%v]:", err)
		return ""
	}
	// 解析信息
	if err = json.Unmarshal(body, &answer); answer.Errcode != 0 {
		log.Errorf("[wx] 解析 aibot 信息失败! Error: %v", answer.Errmsg)
		return ""
	}
	log.Printf("[wx] 解析 aibot 信息成功!")
	if answer.Answer == "" {
		log.Println("[wx] 机器人 回复信息为空")
		return ""
	}
	//log.Printf("[wx] AnsNodeName: [%v], Answer: [%v], Errcode: [%v], Errmsg: [%v]", answer.AnsNodeName, answer.Answer, answer.Errcode, answer.Errmsg)
	return answer.Answer
}
