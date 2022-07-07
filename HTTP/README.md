# HTTP Proxy
## 安装
需 Python >= 3.8
```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install -r requirements.txt
```

## 一键部署
访问 [API 密钥管理](https://console.cloud.tencent.com/cam/capi) 获取 `SecretId` 与 `SecretKey`，填入 `setup.py` 中

运行如下代码查看具体部署方式
```bash
python3 setup.py --help
```

## 手动项目配置

### 函数配置
1. 开通[腾讯云函数服务](https://console.cloud.tencent.com/scf/list)
2. 在 函数服务 > 新建 中使用自定义创建，函数名称及地域任选，运行环境选择 Python3.6。
![函数创建](img/create_function.png)

3. 修改 server.py 中的 `SCF_TOKEN` 为随机值（该值将用于鉴权），并将相同的值填入 client.py 中的 `SCF_TOKEN`，将 server.py 代码复制粘贴到编辑器中。
4. 点击完成
5. 如需更多 ip 数，可重复上述步骤创建函数服务，地域选择不同区域。
![地区列表](img/regions.png)

### 触发器配置
1. 成功创建函数后进入 触发管理，创建触发器
![触发器](img/trigger.jpg)


2. 触发方式选择 API 网关触发，其他保持不变即可
![网关](img/gateway.jpg)

3. 将触发器中的访问路径添加至 client.py 中 `scf_servers` 变量中，以逗号 `,` 分隔。


## 客户端配置
本项目基于 mitmproxy 提供本地代理，为代理 HTTPS 流量需安装证书。
运行 `mitmdump` 命令，证书目录自动生成在在 ~/.mitmproxy 中，安装并信任。

开启代理开始运行：
```bash
mitmdump -s client.py -p 8081 --no-http2
```

如在 VPS 上运行需将 `block_global` 参数设为 false
```bash
mitmdump -s client.py -p 8081 --no-http2 --set block_global=false
```

## 效果
挂上代理获取当前 ip:
![ip](img/ip.png)
查询 ipinfo 为腾讯的服务器:
![tencent](img/tencent.png)

### ip 数量
经测试，单个地区服务器 200 个请求分配 ip 数量在 60-70 左右。

## 限制
1. 请求与响应流量包不能大于 6M
2. 云函数操作最大超时限制默认为 3 秒，可在云函数环境配置中修改执行超时时间
3. 因云函数限制不能进行长连接，仅支持代理 HTTP 流量

## 免责声明
此工具仅供测试和教育使用，请勿用于非法目的。
