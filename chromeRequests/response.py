import json as toJson

class Response:
    def __init__(self, payload):
        self.status_code = payload["StatusCode"]
        self.text = payload["Body"]
        self.cookies = payload["Cookies"]
        self.headers = payload["Headers"]
        self.url = payload["Url"]

    def json(self):
        try:
            return toJson.loads(self.text)
        except Exception as e:
            raise(e)

    def __str__(self):
        return str(self.status_code)
