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