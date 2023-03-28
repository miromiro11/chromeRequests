import json
from typing import OrderedDict
from .response import Response
from .headers import Headers
from .cookies import Cookies
from .utils import pull_from_mem, check_error
import ctypes

library = ctypes.cdll.LoadLibrary("./go/chromeRequests.so")

class Session:
    def __init__(self, proxy="", one_time=False):
        if library is None:
            raise Exception("Library is not loaded")
        self.__request = library.request
        self.__request.restype = ctypes.c_void_p
        self.cookies = Cookies()
        self.headers = Headers()

    def request(self, load):
        response = self.__request(load)
        json_ = json.loads(pull_from_mem(response))
        check_error(json_)
        toReturn = Response(json_)
        cookies = toReturn.cookies
        self.cookies.update(cookies)
        return toReturn

    def createPayload(self, request_type: str, url: str, **kwargs) -> dict:
        payload = OrderedDict({
            "requestType": request_type,
            "url": url,
            "parameters": {
                "url": url,
            }
        })
        allowed_get_params = ['proxy', 'headers', "cookies", "allow_redirects"]
        allowed_post_params = ['json', 'proxy',
                               'headers', 'data', "cookies", "allow_redirects"]
        params = allowed_get_params if request_type == "GET" else allowed_post_params
        for item in kwargs:
            if item not in params:
                raise Exception(f"{item} is not an accepted PARAM")
        payload['parameters']['headers'] = {
            k.lower(): v for k, v in kwargs.get("headers", {}).items()}
        payload['parameters']['headers'].update(self.headers.get_dict())
        payload['parameters']['proxy'] = kwargs.get("proxy", "")
        payload['parameters']['cookies'] = kwargs.get("cookies", {})
        payload['parameters']['redirects'] = kwargs.get(
            "allow_redirects", True)
        payload['parameters']['cookies'].update(self.cookies.get_dict())
        payload['parameters']['headerOrder'] = list(
            payload['parameters']['headers'].keys())
        if request_type == "POST" or request_type == "PUT":
            payload['parameters']['form'] = kwargs.get("data", {})
            payload['parameters']['json'] = json.dumps(
                kwargs.get("json", {}), sort_keys=True)
        return payload

    def get(self, url: str, **kwargs) -> Response:
        payload = self.createPayload("GET", url, **kwargs)
        toSend = json.dumps(
            payload, sort_keys=True).encode('utf-8')
        response = self.request(toSend)
        return response

    def post(self, url: str, **kwargs) -> Response:
        payload = self.createPayload("POST", url, **kwargs)
        print(json.dumps(
            payload, sort_keys=True).encode('utf-8'))
        response = self.request(json.dumps(
            payload, sort_keys=True).encode('utf-8'))
        return response

    def put(self, url: str, **kwargs) -> Response:
        payload = self.createPayload("PUT", url, **kwargs)
        response = self.request(json.dumps(
            payload, sort_keys=True).encode('utf-8'))
        return response

    def close(self):
        self.__closeSession(self.__uuid.encode('utf-8'))

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        if self.__uuid != "":
            self.close()
        return False
