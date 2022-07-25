import ctypes
import json
from urllib import request
import json as toJson


#load library from where the current file is located no matter where it is run from
library = None
request = None

def checkLibrary():
    if library == None:
        raise Exception("Library not loaded")

def loadLibrary(pathToLibrary):
    global library
    global request
    library = ctypes.cdll.LoadLibrary(pathToLibrary)
    request = library.request
    request.restype = ctypes.c_void_p


def pullFromMem(bytes_):
    out = ctypes.string_at(bytes_)
    return out.decode('utf-8')


class Response:
    def __init__(self, payload):
        payload = json.loads(payload)
        self.status_code = payload["StatusCode"]
        self.text = payload["Body"]
        self.cookies = payload["Cookies"]
        self.headers = payload["Headers"]

    def json(self):
        try:
            return toJson.loads(self.text)
        except Exception as e:
            raise(e)

    def __str__(self):
        return str(self.status_code)


def get(url, **kwargs):
    checkLibrary()
    payload = {
        "session": "",
        "requestType": "GET",
        "url": url,
        "paramters": {
            "url": url,
        }
    }
    allowedParams = ['proxy', 'headers', 'data']
    for item in kwargs:
        if not item in allowedParams:
            raise Exception(f"{item} is not an accepted PARAM")
    payload['paramters']['headers'] = kwargs.get("headers", [])
    response = request(json.dumps(payload).encode('utf-8'))
    return Response(pullFromMem(response))


def post(url, **kwargs):
    checkLibrary()
    payload = {
        "session": "",
        "requestType": "POST",
        "paramters": {
            "url": url,
        },
    }
    allowedParams = ['json', 'proxy', 'headers', 'data']
    for item in kwargs:
        if not item in allowedParams:
            raise Exception(f"{item} is not an accepted PARAM")
    payload['paramters']['headers'] = kwargs.get("headers", [])
    payload['paramters']['JSON'] = json.dumps(kwargs.get("json", []))
    response = request(json.dumps(payload).encode('utf-8'))
    return Response(pullFromMem(response))


def put(url, **kwargs):
    checkLibrary()
    payload = {
        "session": "",
        "requestType": "POST",
        "paramters": {
            "url": url,
        },
    }
    allowedParams = ['json', 'proxy', 'headers', 'data']
    for item in kwargs:
        if not item in allowedParams:
            raise Exception(f"{item} is not an accepted PARAM")
    payload['paramters']['headers'] = kwargs.get("headers", [])
    payload['paramters']['JSON'] = json.dumps(kwargs.get("json", []))
    response = request(json.dumps(payload).encode('utf-8'))
    return Response(pullFromMem(response))

class Session:
    def __init__(self, proxy=""):
        checkLibrary()
        self.session = library.createSession
        self.session.restype = ctypes.c_void_p
        self.request = library.request
        self.request.restype = ctypes.c_void_p
        self.changeProxy = library.changeProxy
        self.changeProxy.restype = ctypes.c_void_p
        self.addHeaders = library.addHeaders
        self.addHeaders.restype = ctypes.c_void_p
        self.removeHeaders = library.removeHeaders
        self.removeHeaders.restype = ctypes.c_void_p
        self.addCookies = library.addCookies
        self.addCookies.restype = ctypes.c_void_p
        self.removeCookies = library.removeCookies
        self.removeCookies.restype = ctypes.c_void_p
        self.uuid = pullFromMem(self.session(proxy.encode('utf-8')))

    def setProxy(self, proxy):
        payload = {
            "session": self.uuid,
            "proxy": proxy,
        }
        load = json.dumps(payload).encode('utf-8')
        self.changeProxy(load)

    def setCookies(self,cookies):
        payload = {
            "session": self.uuid,
            "cookies": cookies,
        }
        self.addCookies(json.dumps(payload).encode('utf-8'), self.uuid)

    def delCookies(self,cookies):
        payload = {
            "Session": self.uuid,
            "Cookies": cookies,
        }
        self.removeCookies(json.dumps(payload).encode('utf-8'), self.uuid)

    def setHeaders(self, headers):
        payload = {
            "Session": self.uuid,
            "Headers": headers,
        }
        self.addHeaders(json.dumps(payload).encode('utf-8'), self.uuid)

    def delHeaders(self, headers: dict):
        payload = {
            "Session": self.uuid,
            "Headers": headers,
        }
        self.removeHeaders(json.dumps(payload).encode('utf-8'), self.uuid)

    def get(self, url, **kwargs):
        payload = {
            "session": self.uuid,
            "requestType": "GET",
            "url": url,
            "paramters": {
                "url": url,
            }
        }
        allowedParams = ['proxy', 'headers', 'data']
        for item in kwargs:
            if not item in allowedParams:
                raise Exception(f"{item} is not an accepted PARAM")
        payload['paramters']['headers'] = kwargs.get("headers", [])
        payload['paramters']['proxy'] = kwargs.get("proxy", "")
        response = self.request(json.dumps(payload).encode('utf-8'))
        return Response(pullFromMem(response))

    def post(self, url, **kwargs):
        payload = {
            "session": self.uuid,
            "requestType": "POST",
            "paramters": {
                "url": url,
            },
        }
        allowedParams = ['json', 'proxy', 'headers', 'data']
        for item in kwargs:
            if not item in allowedParams:
                raise Exception(f"{item} is not an accepted PARAM")
        payload['paramters']['headers'] = kwargs.get("headers", {})
        payload['paramters']['Json'] = json.dumps(kwargs.get("json", {}))
        payload['paramters']['FORM'] = kwargs.get("data", [])
        response = self.request(json.dumps(payload).encode('utf-8'))
        return Response(pullFromMem(response))
    
    def put(self, url, **kwargs):
        payload = {
            "session": self.uuid,
            "requestType": "PUT",
            "paramters": {
                "url": url,
            },
        }
        allowedParams = ['json', 'proxy', 'headers', 'data']
        for item in kwargs:
            if not item in allowedParams:
                raise Exception(f"{item} is not an accepted PARAM")
        payload['paramters']['headers'] = kwargs.get("headers", {})
        payload['paramters']['Json'] = json.dumps(kwargs.get("json", {}))
        payload['paramters']['FORM'] = kwargs.get("data", [])
        response = self.request(json.dumps(payload).encode('utf-8'))
        return Response(pullFromMem(response))

