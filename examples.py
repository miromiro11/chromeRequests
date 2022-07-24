import chromeRequests


chromeRequests.loadLibrary("./GoLangSource/library.so")

session = chromeRequests.chromeRequests.Session()


resp = session.get("https://www.google.com")
print(resp.text)
