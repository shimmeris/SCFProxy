# SOCKS5

## 安装
上传 socks_client 文件夹到 VPS 上，执行安装
```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install -r requirements.txt
```

## 项目配置
### 云函数配置
参见 HTTP Proxy 的[代理配置](https://github.com/shimmeris/SCFProxy/tree/main/HTTP)

注意事项
1. 修改 server.py 中的 `bridge_ip` 与 `bridge_port` 变量为自己的 VPS 的 ip 及开启监听的端口
2. 修改云函数超时时间为 900s（因此一个 SOCKS5 长连接最多维持 15m）

## 客户端配置
执行 socks5.py（仅支持 Python >= 3.8)
```bash
$ python3 socks5.py -h
usage: socks5.py [-h] -u SCF_URL [-l ip] -sp SOCKS_PORT -bp BRIDGE_PORT [--user USER] [--passwd PASSWD]
SCF Socks5 Proxy Server
optional arguments:
  -h, --help            show this help message and exit
  -u SCF_URL, --scf-url SCF_URL
                        API Gate Way URL
  -l ip, --listen ip    Bind address to listen, default to 0.0.0.0
  -sp SOCKS_PORT, --socks-port SOCKS_PORT
                        Port accept connections from client
  -bp BRIDGE_PORT, --bridge-port BRIDGE_PORT
                        Port accept connections from SCF
  --user USER           Authentication username
  --passwd PASSWD       Authentication password
```

* `-u` 参数需要填写 API 网关提供的地址，必填
* `-l` 表示本机监听的 ip，默认为 0.0.0.0
* `-sp` 表示 SOCKS5 代理监听的端口，必填
* `bp` 表示用于监听来自云函数连接的端口，与 server.py 中的 `bridge_port` 相同，必填
* `--user` 和 `--passwd` 将用于 SOCKS5 服务器对连接进行身份验证，客户端需配置相应的用户名和密码

常用语法
```bash
python3 socks5.py -u "https://service-xxx.sh.apigw.tencentcs.com/release/xxx" -bp 53203 -sp 53201 --user test --passwd test
```

## 免责声明
此工具仅供测试和教育使用，请勿用于非法目的。
