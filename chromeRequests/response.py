import json as toJson


class Response:
    def __init__(self, payload):
        self.status_code = payload["statusCode"]
        self.text = payload["body"]
        self.cookies = payload["cookies"]
        self.headers = payload["headers"]
        self.url = payload["url"]

    def json(self):
        try:
            return toJson.loads(self.text)
        except Exception as exception:
            raise exception

    def __str__(self):
        return str(self.status_code)
