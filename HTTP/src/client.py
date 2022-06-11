import json
from random import choice
from urllib.parse import urlparse
from base64 import b64encode, b64decode

import mitmproxy
from mitmproxy.version import VERSION


if int(VERSION[0]) > 6:
    from mitmproxy.http import Headers
else:
    from mitmproxy.net.http import Headers


with open('cities.txt', 'r') as f:
    scf_servers = [line.split()[1].strip() for line in f]


css = """pre {
    white-space: pre-wrap;       /* Since CSS 2.1 */
    white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
    white-space: -pre-wrap;      /* Opera 4-6 */
    white-space: -o-pre-wrap;    /* Opera 7 */
    word-wrap: break-word;       /* Internet Explorer 5.5+ */
}"""


def request(flow: mitmproxy.http.HTTPFlow):
    scf_server = choice(scf_servers)
    request = flow.request
    data = {
        "method": request.method,
        "url": request.pretty_url,
        "headers": dict(request.headers),
        "body": b64encode(request.raw_content).decode('utf-8'),
    }
    flow.request = flow.request.make(
        "POST",
        url=scf_server,
        content=json.dumps(data),
        headers={
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Encoding": "gzip, deflate, compress",
            "Accept-Language": "en-us;q=0.8",
            "Cache-Control": "max-age=0",
            "User-Agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
            "Connection": "Keep-Alive",
            "Host": urlparse(scf_server).netloc,
        },
    )


def response(flow: mitmproxy.http.HTTPFlow):
    status = flow.response.status_code

    if status != 200:
        mitmproxy.ctx.log.warn("Error")

    if status == 401:
        flow.response.headers = Headers(content_type="text/html;charset=utf-8")
        return

    elif status == 433:
        flow.response.headers = Headers(content_type="text/html;charset=utf-8")
        flow.response.content="test",
        return

    elif status == 430:
        body = flow.response.content.decode("utf-8")
        data = json.loads(body)
        flow.response.headers = Headers(content_type="text/html;charset=utf-8")
        flow.response.text = f'<style>{css}</style><pre id="json"></pre><script>document.getElementById("json").textContent = JSON.stringify({data}, undefined, 2);</script>'
        return

    elif status == 200:
        body = flow.response.content.decode("utf-8")
        resp = json.loads(body)
        headers = resp["headers"]
        raw_content = b64decode(resp["content"])

        r = flow.response.make(
            status_code=resp["status_code"],
            headers=dict(headers),
        )

        if headers.get("content-encoding", None):
            r.raw_content = raw_content
            if "transfer-encoding" not in r.headers:
                r.headers['content-length'] = str(len(r.raw_content))
        else:
            r.content = raw_content

        flow.response = r
