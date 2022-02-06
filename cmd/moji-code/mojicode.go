package main

import (
	"errors"
	"strconv"
	"unicode"
)

const (
	MaxR16 = '\uFFFF'
)

type Entity struct {
	String string
	Detail string
}

var (
	ASCII = struct {
		All      *unicode.RangeTable
		Control  *unicode.RangeTable
		Number   *unicode.RangeTable
		Symbol   *unicode.RangeTable
		Alphabet struct {
			Upper *unicode.RangeTable
			Lower *unicode.RangeTable
		}
		ControlDetails map[uint16]Entity
	}{
		All: &unicode.RangeTable{
			R16: []unicode.Range16{
				{0x0000, 0x007F, 1},
			},
		},
		Control: &unicode.RangeTable{
			R16: []unicode.Range16{
				{0x0000, 0x0020, 1},
				{0x007F, 0x007F, 1},
			},
		},
		Number: &unicode.RangeTable{
			R16: []unicode.Range16{{0x0030, 0x0039, 1}},
		},
		Symbol: &unicode.RangeTable{
			R16: []unicode.Range16{
				{0x0021, 0x0040, 1},
				{0x005B, 0x0060, 1},
				{0x007B, 0x007E, 1},
			},
		},
		Alphabet: struct {
			Upper *unicode.RangeTable
			Lower *unicode.RangeTable
		}{
			Upper: &unicode.RangeTable{
				R16: []unicode.Range16{{0x0041, 0x005A, 1}},
			},
			Lower: &unicode.RangeTable{
				R16: []unicode.Range16{{0x0061, 0x007A, 1}},
			},
		},
		ControlDetails: map[uint16]Entity{
			0x0000: {"NULL", "Null"},
			0x0001: {"SOH", "Start Of Heading"},
			0x0002: {"STX", "End Of Text"},
			0x0003: {"ETX", "End Of Transmission"},
			0x0004: {"EOT", "End Of Transmission"},
			0x0005: {"ENQ", "Enquiry"},
			0x0006: {"ACK", "Acknowledgement"},
			0x0007: {"BEL", "Bell"},
			0x0008: {"BS", "Back Space"},
			0x0009: {"HT", "Horizontal Tabulation"},
			0x000A: {"LF/NL", "Line Feed/New Line"},
			0x000B: {"VT", "Vertical Tabulation"},
			0x000C: {"FF/NP", "Form Feed/New Page"},
			0x000D: {"CR", "Carriage Return"},
			0x000E: {"SO", "Shift Out"},
			0x000F: {"SI", "Shift In"},
			0x0010: {"DLE", "Data Link Escape"},
			0x0011: {"DC1", "Device Control 1"},
			0x0012: {"DC2", "Device Control 2"},
			0x0013: {"DC3", "Device Control 3"},
			0x0014: {"DC4", "Device Control 4"},
			0x0015: {"NAK", "Negative Acknowledgement"},
			0x0016: {"SYN", "Synchronous idle"},
			0x0017: {"ETB", "End of Transmission Block"},
			0x0018: {"CAN", "Cancel"},
			0x0019: {"EM", "End of Medium"},
			0x001A: {"SUB/EOF", "Substitute/End Of File"},
			0x001B: {"ESC", "Escape"},
			0x001C: {"FS", "File Separator"},
			0x001D: {"GS", "Group Separator"},
			0x001E: {"RS", "Record Separator"},
			0x001F: {"US", "Unit Separator"},
			0x0020: {"SPC", "Space"},
			0x007F: {"DEL", "Delete"},
		},
	}
)

func inspect(r rune) string {
	switch {
	case unicode.IsControl(r):
		return "Control"
	case unicode.IsSpace(r):
		return "Space"
	case unicode.IsGraphic(r):
		switch {
		case unicode.IsSymbol(r) || unicode.IsPunct(r):
			return "Symbol"
		case unicode.IsDigit(r):
			return "Digit"
		default:
			switch {
			case unicode.Is(unicode.Hiragana, r):
				return "Japanese Hiragana"
			case unicode.Is(unicode.Katakana, r):
				return "Japanese Katakana"
			case unicode.Is(unicode.Han, r):
				return "Chinese character"
			case unicode.Is(ASCII.Alphabet.Upper, r):
				return "Alphabet Upper-case"
			case unicode.Is(ASCII.Alphabet.Lower, r):
				return "Alphabet Lower-case"
			default:
				return ""
			}
		}
	default:
		return ""
	}
}

func toEntity(r rune) Entity {
	e := Entity{
		String: string(r),
	}
	if !unicode.IsPrint(r) {
		e.String = "Unprintable"
	}
	if r <= unicode.MaxASCII {
		if entity, ok := ASCII.ControlDetails[uint16(r)]; ok {
			return entity
		}
	}

	e.Detail = inspect(r)
	return e
}

type Callback func(uint32, Entity)

func Each(table *unicode.RangeTable, callback Callback) {
	for _, v := range table.R16 {
		for i := v.Lo; i <= v.Hi; i++ {
			callback(uint32(i), toEntity(rune(i)))
		}
	}
	for _, v := range table.R32 {
		for i := v.Lo; i <= v.Hi; i++ {
			callback(i, toEntity(rune(i)))
		}
	}
}

type Uint32Slice []uint32

func (u Uint32Slice) ToRangeTable() (*unicode.RangeTable, error) {
	t := unicode.RangeTable{}
	l := len(u) - 1
	for i, j := 0, 1; i <= l; i, j = i+2, j+2 {
		if j > l {
			if u[i] <= MaxR16 {
				r16 := unicode.Range16{Lo: uint16(u[i]), Hi: uint16(u[i]), Stride: 1}
				t.R16 = append(t.R16, r16)
			} else {
				r32 := unicode.Range32{Lo: uint32(u[i]), Hi: uint32(u[i]), Stride: 1}
				t.R32 = append(t.R32, r32)
			}
			continue
		}

		if u[i] <= MaxR16 && u[j] <= MaxR16 {
			r16 := unicode.Range16{Lo: uint16(u[i]), Hi: uint16(u[j]), Stride: 1}
			t.R16 = append(t.R16, r16)
			continue
		}
		if u[i] > MaxR16 && u[j] > MaxR16 {
			r32 := unicode.Range32{Lo: uint32(u[i]), Hi: uint32(u[j]), Stride: 1}
			t.R32 = append(t.R32, r32)
			continue
		}

		return nil, errors.New("ToRnageTable: Failed to create a range table")
	}
	return &t, nil
}

func (u Uint32Slice) Each(callback Callback) {
	for _, v := range u {
		callback(v, toEntity(rune(v)))
	}
}

func toUint32Slice(s []string, base int) (Uint32Slice, error) {
	if base < 0 {
		var a []uint32
		for _, v := range s {
			for _, vv := range v {
				a = append(a, uint32(vv))
			}
		}
		return a, nil
	}

	a := make([]uint32, len(s), len(s))
	for i, v := range s {
		u64, err := strconv.ParseUint(v, base, 32)
		if err != nil {
			return nil, err
		}
		a[i] = uint32(u64)
	}
	return a, nil
}
