# chromeRequests

A python package that mimics chromes TLS using golang

## Installation

Instll chromeRequests using pip

```bash
  pip install chromeRequests
```

## API Reference

```py
import chromeRequests

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
response.json()
cookies = response.cookies
headers = response.headers
status = response.status_code
text = response.text

#put request
response = chromeRequests.put(
    "example.com",
    headers = {
        "foo":"bar",
    },
    json = {"foo":"bar"},
    data = {"foo":"bar"}
)
#session

session = chromeRequests.Session()

session.setProxy("http://user:password@host:port")

session.headers.update({
  "foo":"bar"
})

session.headers.get_dict()

session.cookies.set(name,value)

session.cookies.get_dict()

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
