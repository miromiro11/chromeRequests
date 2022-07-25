import chromeRequests


chromeRequests.loadLibrary("./GoLangSource/library.so")

session = chromeRequests.chromeRequests.Session()
session.addCookies({
    "name": "cookie1",
    "value": "value1",
    "domain": "www.google.com",
    "path": "/",
    "secure": False,
})

session.removeCookies({
    "name": "cookie1",
})

resp = session.get("https://www.httpbin.org/get", data = {"test":"test"})

print(resp.text)
