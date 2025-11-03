# 19-00036 por lo tanto, x=0, y=3 , z=6
def suspenso(a, b):
    if b == []:
        yield a
    else:
        yield a + b[0]
        for x in suspenso(b[0], b[1:]):
            yield x

for x in suspenso (9, [0,3,6]):
    print(x)

# 19-00036 por lo tanto, x=0, y=3 , z=6
def misterio(n):
    if n == 0:
        yield [1]
    else:
        for x in misterio(n-1):
            r = []
        for y in suspenso(0, x):
            r = [*r, y]
        yield r

for x in misterio(5):
 print(x)