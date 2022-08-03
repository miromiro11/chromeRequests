class Cookies:
    def __init__(self):
        self.__cookies = {}

    def set(self, name: str, value: str) -> None:
        self.__cookies.update({
            name: value,
        })

    def update(self, cookies: dict):
        self.__cookies.update(cookies)

    def get_dict(self) -> dict:
        return self.__cookies

    def clear(self) -> None:
        self.__cookies = {}

    def __str__(self) -> str:
        return str(self.__cookies)
