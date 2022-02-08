package main

import (
	"errors"
	"strconv"
	"strings"
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
			LatinOffset: 1,
		},
		Control: &unicode.RangeTable{
			R16: []unicode.Range16{
				{0x0000, 0x0020, 1},
				{0x007F, 0x007F, 1},
			},
			LatinOffset: 2,
		},
		Number: &unicode.RangeTable{
			R16:         []unicode.Range16{{0x0030, 0x0039, 1}},
			LatinOffset: 1,
		},
		Symbol: &unicode.RangeTable{
			R16: []unicode.Range16{
				{0x0021, 0x0040, 1},
				{0x005B, 0x0060, 1},
				{0x007B, 0x007E, 1},
			},
			LatinOffset: 1,
		},
		Alphabet: struct {
			Upper *unicode.RangeTable
			Lower *unicode.RangeTable
		}{
			Upper: &unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0041, 0x005A, 1},
				},
				LatinOffset: 1,
			},
			Lower: &unicode.RangeTable{
				R16: []unicode.Range16{
					{0x0061, 0x007A, 1},
				},
				LatinOffset: 1,
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

type Inspector struct {
	s string
	f func(rune) bool
	t *unicode.RangeTable
}

var (
	Inspectors = []Inspector{
		{s: "Control", f: unicode.IsControl},
		{s: "Space", f: unicode.IsSpace},
		{s: "Graphic", f: unicode.IsGraphic},
		{s: "Letter", f: unicode.IsLetter},
		{s: "Title", f: unicode.IsTitle},
		{s: "Mark", f: unicode.IsMark},
		{s: "Symbolic", f: unicode.IsSymbol},
		{s: "Punctuation", f: unicode.IsPunct},
		{s: "Digit", f: unicode.IsDigit},

		{s: "Adlam", t: unicode.Adlam},
		{s: "Ahom", t: unicode.Ahom},
		{s: "Anatolian-Hieroglyphs", t: unicode.Anatolian_Hieroglyphs},
		{s: "Arabic", t: unicode.Arabic},
		{s: "Armenian", t: unicode.Armenian},
		{s: "Avestan", t: unicode.Avestan},
		{s: "Balinese", t: unicode.Balinese},
		{s: "Bamum", t: unicode.Bamum},
		{s: "Bassa-Vah", t: unicode.Bassa_Vah},
		{s: "Batak", t: unicode.Batak},
		{s: "Bengali", t: unicode.Bengali},
		{s: "Bhaiksuki", t: unicode.Bhaiksuki},
		{s: "Bopomofo", t: unicode.Bopomofo},
		{s: "Brahmi", t: unicode.Brahmi},
		{s: "Braille", t: unicode.Braille},
		{s: "Buginese", t: unicode.Buginese},
		{s: "Buhid", t: unicode.Buhid},
		{s: "Canadian-Aboriginal", t: unicode.Canadian_Aboriginal},
		{s: "Carian", t: unicode.Carian},
		{s: "Caucasian-Albanian", t: unicode.Caucasian_Albanian},
		{s: "Chakma", t: unicode.Chakma},
		{s: "Cham", t: unicode.Cham},
		{s: "Cherokee", t: unicode.Cherokee},
		{s: "Chorasmian", t: unicode.Chorasmian},
		{s: "Common", t: unicode.Common},
		{s: "Coptic", t: unicode.Coptic},
		{s: "Cuneiform", t: unicode.Cuneiform},
		{s: "Cypriot", t: unicode.Cypriot},
		{s: "Cyrillic", t: unicode.Cyrillic},
		{s: "Deseret", t: unicode.Deseret},
		{s: "Devanagari", t: unicode.Devanagari},
		{s: "Dives-Akuru", t: unicode.Dives_Akuru},
		{s: "Dogra", t: unicode.Dogra},
		{s: "Duployan", t: unicode.Duployan},
		{s: "Egyptian-Hieroglyphs", t: unicode.Egyptian_Hieroglyphs},
		{s: "Elbasan", t: unicode.Elbasan},
		{s: "Elymaic", t: unicode.Elymaic},
		{s: "Ethiopic", t: unicode.Ethiopic},
		{s: "Georgian", t: unicode.Georgian},
		{s: "Glagolitic", t: unicode.Glagolitic},
		{s: "Gothic", t: unicode.Gothic},
		{s: "Grantha", t: unicode.Grantha},
		{s: "Greek", t: unicode.Greek},
		{s: "Gujarati", t: unicode.Gujarati},
		{s: "Gunjala-Gondi", t: unicode.Gunjala_Gondi},
		{s: "Gurmukhi", t: unicode.Gurmukhi},
		{s: "Han", t: unicode.Han},
		{s: "Hangul", t: unicode.Hangul},
		{s: "Hanifi-Rohingya", t: unicode.Hanifi_Rohingya},
		{s: "Hanunoo", t: unicode.Hanunoo},
		{s: "Hatran", t: unicode.Hatran},
		{s: "Hebrew", t: unicode.Hebrew},
		{s: "Hiragana", t: unicode.Hiragana},
		{s: "Imperial-Aramaic", t: unicode.Imperial_Aramaic},
		{s: "Inherited", t: unicode.Inherited},
		{s: "Inscriptional-Pahlavi", t: unicode.Inscriptional_Pahlavi},
		{s: "Inscriptional-Parthian", t: unicode.Inscriptional_Parthian},
		{s: "Javanese", t: unicode.Javanese},
		{s: "Kaithi", t: unicode.Kaithi},
		{s: "Kannada", t: unicode.Kannada},
		{s: "Katakana", t: unicode.Katakana},
		{s: "Kayah-Li", t: unicode.Kayah_Li},
		{s: "Kharoshthi", t: unicode.Kharoshthi},
		{s: "Khitan-Small-Script", t: unicode.Khitan_Small_Script},
		{s: "Khmer", t: unicode.Khmer},
		{s: "Khojki", t: unicode.Khojki},
		{s: "Khudawadi", t: unicode.Khudawadi},
		{s: "Lao", t: unicode.Lao},
		{s: "Latin", t: unicode.Latin},
		{s: "Lepcha", t: unicode.Lepcha},
		{s: "Limbu", t: unicode.Limbu},
		{s: "Linear-A", t: unicode.Linear_A},
		{s: "Linear-B", t: unicode.Linear_B},
		{s: "Lisu", t: unicode.Lisu},
		{s: "Lycian", t: unicode.Lycian},
		{s: "Lydian", t: unicode.Lydian},
		{s: "Mahajani", t: unicode.Mahajani},
		{s: "Makasar", t: unicode.Makasar},
		{s: "Malayalam", t: unicode.Malayalam},
		{s: "Mandaic", t: unicode.Mandaic},
		{s: "Manichaean", t: unicode.Manichaean},
		{s: "Marchen", t: unicode.Marchen},
		{s: "Masaram-Gondi", t: unicode.Masaram_Gondi},
		{s: "Medefaidrin", t: unicode.Medefaidrin},
		{s: "Meetei-Mayek", t: unicode.Meetei_Mayek},
		{s: "Mende-Kikakui", t: unicode.Mende_Kikakui},
		{s: "Meroitic-Cursive", t: unicode.Meroitic_Cursive},
		{s: "Meroitic-Hieroglyphs", t: unicode.Meroitic_Hieroglyphs},
		{s: "Miao", t: unicode.Miao},
		{s: "Modi", t: unicode.Modi},
		{s: "Mongolian", t: unicode.Mongolian},
		{s: "Mro", t: unicode.Mro},
		{s: "Multani", t: unicode.Multani},
		{s: "Myanmar", t: unicode.Myanmar},
		{s: "Nabataean", t: unicode.Nabataean},
		{s: "Nandinagari", t: unicode.Nandinagari},
		{s: "New-Tai-Lue", t: unicode.New_Tai_Lue},
		{s: "Newa", t: unicode.Newa},
		{s: "Nko", t: unicode.Nko},
		{s: "Nushu", t: unicode.Nushu},
		{s: "Nyiakeng-Puachue-Hmong", t: unicode.Nyiakeng_Puachue_Hmong},
		{s: "Ogham", t: unicode.Ogham},
		{s: "Ol-Chiki", t: unicode.Ol_Chiki},
		{s: "Old-Hungarian", t: unicode.Old_Hungarian},
		{s: "Old-Italic", t: unicode.Old_Italic},
		{s: "Old-North-Arabian", t: unicode.Old_North_Arabian},
		{s: "Old-Permic", t: unicode.Old_Permic},
		{s: "Old-Persian", t: unicode.Old_Persian},
		{s: "Old-Sogdian", t: unicode.Old_Sogdian},
		{s: "Old-South-Arabian", t: unicode.Old_South_Arabian},
		{s: "Old-Turkic", t: unicode.Old_Turkic},
		{s: "Oriya", t: unicode.Oriya},
		{s: "Osage", t: unicode.Osage},
		{s: "Osmanya", t: unicode.Osmanya},
		{s: "Pahawh-Hmong", t: unicode.Pahawh_Hmong},
		{s: "Palmyrene", t: unicode.Palmyrene},
		{s: "Pau-Cin-Hau", t: unicode.Pau_Cin_Hau},
		{s: "Phags-Pa", t: unicode.Phags_Pa},
		{s: "Phoenician", t: unicode.Phoenician},
		{s: "Psalter-Pahlavi", t: unicode.Psalter_Pahlavi},
		{s: "Rejang", t: unicode.Rejang},
		{s: "Runic", t: unicode.Runic},
		{s: "Samaritan", t: unicode.Samaritan},
		{s: "Saurashtra", t: unicode.Saurashtra},
		{s: "Sharada", t: unicode.Sharada},
		{s: "Shavian", t: unicode.Shavian},
		{s: "Siddham", t: unicode.Siddham},
		{s: "SignWriting", t: unicode.SignWriting},
		{s: "Sinhala", t: unicode.Sinhala},
		{s: "Sogdian", t: unicode.Sogdian},
		{s: "Sora-Sompeng", t: unicode.Sora_Sompeng},
		{s: "Soyombo", t: unicode.Soyombo},
		{s: "Sundanese", t: unicode.Sundanese},
		{s: "Syloti-Nagri", t: unicode.Syloti_Nagri},
		{s: "Syriac", t: unicode.Syriac},
		{s: "Tagalog", t: unicode.Tagalog},
		{s: "Tagbanwa", t: unicode.Tagbanwa},
		{s: "Tai-Le", t: unicode.Tai_Le},
		{s: "Tai-Tham", t: unicode.Tai_Tham},
		{s: "Tai-Viet", t: unicode.Tai_Viet},
		{s: "Takri", t: unicode.Takri},
		{s: "Tamil", t: unicode.Tamil},
		{s: "Tangut", t: unicode.Tangut},
		{s: "Telugu", t: unicode.Telugu},
		{s: "Thaana", t: unicode.Thaana},
		{s: "Thai", t: unicode.Thai},
		{s: "Tibetan", t: unicode.Tibetan},
		{s: "Tifinagh", t: unicode.Tifinagh},
		{s: "Tirhuta", t: unicode.Tirhuta},
		{s: "Ugaritic", t: unicode.Ugaritic},
		{s: "Vai", t: unicode.Vai},
		{s: "Wancho", t: unicode.Wancho},
		{s: "Warang-Citi", t: unicode.Warang_Citi},
		{s: "Yezidi", t: unicode.Yezidi},
		{s: "Yi", t: unicode.Yi},
		{s: "Zanabazar-Square", t: unicode.Zanabazar_Square},

		{s: "Lowercase", f: unicode.IsLower},
		{s: "Uppercase", f: unicode.IsUpper},

		{s: "ASCII-Hex-Digit", t: unicode.ASCII_Hex_Digit},
		{s: "Bidi-Control", t: unicode.Bidi_Control},
		{s: "Dash", t: unicode.Dash},
		{s: "Deprecated", t: unicode.Deprecated},
		{s: "Diacritic", t: unicode.Diacritic},
		{s: "Extender", t: unicode.Extender},
		{s: "Hex-Digit", t: unicode.Hex_Digit},
		{s: "Hyphen", t: unicode.Hyphen},
		{s: "IDS-Binary-Operator", t: unicode.IDS_Binary_Operator},
		{s: "IDS-Trinary-Operator", t: unicode.IDS_Trinary_Operator},
		{s: "Ideographic", t: unicode.Ideographic},
		{s: "Join-Control", t: unicode.Join_Control},
		{s: "Logical-Order-Exception", t: unicode.Logical_Order_Exception},
		{s: "Noncharacter-Code-Point", t: unicode.Noncharacter_Code_Point},
		{s: "Other-Alphabetic", t: unicode.Other_Alphabetic},
		{s: "Other-Default-Ignorable-Code-Point", t: unicode.Other_Default_Ignorable_Code_Point},
		{s: "Other-Grapheme-Extend", t: unicode.Other_Grapheme_Extend},
		{s: "Other-ID-Continue", t: unicode.Other_ID_Continue},
		{s: "Other-ID-Start", t: unicode.Other_ID_Start},
		{s: "Other-Lowercase", t: unicode.Other_Lowercase},
		{s: "Other-Math", t: unicode.Other_Math},
		{s: "Other-Uppercase", t: unicode.Other_Uppercase},
		{s: "Pattern-Syntax", t: unicode.Pattern_Syntax},
		{s: "Pattern-White-Space", t: unicode.Pattern_White_Space},
		{s: "Prepended-Concatenation-Mark", t: unicode.Prepended_Concatenation_Mark},
		{s: "Quotation-Mark", t: unicode.Quotation_Mark},
		{s: "Radical", t: unicode.Radical},
		{s: "Regional-Indicator", t: unicode.Regional_Indicator},
		{s: "STerm", t: unicode.STerm},
		{s: "Sentence-Terminal", t: unicode.Sentence_Terminal},
		{s: "Soft-Dotted", t: unicode.Soft_Dotted},
		{s: "Terminal-Punctuation", t: unicode.Terminal_Punctuation},
		{s: "Unified-Ideograph", t: unicode.Unified_Ideograph},
		{s: "Variation-Selector", t: unicode.Variation_Selector},
		{s: "White-Space", t: unicode.White_Space},
	}
)

func classify(r rune) string {
	var a []string
	for _, v := range Inspectors {
		if v.f != nil && v.f(r) {
			a = append(a, v.s)
		}
		if v.t != nil && unicode.Is(v.t, r) {
			a = append(a, v.s)
		}
	}
	return strings.Join(a, " ")
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
			entity.Detail = classify(r) + " " + entity.Detail
			return entity
		}
	}

	e.Detail = classify(r)
	return e
}

type Callback func(uint32, Entity) error

func Each(table *unicode.RangeTable, callback Callback) error {
	for _, v := range table.R16 {
		for i := v.Lo; i <= v.Hi; i++ {
			if err := callback(uint32(i), toEntity(rune(i))); err != nil {
				return err
			}
		}
	}
	for _, v := range table.R32 {
		for i := v.Lo; i <= v.Hi; i++ {
			if err := callback(i, toEntity(rune(i))); err != nil {
			}
		}
	}
	return nil
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

func (u Uint32Slice) Each(callback Callback) error {
	for _, v := range u {
		if err := callback(v, toEntity(rune(v))); err != nil {
			return err
		}
	}
	return nil
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
