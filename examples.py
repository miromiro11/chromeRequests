import chromeRequests


chromeRequests.loadLibrary("./GoLangSource/library.so")

session = chromeRequests.chromeRequests.Session()

resp = session.put("https://www.httpbin.org/put", data = {"test":"test"})

