package myRPC

import (
	"errors"
	"net"
	"reflect"
)

//客户端的结构体
type Client struct {
	conn net.Conn
}

// NewClient creates a new client
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

//Call将函数原型转换成函数
func (c *Client) Call(serviceName string, fptr []interface{}) {
	//调用 reflect.ValueOf 获取变量指针；
	//调用 reflect.Value.Elem 获取指针指向的变量；
	//调用 reflect.Value.SetInt 更新变量的值：
	container := reflect.ValueOf(fptr).Elem()

	//在这个函数中完成对错误的处理，与服务器端建立连接，接收服务器端数据等操作
	//客户端只需要调用该函数即可
	f := func(req []reflect.Value) []reflect.Value {
		//创建连接
		clientTrans := NewTransport(c.conn)

		//注册错误处理机制
		errorHandler := func(err error) []reflect.Value {
			//NumOut returns a function type's output parameter count
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				//Out returns the type of a function type's i'th output parameter.
				//Zero returns a Value representing the zero value for the specified type.
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			//将错误信息放置在包末尾
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		//处理包请求参数
		inArgs := make([]interface{}, 0, len(req)) //req []reflect.Value
		for i := range req {
			//将请求参数作为接口类型加入，方便后续处理
			//reflect.Value.Interface可以从反射对象可以获取 interface{} 变量,想要将其还原成最原始的状态还需要经过显式类型转换
			inArgs = append(inArgs, req[i].Interface())
		}
		err := clientTrans.Send(Data{
			ServiceName: serviceName,
			Args:        inArgs,
		})
		if err != nil {
			return errorHandler(err)
		}
		rsp, err := clientTrans.Receive()
		if err != nil {
			return errorHandler(err)
		}
		if rsp.Err != "" {
			//New returns an error that formats as the given text.
			return errorHandler(errors.New(rsp.Err))
		}
		if len(rsp.Args) == 0 {
			rsp.Args = make([]interface{}, container.Type().NumOut())
		}
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 {
				if rsp.Args[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rsp.Args[i])
				}
			} else {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), f))
}
