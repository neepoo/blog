# pip install pycryptodome
from Crypto.Hash import TupleHash128


def foo1():
    hd = TupleHash128.new(digest_bytes=16)
    hd.update(b'foo')
    hd.update(b'10')
    hd.update(b'10')
    print(hd.hexdigest())


def foo2():
    hd = TupleHash128.new(digest_bytes=16)
    hd.update(b'foo')
    hd.update(b'1')
    hd.update(b'010')
    print(hd.hexdigest())


if __name__ == "__main__":
    foo1()
    foo2()
