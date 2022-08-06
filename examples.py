import chromeRequests
import requests
import time

chromeRequests.load_library("./go/library.so")

session = chromeRequests.session()


def testCase(session):
    start = time.time()
    session.get("https://www.facebook.com")
    end = time.time()
    print("Time taken: ", end - start)


def sesionCreation(libary):
    start = time.time()
    libary.session()
    end = time.time()
    print("Time taken: ", end - start)


print("Testing session creation speed with requests")
sesionCreation(requests)
print("Testing session creation speed with chromeRequests")
sesionCreation(chromeRequests)

print("==========================================================")

requestsSession = requests.Session()
chromeRequestsSession = chromeRequests.session()
print("Testing session speed with requests")
testCase(requestsSession)
print("Testing session speed with chromeRequests")
testCase(chromeRequestsSession)
