# -*- coding: utf8 -*-
import pytz
import pymysql.cursors


push_back_host = ""
db_host = ""
db_user = ""
db_pass = ""
db_port = 123

db = "SCF"
db_table = "Connections"
tz = pytz.timezone("Asia/Shanghai")


def delete_connectionID(connectionID):
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
        sql = f"delete from {db_table} where ConnectionID ='{connectionID}'"
        cursor.execute(sql)
        conn.commit()


def main_handler(event, context):
    if "websocket" not in event.keys():
        return {"errNo": 102, "errMsg": "not found web socket"}

    connectionID = event["websocket"]["secConnectionID"]
    delete_connectionID(connectionID)
    return event
