import asyncio

from utils import print_time
from models import Conn, uid_socket


async def scf_handle(reader: asyncio.StreamReader, writer: asyncio.StreamWriter):
    bridge = Conn("Bridge", reader, writer)
    uid = await bridge.read(4)
    uid = uid.decode("ascii")
    client = uid_socket[uid]
    bridge.target = client.target
    bridge_addr, _ = bridge.writer.get_extra_info("peername")
    print_time(f"Tencent IP:{bridge_addr} <=> {client.target} established")

    await socks5_forward(client, bridge)


async def socks5_forward(client: Conn, target: Conn):
    async def forward(src: Conn, dst: Conn):
        while True:
            try:
                data = await src.read(4096)
                if not data:
                    break
                await dst.write(data)
            except RuntimeError as e:
                print_time(f"RuntimeError occured when connecting to {src.target}")
                print_time(f"Direction: {src.role} => {dst.role}")
                print(e)
            except ConnectionResetError:
                print_time(f"{src.target} sends a ConnectionReset")
                pass

            await asyncio.sleep(0.01)

    tasks = [forward(client, target), forward(target, client)]
    await asyncio.gather(*tasks)
