package util_test

import (
	"github.com/Foxcapades/Argonaut/v0/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsValidShortFlag(t *testing.T) {
	Convey("util.IsValidShortFlag()", t, func() {
		for k, v := range validShortFlags {
			So(util.IsValidShortFlag(k), ShouldEqual, v)
		}
	})
}

func TestIsValidLongFlag(t *testing.T) {
	Convey("util.IsValidLongFlag()", t, func() {
		Convey("Leader character", func() {
			for k, v := range validShortFlags {
				So(util.IsValidLongFlag(string(k)), ShouldEqual, v)
			}
		})

		Convey("Followup characters", func() {
			for k, v := range validLongFlags {
				So(util.IsValidLongFlag(string([]byte{'a', k})), ShouldEqual, v)
			}
		})
	})
}

var validShortFlags = map[byte]bool{
	0: false, 1: false, 2: false, 3: false, 4: false, 5: false,
	6: false, 7: false, 8: false, 9: false, 10: false, 11: false,
	12: false, 13: false, 14: false, 15: false, 16: false, 17: false,
	18: false, 19: false, 20: false, 21: false, 22: false, 23: false,
	24: false, 25: false, 26: false, 27: false, 28: false, 29: false,
	30: false, 31: false, 32: false, 33: false, 34: false, 35: false,
	36: false, 37: false, 38: false, 39: false, 40: false, 41: false,
	42: false, 43: false, 44: false, 45: false, 46: false, 47: false,

	// 0-9
	48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true,
	55: true, 56: true, 57: true,

	// other
	58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false,

	// A-Z
	65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true,
	72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true,
	79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true,
	86: true, 87: true, 88: true, 89: true, 90: true,

	// other
	91: false, 92: false, 93: false, 94: false, 95: false, 96: false,

	// a-z
	97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true,
	104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true,
	111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true,
	118: true, 119: true, 120: true, 121: true, 122: true,

	// other
	123: false, 124: false, 125: false, 126: false, 127: false, 128: false,
	129: false, 130: false, 131: false, 132: false, 133: false, 134: false,
	135: false, 136: false, 137: false, 138: false, 139: false, 140: false,
	141: false, 142: false, 143: false, 144: false, 145: false, 146: false,
	147: false, 148: false, 149: false, 150: false, 151: false, 152: false,
	153: false, 154: false, 155: false, 156: false, 157: false, 158: false,
	159: false, 160: false, 161: false, 162: false, 163: false, 164: false,
	165: false, 166: false, 167: false, 168: false, 169: false, 170: false,
	171: false, 172: false, 173: false, 174: false, 175: false, 176: false,
	177: false, 178: false, 179: false, 180: false, 181: false, 182: false,
	183: false, 184: false, 185: false, 186: false, 187: false, 188: false,
	189: false, 190: false, 191: false, 192: false, 193: false, 194: false,
	195: false, 196: false, 197: false, 198: false, 199: false, 200: false,
	201: false, 202: false, 203: false, 204: false, 205: false, 206: false,
	207: false, 208: false, 209: false, 210: false, 211: false, 212: false,
	213: false, 214: false, 215: false, 216: false, 217: false, 218: false,
	219: false, 220: false, 221: false, 222: false, 223: false, 224: false,
	225: false, 226: false, 227: false, 228: false, 229: false, 230: false,
	231: false, 232: false, 233: false, 234: false, 235: false, 236: false,
	237: false, 238: false, 239: false, 240: false, 241: false, 242: false,
	243: false, 244: false, 245: false, 246: false, 247: false, 248: false,
	249: false, 250: false, 251: false, 252: false, 253: false, 254: false,
	255: false,
}

var validLongFlags = map[byte]bool{
	0: false, 1: false, 2: false, 3: false, 4: false, 5: false,
	6: false, 7: false, 8: false, 9: false, 10: false, 11: false,
	12: false, 13: false, 14: false, 15: false, 16: false, 17: false,
	18: false, 19: false, 20: false, 21: false, 22: false, 23: false,
	24: false, 25: false, 26: false, 27: false, 28: false, 29: false,
	30: false, 31: false, 32: false, 33: false, 34: false, 35: false,
	36: false, 37: false, 38: false, 39: false, 40: false, 41: false,
	42: false, 43: false, 44: false,

	// -
	45: true,

	// other
	46: false, 47: false,

	// 0-9
	48: true, 49: true, 50: true, 51: true, 52: true, 53: true, 54: true,
	55: true, 56: true, 57: true,

	// other
	58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false,

	// A-Z
	65: true, 66: true, 67: true, 68: true, 69: true, 70: true, 71: true,
	72: true, 73: true, 74: true, 75: true, 76: true, 77: true, 78: true,
	79: true, 80: true, 81: true, 82: true, 83: true, 84: true, 85: true,
	86: true, 87: true, 88: true, 89: true, 90: true,

	// other
	91: false, 92: false, 93: false, 94: false,

	// _
	95: true,

	// other
	96: false,

	// a-z
	97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true,
	104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true,
	111: true, 112: true, 113: true, 114: true, 115: true, 116: true, 117: true,
	118: true, 119: true, 120: true, 121: true, 122: true,

	// other
	123: false, 124: false, 125: false, 126: false, 127: false, 128: false,
	129: false, 130: false, 131: false, 132: false, 133: false, 134: false,
	135: false, 136: false, 137: false, 138: false, 139: false, 140: false,
	141: false, 142: false, 143: false, 144: false, 145: false, 146: false,
	147: false, 148: false, 149: false, 150: false, 151: false, 152: false,
	153: false, 154: false, 155: false, 156: false, 157: false, 158: false,
	159: false, 160: false, 161: false, 162: false, 163: false, 164: false,
	165: false, 166: false, 167: false, 168: false, 169: false, 170: false,
	171: false, 172: false, 173: false, 174: false, 175: false, 176: false,
	177: false, 178: false, 179: false, 180: false, 181: false, 182: false,
	183: false, 184: false, 185: false, 186: false, 187: false, 188: false,
	189: false, 190: false, 191: false, 192: false, 193: false, 194: false,
	195: false, 196: false, 197: false, 198: false, 199: false, 200: false,
	201: false, 202: false, 203: false, 204: false, 205: false, 206: false,
	207: false, 208: false, 209: false, 210: false, 211: false, 212: false,
	213: false, 214: false, 215: false, 216: false, 217: false, 218: false,
	219: false, 220: false, 221: false, 222: false, 223: false, 224: false,
	225: false, 226: false, 227: false, 228: false, 229: false, 230: false,
	231: false, 232: false, 233: false, 234: false, 235: false, 236: false,
	237: false, 238: false, 239: false, 240: false, 241: false, 242: false,
	243: false, 244: false, 245: false, 246: false, 247: false, 248: false,
	249: false, 250: false, 251: false, 252: false, 253: false, 254: false,
	255: false,
}
