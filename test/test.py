import socket

HOST = "127.0.0.1"
PORT = 15001

a = """Formal letter writing: block style vs. AMS style
Formal letters—like cover letters, business inquiries, and urgent notifications— are some of the most important letters you’ll ever have to write. Because they’re sometimes used as official documents, formal letters have a very precise structure and particular format. In fact, there are a few different “correct formats” to choose from.

The most common formats for formal letter writing are block style and American Mathematical Society, or AMS, style. In the example below, we use block style, specifically full block style, because it’s the most popular. Block style is characterized by all elements being aligned on the left margin of the page. This includes the first lines of paragraphs, which don’t use indentation. 

AMS is fairly similar, following many of the same rules as block style. There are a few differences, however, which we briefly cover after the next section. 

How to write a formal letter in block style
Step 1: Write the contact information and date 
All formal letters start with the contact information and date. In the full block style, this goes in the upper left-hand corner. 

First, as the sender, type your full name and address aligned to the left side, just as you would when addressing an envelope. This isn’t just a formality, but a useful inclusion so the recipient can easily find your contact information when they want to respond. 

If you’re writing on official company letterhead that already includes this information, you do not need to rewrite the contact information. 

After your address, skip a line and then add the date you’re writing the letter. 

Last, skip a line again and add the recipient’s name and full address. Feel free to include their job title below their name if it’s relevant. Leave a blank line after the contact information before writing the salutation. """
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
