import socket
import time

from util import add_string

HOST = "127.0.0.1"
PORT = 15002
ADDR = (HOST, PORT)

b = "sub".encode("utf-8")
# b = add_string(b, "sub")
print(b)

s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.sendto(b, ADDR)
while 1:
    recv = s.recvfrom(1024)
    msg = recv[0]
    if msg == b"ok":
        print("subscribed")
    elif len(msg) == 1:
        if msg[0] == 1:
            print("ping")
            s.sendto(bytes([1]), ADDR)
        elif msg[0] == 2:
            print("notify")
    else:
        raise Exception(f"wtf? {msg}")

    # if msg:
    #     print(msg)
    #     break
    # time.sleep(1)
s.close()
