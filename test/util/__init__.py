def int_to_255(total) -> list[int]:
    l: list[int] = []
    while total >= 255:
        l.append(total % 255)
        total //= 255
    l.append(total)
    l.append(255)
    return l


def get_len(content: bytes) -> list[int]:
    return int_to_255(len(content))


def add_bytes(old: bytes, add: bytes) -> bytes:
    tl = get_len(add)
    return old + bytes(tl) + add


def add_string(old: bytes, content: str) -> bytes:
    t = content.encode()
    return add_bytes(old, t)
