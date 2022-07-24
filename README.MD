# chromeRequests

A python package that mimics chromes TLS using golang

## Installation

Install my-project with npm

```bash
  pip install chromeRequests
```

## API Reference

```py
from chromeRequests import chromeRequests

chromeRequests.loadLibrary("./library.so") #this is necessary, download from GoLangSource folder.

#get request

response = chromeRequests.get(
    "example.com",
    headers = {
        "foo":"bar",
    },
    proxy = "http://user:password@host:port" #never tested with passwordless.
)

#post request

response = chromeRequests.post(
    "example.com",
    headers = {
        "foo":"bar",
    },
    json = {"foo":"bar"},
    data = {"foo":"bar"}
)

cookies = response.cookies
headers = response.headers
status = response.status_code
text = response.text

#session

session = chromeRequests.Session()

#sessions functions

session.setHeaders({"foo":"bar"})

session.setProxy("http://user:password@host:port")

session.delHeaders({"test":""})

```

## Acknowledgements

- [Mimic By Sauce](https://github.com/saucesteals/mimic)

## Authors

- [@Miro](https://github.com/miromiro11)

## To do

- [ ] Speed Improvement
- [ ] Same syntax as pythons [Request Library](https://pypi.org/project/requests/)

## Contributing

Contributions are always welcome!

Please adhere to this project's `code of conduct`.
