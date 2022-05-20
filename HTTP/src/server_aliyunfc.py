# -*- coding: utf8 -*-
import json
import pickle
from base64 import b64decode, b64encode

import requests


SCF_TOKEN = "Token"


def handler(environ: dict, start_response):
    try:
        token = environ["HTTP_SCF_TOKEN"]
        assert token == SCF_TOKEN, "Invalid token."
    except:
        status = '403 Forbidden'
        response_headers = [('Content-type', 'text/json')]
        start_response(status, response_headers)
        return []

    try:
        request_body_size = int(environ.get('CONTENT_LENGTH', 0))
    except (ValueError):
        request_body_size = 0
    request_body = environ['wsgi.input'].read(request_body_size)

    kwargs = json.loads(request_body.decode("utf-8"))
    kwargs['data'] = b64decode(kwargs['data'])
    # Prohibit automatic redirect to avoid network errors such as connection reset
    r = requests.request(**kwargs, verify=False, allow_redirects=False)

    # TODO: REFACTOR NEEDED. Return response headers and body directly.
    # There are many errors occured when setting headers to r.headers with some aujustments(https://cloud.tencent.com/document/product/583/12513).
    # and the response `r.content`/`r.raw.read()` to body.(like gzip error)
    serialized_resp = pickle.dumps(r)

    status = '200 OK'
    response_headers = [('Content-type', 'text/json')]
    start_response(status, response_headers)
    return [b64encode(serialized_resp)]
