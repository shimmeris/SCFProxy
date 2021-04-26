import pickle
from urllib.parse import urlparse
from base64 import b64encode, b64decode

import requests
import mitmproxy
from mitmproxy.net.http import Headers


scf_server = ""
SCF_TOKEN = "Token"


def request(flow: mitmproxy.http.HTTPFlow):
    # TODO: Just send request kwargs rather than serialized `PreparedRequest` object.
    r = flow.request
    req = requests.Request(
        method=r.method,
        url=r.pretty_url,
        headers=r.headers,
        cookies=r.cookies,
        params=r.query,
        data=r.raw_content,
    )
    prepped = req.prepare()
    serialized_req = pickle.dumps(prepped)

    flow.request = flow.request.make(
        "POST",
        url=scf_server,
        content=b64encode(serialized_req),
        headers={
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Encoding": "gzip, deflate, compress",
            "Accept-Language": "en-us;q=0.8",
            "Cache-Control": "max-age=0",
            "User-Agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
            "Connection": "close",
            "Host": urlparse(scf_server).netloc,
            "SCF-Token": SCF_TOKEN,
        },
    )


def response(flow: mitmproxy.http.HTTPFlow):
    if flow.response.status_code != 200:
        mitmproxy.ctx.log.warn("Error")

    if flow.response.status_code == 401:
        flow.response.headers = Headers(content_type="text/html;charset=utf-8")
        return

    if flow.response.status_code == 433:
        flow.response.headers = Headers(content_type="text/html;charset=utf-8")
        flow.response.content="<html><body>操作已超过云函数服务最大时间限制，可在函数配置中修改执行超时时间</body></html>",
        return


    if flow.response.status_code == 200:
        body = flow.response.content.decode("utf-8")
        resp = pickle.loads(b64decode(body))

        r = flow.response.make(
            status_code=resp.status_code,
            headers=dict(resp.headers),
            content=resp.content,
        )
        flow.response = r
