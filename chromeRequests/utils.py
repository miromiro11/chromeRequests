import ctypes

def pullFromMem(bytes_) -> str:
    out = ctypes.string_at(bytes_)
    return out.decode('utf-8')

def checkError(json_):
    if json_.get("Error",False):
        raise Exception(json_['Error'])
    return True