package myRPC

import (
	"bytes"
	"encoding/gob"
)

//数据包的定义
type Data struct {
	ServiceName string        //服务名称
	Args        []interface{} //传递的参数
	Err         string        //socket的错误
}

//数据包的序列化
func encode(data Data) ([]byte, error) {
	var buf bytes.Buffer            //定义一个空的字节缓冲区
	encoder := gob.NewEncoder(&buf) //NewEncoder returns a new encoder that will transmit on the &buf.
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//数据包的反序列化
func decode(b []byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	return data, nil
}
