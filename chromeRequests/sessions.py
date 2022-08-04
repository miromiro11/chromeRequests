import json
from .response import Response
from .headers import Headers
from .cookies import Cookies
from .utils import *

library = None
request = None


def load_library(path_to_lib):
    library = ctypes.cdll.LoadLibrary(path_to_lib)
    return library


class Session:
    def __init__(self, proxy="", one_time=False, library=None):
        if library is None:
            raise Exception("Library is not loaded")
        self.__session = library.createSession
        self.__session.restype = ctypes.c_void_p
        self.__request = library.request
        self.__request.restype = ctypes.c_void_p
        self.__changeProxy = library.changeProxy
        self.__changeProxy.restype = ctypes.c_void_p
        self.__uuid = json.loads(pull_from_mem(self.__session(proxy.encode('utf-8'))) if not one_time else "")['sessionId']
        self.__closeSession = library.closeSession
        self.__closeSession.restype = ctypes.c_void_p
        self.cookies = Cookies()
        self.headers = Headers()

    def set_proxy(self, proxy: str):
        payload = {
            "session": self.__uuid,
            "proxy": proxy,
        }
        load = json.dumps(payload).encode('utf-8')
        self.__changeProxy(load)

    def request(self, load):
        response = self.__request(load)
        json_ = json.loads(pull_from_mem(response))
        check_error(json_)
        toReturn = Response(json_)
        cookies = toReturn.cookies
        self.cookies.update(cookies)
        return toReturn

    def createPayload(self, request_type: str, url: str, **kwargs) -> dict:
        payload = {
            "session": self.__uuid,
            "requestType": request_type,
            "url": url,
            "parameters": {
                "url": url,
            }
        }
        allowed_get_params = ['proxy', 'headers', "cookies", "allow_redirects"]
        allowed_post_params = ['json', 'proxy', 'headers', 'data', "cookies", "allow_redirects"]
        params = allowed_get_params if request_type == "GET" else allowed_post_params
        for item in kwargs:
            if not item in params:
                raise Exception(f"{item} is not an accepted PARAM")
        payload['parameters']['headers'] = kwargs.get("headers", {})
        payload['parameters']['headers'].update(self.headers.get_dict())
        payload['parameters']['proxy'] = kwargs.get("proxy", "")
        payload['parameters']['cookies'] = kwargs.get("cookies", {})
        payload['parameters']['redirects'] = kwargs.get("allow_redirects", True)
        payload['parameters']['cookies'].update(self.cookies.get_dict())
        if request_type == "POST" or request_type == "PUT":
            payload['parameters']['form'] = kwargs.get("data", [])
            payload['parameters']['json'] = json.dumps(kwargs.get("json", {}))
        return payload

    def get(self, url: str, **kwargs) -> Response:
        payload = self.createPayload("GET", url, **kwargs)
        response = self.request(json.dumps(payload).encode('utf-8'))
        return response

    def post(self, url: str, **kwargs) -> Response:
        payload = self.createPayload("POST", url, **kwargs)
        response = self.request(json.dumps(payload).encode('utf-8'))
        return response

    def put(self, url: str, **kwargs) -> Response:
        payload = self.createPayload("PUT", url, **kwargs)
        response = self.request(json.dumps(payload).encode('utf-8'))
        return response

    def close(self):
        self.__closeSession(self.__uuid.encode('utf-8'))

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        if self.__uuid != "":
            self.close()
        return False
