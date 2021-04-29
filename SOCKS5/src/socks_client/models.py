import asyncio
from typing import Union
from collections import OrderedDict

import aiohttp


class Conn:
    def __init__(
        self,
        role: str,
        reader: asyncio.StreamReader,
        writer: asyncio.StreamWriter,
    ) -> None:
        self.target = None
        self.role = role
        self.reader = reader
        self.writer = writer

    async def read(self, size: int):
        return await self.reader.read(size)

    async def write(self, data: Union[str, bytes]):
        self.writer.write(data)
        await self.writer.drain()

    def close(self):
        self.writer.close()


class LRUDict(OrderedDict):
    def __init__(self, capacity):
        self.capacity = capacity
        self.cache = OrderedDict()

    def get(self, key):
        value = self.cache.pop(key)
        self.cache[key] = value
        return value

    def set(self, key, value):
        if key in self.cache:
            self.cache.pop(key)
        elif len(self.cache) == self.capacity:
            self.cache.popitem(last=True)
        self.cache[key] = value


class Request:
    def __init__(self):
        self._session = None

    async def init_session(self):
        self._session = aiohttp.ClientSession()

    async def request(self, method, url, bypass_cf=False, **kwargs):
        await self._session.request(method=method, url=url, **kwargs)

    async def post(self, url, **kwargs):
        return await self.request("POST", url, **kwargs)

    async def close(self):
        await self._session.close()


http = Request()
uid_socket = LRUDict(150)
