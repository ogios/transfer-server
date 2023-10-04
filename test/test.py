import socket

HOST = "127.0.0.1"
PORT = 15001

a = "a"
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

t = "text".encode()
tl = get_len(t)
b = b + bytes(tl) + t

c = a.encode()
cl = get_len(c)
b = b + bytes(cl) + c

print(b)



s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.connect((HOST, PORT))
s.send(b)
s.close()
