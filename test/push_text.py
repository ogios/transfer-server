import socket

HOST = "127.0.0.1"
PORT = 15001

b = b""

def get_len(content: bytes) -> list[int]:
    l: list[int] = []
    total = len(content)
    while total >= 255:
        l.append(total%255)
        total //=255
    l.append(total)
    if len(l) > 8:
        raise Exception("Length over 8")
    while len(l) < 8:
        l.append(0)
    return l

def add_bytes(old: bytes, add: bytes) -> bytes:
    tl = get_len(add)
    return old + bytes(tl) + add


def add_string(old: bytes, content: str) -> bytes:
    t = content.encode()
    return add_bytes(old, t)

b = add_string(b, "text")
b = add_string(b, "a")
print(b)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))
s.send(b)
s.close()
