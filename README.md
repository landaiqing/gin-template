# gin 模板代码

## 环境准备

- go 1.19+

## 功能说明

### 注册登录

- 手机号注册登录
- 账号密码登录
- 微信公众号扫码登录
- 忘记密码
- 第三方登录（QQ、Gitee、GitHub）

### 评论

- 评论列表
- 评论回复
- 评论点赞

### 具体技术栈

- gin 框架
- gorm 数据库ORM框架
- jwt 双token验证
- redis 缓存
- powerWechat 微信公众号扫码登录
- websocket 扫码登录监听状态
- 第三方登录(QQ、Gitee、GitHub)
- 短信发送 (阿里云短信，短信宝)
- mongodb 数据库存储评论图片
- 评论图片nsfw检测
- go i18n 国际化
- casbin 权限管理
- ip2region 实现ip定位
- nsq 消息队列(异步处理评论点赞)
- 日志记录(logrus 按天记录日志,按照等级分文件记录日志)
- 评论敏感词检测
- go-captcha 验证码(滑动验证,旋转验证,点选验证)


 


    