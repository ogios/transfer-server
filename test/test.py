from typing import Any


a = """123 34 116 111 116 97 108 34 58 49 44 34 100 97 116 97 34 58 91 123 34 116 121 112 101 34 58 49 44 34 116 105 109 101 34 58 49 54 57 55 52 51 57 48 52 56 49 51 51 44 34 105 100 34 58 34 56 56 110 120 88 34 44 34 100 97 116 97 34 58 34 229 136 134 230 146 146 229 143 145 229 163 176 230 179 149 230 152 175 229 144 166 228 188 154 230 146 146 232 176 142 231 166 143 229 174 137 229 184 130 230 179 149 229 141 142 229 175 186 231 154 132 231 154 132 230 146 146 231 154 132 232 175 157 229 136 134 230 146 146 229 136 134 230 146 146 34 125 93 125"""

b: Any = a.split(" ")
for i in range(len(b)):
    b[i] = int(b[i])

c = bytes(b)
print(c)
print(c.decode())
