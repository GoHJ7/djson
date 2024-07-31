package djson

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/friendsofgo/errors"
)

const (
	V_TYPE_NULL int = iota
	V_TYPE_INT
	V_TYPE_FLOAT
	V_TYPE_NUMBER
	V_TYPE_STRING
	V_TYPE_BOOL
	V_TYPE_OBJECT
	V_TYPE_ARRAY
	V_TYPE_MULTI
)

var MAX_STR_LEN int64 = 2097152 // .NET JavaScriptSerializer.MaxJsonLength

func SetMaxJsonLength(maxLen int64) {
	MAX_STR_LEN = maxLen
}

var CountryCodes = []string{
	"GH", "GA", "GY", "GM", "GG", "GP", "GT", "GU", "GD", "GR", "GL", "GW", "GN", "NA", "NR", "NG", "AQ", "SS", "ZA", "AN", "NL",
	"NP", "NO", "NF", "NC", "NZ", "NU", "NE", "NI", "KR", "DK", "DO", "DM", "DE", "TL", "LA", "LR", "LV", "RU", "LB", "LS", "RE",
	"RO", "LU", "RW", "LY", "LT", "LI", "MG", "MQ", "MH", "YT", "MO", "MW", "MY", "ML", "IM", "MX", "MC", "MA", "MU", "MR", "MZ",
	"ME", "MS", "MD", "MV", "MT", "MN", "UM", "VI", "US", "MM", "FM", "VU", "BH", "BB", "VA", "BS", "BD", "BM", "BJ", "VE", "VN",
	"BE", "BY", "BZ", "BA", "BW", "BO", "BI", "BF", "BV", "BT", "MP", "MK", "BG", "BR", "BN", "WS", "SA", "GS", "SM", "ST", "PM",
	"EH", "SN", "RS", "SC", "LC", "VC", "KN", "SH", "SO", "SB", "SD", "SR", "LK", "SJ", "SE", "CH", "ES", "SK", "SI", "SY", "SL",
	"SX", "SG", "AE", "AW", "AM", "AR", "AS", "IS", "HT", "IE", "AZ", "AF", "AD", "AL", "DZ", "AO", "AG", "AI", "ER", "SZ", "EE",
	"EC", "ET", "SV", "VG", "IO", "GB", "YE", "OM", "AU", "AT", "HN", "AX", "WF", "JO", "UG", "UY", "UZ", "UA", "IQ", "IR", "IL",
	"EG", "IT", "ID", "IN", "JP", "JM", "ZM", "JE", "GQ", "KP", "GE", "CN", "CF", "DJ", "GI", "ZW", "TD", "CZ", "CL", "CM", "CV",
	"KZ", "QA", "KH", "CA", "KE", "KY", "KM", "CR", "CC", "CI", "CO", "CG", "CD", "CU", "KW", "CK", "HR", "CX", "KG", "KI", "CY",
	"TW", "TJ", "TZ", "TH", "TC", "TR", "TG", "TK", "TO", "TM", "TV", "TN", "TT", "PA", "PY", "PK", "PG", "PW", "PS", "FO", "PE",
	"PT", "FK", "PL", "PR", "GF", "TF", "PF", "FR", "FJ", "FI", "PH", "PN", "HM", "HU", "HK",
}

var HexRegExp *regexp.Regexp
var TimestampRegExp *regexp.Regexp
var YYYYMMDDRegExp *regexp.Regexp
var YYMMDDRegExp *regexp.Regexp
var HHMMSSRegExp *regexp.Regexp
var HHMMRegExp *regexp.Regexp
var EmailRegExp *regexp.Regexp
var UUIDRegExp *regexp.Regexp
var TelRegExp *regexp.Regexp
var BinRegExp *regexp.Regexp
var DecRegExp *regexp.Regexp

func CheckFuncHex(ts string, vi ...int64) bool {
	return HexRegExp.Match([]byte(ts))
}

func CheckFuncTimestamp(ts string, vi ...int64) bool {
	return TimestampRegExp.Match([]byte(ts))
}

func CheckFuncYYYYMMDD(ts string, vi ...int64) bool {
	return YYYYMMDDRegExp.Match([]byte(ts))
}

func CheckFuncYYMMDD(ts string, vi ...int64) bool {
	return YYMMDDRegExp.Match([]byte(ts))
}

func CheckFuncHHMMSS(ts string, vi ...int64) bool {
	return HHMMSSRegExp.Match([]byte(ts))
}

func CheckFuncHHMM(ts string, vi ...int64) bool {
	return HHMMRegExp.Match([]byte(ts))
}

func CheckFuncEmail(ts string, vi ...int64) bool {
	return EmailRegExp.Match([]byte(ts))
}

func CheckFuncIntString(ts string, vi ...int64) bool {
	_, err := strconv.Atoi(ts)
	return err == nil
}

func CheckFuncFloatString(ts string, vi ...int64) bool {
	_, err := strconv.ParseFloat(ts, 64)
	return err == nil
}

func CheckFuncUUID(ts string, vi ...int64) bool {
	return UUIDRegExp.Match([]byte(ts))
}

func CheckISO31661A2(val string, vi ...int64) bool {
	if len(val) != 2 {
		return false
	}

	valUpper := strings.ToUpper(val)

	for idx := range CountryCodes {
		if valUpper == CountryCodes[idx] {
			return true
		}
	}

	return false
}

func CheckBase64(ts string, vi ...int64) bool {
	_, err := base64.StdEncoding.DecodeString(ts)
	return err == nil
}

func CheckTelephone(ts string, vi ...int64) bool {
	return TelRegExp.Match([]byte(ts))
}

// ISO 3166-2 : KR-XX, GH-XX, ...
func CheckISO31662(val string, vi ...int64) bool {
	if len(val) < 4 {
		return false
	}

	if val[2:3] != "-" {
		return false
	}

	return CheckISO31661A2(val[0:2])
}

func CheckFuncBoolString(ts string, vi ...int64) bool {
	tslower := strings.ToLower(ts)
	return tslower == "true" || tslower == "false"
}

func CheckHexIfExist(ts string, vi ...int64) bool {
	if CheckFuncMinMaxString(ts, vi...) {
		return CheckFuncHex(ts)
	}
	return false
}

func CheckFuncMinMaxString(ts string, vi ...int64) bool {
	if len(vi) >= 2 {
		return len(ts) == int(vi[0]) || len(ts) == int(vi[1])
	}

	return false
}

func CheckFuncBin(ts string, vi ...int64) bool {
	return BinRegExp.Match([]byte(ts))
}

func CheckFuncDec(ts string, vi ...int64) bool {
	return DecRegExp.Match([]byte(ts))
}

func init() {
	HexRegExp = regexp.MustCompile(`^([A-Fa-f0-9]{2})*$`)
	TimestampRegExp = regexp.MustCompile(`^[0-9]{9,11}$`)
	YYYYMMDDRegExp = regexp.MustCompile(`^[1-2][0-9]{3}-{0,1}(0[1-9]|1[0-2])-{0,1}(0[1-9]|[1-2][0-9]|3[0-1])$`)
	YYMMDDRegExp = regexp.MustCompile(`^[0-9]{2}-{0,1}(0[1-9]|1[0-2])-{0,1}(0[1-9]|[1-2][0-9]|3[0-1])$`)
	HHMMSSRegExp = regexp.MustCompile(`^([0-1][0-9]|2[0-3])\:{0,1}([0-5][0-9])\:{0,1}([0-5][0-9])$`)
	HHMMRegExp = regexp.MustCompile(`^([0-1][0-9]|2[0-3])\:{0,1}([0-5][0-9])$`)
	EmailRegExp = regexp.MustCompile(`^([\w\.\_\-])*[a-zA-Z0-9]+([\w\.\_\-])*([a-zA-Z0-9])+([\w\.\_\-])+@([a-zA-Z0-9]+\.)+[a-zA-Z0-9]{2,8}$`)
	UUIDRegExp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89AB][0-9a-f]{3}-[0-9a-f]{12}$`)
	TelRegExp = regexp.MustCompile(`^((|\+\d{1,2})(|[-.\s])\d{2}|\d{2,3}|\(\d{2,3}\))(|[-.\s])\d{3,4}(|[-.\s])\d{4}$`)
	BinRegExp = regexp.MustCompile(`^[0-1]*$`)
	DecRegExp = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)
}

type VItem struct {
	Type      int
	Name      string
	Max       int64
	Min       int64
	MaxFloat  float64
	MinFloat  float64
	Size      int64
	IsRequred bool
	SubItems  []*VItem
	CheckFunc func(string, ...int64) bool
	RegExp    *regexp.Regexp
}

type Validator struct {
	Syntax    *JSON
	RootItems []*VItem
}

func NewValidator() *Validator {
	return &Validator{
		Syntax: New(),
	}
}

func (m *Validator) Compile(syntax string) bool {
	m.Syntax.Parse(syntax)

	if !m.Syntax.IsObject() && !m.Syntax.IsString() && !m.Syntax.IsArray() {
		m.Syntax = New()
		return false
	}

	m.RootItems = make([]*VItem, 0)

	if m.Syntax.IsObject() || m.Syntax.IsString() {
		vItem := GetVItem("__root__", m.Syntax)
		if vItem != nil {
			m.RootItems = append(m.RootItems, vItem)
		}

	} else if m.Syntax.IsArray() {
		m.Syntax.Seek()

		for m.Syntax.Next() {
			es := m.Syntax.Scan()
			vi := GetVItem("__root__", es)
			if vi != nil {
				m.RootItems = append(m.RootItems, vi)
			}
		}
	}

	return true
}

func GetVItem(name string, ejson *JSON) *VItem {
	eitem := new(VItem)
	eitem.Name = name
	etype := ""

	if ejson.IsString() {
		etype = ejson.String()

		switch etype {
		case "INT":
			eitem.Type = V_TYPE_INT
			eitem.Min = int64(-9007199254740991)
			eitem.Max = int64(9007199254740991)
		case "UNIXTIME", "UINT":
			eitem.Type = V_TYPE_INT
			eitem.Min = 0
			eitem.Max = int64(9007199254740991)
		case "FLOAT":
			eitem.Type = V_TYPE_FLOAT
			eitem.MinFloat = float64(-1.7976931348623157e+308)
			eitem.MaxFloat = float64(1.7976931348623157e+308)
		case "NUMBER":
			eitem.Type = V_TYPE_NUMBER
			eitem.MinFloat = float64(-1.7976931348623157e+308)
			eitem.MaxFloat = float64(1.7976931348623157e+308)
		case "STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = MAX_STR_LEN
		case "OBJECT":
			eitem.Type = V_TYPE_OBJECT
		case "ARRAY":
			eitem.Type = V_TYPE_ARRAY
			eitem.Min = 0
			eitem.Max = int64(9007199254740991)
		case "EMPTY.STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = 0
		case "NONEMPTY.STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 1
			eitem.Max = MAX_STR_LEN
		case "NONEMPTY.ARRAY":
			eitem.Type = V_TYPE_ARRAY
			eitem.Min = 1
			eitem.Max = int64(9007199254740991)
		case "BIN":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = MAX_STR_LEN
			eitem.CheckFunc = CheckFuncBin
		case "DEC":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = MAX_STR_LEN
			eitem.CheckFunc = CheckFuncDec
		case "HEX":
			eitem.Type = V_TYPE_STRING
			eitem.Min = 0
			eitem.Max = MAX_STR_LEN
			eitem.CheckFunc = CheckFuncHex
		case "BOOL":
			eitem.Type = V_TYPE_BOOL
		}

	} else if ejson.IsArray() {

		eitem.Type = V_TYPE_MULTI
		ejson.Seek()
		for ejson.Next() {
			es := ejson.Scan()
			vi := GetVItem(name, es)
			if vi != nil {
				eitem.SubItems = append(eitem.SubItems, vi)
			}
		}

	} else if ejson.IsObject() {

		etype = ejson.String("type")
		eitem.IsRequred = ejson.Bool("required")
		if ejson.String("regexp") != "" {
			eitem.RegExp, _ = regexp.Compile(ejson.String("regexp"))
		}

		switch etype {
		case "INT":
			eitem.Type = V_TYPE_INT
			eitem.Min = ejson.Int("min", int64(-9007199254740991))
			eitem.Max = ejson.Int("max", int64(9007199254740991))
		case "UNIXTIME", "UINT":
			eitem.Type = V_TYPE_INT
			eitem.Min = ejson.Int("min", 0)
			eitem.Max = ejson.Int("max", int64(9007199254740991))

			if eitem.Min < 0 {
				eitem.Min = 0
			}
		case "FLOAT":
			eitem.Type = V_TYPE_FLOAT
			eitem.MinFloat = ejson.Float("min", float64(-1.7976931348623157e+308))
			eitem.MaxFloat = ejson.Float("max", float64(1.7976931348623157e+308))
		case "NUMBER":
			eitem.Type = V_TYPE_NUMBER
			eitem.MinFloat = ejson.Float("min", float64(-1.7976931348623157e+308))
			eitem.MaxFloat = ejson.Float("max", float64(1.7976931348623157e+308))
		case "STRING":
			eitem.Type = V_TYPE_STRING
			if ejson.IsInt("size") {
				eitem.Min = ejson.Int("size")
				eitem.Max = eitem.Min
			} else {
				eitem.Min = ejson.Int("min", 0)
				eitem.Max = ejson.Int("max", MAX_STR_LEN)
			}
		case "MIN.MAX.STRING":
			eitem.Type = V_TYPE_STRING
			eitem.Min = ejson.Int("min", 0)
			eitem.Max = ejson.Int("max", MAX_STR_LEN)
			eitem.CheckFunc = CheckFuncMinMaxString
		case "OBJECT":
			subJson, ok := ejson.Object("object")
			if ok {
				eitem.Type = V_TYPE_OBJECT
				ks := subJson.GetKeys()
				for _, ek := range ks {
					ejson, ok := subJson.Get(ek)
					if ok {
						vItem := GetVItem(ek, ejson)
						if vItem != nil {
							eitem.SubItems = append(eitem.SubItems, vItem)
						}
					}

				}
			}
		case "NONEMPTY.STRING":
			eitem.Type = V_TYPE_STRING
			if ejson.IsInt("size") {
				eitem.Min = ejson.Int("size")
				eitem.Max = eitem.Min
			} else {
				eitem.Min = ejson.Int("min", 1)
				eitem.Max = ejson.Int("max", MAX_STR_LEN)
			}

			if eitem.Min < 1 {
				eitem.Min = 1
			}

		case "ARRAY":
			if ejson.IsInt("size") {
				eitem.Min = ejson.Int("size")
				eitem.Max = eitem.Min
			} else {
				eitem.Min = ejson.Int("min", 0)
				eitem.Max = ejson.Int("max", int64(9007199254740991))
			}

			if eitem.Min < 0 {
				eitem.Min = 0
			}

		case "NONEMPTY.ARRAY":
			if ejson.IsInt("size") {
				eitem.Min = ejson.Int("size")
				eitem.Max = eitem.Min
			} else {
				eitem.Min = ejson.Int("min", 1)
				eitem.Max = ejson.Int("max", int64(9007199254740991))
			}

			if eitem.Min < 1 {
				eitem.Min = 1
			}
		case "BIN":
			eitem.CheckFunc = CheckFuncBin
		case "DEC":
			eitem.CheckFunc = CheckFuncDec
		case "HEX":
			eitem.CheckFunc = CheckFuncHex
		case "BOOL":
			eitem.Type = V_TYPE_BOOL
		}

		if etype == "BIN" || etype == "DEC" || etype == "HEX" {
			eitem.Type = V_TYPE_STRING
			if ejson.IsInt("size") {
				eitem.Min = ejson.Int("size")
				eitem.Max = eitem.Min
			} else {
				eitem.Min = ejson.Int("min", 0)
				eitem.Max = ejson.Int("max", MAX_STR_LEN)
			}
		}

		if etype == "ARRAY" || etype == "NONEMPTY.ARRAY" {
			eitem.Type = V_TYPE_ARRAY
			eitem.Max = ejson.Int("max", int64(9007199254740991))
			oa, ok := ejson.Get("array") // type of element
			if ok {
				eitem.SubItems = make([]*VItem, 0)
				if oa.IsArray() {
					oa.Seek()
					for oa.Next() {
						es := oa.Scan()
						vi := GetVItem("__array__", es)
						if vi != nil {
							eitem.SubItems = append(eitem.SubItems, vi)
						}
					}
				} else if oa.IsString() || oa.IsObject() {
					vi := GetVItem("__array__", oa)
					eitem.SubItems = append(eitem.SubItems, vi)
				}
			}
		}

	}

	switch etype {
	case "TIMESTAMP":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = 10
		eitem.CheckFunc = CheckFuncTimestamp
	case "YYYYMMDD":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 8
		eitem.Max = 10
		eitem.CheckFunc = CheckFuncYYYYMMDD
	case "YYMMDD":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 6
		eitem.Max = 8
		eitem.CheckFunc = CheckFuncYYMMDD
	case "HHMMSS":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 6
		eitem.Max = 8
		eitem.CheckFunc = CheckFuncHHMMSS
	case "HHMM":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 4
		eitem.Max = 5
		eitem.CheckFunc = CheckFuncHHMM
	case "EMAIL":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 3
		eitem.Max = 255
		eitem.CheckFunc = CheckFuncEmail
	case "INT.STRING", "INT_STRING":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 1
		eitem.Max = 17
		eitem.CheckFunc = CheckFuncIntString
	case "FLOAT.STRING", "FLOAT_STRING":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 1
		eitem.Max = 24
		eitem.CheckFunc = CheckFuncFloatString
	case "BOOL.STRING", "BOOL_STRING":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 4
		eitem.Max = 5
		eitem.CheckFunc = CheckFuncBoolString
	case "UUID":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 36
		eitem.Max = 36
		eitem.CheckFunc = CheckFuncUUID
	case "ISO31661A2":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 2
		eitem.Max = 2
		eitem.CheckFunc = CheckISO31661A2
	case "ISO31662":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 5
		eitem.Max = 5
		eitem.CheckFunc = CheckISO31662
	case "BASE64":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = MAX_STR_LEN
		eitem.CheckFunc = CheckBase64
	case "TELEPHONE":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 4
		eitem.Max = 20
		eitem.CheckFunc = CheckTelephone
	case "HEX64.IF.EXIST":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = 16
		eitem.CheckFunc = CheckHexIfExist
	case "HEX128.IF.EXIST":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = 32
		eitem.CheckFunc = CheckHexIfExist
	case "HEX256.IF.EXIST":
		eitem.Type = V_TYPE_STRING
		eitem.Min = 0
		eitem.Max = 64
		eitem.CheckFunc = CheckHexIfExist
	}

	return eitem
}

var ErrEmptyjson = errors.New("empty validator")
var ErrArrayjson = errors.New("No valid element while root is an array or string")

func (m *Validator) IsValidWithError(tjson *JSON) (error, bool) {
	if tjson == nil {
		if len(m.RootItems) == 0 {
			return nil, true
		} else {
			return ErrEmptyjson, false
		}
	}

	if m.Syntax.IsObject() { // json must be valid one

		for _, vitem := range m.RootItems {
			return CheckVItemWithError(vitem, tjson)
		}

	} else if m.Syntax.IsArray() || m.Syntax.IsString() {
		// each element must be valid for one of vitems

		if len(m.RootItems) == 0 {
			// log.Println("empty rootitems")
			return nil, true
		}

		for _, vitem := range m.RootItems {
			if _, check := CheckVItemWithError(vitem, tjson); check {
				return nil, true
			}
		}

		return ErrArrayjson, false
	}

	return nil, true

}

func typeToReadable(t int) string {
	switch t {
	case 0:
		return "null"
	case 1:
		return "int"
	case 2:
		return "float"
	case 3:
		return "number"
	case 4:
		return "string"
	case 5:
		return "bool"
	case 6:
		return "object"
	case 7:
		return "array"
	case 8:
		return "multi"
	default:
		return "something wrong with type"
	}
}

func CheckVItemWithError(vi *VItem, tjson *JSON) (error, bool) {
	if vi.Name == "" {
		err := fmt.Errorf("%s >> empty name", vi.Name)
		return err, false
	}

	var vtype string

	if vi.Name == "__root__" || vi.Name == "__array__" {
		vtype = tjson.Type()
	} else {
		vtype = tjson.Type(vi.Name)
	}

	if vtype == "" && !vi.IsRequred {
		return nil, true
	}

	switch vi.Type {
	case V_TYPE_INT:
		if vtype != "int" {
			err := fmt.Errorf("%s >> your type must be %s but is %s", vi.Name, typeToReadable(vi.Type), vtype)
			return err, false
		}

		var si int64

		if vi.Name == "__root__" || vi.Name == "__array__" {
			si = tjson.Int()
		} else {
			si = tjson.Int(vi.Name)
		}

		if vi.Max < si {
			err := fmt.Errorf("%s >> Out of range. Must be under %d but is %d", vi.Name, vi.Max, si)
			return err, false
		}

		if vi.Min > si {
			err := fmt.Errorf("%s >> Out of range. Must be over %d but is %d", vi.Name, vi.Min, si)
			return err, false
		}

	case V_TYPE_FLOAT:
		if vtype != "float" {
			err := fmt.Errorf("%s >> your type must be %s but is %s", vi.Name, typeToReadable(vi.Type), vtype)
			return err, false
		}

		fallthrough

	case V_TYPE_NUMBER:
		if !(vtype == "float" || vtype == "int") {
			err := fmt.Errorf("%s >> your type must be %s but is %s", vi.Name, typeToReadable(vi.Type), vtype)
			return err, false
		}

		var sf float64

		if vi.Name == "__root__" || vi.Name == "__array__" {
			sf = tjson.Float()
		} else {
			sf = tjson.Float(vi.Name)
		}

		if vi.MaxFloat < sf {
			err := fmt.Errorf("%s >> Out of range. Must be under %f but is %f", vi.Name, vi.MaxFloat, sf)
			return err, false
		}

		if vi.MinFloat > sf {
			err := fmt.Errorf("%s >> Out of range. Must be over %f but is %f", vi.Name, vi.MinFloat, sf)
			return err, false
		}
	case V_TYPE_STRING:
		if vtype != "string" {
			err := fmt.Errorf("%s >> your type must be %s but is %s", vi.Name, typeToReadable(vi.Type), vtype)
			return err, false
		}

		var ss string

		if vi.Name == "__root__" || vi.Name == "__array__" {
			ss = tjson.String()
		} else {
			ss = tjson.String(vi.Name)
		}

		lenv := int64(len(ss))

		if lenv > vi.Max {
			err := fmt.Errorf("%s >> Length out of range. Must be under %d characters but is %d characters", vi.Name, vi.Max, lenv)
			return err, false
		}
		if lenv < vi.Min {
			err := fmt.Errorf("%s >> Length out of range. Must be at least %d characters but is %d characters", vi.Name, vi.Min, lenv)
			return err, false
		}

		if vi.RegExp != nil {
			err := fmt.Errorf("%s >> regex compilation fail", vi.Name)
			return err, vi.RegExp.Match([]byte(ss))
		}

		if vi.CheckFunc != nil {
			err := fmt.Errorf("%s >> check func fail", vi.Name)
			return err, vi.CheckFunc(ss, vi.Min, vi.Max)
		}

	case V_TYPE_OBJECT:
		if vi.Name == "__root__" && vtype != "object" {
			err := fmt.Errorf("%s >> your type must be %s but is %s", vi.Name, typeToReadable(vi.Type), vtype)
			return err, false
		}

		var so *JSON
		var ok bool

		if vi.Name == "__root__" || vi.Name == "__array__" {
			so, ok = tjson.Object()
			if !ok {
				err := fmt.Errorf("%s >> object conversion fail", vi.Name)
				return err, false
			}
		} else {
			so, ok = tjson.Object(vi.Name)
		}

		if vi.IsRequred && !ok {
			err := fmt.Errorf("%s >> required", vi.Name)
			return err, false
		}

		if !ok {
			return nil, true
		}

		for _, svi := range vi.SubItems {
			if err, check := CheckVItemWithError(svi, so); !check {
				err = fmt.Errorf("%s.%s", vi.Name, err.Error())
				return err, false
			}
		}

	case V_TYPE_ARRAY:
		if vi.Name == "__root__" && vtype != "array" {
			err := fmt.Errorf("%s >> your type must be %s but is %s", vi.Name, typeToReadable(vi.Type), vtype)
			return err, false
		}

		var sa *JSON
		var ok bool

		if vi.Name == "__root__" || vi.Name == "__array__" {
			sa, ok = tjson.Array()
			if !ok {
				err := fmt.Errorf("%s >> array conversion fail", vi.Name)
				return err, false
			}
		} else {
			sa, ok = tjson.Array(vi.Name)
		}

		if vi.IsRequred && !ok {
			err := fmt.Errorf("%s >> required", vi.Name)
			return err, false
		}

		if !ok {
			return nil, true
		}

		lenv := int64(sa.Len())
		if lenv > vi.Max {
			err := fmt.Errorf("%s >> array length out of range. Must be at most %d but is %d", vi.Name, vi.Max, lenv)
			return err, false
		}

		if lenv < vi.Min {
			err := fmt.Errorf("%s >> array length out of range. Must be at least %d but is %d", vi.Name, vi.Min, lenv)
			return err, false
		}

		if ok {
			if len(vi.SubItems) == 0 {
				return nil, true
			}

			sa.Seek() // valid element type
			for sa.Next() {
				ssa := sa.Scan()
				isValid := false
				for _, svi := range vi.SubItems {
					if _, check := CheckVItemWithError(svi, ssa); check {
						isValid = true
						break
					}
				}

				// error of array can't be determined
				if !isValid {
					err := fmt.Errorf("%s >> no valid element in the array", vi.Name)
					return err, false
				}
			}
		}

	case V_TYPE_BOOL:

		var ok bool
		if vi.Name == "__root__" || vi.Name == "__array__" {
			ok = tjson.IsBool()
		} else {
			ok = tjson.IsBool(vi.Name)
		}

		if vi.IsRequred && !ok {
			err := fmt.Errorf("%s >> required", vi.Name)
			return err, false
		}

	case V_TYPE_MULTI:

		isValid := false

		for _, svi := range vi.SubItems {
			if _, check := CheckVItemWithError(svi, tjson); check {
				isValid = true
				break
			}
		}

		if !isValid {
			err := fmt.Errorf("%s >> no valid element in the multi type", vi.Name)
			return err, false
		}
	}
	return nil, true
}

func (m *Validator) IsValid(tjson *JSON) bool {
	if tjson == nil {
		return len(m.RootItems) == 0
	}

	if m.Syntax.IsObject() { // json must be valid one

		for _, vitem := range m.RootItems {
			return CheckVItem(vitem, tjson)
		}

	} else if m.Syntax.IsArray() || m.Syntax.IsString() {
		// each element must be valid for one of vitems

		if len(m.RootItems) == 0 {
			// log.Println("empty rootitems")
			return true
		}

		for _, vitem := range m.RootItems {
			if CheckVItem(vitem, tjson) {
				return true
			}
		}

		return false
	}

	return true

}

func CheckVItem(vi *VItem, tjson *JSON) bool {
	if vi.Name == "" {
		return false
	}

	var vtype string

	if vi.Name == "__root__" || vi.Name == "__array__" {
		vtype = tjson.Type()
	} else {
		vtype = tjson.Type(vi.Name)
	}

	if vtype == "" && !vi.IsRequred {
		return true
	}

	switch vi.Type {
	case V_TYPE_INT:
		if vtype != "int" {
			return false
		}

		var si int64

		if vi.Name == "__root__" || vi.Name == "__array__" {
			si = tjson.Int()
		} else {
			si = tjson.Int(vi.Name)
		}

		if vi.Max < si || vi.Min > si {
			return false
		}

	case V_TYPE_NUMBER:
		if vtype != "float" && vtype != "int" {
			return false
		}

		fallthrough

	case V_TYPE_FLOAT:
		if vtype != "float" {
			return false
		}

		var sf float64

		if vi.Name == "__root__" || vi.Name == "__array__" {
			sf = tjson.Float()
		} else {
			sf = tjson.Float(vi.Name)
		}

		if vi.MaxFloat < sf || vi.MinFloat > sf {
			return false
		}
	case V_TYPE_STRING:
		if vtype != "string" {
			return false
		}

		var ss string

		if vi.Name == "__root__" || vi.Name == "__array__" {
			ss = tjson.String()
		} else {
			ss = tjson.String(vi.Name)
		}

		lenv := int64(len(ss))

		if lenv > vi.Max || lenv < vi.Min {
			return false
		}

		if vi.RegExp != nil {
			return vi.RegExp.Match([]byte(ss))
		}

		if vi.CheckFunc != nil {
			return vi.CheckFunc(ss, vi.Min, vi.Max)
		}

	case V_TYPE_OBJECT:
		if vi.Name == "__root__" && vtype != "object" {
			return false
		}

		var so *JSON
		var ok bool

		if vi.Name == "__root__" || vi.Name == "__array__" {
			so, ok = tjson.Object()
			if !ok {
				return false
			}
		} else {
			so, ok = tjson.Object(vi.Name)
		}

		if vi.IsRequred && !ok {
			return false
		}

		if !ok {
			return true
		}

		for _, svi := range vi.SubItems {
			if !CheckVItem(svi, so) {
				return false
			}
		}

	case V_TYPE_ARRAY:
		if vi.Name == "__root__" && vtype != "array" {
			return false
		}

		var sa *JSON
		var ok bool

		if vi.Name == "__root__" || vi.Name == "__array__" {
			sa, ok = tjson.Array()
			if !ok {
				return false
			}
		} else {
			sa, ok = tjson.Array(vi.Name)
		}

		if vi.IsRequred && !ok {
			return false
		}

		if !ok {
			return true
		}

		lenv := int64(sa.Len())
		if lenv > vi.Max || lenv < vi.Min {
			return false
		}

		if ok {
			if len(vi.SubItems) == 0 {
				return true
			}

			sa.Seek() // valid element type
			for sa.Next() {
				ssa := sa.Scan()
				isValid := false
				for _, svi := range vi.SubItems {
					if CheckVItem(svi, ssa) {
						isValid = true
						break
					}
				}
				if !isValid {
					return false
				}
			}
		}

	case V_TYPE_BOOL:

		var ok bool
		if vi.Name == "__root__" || vi.Name == "__array__" {
			ok = tjson.IsBool()
		} else {
			ok = tjson.IsBool(vi.Name)
		}

		if vi.IsRequred && !ok {
			return false
		}

	case V_TYPE_MULTI:

		isValid := false

		for _, svi := range vi.SubItems {
			if CheckVItem(svi, tjson) {
				isValid = true
				break
			}
		}

		if !isValid {
			return false
		}

	}

	return true
}
