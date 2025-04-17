package account

//路由定义，实现接口的动态调用

import (
	"fmt"
	"reflect"
)

// 创建一个映射来存储函数名和它们的反射值
var funcs = map[string]interface{}{
	"accouont_login": accouont_login,
	"foo":            foo,
	"bar":            bar,
	"baz":            baz,
}

// 定义一个函数来动态调用映射中的函数
var callFunc = func(name string, params ...interface{}) (result []reflect.Value, err error) {
	f, ok := funcs[name]
	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	fn := reflect.ValueOf(f)
	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("%s is not a function", name)
	}

	if len(params) != fn.Type().NumIn() {
		return nil, fmt.Errorf("wrong number of arguments for %s", name)
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	result = fn.Call(in)
	return
}
