// +build wireinject

package main

import (
	"github.com/google/wire"
)

// InitializeEvent 声明injector的函数签名
func InitializeEvent(msg string) Event {
	wire.Build(NewGreeter, NewMessage, NewEvent)
	return Event{} //返回值没有实际意义，只需符合函数签名即可
}

type Foo int
type Bar int

func ProvideFoo() Foo {
	return 0
}

func ProvideBar() Bar {
	return 1
}

type FooBar struct {
	MyFoo Foo
	MyBar Bar
}

var Set = wire.NewSet(
	ProvideFoo,
	ProvideBar,
	wire.Struct(new(FooBar), "MyFoo", "MyBar"))

func initFooBar() FooBar {
	wire.Build(Set)
	return FooBar{}
}
