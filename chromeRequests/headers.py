from typing import OrderedDict


class Headers:
    def __init__(self):
        self.__headers = OrderedDict()

    def update(self, headers: dict) -> None:
        self.__headers.update(headers)

    def get_dict(self) -> dict:
        ordered = OrderedDict()
        for key in self.__headers:
            ordered[key.lower()] = self.__headers[key]
        return ordered

    def clear(self) -> None:
        self.__headers = {}

    def __str__(self) -> str:
        return str(self.__headers)
