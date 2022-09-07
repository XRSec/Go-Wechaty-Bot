package AutoReply

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	. "wechatBot/General"
	. "wechatBot/Plug"
)

func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	plug.OnMessage(onMessage)
	return plug
}

func onMessage(context *wechaty.Context, message *user.Message) {
	m, ok := (context.GetData("msgInfo")).(MessageInfo)
	if !ok {
		log.Errorf("Conversion Failed CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Pass {
		log.Infof("Pass CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if m.Reply {
		log.Infof("Reply CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if !m.AtMe {
		log.Infof("AtMe CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Type() != schemas.MessageTypeText {
		log.Infof("Type: [%v] CoptRight: [%v]", message.Type().String(), Copyright(make([]uintptr, 1)))
		return
	}
	if message.Self() {
		log.Infof("Self CoptRight: [%s]", Copyright(make([]uintptr, 1)))
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Infof("Age: [%v] CoptRight: [%v]", message.Age()/(60*time.Second), Copyright(make([]uintptr, 1)))
		return
	}

	/*
	 TODO Debug
	*/
	//if message.Room() == nil {
	//	return
	//}
	//if !strings.Contains(message.Room().Topic(), "Debug") {
	//	return
	//}

	var msg string
	if message.MentionText() == "" {
		msg = "你想表达什么[破涕为笑]?"
		goto labelSay
	}
	if msg = wxApi(message); msg != "" {
		goto labelSay
	}
	if msg = qingYunKe(message.MentionText()); msg != "" {
		goto labelSay
	}
	if msg = tuLingApi(message.MentionText()); msg != "" {
		msg = "我又不会了 [捂脸]"
	}
labelSay:
	//time.Sleep(5 * time.Second)
	SayMessage(context, message, msg)
}

func wxApi(message *user.Message) string {
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
		err       error
	)
	// 微信鉴权
	if resp, err = http.Post(
		viper.GetString("wxOpenAi.signUrl")+
			viper.GetString("wxOpenAi.TOKEN"), "application/json",
		strings.NewReader(
			fmt.Sprintf(`{"username":"%v","userid": "%v"}`, message.Talker().Name(), message.Talker().ID()),
		)); err != nil {
		log.Errorf("[wx] 请求 signature 接口失败! Error: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 关闭鉴权请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("Error: [%v] CoptRight: [%s]", err.Error(), Copyright(make([]uintptr, 1)))
		}
	}(resp.Body)

	// 读取鉴权信息
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("[wx] 读取 signature 信息失败! Error: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 解析鉴权信息
	if err = json.Unmarshal(body, &wxSession); wxSession.ExpiresIn == 0 {
		log.Errorf("[wx] 解析 signature 信息失败! Error: %v CoptRight: [%v]", wxSession.Errmsg, Copyright(make([]uintptr, 1)))
		return ""
	}
	//log.Infof("[wx] 解析 signature 信息成功! Copyright: [%v]", Copyright(make([]uintptr, 1)))
	if resp, err = http.Post(viper.GetString("wxOpenAi.url")+
		viper.GetString("wxOpenAi.TOKEN"),
		"application/json", strings.NewReader(
			fmt.Sprintf(`{"signature": "%v", "query": "%v","env": "%v"}`, wxSession.Signature, message.MentionText(), viper.GetString("WXopenai.ENV")))); err != nil {
		//fmt.Sprintf(`{"signature": "%v", "query": "%v","env": "%v"}`, wxSession.Signature, url.QueryEscape(message.MentionText()), viper.GetString("WXopenai.ENV")))); err != nil {
		log.Errorf("[wx] 请求 aibot 接口失败! Error: [%v] CoptRight: [%v]:", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("Error: [%v] CoptRight: [%s]", err.Error(), Copyright(make([]uintptr, 1)))
		}
	}(resp.Body)

	// 读取信息
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("[wx] 读取 aibot 信息失败! Error: [%v] CoptRight: [%v]:", err, Copyright(make([]uintptr, 1)))
		return ""
	}

	// 解析信息
	if err = json.Unmarshal(body, &answer); answer.Errcode != 0 {
		log.Errorf("[wx] 解析 aibot 信息失败! Error: %v CoptRight: [%v]", answer.Errmsg, Copyright(make([]uintptr, 1)))
		return ""
	}

	//log.Infof("[wx] 解析 aibot 信息成功!")
	// log.Printf("[wx] msg: [%v], Answer: [%v], Confidence: [%v], Errcode: [%v], Errmsg: [%v]", message.MentionText(), answer.Answer, answer.Confidence, answer.Errcode, answer.Errmsg)
	if answer.Answer == "" {
		return ""
	}
	return answer.Answer
}

func qingYunKe(msg string) string {
	type (
		qingYunKeResult struct {
			Result  int
			Content string
		}
	)
	var (
		q    qingYunKeResult
		resp *http.Response
		err  error
	)
	if resp, err = http.Get(viper.GetString("QingYunKe.url") + msg); err != nil {
		log.Errorf("[青云客] 机器人请求错误: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("[青云客] Close Body Error: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&q); err != nil {
		log.Errorf("[图青云客灵] 机器人解析错误: [%v]  CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// CoptRight: [%v]
	// 判断响应码
	if q.Result != 0 {
		log.Infof("[青云客] 机器人返回错误: [%v] CoptRight: [%v]", q.Content, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 输出结果
	log.Infof("[青云客] 机器人 回复信息: %v CoptRight: [%v]", q.Content, Copyright(make([]uintptr, 1)))
	return q.Content
}

func tuLingApi(msg string) string {
	type (
		TulingBotResult struct {
			Code int    `json:"code"`
			Text string `json:"text"`
		}
	)
	var (
		t    TulingBotResult
		resp *http.Response
		err  error
	)
	// 发送请求
	tuLingWebhook := viper.GetString("TuLing.URL") + viper.GetString("TuLing.TOKEN")
	if resp, err = http.Get(tuLingWebhook + url.QueryEscape(msg)); err != nil {
		log.Errorf("[图灵] 机器人请求错误: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("[图灵] Close Body Error: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		log.Errorf("[图灵] 机器人解析错误: [%v] CoptRight: [%v]", err, Copyright(make([]uintptr, 1)))
		return ""
	}
	// 判断响应码
	// TODO 添加自动更换TOKEN
	if t.Code != 100000 {
		if strings.Contains(t.Text, "当天请求次数已用完") {
			if token2 := viper.GetString("TuLing.token2"); token2 != "" {
				viper.Set("TuLing.token2", viper.GetString("TuLing.token"))
				viper.Set("TuLing.token", token2)
				return ""
			}
		}
		if t.Text != "你想和我说什么呢?" {
			return ""
		}
		log.Infof("[图灵] 机器人返回错误: [%v] CoptRight: [%v]", t.Text, Copyright(make([]uintptr, 1)))
	}
	// 输出结果
	log.Infof("[图灵] 机器人 回复信息: %v CoptRight: [%v]", t.Text, Copyright(make([]uintptr, 1)))
	return t.Text
}
