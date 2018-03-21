import sys
import base64
import struct
from pyasn1.type import univ
from pyasn1.codec.der import encoder as der_encoder

keydata = base64.b64decode(open('test_rsa.pub').read().split(None)[1])
keydata = base64.b64decode(open('test_rsa.pub').read().split(None)[1])

parts = []

while keydata:
        dlen = struct.unpack('>I', keydata[:4])[0]
        data, keydata = keydata[4:dlen+4], keydata[4+dlen:]
        parts.append(data)
#print(parts)
#print(parts[1])
e_val = eval('0x' + ''.join(['%02X' % struct.unpack('B', x.to_bytes(1, 'big'))[0] for x in parts[1]]))
n_val = eval('0x' + ''.join(['%02X' % struct.unpack('B', x.to_bytes(1, 'big'))[0] for x in parts[2]]))

#print(e_val)
#print(n_val)

pkcs1_seq = univ.Sequence()
pkcs1_seq.setComponentByPosition(0, univ.Integer(n_val))
pkcs1_seq.setComponentByPosition(1, univ.Integer(e_val))

print('-----BEGIN RSA PUBLIC KEY-----')
print(base64.encodestring(der_encoder.encode(pkcs1_seq)).decode())
print('-----END RSA PUBLIC KEY-----')