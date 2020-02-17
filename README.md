# paramParse
用于http.Request 集中解析参数， 并返回错误

## Exmaple
```
func Hello(w http.ResponseWriter, r *http.Request)  {
	var (
		paramInt64   int64
		paramInt32 int32
		paramStr string
	)

	P := paramParse.NewParseRequest(r)
	//入参：变量指针、 表单key、默认值， 默认值为nil时， 相当于required

	P.Int64Val(&paramInt64, "paramInt64", nil)
	P.Int32Val(&paramInt32, "paramInt32", "0")
	P.Str32Val(&paramStr, "paramInt32", nil)
	if ok := P.DoParse(); !ok {
		fmt.Printf("error parse form fail :%s", P.Err.Error())
		return
	}
}


```