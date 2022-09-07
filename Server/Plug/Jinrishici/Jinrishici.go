package Jinrishici

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wechaty/go-wechaty/wechaty"
	"io"
	"net/http"
)

var (
	err error
)

/*
	Health()
	健康监测
*/
func New() *wechaty.Plugin {
	plug := wechaty.NewPlugin()
	log.Errorf("今日诗词模块不可当做模块导入,请直接使用!")
	return plug
}

func Do() string {
	type Jin struct {
		Content  string
		origin   string
		Author   string
		Category string
	}

	var (
		jin  Jin
		resp *http.Response
		err  error
	)
	// 发起请求
	resp, err = http.Get("https://v1.jinrishici.com/shuqing")
	if err != nil {
		log.Errorf("今日诗词接口请求错误: [%v] ", err)
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Errorf("今日诗词请求 Close Body Error: [%v]", err)
		}
	}(resp.Body)
	// 解析响应
	if err = json.NewDecoder(resp.Body).Decode(&jin); err != nil {
		log.Errorf("今日诗词接口解析错误: [%v]", err)
	}
	if jin.Content == "" {
		jin.Content = "情如之何，暮涂为客，忍堪送君。"
	}
	return jin.Content
}
