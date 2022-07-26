import chromeRequests
import requests
import time

chromeRequests.loadLibrary("./GoLangSource/library.so")     

session = chromeRequests.Session()

def testCase(session):
    start = time.time()
    session.get("https://www.facebook.com")
    end = time.time()
    print("Time taken: ", end - start)

def sesionCreation(libary):
    start = time.time()
    libary.Session()
    end = time.time()
    print("Time taken: ", end - start)

print("Testing session creation speed with requests")
sesionCreation(requests)
print("Testing session creation speed with chromeRequests")
sesionCreation(chromeRequests)

print("==========================================================")

requestsSession = requests.Session()
chromeRequestsSession = chromeRequests.Session()
print("Testing session speed with requests")
testCase(requestsSession)
print("Testing session speed with chromeRequests")
testCase(chromeRequestsSession)
