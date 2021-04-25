import json
import socket
import select


bridge_ip = ""
bridge_port = 1234


def main_handler(event, context):
    data = json.loads(event["body"])
    out = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    out.connect((data["host"], data["port"]))

    bridge = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    bridge.connect((bridge_ip, bridge_port))
    bridge.send(data["uid"].encode("ascii"))

    while True:
        readable, _, _ = select.select([out, bridge], [], [])
        if out in readable:
            data = out.recv(4096)
            bridge.send(data)
        if bridge in readable:
            data = bridge.recv(4096)
            out.send(data)
