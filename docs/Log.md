## 04/18

- 完善 `加好友、踢人、邀请进群` 的功能
- 添加`退群、倒计时` 的功能
- 关闭 `OnRoomInvite` 的功能
- 修改 `DingMessage` 的传参方式
- 修复 `@ 异常` 的问题
- 添加 `管理员 夜间模式不受限制`
- 加强 `OnError` 的回溯功能
- 移除关键字模式，改为 `微信对话开放平台`
- 修复 自动更新文档 bug

## 04/10

- 添加日志输出 json 格式

- ```diff
  -- AutoInfo: "用户ID: [" + UserID + "] 用户名称: [" + UserName + "]" + message.Text() +"]",
  ++ AutoInfo: message.Text(),
  ```

- 添加 Tuling 图灵机器人
- 添加微信对话平台机器人
- 添加夜间模式
- 调整 消息保存格式

## 04/09

- 添加失败重试机制, 暂时没有奇怪的错误验证
- 添加退出保存配置的功能
- 添加自动编译功能

## 04/07

- 添加 `Atme` 方法 用来替代官方的 `message.MentionSelf()` 方法,需要填写机器人 名称 bot: name: 随缘

## 04/05

- 重构架构，所有功能在Plug文件夹，正在计划性删除 其他文件夹
- Keyword 自定义内容回复
- @机器人的事件只能捕捉到 `@机器人名字` 已经发起[issues](https://github.com/wechaty/puppet-xp/issues/97)