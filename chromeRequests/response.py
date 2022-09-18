import json as toJson


class Response:
    def __init__(self, payload):
        self.status_code = payload.get("statusCode", 0)
        self.text = payload.get("body", "")
        self.cookies = payload.get("cookies", {})
        self.headers = payload.get("headers", {})
        self.url = payload.get("url", "")

    def json(self):
        try:
            return toJson.loads(self.text)
        except Exception as exception:
            raise exception

    def __str__(self):
        return str(self.status_code)
