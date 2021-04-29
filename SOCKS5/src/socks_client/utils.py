import sys
import asyncio
import argparse
from datetime import datetime, timezone, timedelta


timezone(timedelta(hours=8))


def print_time(data):
    print(f'{datetime.now().strftime("%Y-%m-%d %H:%M:%S")} {data}')


def parse_error(errmsg):
    print("Usage: python " + sys.argv[0] + " [Options] use -h or --help for help")
    sys.exit()


def parse_args():
    parser = argparse.ArgumentParser(description="SCF Socks5 Proxy Server")
    parser.error = parse_error

    parser.add_argument(
        "-u", "--scf-url", type=str, help="API Gate Way URL", required=True
    )
    parser.add_argument(
        "-l",
        "--listen",
        default="0.0.0.0",
        metavar="ip",
        help="Bind address to listen, default to 0.0.0.0",
    )
    parser.add_argument(
        "-sp",
        "--socks-port",
        type=int,
        help="Port accept connections from client",
        required=True,
    )
    parser.add_argument(
        "-bp",
        "--bridge-port",
        type=int,
        help="Port accept connections from SCF",
        required=True,
    )
    parser.add_argument("--user", type=str, help="Authentication username")
    parser.add_argument("--passwd", type=str, help="Authentication password")
    args = parser.parse_args()
    return args


def cancel_task(msg):
    print_time(f"[ERROR] {msg}")
    task = asyncio.current_task()
    task.cancel()
