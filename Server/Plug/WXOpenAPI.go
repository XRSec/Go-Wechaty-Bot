package Plug

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
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
			AnsNodeName string  `json:"ans_node_name"`
			Answer      string  `json:"answer"`
			Confidence  float64 `json:"confidence"`
			Errcode     int     `json:"errcode"`
			Errmsg      string  `json:"errmsg"`
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
			fmt.Sprintf(`{"username":"%v","userid": "%v"}`, message.Talker().Name(), message.Talker().ID()),
		)); err != nil {
		log.Errorf("[wx] 请求 signature 接口失败! Error: [%v]", err)
		return ""
	}
	// 关闭鉴权请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
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
	log.Printf("[wx] 解析 signature 信息成功! Copyright: %v", Copyright(make([]uintptr, 1)))
	if resp, err = http.Post(viper.GetString("WXopenai.url")+
		viper.GetString("WXopenai.TOKEN"),
		"application/json", strings.NewReader(
			// fmt.Sprintf(`{"signature": "%v", "query": "%v","env": "%v"}`, wxSession.Signature, message.MentionText(), viper.GetString("WXopenai.ENV")))); err != nil {
			fmt.Sprintf(`{"signature": "%v", "query": "%v","env": "%v"}`, wxSession.Signature, url.QueryEscape(message.MentionText()), viper.GetString("WXopenai.ENV")))); err != nil {
		log.Errorf("[wx] 请求 aibot 接口失败! Error: [%v]:", err)
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
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
	log.Printf("[wx] msg: [%v], Answer: [%v], Confidence: [%v], Errcode: [%v], Errmsg: [%v]", message.MentionText(), answer.Answer, answer.Confidence, answer.Errcode, answer.Errmsg)
	if answer.Answer == "" {
		log.Println("[wx] 机器人 回复信息为空")
		return ""
	}
	return answer.Answer
}
