package General

import "log"

func init() {
	log.SetPrefix("\x1b[1;32m[wechatBot] \x1b[0m")
	// \x1b[%dm%s\x1b[0m
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
