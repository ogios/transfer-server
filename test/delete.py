import socket
import time

from util import add_string, add_bytes, int_to_255

HOST = "127.0.0.1"
PORT = 15001

b = b""
b = add_string(b, "delete")
b = add_string(b, "LZDNA")
print(b)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((HOST, PORT))
s.send(b)
while 1:
    msg = s.recv(1024)
    if msg:
        print(msg)
        break
    time.sleep(1)
s.close()
