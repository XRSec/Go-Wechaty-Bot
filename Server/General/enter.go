package General

import (
	"github.com/wechaty/go-wechaty/wechaty"
	"os/exec"
)

var (
	err                   error
	GatewayCmd            *exec.Cmd
	globleServiceInstance = &GlobleService{}
)

type GlobleService struct {
	Bot *wechaty.Wechaty
}
