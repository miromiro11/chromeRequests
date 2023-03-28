from . import sessions
from .response import Response
import ctypes

globalLibrary = None


def load_library(path_to_lib):
    global globalLibrary
    globalLibrary = ctypes.cdll.LoadLibrary(path_to_lib)
    print("Library Loaded Successfully")


def get(url: str, **kwargs) -> Response:
    with sessions.Session(one_time=True) as session:
        return session.get(url, **kwargs)


def post(url: str, **kwargs) -> Response:
    with sessions.Session(one_time=True) as session:
        return session.post(url, **kwargs)


def put(url: str, **kwargs) -> Response:
    with sessions.Session(one_time=True) as session:
        return session.put(url, **kwargs)


def session(proxy="") -> sessions.Session:
    return sessions.Session(proxy)
