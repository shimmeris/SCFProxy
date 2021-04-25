# -*- coding: utf8 -*-
from os import close
import pytz
import requests
import pymysql.cursors


push_back_host = ""
db_host = ""
db_user = ""
db_pass = ""
db_port = 123
PASSWORD = "test"


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


def get_connectionIDs(conn):
    with conn.cursor() as cursor:
        sql = f"use {db}"
        cursor.execute(sql)
        sql = f"select * from {db_table}"
        cursor.execute(sql)
        result = cursor.fetchall()
        connectionIDs = {c["ConnectionID"]: c["is_user"] for c in result}
    return connectionIDs


def update_user_type(conn, connectionID):
    with conn.cursor() as cursor:
        sql = f"use {db}"
        cursor.execute(sql)
        sql = f"update {db_table} set is_user=True where ConnectionID='{connectionID}'"
        cursor.execute(sql)
        conn.commit()


def main_handler(event, context):
    if "websocket" not in event.keys():
        return {"errNo": 102, "errMsg": "not found web socket"}
    data = event["websocket"]["data"].strip()
    current_connectionID = event["websocket"]["secConnectionID"]

    if data == "close":
        send(current_connectionID, "[INFO] current connection closed")
        close_ws(current_connectionID)
        return

    if data == "help":
        msg = """Commands
        auth PASSWORD - provide a password to set current connection to be a user
        close - close curren websocket connection
        closeall - close all websocket connections
        help - show this help message
        """
        send(current_connectionID, msg)
        return

    conn = pymysql.connect(
        host=db_host,
        user=db_user,
        password=db_pass,
        port=db_port,
        db=db,
        charset="utf8",
        cursorclass=pymysql.cursors.DictCursor,
    )
    connectionIDs = get_connectionIDs(conn)

    if data[:5] == "auth ":
        try:
            password = data.split()[1]
        except IndexError:
            password = None
        if password == PASSWORD:
            send(current_connectionID, "[INFO] AUTH SUCCESS")
            update_user_type(conn, current_connectionID)
        else:
            send(current_connectionID, "[ERROR] AUTH FAILED")
    if data == "closeall":
        send(current_connectionID, "[INFO] all connections closed")
        for ID in connectionIDs.keys():
            close_ws(ID)
        return

    is_current_user = connectionIDs.pop(current_connectionID)
    for ID, is_user in connectionIDs.items():
        if is_current_user:
            send(ID, data)
        elif is_user:
            send(ID, data)

    return "send success"
