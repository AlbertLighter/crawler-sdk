'''
Python implementation of the sign.js logic.
'''
import time
import math


class SignGenerator:
    '''
    Translates the signing logic from sign.js to Python.
    It can be configured with fixed values for timestamp and random numbers
    to allow for deterministic testing against the JavaScript implementation.
    '''

    def __init__(self, fixed_timestamp=None, random_values=None):
        self._fixed_timestamp = fixed_timestamp
        self._random_values = random_values if random_values is not None else [0.123, 0.456, 0.789]

    def _get_timestamp(self):
        if self._fixed_timestamp is not None:
            return self._fixed_timestamp
        return int(time.time() * 1000)

    def rc4_encrypt(self, plaintext, key):
        s = list(range(256))
        j = 0
        key_bytes = key.encode('latin-1') if isinstance(key, str) else bytes(key)
        plaintext_bytes = plaintext.encode('latin-1') if isinstance(plaintext, str) else bytes(plaintext)

        for i in range(256):
            j = (j + s[i] + key_bytes[i % len(key_bytes)]) % 256
            s[i], s[j] = s[j], s[i]

        i = 0
        j = 0
        cipher = []
        for k in range(len(plaintext_bytes)):
            i = (i + 1) % 256
            j = (j + s[i]) % 256
            s[i], s[j] = s[j], s[i]
            t = (s[i] + s[j]) % 256
            cipher.append(s[t] ^ plaintext_bytes[k])
        return "".join(map(chr, cipher))

    def _le(self, e, r):
        r %= 32
        return ((e << r) | (e >> (32 - r))) & 0xFFFFFFFF

    def _de(self, e):
        if 0 <= e < 16: return 2043430169
        if 16 <= e < 64: return 2055708042
        raise ValueError("invalid j for constant Tj")

    def _pe(self, e, r, t, n):
        if 0 <= e < 16: return (r ^ t ^ n) & 0xFFFFFFFF
        if 16 <= e < 64: return ((r & t) | (r & n) | (t & n)) & 0xFFFFFFFF
        raise ValueError('invalid j for bool function FF')

    def _he(self, e, r, t, n):
        if 0 <= e < 16: return (r ^ t ^ n) & 0xFFFFFFFF
        if 16 <= e < 64: return ((r & t) | (~r & n)) & 0xFFFFFFFF
        raise ValueError('invalid j for bool function GG')

    class SM3:
        def __init__(self, parent):
            self._parent = parent
            self.reg = [0] * 8
            self.chunk = bytearray()
            self.size = 0
            self.reset()

        def reset(self):
            self.reg = [1937774191, 1226093241, 388252375, 3666478592, 2842636476, 372324522, 3817729613, 2969243214]
            self.chunk = bytearray()
            self.size = 0

        def write(self, data):
            if isinstance(data, str):
                a = data.encode('utf-8')
            else:
                a = bytes(data)
            self.size += len(a)
            self.chunk.extend(a)

        def sum(self, data=None, t='bytes'):
            if data is not None:
                self.reset()
                self.write(data)
            
            chunk_copy = self.chunk[:]
            self._fill(chunk_copy)
            
            for i in range(0, len(chunk_copy), 64):
                self._compress(chunk_copy[i:i+64])

            if t == 'hex':
                return "".join([f'{val:08x}' for val in self.reg])
            else:
                i = []
                for val in self.reg:
                    i.extend(val.to_bytes(4, 'big'))
                return i

        def _compress(self, t):
            r = [0] * 68
            for i in range(16):
                r[i] = int.from_bytes(bytes(t[i*4:(i+1)*4]), 'big')

            for n in range(16, 68):
                a = r[n - 16] ^ r[n - 9] ^ self._parent._le(r[n - 3], 15)
                a = a ^ self._parent._le(a, 15) ^ self._parent._le(a, 23)
                r[n] = (a ^ self._parent._le(r[n - 13], 7) ^ r[n - 6]) & 0xFFFFFFFF

            w_prime = [r[j] ^ r[j + 4] for j in range(64)]

            i_reg = self.reg[:]
            for c in range(64):
                o = self._parent._le(i_reg[0], 12) + i_reg[4] + self._parent._le(self._parent._de(c), c)
                o = self._parent._le(o & 0xFFFFFFFF, 7)
                s = (o ^ self._parent._le(i_reg[0], 12)) & 0xFFFFFFFF
                
                u = self._parent._pe(c, i_reg[0], i_reg[1], i_reg[2])
                u = (u + i_reg[3] + s + w_prime[c]) & 0xFFFFFFFF
                
                b = self._parent._he(c, i_reg[4], i_reg[5], i_reg[6])
                b = (b + i_reg[7] + o + r[c]) & 0xFFFFFFFF
                
                i_reg[3] = i_reg[2]
                i_reg[2] = self._parent._le(i_reg[1], 9)
                i_reg[1] = i_reg[0]
                i_reg[0] = u
                i_reg[7] = i_reg[6]
                i_reg[6] = self._parent._le(i_reg[5], 19)
                i_reg[5] = i_reg[4]
                i_reg[4] = (b ^ self._parent._le(b, 9) ^ self._parent._le(b, 17)) & 0xFFFFFFFF

            for l in range(8):
                self.reg[l] = (self.reg[l] ^ i_reg[l]) & 0xFFFFFFFF

        def _fill(self, chunk):
            a = self.size * 8
            chunk.append(128)
            
            while (len(chunk) % 64) != 56:
                chunk.append(0)
            
            chunk.extend(a.to_bytes(8, 'big'))

    def get_sm3_instance(self):
        return self.SM3(self)

    def result_encrypt(self, long_str, num=None):
        s_obj = {
            "s0": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
            "s1": "Dkdpgh4ZKsQB80/Mfvw36XI1R25+WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
            "s2": "Dkdpgh4ZKsQB80/Mfvw36XI1R25-WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
            "s3": "ckdp1h4ZKsUB80/Mfvw36XIgR25+WQAlEi7NLboqYTOPuzmFjJnryx9HVGDaStCe",
            "s4": "Dkdpgh2ZmsQB80/MfvV36XI1R45-WUAlEixNLwoqYTOPuzKFjJnry79HbGcaStCe"
        }
        alphabet = s_obj[num]
        result = []
        long_str_bytes = long_str.encode('latin-1')
        
        for i in range(0, len(long_str_bytes), 3):
            chunk = long_str_bytes[i:i+3]
            b1 = chunk[0]
            b2 = chunk[1] if len(chunk) > 1 else 0
            b3 = chunk[2] if len(chunk) > 2 else 0
            
            long_int = (b1 << 16) | (b2 << 8) | b3

            result.append(alphabet[(long_int >> 18) & 0x3F])
            result.append(alphabet[(long_int >> 12) & 0x3F])
            
            if len(chunk) > 1:
                result.append(alphabet[(long_int >> 6) & 0x3F])
            if len(chunk) > 2:
                result.append(alphabet[long_int & 0x3F])

        return "".join(result)

    def _gener_random(self, random_val, option):
        random_val = int(random_val)
        return [
            (random_val & 255 & 170) | (option[0] & 85),
            (random_val & 255 & 85) | (option[0] & 170),
            ((random_val >> 8) & 255 & 170) | (option[1] & 85),
            ((random_val >> 8) & 255 & 85) | (option[1] & 170),
        ]

    def generate_random_str(self):
        random_str_list = []
        random_str_list.extend(self._gener_random(self._random_values[0] * 10000, [3, 45]))
        random_str_list.extend(self._gener_random(self._random_values[1] * 10000, [1, 0]))
        random_str_list.extend(self._gener_random(self._random_values[2] * 10000, [1, 5]))
        return "".join(map(chr, random_str_list))

    def generate_rc4_bb_str(self, url_search_params, user_agent, window_env_str, suffix="cus", Arguments=[0, 1, 14]):
        start_time = self._get_timestamp()

        # url_search_params_list
        sm3_1 = self.get_sm3_instance()
        sm3_1.write(url_search_params + suffix)
        hash1 = sm3_1.sum()
        url_search_params_list = self.get_sm3_instance().sum(hash1)

        # cus
        sm3_2 = self.get_sm3_instance()
        sm3_2.write(suffix)
        hash2 = sm3_2.sum()
        cus = self.get_sm3_instance().sum(hash2)

        # ua
        rc4_key_bytes = [int(0.00390625), 1, Arguments[2]]
        ua_rc4 = self.rc4_encrypt(user_agent, bytes(rc4_key_bytes))
        ua_encrypted = self.result_encrypt(ua_rc4, "s3")
        ua = self.get_sm3_instance().sum(ua_encrypted)

        end_time = self._get_timestamp()

        b = {}
        b[8] = 3
        b[10] = end_time
        b[15] = {
            "aid": 6383,
            "pageId": 6241,
            "boe": False,
            "ddrt": 7,
            "paths": {
                "include": [
                    {},
                    {},
                    {},
                    {},
                    {},
                    {},
                    {}
                ],
                "exclude": []
            },
            "track": {
                "mode": 0,
                "delay": 300,
                "paths": []
            },
            "dump": True,
            "rpU": ""
        }
        b[16] = start_time
        b[18] = 44

        b[20] = (b[16] >> 24) & 255
        b[21] = (b[16] >> 16) & 255
        b[22] = (b[16] >> 8) & 255
        b[23] = b[16] & 255
        b[24] = (b[16] // (256**4)) & 255
        b[25] = (b[16] // (256**5)) & 255

        b[26] = (Arguments[0] >> 24) & 255
        b[27] = (Arguments[0] >> 16) & 255
        b[28] = (Arguments[0] >> 8) & 255
        b[29] = Arguments[0] & 255

        b[30] = (Arguments[1] // 256) & 255
        b[31] = (Arguments[1] % 256) & 255
        b[32] = (Arguments[1] >> 24) & 255
        b[33] = (Arguments[1] >> 16) & 255

        b[34] = (Arguments[2] >> 24) & 255
        b[35] = (Arguments[2] >> 16) & 255
        b[36] = (Arguments[2] >> 8) & 255
        b[37] = Arguments[2] & 255

        b[38] = url_search_params_list[21]
        b[39] = url_search_params_list[22]

        b[40] = cus[21]
        b[41] = cus[22]

        b[42] = ua[23]
        b[43] = ua[24]

        b[44] = (b[10] >> 24) & 255
        b[45] = (b[10] >> 16) & 255
        b[46] = (b[10] >> 8) & 255
        b[47] = b[10] & 255
        b[48] = b[8]
        b[49] = (b[10] // (256**4)) & 255
        b[50] = (b[10] // (256**5)) & 255

        b[52] = (b[15]['pageId'] >> 24) & 255
        b[53] = (b[15]['pageId'] >> 16) & 255
        b[54] = (b[15]['pageId'] >> 8) & 255
        b[55] = b[15]['pageId'] & 255

        b[57] = b[15]['aid'] & 255
        b[58] = (b[15]['aid'] >> 8) & 255
        b[59] = (b[15]['aid'] >> 16) & 255
        b[60] = (b[15]['aid'] >> 24) & 255

        window_env_list = list(window_env_str.encode('utf-8'))
        b[64] = len(window_env_list)
        b[65] = b[64] & 255
        b[66] = (b[64] >> 8) & 255

        b[69] = 0
        b[70] = b[69] & 255
        b[71] = (b[69] >> 8) & 255

        b[72] = (b[18] ^ b[20] ^ b[26] ^ b[30] ^ b[38] ^ b[40] ^ b[42] ^ b[21] ^ b[27] ^ b[31] ^ b[35] ^ b[39] ^ b[41] ^ b[43] ^ b[22] ^
                 b[28] ^ b[32] ^ b[36] ^ b[23] ^ b[29] ^ b[33] ^ b[37] ^ b[44] ^ b[45] ^ b[46] ^ b[47] ^ b[48] ^ b[49] ^ b[50] ^ b[24] ^
                 b[25] ^ b[52] ^ b[53] ^ b[54] ^ b[55] ^ b[57] ^ b[58] ^ b[59] ^ b[60] ^ b[65] ^ b[66] ^ b[70] ^ b[71])

        bb = [
            b[18], b[20], b[52], b[26], b[30], b[34], b[58], b[38], b[40], b[53], b[42], b[21], b[27], b[54], b[55], b[31],
            b[35], b[57], b[39], b[41], b[43], b[22], b[28], b[32], b[60], b[36], b[23], b[29], b[33], b[37], b[44], b[45],
            b[59], b[46], b[47], b[48], b[49], b[50], b[24], b[25], b[65], b[66], b[70], b[71]
        ]
        bb.extend(window_env_list)
        bb.append(b[72])
        
        bb_str = "".join(map(chr, bb))
        return self.rc4_encrypt(bb_str, "\x79")

    def sign(self, url_search_params, user_agent, arguments):
        result_str = self.generate_random_str() + self.generate_rc4_bb_str(
            url_search_params,
            user_agent,
            "1536|747|1536|834|0|30|0|0|1536|834|1536|864|1525|747|24|24|Win32",
            "cus",
            arguments
        )
        return self.result_encrypt(result_str, "s4") + "="

    def sign_datail(self, params, userAgent):
        return self.sign(params, userAgent, [0, 1, 14])

    def sign_reply(self, params, userAgent):
        return self.sign(params, userAgent, [0, 1, 8])
