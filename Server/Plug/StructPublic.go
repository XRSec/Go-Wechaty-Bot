package Plug

type (
	MessageInfo struct {
		Date        string `json:"data"`         // 时间
		Status      bool   `json:"status"`       // 群聊属性
		AtMe        bool   `json:"atme"`         // 是否@我
		RoomName    string `json:"room_name"`    // 群聊名称
		RoomID      string `json:"room_id"`      // 群聊ID
		UserName    string `json:"user_name"`    // 用户名称
		UserID      string `json:"userid"`       // 用户ID
		Content     string `json:"content"`      // 聊天内容
		AutoInfo    string `json:"auto_info"`    // 信息一览
		ReplyResult string `json:"reply_result"` // 自动回复
		Reply       bool   `json:"reply"`        // 自动回复状态
		PassResult  string `json:"pass_result"`  // pass 原因
		Pass        bool   `json:"pass"`         // pass 状态
	}
)
