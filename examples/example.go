package examples

import (
	"fmt"
	"paramParse"
	"net/http"
)

func main() {
	http.HandleFunc("/ptr", Hello)
	_ = http.ListenAndServe(":12345", nil)

}

func Hello(w http.ResponseWriter, req *http.Request) {
	P := paramParse.NewParseRequest(req)

	//返回指针
	s1 := P.String("string1", "fddd")
	s2 := P.Int64("int1", 2)
	s3 := P.Int64("int2", 64)
	s4 := P.Int("int3", nil)
	P.DoParse()
	if ok := P.DoParse(); !ok {
		fmt.Println(P.Err.Error())
	}
	fmt.Println(*s1)
	fmt.Println(*s2)
	fmt.Println(*s3)
	fmt.Println(*s4)

	//返回值
	//P := paramParse.NewParseRequest(req)
	//var (
	//	s1 string
	//	i1 uint64
	//)
	//P.StringVal(&s1, "string2","fff")
	//P.Uint64Val(&i1, "int1",nil)
	//if ok := P.DoParse(); !ok{
	//	fmt.Println(P.Err.Error())
	//}
	//fmt.Println(s1)
	//fmt.Println(i1)

	w.Write([]byte("hello"))
	return
}
