package Plug

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Qingyunke(msg string) string {
	type (
		QingyunkeResult struct {
			Result  int
			Content string
		}
	)
	var (
		qingyunkeResult QingyunkeResult
		resp            *http.Response
	)
	// 发送请求
	// apiUrl := fmt.Sprintf("%s?key=%s&appid=%v&msg=%s", viper.GetString("Qingyunke.url"), viper.GetString("Qingyunke.key"), viper.GetInt("Qingyunke.appid"), url.QueryEscape(msg))
	// apiUrl := fmt.Sprintf("%s'%s'", viper.GetString("Qingyunke.url"), url.QueryEscape(msg))

	if resp, err = http.Get(viper.GetString("Qingyunke.url") + msg); err != nil {
		log.Errorf("[青云客] 机器人请求错误: [%v]", err)
		return ""
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println("[青云客] Close body error:", err)
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&qingyunkeResult); err != nil {
		log.Errorf("[图青云客灵] 机器人解析错误: [%v]", err)
		return ""
	}
	// 判断响应码
	// TODO 添加自动更换TOKEN
	if qingyunkeResult.Result != 0 {
		log.Errorf("[青云客] 机器人返回错误: [%v]", qingyunkeResult.Content)
		return ""
	}
	// 输出结果
	log.Printf("[青云客] 机器人 回复信息: %+v", qingyunkeResult.Content)
	return qingyunkeResult.Content
}
