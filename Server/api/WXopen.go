package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	. "wechatBot/data"
)

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

func WXAPI(messages MessageInfo) MessageInfo {
	/*
		WXopenai.TOKEN
		WXopenai.ENV
	*/
	if NightMode(messages.UserID) {
		log.Println("现在处于夜间模式，请在白天使用")
		return messages
	} else {
		var (
			res       *http.Response
			err       error
			body      []byte
			wxSession = WxSession{}
			answer    = Answer{}
		)
		if res, err = http.Post("https://openai.weixin.qq.com/openapi/sign/"+viper.GetString("WXopenai.TOKEN"), "application/json", strings.NewReader(fmt.Sprintf(`{"username":"%s","userid": "%s"}`, messages.UserName, messages.UserID))); err != nil {
			log.Println("请求 signature 接口失败! 错误:", err)
		} else {
			if body, err = ioutil.ReadAll(res.Body); err != nil {
				log.Println("读取 signature 信息失败! 错误:", err)
			} else {
				if err = json.Unmarshal(body, &wxSession); wxSession.ExpiresIn == 0 {
					log.Printf("解析 signature 信息失败! Error: %+v", wxSession.Errmsg)
				} else {
					log.Printf("解析 signature 信息成功!")
					if res, err = http.Post("https://openai.weixin.qq.com/openapi/aibot/"+viper.GetString("WXopenai.TOKEN"), "application/json", strings.NewReader(fmt.Sprintf(`{"signature": "%s", "query": "%s","env": "%s"}`, wxSession.Signature, messages.Content, viper.GetString("WXopenai.ENV")))); err != nil {
						log.Println("请求 aibot 接口失败! 错误:", err)
					} else {
						if body, err = ioutil.ReadAll(res.Body); err != nil {
							log.Println("读取 aibot 信息失败! 错误:", err)
						} else {
							if err = json.Unmarshal(body, &answer); answer.Errcode != 0 {
								log.Printf("解析 aibot 信息失败! Error: %+v", answer.Errmsg)
							} else {
								log.Printf("解析 aibot 信息成功!")
								if answer.Answer != "" {
									messages.Reply = answer.Answer
									log.Printf("wx 机器人 回复信息: %+v", messages.Reply)
									messages.AutoInfo += " 回复: [" + messages.Reply + "]"
									//log.Printf("WXopenai.TOKEN:[%s] messages.UserName:[%s], messages.UserID:[%s] wxSession.Signature:[%s] messages.Content:[%v] WXopenai.ENV:[%s] answer:[%v]", viper.GetString("WXopenai.TOKEN"), messages.UserName, messages.UserID, wxSession.Signature, messages.Content, viper.GetString("WXopenai.ENV"), answer)
								} else {
									log.Println("wx 机器人 回复信息为空")
								}
							}
						}
					}
					defer func(Body io.ReadCloser) {
						err := Body.Close()
						if err != nil {
							fmt.Println(err)
						}
					}(res.Body)
				}
			}
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(res.Body)
		return messages
	}
}
