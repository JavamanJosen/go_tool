package tool_bzip

import (
	"bytes"
	"compress/zlib"
	"io"
)

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//func main() {
	//buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
	//	47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
	//b := bytes.NewReader(buff)
	//r, err := zlib.NewReader(b)
	//if err != nil {
	//	panic(err)
	//}
	//io.Copy(os.Stdout, r)
	//r.Close()
	//
	//zip := DoZlibCompress([]byte("hello, world\n"))
	//fmt.Println(zip)
	//fmt.Println(string(DoZlibUnCompress(zip)))
//}
