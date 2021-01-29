// +build linux

package tc

import (
	"testing"

	"github.com/dennisafa/go-tc/internal/unix"
)

var (
	rate_1kbit_burst_40_mtu_9k = []uint32{0x00127a00, 0x0024f400, 0x00366e01, 0x0048e801, 0x005a6202,
		0x006cdc02, 0x007e5603, 0x0090d003, 0x00a24a04, 0x00b4c404, 0x00c63e05,
		0x00d8b805, 0x00ea3206, 0x00fcac06, 0x000e2707, 0x0020a107, 0x00321b08,
		0x00449508, 0x00560f09, 0x00688909, 0x007a030a, 0x008c7d0a, 0x009ef70a,
		0x00b0710b, 0x00c2eb0b, 0x00d4650c, 0x00e6df0c, 0x00f8590d, 0x000ad40d,
		0x001c4e0e, 0x002ec80e, 0x0040420f, 0x0052bc0f, 0x00643610, 0x0076b010,
		0x00882a11, 0x009aa411, 0x00ac1e12, 0x00be9812, 0x00d01213, 0x00e28c13,
		0x00f40614, 0x00068114, 0x0018fb14, 0x002a7515, 0x003cef15, 0x004e6916,
		0x0060e316, 0x00725d17, 0x0084d717, 0x00965118, 0x00a8cb18, 0x00ba4519,
		0x00ccbf19, 0x00de391a, 0x00f0b31a, 0x00022e1b, 0x0014a81b, 0x0026221c,
		0x00389c1c, 0x004a161d, 0x005c901d, 0x006e0a1e, 0x0080841e, 0x0092fe1e,
		0x00a4781f, 0x00b6f21f, 0x00c86c20, 0x00dae620, 0x00ec6021, 0x00feda21,
		0x00105522, 0x0022cf22, 0x00344923, 0x0046c323, 0x00583d24, 0x006ab724,
		0x007c3125, 0x008eab25, 0x00a02526, 0x00b29f26, 0x00c41927, 0x00d69327,
		0x00e80d28, 0x00fa8728, 0x000c0229, 0x001e7c29, 0x0030f629, 0x0042702a,
		0x0054ea2a, 0x0066642b, 0x0078de2b, 0x008a582c, 0x009cd22c, 0x00ae4c2d,
		0x00c0c62d, 0x00d2402e, 0x00e4ba2e, 0x00f6342f, 0x0008af2f, 0x001a2930,
		0x002ca330, 0x003e1d31, 0x00509731, 0x00621132, 0x00748b32, 0x00860533,
		0x00987f33, 0x00aaf933, 0x00bc7334, 0x00ceed34, 0x00e06735, 0x00f2e135,
		0x00045c36, 0x0016d636, 0x00285037, 0x003aca37, 0x004c4438, 0x005ebe38,
		0x00703839, 0x0082b239, 0x00942c3a, 0x00a6a63a, 0x00b8203b, 0x00ca9a3b,
		0x00dc143c, 0x00ee8e3c, 0x0000093d, 0x0012833d, 0x0024fd3d, 0x0036773e,
		0x0048f13e, 0x005a6b3f, 0x006ce53f, 0x007e5f40, 0x0090d940, 0x00a25341,
		0x00b4cd41, 0x00c64742, 0x00d8c142, 0x00ea3b43, 0x00fcb543, 0x000e3044,
		0x0020aa44, 0x00322445, 0x00449e45, 0x00561846, 0x00689246, 0x007a0c47,
		0x008c8647, 0x009e0048, 0x00b07a48, 0x00c2f448, 0x00d46e49, 0x00e6e849,
		0x00f8624a, 0x000add4a, 0x001c574b, 0x002ed14b, 0x00404b4c, 0x0052c54c,
		0x00643f4d, 0x0076b94d, 0x0088334e, 0x009aad4e, 0x00ac274f, 0x00bea14f,
		0x00d01b50, 0x00e29550, 0x00f40f51, 0x00068a51, 0x00180452, 0x002a7e52,
		0x003cf852, 0x004e7253, 0x0060ec53, 0x00726654, 0x0084e054, 0x00965a55,
		0x00a8d455, 0x00ba4e56, 0x00ccc856, 0x00de4257, 0x00f0bc57, 0x00023758,
		0x0014b158, 0x00262b59, 0x0038a559, 0x004a1f5a, 0x005c995a, 0x006e135b,
		0x00808d5b, 0x0092075c, 0x00a4815c, 0x00b6fb5c, 0x00c8755d, 0x00daef5d,
		0x00ec695e, 0x00fee35e, 0x00105e5f, 0x0022d85f, 0x00345260, 0x0046cc60,
		0x00584661, 0x006ac061, 0x007c3a62, 0x008eb462, 0x00a02e63, 0x00b2a863,
		0x00c42264, 0x00d69c64, 0x00e81665, 0x00fa9065, 0x000c0b66, 0x001e8566,
		0x0030ff66, 0x00427967, 0x0054f367, 0x00666d68, 0x0078e768, 0x008a6169,
		0x009cdb69, 0x00ae556a, 0x00c0cf6a, 0x00d2496b, 0x00e4c36b, 0x00f63d6c,
		0x0008b86c, 0x001a326d, 0x002cac6d, 0x003e266e, 0x0050a06e, 0x00621a6f,
		0x0074946f, 0x00860e70, 0x00988870, 0x00aa0271, 0x00bc7c71, 0x00cef671,
		0x00e07072, 0x00f2ea72, 0x00046573, 0x0016df73, 0x00285974, 0x003ad374,
		0x004c4d75, 0x005ec775, 0x00704176, 0x0082bb76, 0x00943577, 0x00a6af77,
		0x00b82978, 0x00caa378, 0x00dc1d79, 0x00ee9779, 0x0000127a}
	rate_1mbit_burst_100k = []uint32{0xe8030000, 0xd0070000, 0xb80b0000, 0xa00f0000, 0x88130000,
		0x70170000, 0x581b0000, 0x401f0000, 0x28230000, 0x10270000, 0xf82a0000,
		0xe02e0000, 0xc8320000, 0xb0360000, 0x983a0000, 0x803e0000, 0x68420000,
		0x50460000, 0x384a0000, 0x204e0000, 0x08520000, 0xf0550000, 0xd8590000,
		0xc05d0000, 0xa8610000, 0x90650000, 0x78690000, 0x606d0000, 0x48710000,
		0x30750000, 0x18790000, 0x007d0000, 0xe8800000, 0xd0840000, 0xb8880000,
		0xa08c0000, 0x88900000, 0x70940000, 0x58980000, 0x409c0000, 0x28a00000,
		0x10a40000, 0xf8a70000, 0xe0ab0000, 0xc8af0000, 0xb0b30000, 0x98b70000,
		0x80bb0000, 0x68bf0000, 0x50c30000, 0x38c70000, 0x20cb0000, 0x08cf0000,
		0xf0d20000, 0xd8d60000, 0xc0da0000, 0xa8de0000, 0x90e20000, 0x78e60000,
		0x60ea0000, 0x48ee0000, 0x30f20000, 0x18f60000, 0x00fa0000, 0xe8fd0000,
		0xd0010100, 0xb8050100, 0xa0090100, 0x880d0100, 0x70110100, 0x58150100,
		0x40190100, 0x281d0100, 0x10210100, 0xf8240100, 0xe0280100, 0xc82c0100,
		0xb0300100, 0x98340100, 0x80380100, 0x683c0100, 0x50400100, 0x38440100,
		0x20480100, 0x084c0100, 0xf04f0100, 0xd8530100, 0xc0570100, 0xa85b0100,
		0x905f0100, 0x78630100, 0x60670100, 0x486b0100, 0x306f0100, 0x18730100,
		0x00770100, 0xe87a0100, 0xd07e0100, 0xb8820100, 0xa0860100, 0x888a0100,
		0x708e0100, 0x58920100, 0x40960100, 0x289a0100, 0x109e0100, 0xf8a10100,
		0xe0a50100, 0xc8a90100, 0xb0ad0100, 0x98b10100, 0x80b50100, 0x68b90100,
		0x50bd0100, 0x38c10100, 0x20c50100, 0x08c90100, 0xf0cc0100, 0xd8d00100,
		0xc0d40100, 0xa8d80100, 0x90dc0100, 0x78e00100, 0x60e40100, 0x48e80100,
		0x30ec0100, 0x18f00100, 0x00f40100, 0xe8f70100, 0xd0fb0100, 0xb8ff0100,
		0xa0030200, 0x88070200, 0x700b0200, 0x580f0200, 0x40130200, 0x28170200,
		0x101b0200, 0xf81e0200, 0xe0220200, 0xc8260200, 0xb02a0200, 0x982e0200,
		0x80320200, 0x68360200, 0x503a0200, 0x383e0200, 0x20420200, 0x08460200,
		0xf0490200, 0xd84d0200, 0xc0510200, 0xa8550200, 0x90590200, 0x785d0200,
		0x60610200, 0x48650200, 0x30690200, 0x186d0200, 0x00710200, 0xe8740200,
		0xd0780200, 0xb87c0200, 0xa0800200, 0x88840200, 0x70880200, 0x588c0200,
		0x40900200, 0x28940200, 0x10980200, 0xf89b0200, 0xe09f0200, 0xc8a30200,
		0xb0a70200, 0x98ab0200, 0x80af0200, 0x68b30200, 0x50b70200, 0x38bb0200,
		0x20bf0200, 0x08c30200, 0xf0c60200, 0xd8ca0200, 0xc0ce0200, 0xa8d20200,
		0x90d60200, 0x78da0200, 0x60de0200, 0x48e20200, 0x30e60200, 0x18ea0200,
		0x00ee0200, 0xe8f10200, 0xd0f50200, 0xb8f90200, 0xa0fd0200, 0x88010300,
		0x70050300, 0x58090300, 0x400d0300, 0x28110300, 0x10150300, 0xf8180300,
		0xe01c0300, 0xc8200300, 0xb0240300, 0x98280300, 0x802c0300, 0x68300300,
		0x50340300, 0x38380300, 0x203c0300, 0x08400300, 0xf0430300, 0xd8470300,
		0xc04b0300, 0xa84f0300, 0x90530300, 0x78570300, 0x605b0300, 0x485f0300,
		0x30630300, 0x18670300, 0x006b0300, 0xe86e0300, 0xd0720300, 0xb8760300,
		0xa07a0300, 0x887e0300, 0x70820300, 0x58860300, 0x408a0300, 0x288e0300,
		0x10920300, 0xf8950300, 0xe0990300, 0xc89d0300, 0xb0a10300, 0x98a50300,
		0x80a90300, 0x68ad0300, 0x50b10300, 0x38b50300, 0x20b90300, 0x08bd0300,
		0xf0c00300, 0xd8c40300, 0xc0c80300, 0x98cc0300, 0x90d00300, 0x68d40300,
		0x60d80300, 0x48dc0300, 0x30e00300, 0x18e40300, 0x00e80300}
	rate_8kbit_burst_5kb_peakrate_12kbit_mpu_64_mtu_1464_drop = []uint32{0x40420f00, 0x40420f00, 0x40420f00, 0x40420f00, 0x40420f00,
		0x40420f00, 0x40420f00, 0x40420f00, 0x882a1100, 0xd0121300, 0x18fb1400,
		0x60e31600, 0xa8cb1800, 0xf0b31a00, 0x389c1c00, 0x80841e00, 0xc86c2000,
		0x10552200, 0x583d2400, 0xa0252600, 0xe80d2800, 0x30f62900, 0x78de2b00,
		0xc0c62d00, 0x08af2f00, 0x50973100, 0x987f3300, 0xe0673500, 0x28503700,
		0x70383900, 0xb8203b00, 0x00093d00, 0x48f13e00, 0x90d94000, 0xd8c14200,
		0x20aa4400, 0x68924600, 0xb07a4800, 0xf8624a00, 0x404b4c00, 0x88334e00,
		0xd01b5000, 0x18045200, 0x60ec5300, 0xa8d45500, 0xf0bc5700, 0x38a55900,
		0x808d5b00, 0xc8755d00, 0x105e5f00, 0x58466100, 0xa02e6300, 0xe8166500,
		0x30ff6600, 0x78e76800, 0xc0cf6a00, 0x08b86c00, 0x50a06e00, 0x98887000,
		0xe0707200, 0x28597400, 0x70417600, 0xb8297800, 0x00127a00, 0x48fa7b00,
		0x90e27d00, 0xd8ca7f00, 0x20b38100, 0x689b8300, 0xb0838500, 0xf86b8700,
		0x40548900, 0x883c8b00, 0xd0248d00, 0x180d8f00, 0x60f59000, 0xa8dd9200,
		0xf0c59400, 0x38ae9600, 0x80969800, 0xc87e9a00, 0x10679c00, 0x584f9e00,
		0xa037a000, 0xe81fa200, 0x3008a400, 0x78f0a500, 0xc0d8a700, 0x08c1a900,
		0x50a9ab00, 0x9891ad00, 0xe079af00, 0x2862b100, 0x704ab300, 0xb832b500,
		0x001bb700, 0x4803b900, 0x90ebba00, 0xd8d3bc00, 0x20bcbe00, 0x68a4c000,
		0xb08cc200, 0xf874c400, 0x405dc600, 0x8845c800, 0xd02dca00, 0x1816cc00,
		0x60fecd00, 0xa8e6cf00, 0xf0ced100, 0x38b7d300, 0x809fd500, 0xc887d700,
		0x1070d900, 0x5858db00, 0xa040dd00, 0xe828df00, 0x3011e100, 0x78f9e200,
		0xc0e1e400, 0x08cae600, 0x50b2e800, 0x989aea00, 0xe082ec00, 0x286bee00,
		0x7053f000, 0xb83bf200, 0x0024f400, 0x480cf600, 0x90f4f700, 0xd8dcf900,
		0x20c5fb00, 0x68adfd00, 0xb095ff00, 0xf87d0101, 0x40660301, 0x884e0501,
		0xd0360701, 0x181f0901, 0x60070b01, 0xa8ef0c01, 0xf0d70e01, 0x38c01001,
		0x80a81201, 0xc8901401, 0x10791601, 0x58611801, 0xa0491a01, 0xe8311c01,
		0x301a1e01, 0x78022001, 0xc0ea2101, 0x08d32301, 0x50bb2501, 0x98a32701,
		0xe08b2901, 0x28742b01, 0x705c2d01, 0xb8442f01, 0x002d3101, 0x48153301,
		0x90fd3401, 0xd8e53601, 0x20ce3801, 0x68b63a01, 0xb09e3c01, 0xf8863e01,
		0x406f4001, 0x88574201, 0xd03f4401, 0x18284601, 0x60104801, 0xa8f84901,
		0xf0e04b01, 0x38c94d01, 0x80b14f01, 0xc8995101, 0x10825301, 0x586a5501,
		0xa0525701, 0xe83a5901, 0x30235b01, 0x780b5d01, 0xc0f35e01, 0x08dc6001,
		0x50c46201, 0x98ac6401, 0xe0946601, 0x287d6801, 0x70656a01, 0xb84d6c01,
		0x00366e01, 0x481e7001, 0x90067201, 0xd8ee7301, 0x20d77501, 0x68bf7701,
		0xb0a77901, 0xf88f7b01, 0x40787d01, 0x88607f01, 0xd0488101, 0x18318301,
		0x60198501, 0xa8018701, 0xf0e98801, 0x38d28a01, 0x80ba8c01, 0xc8a28e01,
		0x108b9001, 0x58739201, 0xa05b9401, 0xe8439601, 0x302c9801, 0x78149a01,
		0xc0fc9b01, 0x08e59d01, 0x50cd9f01, 0x98b5a101, 0xe09da301, 0x2886a501,
		0x706ea701, 0xb856a901, 0x003fab01, 0x4827ad01, 0x900faf01, 0xd8f7b001,
		0x20e0b201, 0x68c8b401, 0xb0b0b601, 0xf898b801, 0x4081ba01, 0x8869bc01,
		0xd051be01, 0x183ac001, 0x6022c201, 0xa80ac401, 0xf0f2c501, 0x38dbc701,
		0x80c3c901, 0xc8abcb01, 0x1094cd01, 0x587ccf01, 0xa064d101, 0xe84cd301,
		0x3035d501, 0x781dd701, 0xc005d901, 0x08eeda01, 0x50d6dc01, 0x98bede01,
		0xe0a6e001, 0x288fe201, 0x7077e401, 0xb85fe601, 0x0048e801}
)

func TestGenerateRateTable(t *testing.T) {
	tests := map[string]struct {
		pol    *Policy
		expect []uint32
	}{
		"police rate 1kbit burst 40 mtu 9k": {
			pol: &Policy{
				Mtu: 9216,
				Rate: RateSpec{
					Rate:      125,
					Linklayer: unix.LINKLAYER_ETHERNET,
				},
			},
			expect: rate_1kbit_burst_40_mtu_9k,
		},
		"police rate 1mbit burst 100k": {
			pol: &Policy{
				Rate: RateSpec{
					Rate:      125000,
					Linklayer: unix.LINKLAYER_ETHERNET,
				},
			},
			expect: rate_1mbit_burst_100k,
		},
		"police rate 8kbit burst 5kb peakrate 12kbit mpu 64 mtu 1464 drop": {
			pol: &Policy{
				Mtu: 1464,
				PeakRate: RateSpec{
					Rate:      1000,
					Mpu:       64,
					Linklayer: unix.LINKLAYER_ETHERNET,
				},
			},
			expect: rate_8kbit_burst_5kb_peakrate_12kbit_mpu_64_mtu_1464_drop,
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := generateRateTable(testcase.pol)
			if err != nil {
				t.Fatalf("could not generate rate table")
			}
			for i := 0; i < 256; i++ {
				tmp := uint32(data[i*4+3]) | uint32(data[i*4+2])<<8 | uint32(data[i*4+1])<<16 | uint32(data[i*4+0])<<24
				if tmp != testcase.expect[i] {
					t.Fatalf("\n%d:\t0x%08x 0x%08x", i, tmp, testcase.expect[i])
				}
			}
		})
	}

}
