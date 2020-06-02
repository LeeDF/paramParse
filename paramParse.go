package paramParse

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"unsafe"
)

const (
	STR              = "string"
	INT              = "int"
	UINT             = "uint"
	INT32            = "int32"
	UINT32           = "uint32"
	INT64            = "int64"
	UINT64           = "uint64"
	FLOAT64          = "float64"
	defaultMaxMemory = 32 << 20 // 32 MB

)

//collect parse error
type ParamError struct {
	Key string
	Msg string
}

type Parse interface {
	ParseData() // parse data to mData
	DoParse() bool
}

type Parser struct {
	Parse
	Err    ParamError
	mData  map[string]interface{}
	mParam map[string]*Param
}

//usage：
//	P := paramParse.NewParseRequest(req)
//
//		var (
//				s1 string
//				i1 uint64
//			)
//			P.StringVal(&s1, "string2","fff")
//			P.Uint64Val(&i1, "int1",nil)
//			if ok := P.DoParse(); !ok{
//				fmt.Println(P.Err.Error())
//			}
//			fmt.Println(s1)
//			fmt.Println(i1)
//	or:
//		s1 := P.String("string1","fddd")
//		s2 := P.Int64("int1", 2)
//		s3 := P.Int64("int2", 64)
//		s4 := P.Int("int3", nil)
//		if ok := P.DoParse(); !ok{
//			fmt.Println(P.Err.Error())
//		}
//
//		fmt.Println(*s1)
//		fmt.Println(*s2)
//		fmt.Println(*s3)
//		fmt.Println(*s4)

type ParseRequest struct {
	Parser
	request *http.Request
}

//var info
type Param struct {
	name         string
	ptr          unsafe.Pointer
	ty           string
	defaultValue interface{}
}

func (e *ParamError) Error() string {
	return fmt.Sprintf("key:%s -- msg:%s", e.Key, e.Msg)
}

func NewParseRequest(req *http.Request) *ParseRequest {
	mData := make(map[string]interface{})
	mParam := make(map[string]*Param)
	p := Parser{mData: mData, mParam: mParam}
	return &ParseRequest{request: req, Parser: p}
}

func (p *Parser) SetMData(name string, val interface{}) {
	p.mData[name] = val
}

//request.Form to mData
//先解析post，再解析url，当键重复时取第一个
func (p *ParseRequest) ParseData() {
	//p.request.Form

	if p.request.Form == nil {
		p.request.ParseForm()
		p.request.ParseMultipartForm(defaultMaxMemory)
	}
	for k, v := range p.request.Form {
		if len(v) > 0 {
			p.SetMData(k, v[0])
		}
	}

}

func (p *Parser) String(name string, defaultVal interface{}) *string {
	var s string
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&s), ty: STR, defaultValue: defaultVal}
	return &s
}

func (p *Parser) StringVal(ptr *string, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: STR, defaultValue: defaultVal}
}

func (p *Parser) Int64(name string, defaultVal interface{}) *int64 {
	var i int64
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&i), ty: INT64, defaultValue: defaultVal}
	return &i
}

func (p *Parser) Int64Val(ptr *int64, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: INT64, defaultValue: defaultVal}
}

func (p *Parser) Uint64(name string, defaultVal interface{}) *uint64 {
	var i uint64
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&i), ty: UINT64, defaultValue: defaultVal}
	return &i
}

func (p *Parser) Uint64Val(ptr *uint64, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: UINT64, defaultValue: defaultVal}
}

func (p *Parser) Int32(name string, defaultVal interface{}) *int32 {
	var i int32
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&i), ty: INT32, defaultValue: defaultVal}
	return &i
}

func (p *Parser) Int32Val(ptr *int32, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: INT32, defaultValue: defaultVal}
}

func (p *Parser) Uint32(name string, defaultVal interface{}) *uint32 {
	var i uint32
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&i), ty: UINT32, defaultValue: defaultVal}
	return &i
}

func (p *Parser) Uint32Val(ptr *uint32, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: INT32, defaultValue: defaultVal}
}

func (p *Parser) Int(name string, defaultVal interface{}) *int {
	var i int
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&i), ty: INT, defaultValue: defaultVal}
	return &i
}

func (p *Parser) IntVal(ptr *int, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: INT, defaultValue: defaultVal}
}

func (p *Parser) Float64(name string, defaultVal interface{}) *float64 {
	var i float64
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(&i), ty: FLOAT64, defaultValue: defaultVal}
	return &i
}

func (p *Parser) Float64Val(ptr *float64, name string, defaultVal interface{}) {
	p.mParam[name] = &Param{name: name, ptr: unsafe.Pointer(ptr), ty: FLOAT64, defaultValue: defaultVal}
}

func getFloat64(unk interface{}) (float64, error) {
	//var floatType = reflect.TypeOf(float64(0))
	//var stringType = reflect.TypeOf("")
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		return 0, errors.New("Can't convert %v to float64")
		//v := reflect.ValueOf(unk)
		//v = reflect.Indirect(v)
		//if v.Type().ConvertibleTo(floatType) {
		//	fv := v.Convert(floatType)
		//	return fv.Float(), nil
		//} else if v.Type().ConvertibleTo(stringType) {
		//	sv := v.Convert(stringType)
		//	s := sv.String()
		//	return strconv.ParseFloat(s, 64)
		//} else {
		//	return math.NaN(), fmt.Errorf("Can't convert %v to float64", v.Type())
		//}
	}
}

func getInt64(unk interface{}) (val int64, err error) {

	switch i := unk.(type) {
	case float64:
		return int64(i), nil
	case float32:
		return int64(i), nil
	case int64:
		return i, nil
	case int32:
		return int64(i), nil
	case int:
		return int64(i), nil
	case uint64:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint:
		return int64(i), nil
	case string:
		if unk == "" {
			return 0,nil
		}
		return strconv.ParseInt(i, 0, 64)
	default:
		return 0, errors.New("Can't convert to int")

	}
}

func getUInt64(unk interface{}) (val uint64, err error) {

	switch i := unk.(type) {
	case float64:
		return uint64(i), nil
	case float32:
		return uint64(i), nil
	case int64:
		if i < 0 {
			return 0, errors.New("< 0")
		}
		return uint64(i), nil
	case int32:
		if i < 0 {
			return 0, errors.New("< 0")
		}
		return uint64(i), nil
	case int:
		if i < 0 {
			return 0, errors.New("< 0")
		}
		return uint64(i), nil
	case uint64:
		return i, nil
	case uint32:
		return uint64(i), nil
	case uint:
		return uint64(i), nil
	case string:
		if unk == "" {
			return 0,nil
		}
		return strconv.ParseUint(i, 0, 64)
	default:
		return 0, errors.New("Can't convert to int")

	}
}

func getInt32(unk interface{}) (val int32, err error) {
	var i64 int64
	switch i := unk.(type) {
	case float64:
		return int32(i),nil
	case float32:
		return int32(i),nil
	case int64:
		s := strconv.FormatInt(i, 10)
		i64, err = strconv.ParseInt(s, 0, 32)
		return int32(i64), err
	case int32:
		return i, nil
	case int:
		return int32(i), nil
	case uint64:
		s := strconv.FormatUint(i, 10)
		i64, err = strconv.ParseInt(s, 0, 32)
		return int32(i64), err
	case uint32:
		s := strconv.FormatUint(uint64(i), 10)
		i64, err = strconv.ParseInt(s, 0, 32)
		return int32(i64), err
	case uint:
		s := strconv.FormatUint(uint64(i), 10)
		i64, err = strconv.ParseInt(s, 0, 32)
		return int32(i64), err
	case string:
		if unk == "" {
			return 0,nil
		}
		i64, err = strconv.ParseInt(i, 0, 32)
		return int32(i64), err
	default:
		return 0, errors.New("Can't convert to int")

	}
}

func getUInt32(unk interface{}) (val uint32, err error) {
	var tmp uint64
	switch i := unk.(type) {
	case float64:
		return uint32(i),nil
	case float32:
		return uint32(i),nil
	case int64:
		if i < 0 {
			return 0, errors.New("< 0")
		}
		s := strconv.FormatInt(i, 10)
		tmp, err = strconv.ParseUint(s, 0, 32)
		return uint32(tmp), nil
	case int32:
		if i < 0 {
			return 0, errors.New("< 0")
		}

		return uint32(i), nil
	case int:
		if i < 0 {
			return 0, errors.New("< 0")
		}
		s := strconv.FormatInt(int64(i), 10)
		tmp, err = strconv.ParseUint(s, 0, 32)
		return uint32(tmp), nil
	case uint64:
		s := strconv.FormatUint(i, 10)
		tmp, err = strconv.ParseUint(s, 0, 32)
		return uint32(tmp), nil
	case uint32:
		return i, nil
	case uint:
		return uint32(i), nil
	case string:
		if unk == "" {
			return 0,nil
		}
		tmp, err = strconv.ParseUint(i, 0, 32)
		return uint32(tmp), nil
	default:
		return 0, errors.New("Can't convert to int")

	}
}

func getInt(unk interface{}) (val int, err error) {

	switch i := unk.(type) {
	case float64:
		return int(i),nil
	case float32:
		return int(i),nil
	case int64:
		s := strconv.FormatInt(i, 10)
		i64, err := strconv.ParseInt(s, 0, strconv.IntSize)
		return int(i64), err
	case int32:
		return int(i), nil
	case int:
		return i, nil
	case uint64:
		if i > uint64(math.MaxInt64) {
			return 0, errors.New("out range")
		}
		s := strconv.FormatInt(int64(i), 10)
		i64, err := strconv.ParseInt(s, 0, strconv.IntSize)
		return int(i64), err
	case uint32:
		s := strconv.FormatInt(int64(i), 10)
		i64, err := strconv.ParseInt(s, 0, strconv.IntSize)
		return int(i64), err
	case uint:
		if strconv.IntSize == 64 && uint64(i) > uint64(math.MaxInt64) {
			return 0, errors.New("out range")
		}
		return int(i), nil
	case string:
		if unk == "" {
			return 0,nil
		}
		return strconv.Atoi(i)
	default:
		return 0, errors.New("Can't convert to int")

	}
}

func getString(unk interface{}) (val string, err error) {

	return fmt.Sprintf("%v", unk), nil

}

func (p *ParseRequest) DoParse() bool {
	p.ParseData()
	for name, v := range p.mParam {
		var (
			err error
		)
		if p.Err.Msg != "" {
			return false
		}
		var (
			val interface{}
			ok  bool
		)
		if val, ok = p.mData[name]; !ok {
			if v.defaultValue == nil {
				p.Err = ParamError{Key: name, Msg: fmt.Sprintf("[%s] not found", name)}
				return false
			} else {
				val = v.defaultValue
			}
		}
		switch v.ty {
		case STR:
			*(*string)(v.ptr), err = getString(val)
		case INT:
			*(*int)(v.ptr), err = getInt(val)
		case INT32:
			*(*int32)(v.ptr), err = getInt32(val)
		case UINT32:
			*(*uint32)(v.ptr), err = getUInt32(val)
		case INT64:
			*(*int64)(v.ptr), err = getInt64(val)
		case UINT64:
			*(*uint64)(v.ptr), err = getUInt64(val)
		case FLOAT64:
			*(*float64)(v.ptr), err = getFloat64(val)
		default:
			err = errors.New(fmt.Sprintf("%s not support", v.ty))
		}

		if err != nil {
			p.Err = ParamError{Key: name, Msg: fmt.Sprintf("[%s] %s:", name, err.Error())}
			return false
		}
	}
	return true
}
