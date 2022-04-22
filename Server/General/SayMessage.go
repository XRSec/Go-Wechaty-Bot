package General

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type (
	TulingBotResult struct {
		Code int    `json:"code"`
		Text string `json:"text"`
	}
	DingBotResult struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
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
	tulingBotResult TulingBotResult
	dingBotResult   DingBotResult
	resp            *http.Response
	lastDate        time.Time
)

func NightMode(userID string) bool {
	//当前时间
	startTimeStr := "00:00:00"
	endTimeStr := "06:00:00"
	now := time.Now()
	//当前时间转换为"年-月-日"的格式
	format := now.Format("2006-01-02")
	//转换为time类型需要的格式
	layout := "2006-01-02 15:04:05"
	//将开始时间拼接“年-月-日 ”转换为time类型
	timeStart, _ := time.ParseInLocation(layout, format+" "+startTimeStr, time.Local)
	//将结束时间拼接“年-月-日 ”转换为time类型
	timeEnd, _ := time.ParseInLocation(layout, format+" "+endTimeStr, time.Local)
	//使用time的Before和After方法，判断当前时间是否在参数的时间范围
	if userID == viper.GetString("bot.adminid") {
		log.Println("[NightMode] 管理员")
		return true
	} else {
		return !(now.Before(timeEnd) && now.After(timeStart))
	}
}
func ChatTimeLimit(date string) bool {
	//当前时间
	if date == "" {
		return true
	}
	now := time.Now()
	if lastDate, err = time.Parse("2006-01-02 15:04:05", date); err != nil {
		log.Errorf("[ChatTimeLimit] 时间转换错误, Error: [%s], Date: [%s]", err, date)
		return false
	}
	//计算两个时间相差的秒数
	second := int(now.Sub(lastDate).Seconds())
	if second < 30 {
		log.Println("[ChatTimeLimit] 时间相差不足")
		return false
	}
	return true
}

func TulingMessage(msg string) string {
	// 发送请求
	tulingWebhook := viper.GetString("Tuling.URL") + viper.GetString("Tuling.TOKEN")
	if resp, err = http.Get(tulingWebhook + url.QueryEscape(msg)); err != nil {
		log.Errorf("[图灵] 机器人请求错误: [%s]", err)
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println("[图灵] Close body error:", err)
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&tulingBotResult); err != nil {
		log.Errorf("[图灵] 机器人解析错误: [%s]", err)
		return ""
	}
	// 判断响应码
	if tulingBotResult.Code != 100000 {
		log.Errorf("[图灵] 机器人返回错误: [%s]", tulingBotResult.Text)
		return ""
	}
	// 输出结果
	log.Printf("[图灵] 机器人 回复信息: %+v", tulingBotResult.Text)
	return tulingBotResult.Text
}

func DingMessage(msg string) {
	dingWebHook := viper.GetString("Ding.URL") + viper.GetString("Ding.TOKEN")
	content := fmt.Sprintf(" {\"msgtype\": \"text\",\"text\": {\"content\": \"%s %s\"}}", viper.GetString("Ding.KEYWORD"), msg)
	// 发送请求
	if resp, err = http.Post(dingWebHook, "application/json; charset=utf-8", strings.NewReader(content)); err != nil {
		log.Errorf("[Ding] 机器人请求错误: %s", err)
		return
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("[Ding] 关闭请求错误: [%s]", err)
		}
	}(resp.Body)
	if err = json.NewDecoder(resp.Body).Decode(&dingBotResult); err != nil {
		log.Errorf("[Ding] 机器人请求错误: %s", err)
		return
	}
	if dingBotResult.Errcode == 0 {
		log.Println("[Ding] 消息发送成功!")
	} else {
		log.Errorf("[Ding] 消息发送失败: [%s]", err)
	}
}

func WXAPI(message *user.Message, msg string) string {
	/*
		WXopenai.TOKEN
		WXopenai.ENV
	*/
	var (
		res       *http.Response
		body      []byte
		wxSession = WxSession{}
		answer    = Answer{}
	)
	// 微信鉴权
	if res, err = http.Post(
		viper.GetString("WXopenai.signUrl")+
			viper.GetString("WXopenai.TOKEN"), "application/json",
		strings.NewReader(
			fmt.Sprintf(`{"username":"%s","userid": "%s"}`, message.From().Name(), message.From().ID()),
		)); err != nil {
		log.Errorf("[wx] 请求 signature 接口失败! Error: [%s]", err)
		return ""
	}
	// 关闭鉴权请求
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf(err.Error())
		}
	}(res.Body)

	// 读取鉴权信息
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		log.Errorf("[wx] 读取 signature 信息失败! Error: [%s]", err)
		return ""
	}
	// 解析鉴权信息
	if err = json.Unmarshal(body, &wxSession); wxSession.ExpiresIn == 0 {
		log.Printf("[wx] 解析 signature 信息失败! Error: %+v", wxSession.Errmsg)
		return ""
	}
	log.Printf("[wx] 解析 signature 信息成功!")
	if res, err = http.Post(viper.GetString("WXopenai.url")+
		viper.GetString("WXopenai.TOKEN"),
		"application/json", strings.NewReader(
			fmt.Sprintf(`{"signature": "%s", "query": "%s","env": "%s"}`, wxSession.Signature, msg, viper.GetString("WXopenai.ENV")))); err != nil {
		log.Errorf("[wx] 请求 aibot 接口失败! Error: [%s]:", err)
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf(err.Error())
		}
	}(res.Body)

	// 读取信息
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		log.Errorf("[wx] 读取 aibot 信息失败! Error: [%s]:", err)
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
	//log.Printf("[wx] AnsNodeName: [%s], Answer: [%s], Errcode: [%v], Errmsg: [%s]", answer.AnsNodeName, answer.Answer, answer.Errcode, answer.Errmsg)
	return answer.Answer
}

func SayMessage(message *user.Message, msg string) {
	if !NightMode(message.From().ID()) { // 夜间模式
		return
	}
	if msg == "" {
		if msg = WXAPI(message, message.MentionText()); msg == "" {
			if msg = TulingMessage(message.MentionText()); msg == "" {
				msg = "我也不知道呀!"
			}
		}
	}
	//if _, err = message.Say(msg); err != nil {
	//	log.Errorf("[SayMessage] [%s], error: %v", msg, err)
	//	return
	//}
	// TODO 0.79 私聊有问题
	_, _ = message.Say(msg)
	Messages.Reply = msg
	Messages.ReplyStatus = true
	Messages.AutoInfo = Messages.AutoInfo + "[" + msg + "]"
	viper.Set(fmt.Sprintf("Chat.%v.Date", message.From().ID()), Messages.Date)
	DingMessage(Messages.AutoInfo)
	return
}
