# -*- coding: utf8 -*-
import pickle
from base64 import b64decode, b64encode

from requests import Session


SCF_TOKEN = "Token"


def authorization():
    return {
        "isBase64Encoded": False,
        "statusCode": 401,
        "headers": {},
        "body": "Please provide correct SCF-Token",
    }


def main_handler(event: dict, context: dict):
    # Tencent cloud has its own authorization system https://console.cloud.tencent.com/cam/capi
    # But it may be a little overqualified for a simple usage like this
    try:
        token = event["headers"]["scf-token"]
    except KeyError:
        return authorization()

    if token != SCF_TOKEN:
        return authorization()

    data = event["body"]
    prepped = pickle.loads(b64decode(data))
    with Session() as s:
        # Prohibit automatic redirect to avoid network errors such as connection reset
        r = s.send(prepped, verify=False, allow_redirects=False)

    # TODO: REFACTOR NEEDED. Return response headers and body directly.
    # There are many errors occured when setting headers to r.headers with some aujustments(https://cloud.tencent.com/document/product/583/12513).
    # and the response `r.content`/`r.raw.read()` to body.(like gzip error)
    serialized_resp = pickle.dumps(r)

    return {
        "isBase64Encoded": False,
        "statusCode": 200,
        "headers": {},
        "body": b64encode(serialized_resp).decode("utf-8"),
    }
