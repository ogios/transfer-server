import socket

HOST = "127.0.0.1"
PORT = 15001

b = b""


def int_to_255(total) -> list[int]:
    l: list[int] = []
    while total >= 255:
        l.append(total % 255)
        total //= 255
    l.append(total)
    return l


def get_len(content: bytes) -> list[int]:
    l = int_to_255(len(content))
    l.append(255)
    return l


def add_bytes(old: bytes, add: bytes) -> bytes:
    tl = get_len(add)
    return old + bytes(tl) + add


def add_string(old: bytes, content: str) -> bytes:
    t = content.encode("utf-8")
    return add_bytes(old, t)


b = add_string(b, "text")
b = add_string(b, "分撒发声法是否会撒谎福安市法华寺的的撒的话分撒分撒")
print(b)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))
s.send(b)
s.close()
