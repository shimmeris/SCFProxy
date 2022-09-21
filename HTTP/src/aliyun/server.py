# -*- coding: utf8 -*-
import json
from base64 import b64decode, b64encode

import urllib3


def handler(environ: dict, start_response):
    try:
        request_body_size = int(environ.get('CONTENT_LENGTH', 0))
    except (ValueError):
        request_body_size = 0
    request_body = environ['wsgi.input'].read(request_body_size)

    kwargs = json.loads(request_body.decode("utf-8"))
    kwargs['body'] = b64decode(kwargs['body'])

    http = urllib3.PoolManager(cert_reqs="CERT_NONE", assert_hostname=False)
    # Prohibit automatic redirect to avoid network errors such as connection reset
    r = http.request(**kwargs, retries=False, decode_content=False)

    response = {
        "headers": {k.lower(): v.lower() for k, v in r.headers.items()},
        "status_code": r.status,
        "content": b64encode(r._body).decode('utf-8')
    }

    status = '200 OK'
    response_headers = [('Content-type', 'text/json')]
    start_response(status, response_headers)
    return [json.dumps(response)]
