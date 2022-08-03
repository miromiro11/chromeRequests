import ctypes


def pull_from_mem(bytes_) -> str:
    out = ctypes.string_at(bytes_)
    return out.decode('utf-8')


def check_error(json_):
    if json_.get("Error", False):
        raise Exception(json_['Error'])
    return True
