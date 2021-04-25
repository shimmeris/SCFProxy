# -*- coding: utf8 -*-
import pytz
import datetime
import requests
import pymysql.cursors


push_back_host = ""
db_host = ""
db_user = ""
db_pass = ""
db_port = 123

db = "SCF"
db_table = "Connections"
tz = pytz.timezone("Asia/Shanghai")


def send(connectionID, data):
    retmsg = {
        "websocket": {
            "action": "data send",
            "secConnectionID": connectionID,
            "dataType": "text",
            "data": data,
        }
    }
    requests.post(push_back_host, json=retmsg)


def close_ws(connectionID):
    msg = {"websocket": {"action": "closing", "secConnectionID": connectionID}}
    requests.post(push_back_host, json=msg)


def record_connectionID(connectionID):
    try:
        conn = pymysql.connect(
            host=db_host,
            user=db_user,
            password=db_pass,
            port=db_port,
            db=db,
            charset="utf8",
            cursorclass=pymysql.cursors.DictCursor,
        )
        with conn.cursor() as cursor:
            sql = f"use {db}"
            cursor.execute(sql)
            time = datetime.datetime.now(tz).strftime("%Y-%m-%d %H:%M:%S")
            sql = f"insert INTO {db_table} (`ConnectionID`, `is_user`, `Date`) VALUES ('{str(connectionID)}', 0, '{time}')"
            cursor.execute(sql)
            conn.commit()
    except Exception as e:
        send(connectionID, f"[Error]: {e}")
        close_ws(connectionID)
    finally:
        conn.close()


def main_handler(event, context):
    if "requestContext" not in event.keys():
        return {"errNo": 101, "errMsg": "not found request context"}
    if "websocket" not in event.keys():
        return {"errNo": 102, "errMsg": "not found web socket"}

    connectionID = event["websocket"]["secConnectionID"]
    retmsg = {
        "errNo": 0,
        "errMsg": "ok",
        "websocket": {"action": "connecting", "secConnectionID": connectionID},
    }
    record_connectionID(connectionID)
    return retmsg
