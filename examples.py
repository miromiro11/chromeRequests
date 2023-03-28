import chromeRequests
from pprint import pprint

session = chromeRequests.session()

response = session.post(
    "https://tls.peet.ws/api/all",
    json = {
        "test":{
            "test":1123
        },
    },
    headers = {
    "test1":"1",
    "test2":"2",
    }
)


pprint(response.json())