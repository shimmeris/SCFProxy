# SCFProxy
一个利用云函数实现各种功能的工具。

# 功能
* HTTP 代理
  * 腾讯云
  * 阿里云 （感谢 [@lyc8503](https://github.com/lyc8503) 提供）
* SOCKS5 代理
* 接收反弹 shell
* C2 域名隐藏

## 原理
详细原理参见文章[浅谈云函数的利用面](https://xz.aliyun.com/t/9502)

## TODO
* 使用 Go 重构
* 支持其他云厂商
* 整合多个 HTTP 代理厂商，实现多厂商的一键部署
