from . import sessions 
import json
from .response import Response
from .utils import pullFromMem
import ctypes

globalLibrary = None

def loadLibrary(pathToLibrary) -> None:
    global globalLibrary
    globalLibrary = ctypes.cdll.LoadLibrary(pathToLibrary)
    print("Library Loaded Successfully")

def get(url: str, **kwargs) -> Response:
    with sessions.Session(oneTime = True, library = globalLibrary) as session:
        response = session.get(url, **kwargs)
        return response

def post(url: str, **kwargs) -> Response:
    with sessions.Session(oneTime = True, library = globalLibrary) as session:
        response = session.post(url, **kwargs)
        return response

def put(url: str, **kwargs) -> Response:
    with sessions.Session(oneTime = True, library = globalLibrary) as session:
        response = session.put(url, **kwargs)
        return response
    
def Session(proxy="") -> sessions.Session:
    return sessions.Session(proxy,library = globalLibrary)