import ctypes

def pullFromMem(bytes_) -> str:
    out = ctypes.string_at(bytes_)
    return out.decode('utf-8')