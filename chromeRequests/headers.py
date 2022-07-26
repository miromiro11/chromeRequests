

class Headers:
    def __init__(self):
        self.__headers = {}
    
    def update(self,headers: dict) -> None: 
        self.__headers.update(headers)
    
    def get_dict(self) -> dict:
        return self.__headers

    def clear(self) -> None:
        self.__headers = {}

    def __str__(self) -> str:
        return str(self.__headers)
