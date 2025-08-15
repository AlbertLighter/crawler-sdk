package xhs

import "encoding/hex"

const (
	// Max32Bit is the maximum value for a 32-bit unsigned integer.
	Max32Bit uint32 = 0xFFFFFFFF
	// MaxSigned32Bit is the maximum value for a 32-bit signed integer.
	MaxSigned32Bit int32 = 0x7FFFFFFF

	// Base58Alphabet is the character set for Base58 encoding.
	Base58Alphabet = "NOPQRStuvwxWXYZabcyz012DEFTKLMdefghijkl4563GHIJBC7mnop89+/AUVqrsOPQefghijkABCDEFGuvwz0123456789xy"
	// StandardBase64Alphabet is the standard Base64 character set.
	StandardBase64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	// CustomBase64Alphabet is the custom Base64 character set used by XHS.
	CustomBase64Alphabet = "ZmserbBoHQtNP+wOcza/LpngG8yJq42KWYj0DSfdikx3VT16IlUAFM97hECvuRX5"
	// Base58Base is the base for Base58 encoding.
	Base58Base = 58
	// ByteSize is the size of a byte.
	ByteSize = 256

	// TimestampBytesCount is the length of the timestamp byte array.
	TimestampBytesCount = 16
	// TimestampXORKey is the XOR key for timestamp encoding.
	TimestampXORKey = 41
	// StartupTimeOffsetMin is the minimum offset for startup time.
	StartupTimeOffsetMin = 1000
	// StartupTimeOffsetMax is the maximum offset for startup time.
	StartupTimeOffsetMax = 4000

	// ExpectedHexLength is the expected length of the hex parameter.
	ExpectedHexLength = 32
	// OutputByteCount is the number of output bytes after processing.
	OutputByteCount = 8
	// HexChunkSize is the size of a hex character chunk.
	HexChunkSize = 2

	// RandomByteCount is the number of random bytes to generate.
	RandomByteCount = 4
	// FixedIntValue1 is a fixed integer value used in the payload.
	FixedIntValue1 = 15
	// FixedIntValue2 is another fixed integer value used in the payload.
	FixedIntValue2 = 1291

	// X3Prefix is the prefix for the x3 signature field.
	X3Prefix = "mns0101_"
	// XYSPrefix is the prefix for the final signature string.
	XYSPrefix = "XYS_"
)

var (
	// HexKey is the XOR key for the payload transformation.
	HexKey, _ = hex.DecodeString("af572b95ca65b2d9ec76bb5d2e97cb653299cc663399cc663399cce673399cce6733190c06030100000000008040209048241289c4e271381c0e0703018040a05028148ac56231180c0683c16030984c2693c964b259ac56abd5eaf5fafd7e3f9f4f279349a4d2e9743a9d4e279349a4d2e9f47a3d1e8f47239148a4d269341a8d4623110884422190c86432994ca6d3e974baddee773b1d8e47a35128148ac5623198cce6f3f97c3e1f8f47a3d168b45aad562b158ac5e2f1f87c3e9f4f279349a4d269b45aad56")

	// VersionBytes are the version identifier bytes.
	VersionBytes = []byte{119, 104, 96, 41}
	// FixedSeparatorBytes are fixed separator bytes.
	FixedSeparatorBytes = []byte{16, 0, 0, 0, 15, 5, 0, 0, 47, 1, 0, 0}

	// EnvStaticBytes are static bytes for the environment information.
	EnvStaticBytes = []byte{1, 249, 83, 102, 103, 201, 181, 131, 99, 94, 7, 68, 250, 132, 21}
)

// SignatureDataTemplate is the template for the signature data structure.
type SignatureDataTemplate struct {
	X0 string `json:"x0"`
	X1 string `json:"x1"`
	X2 string `json:"x2"`
	X3 string `json:"x3"`
	X4 string `json:"x4"`
}

// NewSignatureDataTemplate creates a new signature data template with default values.
func NewSignatureDataTemplate() *SignatureDataTemplate {
	return &SignatureDataTemplate{
		X0: "4.2.2",
		X1: "xhs-pc-web",
		X2: "Windows",
		X4: "object",
	}
}

// Lookup table for custom base64 encoding
var Lookup = []byte{
	'Z', 'm', 's', 'e', 'r', 'b', 'B', 'o', 'H', 'Q', 't', 'N', 'P', '+', 'w', 'O', 'c', 'z', 'a', '/', 'L', 'p', 'n',
	'g', 'G', '8', 'y', 'J', 'q', '4', '2', 'K', 'W', 'Y', 'j', '0', 'D', 'S', 'f', 'd', 'i', 'k', 'x', '3', 'V', 'T',
	'1', '6', 'I', 'l', 'U', 'A', 'F', 'M', '9', '7', 'h', 'E', 'C', 'v', 'u', 'R', 'X', '5',
}

// xn string for custom base64 encoding
const XN = "A4NjFqYu5wPHsO0XTdDgMa2r1ZQocVte9UJBvk6/7=yRnhISGKblCWi+LpfE8xzm3"
const XN64 = '=' // xn[64]

// ie array for custom CRC algorithm
var IE = []uint32{
	0, 1996959894, 3993919788, 2567524794, 124634137, 1886057615, 3915621685, 2657392035, 249268274, 2044508324,
	3772115230, 2547177864, 162941995, 2125561021, 3887607047, 2428444049, 498536548, 1789927666, 4089016648,
	2227061214, 450548861, 1843258603, 4107580753, 2211677639, 325883990, 1684777152, 4251122042, 2321926636,
	335633487, 1661365465, 4195302755, 2366115317, 997073096, 1281953886, 3579855332, 2724688242, 1006888145,
	1258607687, 3524101629, 2768942443, 901097722, 1119000684, 3686517206, 2898065728, 853044451, 1172266101,
	3705015759, 2882616665, 651767980, 1373503546, 3369554304, 3218104598, 565507253, 1454621731, 3485111705,
	3099436303, 671266974, 1594198024, 3322730930, 2970347812, 795835527, 1483230225, 3244367275, 3060149565,
	1994146192, 31158534, 2563907772, 4023717930, 1907459465, 112637215, 2680153253, 3904427059, 2013776290,
	251722036, 2517215374, 3775830040, 2137656763, 141376813, 2439277719, 3865271297, 1802195444, 476864866,
	2238001368, 4066508878, 1812370925, 453092731, 2181625025, 4111451223, 1706088902, 314042704, 2344532202,
	4240017532, 1658658271, 366619977, 2362670323, 4224994405, 1303535960, 984961486, 2747007092, 3569037538,
	1256170817, 1037604311, 2765210733, 3554079995, 1131014506, 879679996, 2909243462, 3663771856, 1141124467,
	855842277, 2852801631, 3708648649, 1342533948, 654459306, 3188396048, 3373015174, 1466479909, 544179635,
	3110523913, 3462522015, 1591671054, 702138776, 2966460450, 3352799412, 1504918807, 783551873, 3082640443,
	3233442989, 3988292384, 2596254646, 62317068, 1957810842, 3939845945, 2647816111, 81470997, 1943803523,
	3814918930, 2489596804, 225274430, 2053790376, 3826175755, 2466906013, 167816743, 2097651377, 4027552580,
	2265490386, 503444072, 1762050814, 4150417245, 2154129355, 426522225, 1852507879, 4275313526, 2312317920,
	282753626, 1742555852, 4189708143, 2394877945, 397917763, 1622183637, 3604390888, 2714866558, 953729732,
	1340076626, 3518719985, 2797360999, 1068828381, 1219638859, 3624741850, 2936675148, 906185462, 1090812512,
	3747672003, 2825379669, 829329135, 1181335161, 3412177804, 3160834842, 628085408, 1382605366, 3423369109,
	3138078467, 570562233, 1426400815, 3317316542, 2998733608, 733239954, 1555261956, 3268935591, 3050360625,
	752459403, 1541320221, 2607071920, 3965973030, 1969922972, 40735498, 2617837225, 3943577151, 1913087877,
	83908371, 2512341634, 3803740692, 2075208622, 213261112, 2463272603, 3855990285, 2094854071, 198958881,
	2262029012, 4057260610, 1759359992, 534414190, 2176718541, 4139329115, 1873836001, 414664567, 2282248934,
	4279200368, 1711684554, 285281116, 2405801727, 4167216745, 1634467795, 376229701, 2685067896, 3608007406,
	1308918612, 956543938, 2808555105, 3495958263, 1231636301, 1047427035, 2932959818, 3654703836, 1088359270,
	936918000, 2847714899, 3736837829, 1202900863, 817233897, 3183342108, 3401237130, 1404277552, 615818150,
	3134207493, 3453421203, 1423857449, 601450431, 3009837614, 3294710456, 1567103746, 711928724, 3020668471,
	3272380065, 1510334235, 755167117,
}
