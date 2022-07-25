import chromeRequests


chromeRequests.loadLibrary("./GoLangSource/library.so")

session = chromeRequests.chromeRequests.Session()
session.setCookies({
    "name": "cookie1",
    "value": "value1",
    "domain": "www.google.com",
    "path": "/",
    "secure": False,
})

session.delCookies({
    "name": "cookie1",
})

resp = session.get("https://www.httpbin.org/get", data = {"test":"test"})

print(resp.text)
