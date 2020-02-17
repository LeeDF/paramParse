package paramParse

import (
	"net/http/httptest"
	"testing"
)

func TestParser(t *testing.T) {

	inMap := map[string]string{
		"s1":      "abc",
		"int":     "10",
		"int32":   "-20",
		"uint32":  "30",
		"float64": "1.2",
	}
	req := httptest.NewRequest("GET", "/hello", nil)
	P := NewParseRequest(req)
	for k, v := range inMap {
		P.SetMData(k, v)
	}

	var (
		s1   string
		i    int
		i32  int32
		ui32 uint32
		fl64 float64
	)
	P.StringVal(&s1, "s1", nil)
	P.IntVal(&i, "int", nil)
	P.Int32Val(&i32, "int32", nil)
	P.Uint32Val(&ui32, "uint32", nil)
	P.Float64Val(&fl64, "float64", nil)
	if ok := P.DoParse(); !ok {
		t.Errorf("ParseError:%s", P.Err.Error())
	}
	if s1 != "abc" {
		t.Errorf("str error")
	}
	if i != 10 {
		t.Errorf("int error")
	}
	if i32 != -20 {
		t.Errorf("int32 error")
	}
	if ui32 != 30 {
		t.Errorf("uint32 error")
	}
	if fl64 != 1.2 {
		t.Errorf("float64 error")
	}
}

func TestParser_defaultVal(t *testing.T) {

	req := httptest.NewRequest("GET", "/hello", nil)
	P := NewParseRequest(req)
	var (
		s1 string
		i  int
	)
	P.StringVal(&s1, "s1", "abc")
	P.IntVal(&i, "int", "10")
	if ok := P.DoParse(); !ok {
		t.Errorf("ParseError:%s", P.Err.Error())
	}
	if s1 != "abc" {
		t.Errorf("str error")
	}
	if i != 10 {
		t.Errorf("int error")
	}
}

func TestParser_defaultValError(t *testing.T) {

	req := httptest.NewRequest("GET", "/hello", nil)
	P := NewParseRequest(req)
	var (
		s1 string
		i  int
	)
	P.StringVal(&s1, "s1", "abc")
	P.IntVal(&i, "int", nil)
	if ok := P.DoParse(); ok {
		t.Errorf("defaultValError")
	}
	t.Logf("ParseError:%s", P.Err.Error())

}

func TestParser_errorOut(t *testing.T) {

	inMap := map[string]string{
		"int32": "111111111111111111111111111111",
	}
	req := httptest.NewRequest("GET", "/hello", nil)
	P := NewParseRequest(req)
	var (
		i32 int32
	)
	for k, v := range inMap {
		P.SetMData(k, v)
	}
	P.Int32Val(&i32, "int32", nil)
	if ok := P.DoParse(); ok {
		t.Errorf("defaultValError")
	}
	t.Logf("ParseError:%s", P.Err.Error())
}

func TestParser_errorStringtoInt(t *testing.T) {

	inMap := map[string]string{
		"int": "abc",
	}
	req := httptest.NewRequest("GET", "/hello", nil)
	P := NewParseRequest(req)
	var (
		i int
	)
	for k, v := range inMap {
		P.SetMData(k, v)
	}
	P.IntVal(&i, "int", nil)
	if ok := P.DoParse(); ok {
		t.Errorf("defaultValError")
	}
	t.Logf("ParseError:%s", P.Err.Error())
}
