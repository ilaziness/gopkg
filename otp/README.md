
## OTP的应用

otp全称One Time Password，一次性密码，通常用来做账户的二次认证

分两种类型的算法:

- 基于时间TOTP, TOTP比较常用
- 和基于次数的HTOP


### TOTP的使用流程

1. 生成TOTP key `key,_ := totp.Generate(...)`
2. 输出二维码用户app扫描，`key.Image(...)`
3. 验证用户输入的验证码和帐号密码，`totp.Validate(...)`
4. 保存密钥到数据库和用户关联起来，密钥`key.Secret()`
5. 可不做，展示一组恢复代码，代替二次验证的验证码
   恢复代码就是随机生成的字符串，用来替代TOTP验证码，用户需要下载保存好，不可以泄漏

### 手机客户端

- Google身份验证器
- 微软的Authenticator，经测试只支持TOTP
