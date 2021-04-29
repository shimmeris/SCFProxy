import asyncio
import argparse
from socket import inet_ntoa
from functools import partial

import uvloop
import shortuuid

from bridge import scf_handle
from models import Conn, http, uid_socket
from utils import print_time, parse_args, cancel_task


async def socks_handle(
    args: argparse.Namespace, reader: asyncio.StreamReader, writer: asyncio.StreamWriter
):
    client = Conn("Client", reader, writer)

    await socks5_auth(client, args)
    remote_addr, port = await socks5_connect(client)

    client.target = f"{remote_addr}:{port}"
    uid = shortuuid.ShortUUID().random(length=4)
    uid_socket[uid] = client

    data = {"host": remote_addr, "port": port, "uid": uid}
    await http.post(args.scf_url, json=data)


async def socks5_auth(client: Conn, args: argparse.Namespace):
    ver, nmethods = await client.read(2)

    if ver != 0x05:
        client.close()
        cancel_task(f"Invalid socks5 version: {ver}")

    methods = await client.read(nmethods)

    if args.user and b"\x02" not in methods:
        cancel_task(
            f"Unauthenticated access from {client.writer.get_extra_info('peername')[0]}"
        )

    if b"\x02" in methods:
        await client.write(b"\x05\x02")
        await socks5_user_auth(client, args)
    else:
        await client.write(b"\x05\x00")


async def socks5_user_auth(client: Conn, args: argparse.Namespace):
    ver, username_len = await client.read(2)
    if ver != 0x01:
        client.close()
        cancel_task(f"Invalid socks5 user auth version: {ver}")

    username = (await client.read(username_len)).decode("ascii")
    password_len = ord(await client.read(1))
    password = (await client.read(password_len)).decode("ascii")

    if username == args.user and password == args.passwd:
        await client.write(b"\x01\x00")
    else:
        await client.write(b"\x01\x01")
        cancel_task(
            f"Wrong user/passwd connection from {client.writer.get_extra_info('peername')[0]}"
        )


async def socks5_connect(client: Conn):
    ver, cmd, _, atyp = await client.read(4)
    if ver != 0x05:
        client.close()
        cancel_task(f"Invalid socks5 version: {ver}")
    if cmd != 1:
        client.close()
        cancel_task(f"Invalid socks5 cmd type: {cmd}")

    if atyp == 1:
        address = await client.read(4)
        remote_addr = inet_ntoa(address)
    elif atyp == 3:
        addr_len = await client.read(1)
        address = await client.read(ord(addr_len))
        remote_addr = address.decode("ascii")
    elif atyp == 4:
        cancel_task("IPv6 not supported")
    else:
        cancel_task("Invalid address type")

    port = int.from_bytes(await client.read(2), byteorder="big")

    # Should return bind address and port, but it's ok to just return 0.0.0.0
    await client.write(b"\x05\x00\x00\x01\x00\x00\x00\x00\x00\x00")
    return remote_addr, port


async def main():
    args = parse_args()
    handle = partial(socks_handle, args)

    if not args.user:
        print_time("[ALERT] Socks server runs without authentication")

    await http.init_session()
    socks_server = await asyncio.start_server(handle, args.listen, args.socks_port)
    print_time(f"SOCKS5 Server listening on: {args.listen}:{args.socks_port}")
    await asyncio.start_server(scf_handle, args.listen, args.bridge_port)
    print_time(f"Bridge Server listening on: {args.listen}:{args.bridge_port}")

    try:
        await socks_server.serve_forever()
    except asyncio.CancelledError:
        await http.close()


if __name__ == "__main__":
    uvloop.install()
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print_time("[INFO] User stoped server")
