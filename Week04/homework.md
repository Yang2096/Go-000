# 作业

> 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

---
我对这个作业怎么写没有什么概念，只能生搬硬套，把课上说的目录结构和API生成写一下。

1. 2020.12.20 下载了 protoc，写了一些 message，进行代码的生成。
2. 2020.12.21 编写了 rpc 方法，又生成了一些代码。看 gRPC 的代码样例，抄了一点服务端的代码。
3. 2020.12.25 学习 wire 的用法


## wire 
### 普通的依赖注入
- provider  
    即普通的 go 代码，用于生成一个对象。  
    形如 `func NewXXX(initParams) (value, error)`  
    provider 设置了第二个返回值为 error 时，injector 也需要设置第二个返回值为 error

- injector
    使用了 wire 的函数，用于生成依赖链上最终的那个对象。  
    在其中使用 `wire.Build(...)` 传入 provider 并返回最终的对象。  
    **Build 参数的顺序不重要**  
    返回值不重要，只需符合函数返回类型即可。

- provider 也可以使用 wire.Build 成为其他对象的 injector

### 清理函数
类似 `context.WithCancel` 会返回的第二个参数 `cancel func()`，用于释放连接等操作。

> wire对provider的返回值个数和顺序有所规定：  
>   第一个参数是需要生成的依赖对象  
>   如果返回2个返回值，第二个参数必须是func()或者error  
>   如果返回3个返回值，第二个参数必须是func()，第三个参数则必须是error  

### 将多个 provider 组装到一起
某些 provider 经常一起使用，可以使用 provider set 将他们放在一起。
```
var ZZZSet  = wire.NewSet(NewZZZ, NewYYY, NewXXX)
```

### 对接口的绑定
依赖注入，依赖的应该是抽象的接口而非具体的某个实现。  
即某个 provider 的输入参数是一个接口，而 go 的 `NewXXX` 函数最佳实践是返回具体实现。
但是 wire 不能自动将实现和接口相关联，需要显式地将两者关联起来。

```go
// proider 
var xxxSet = wire.NewSet(NewXXX, wire.Bind(new(interfaceXXX), new(*implXXX)))

// YYY 的构造函数 NewYYY 依赖的是接口 interfaceXXX, 而 NewXXX 返回的是 *implXXX

// injector
func initYYY() *YYY {
    wire.Build(NewYYY, xxxSet)
}
```

### 对结构体应用注入
结构体中的字段初始化也可以使用 wire

```go
type Foo int
type Bar int

func ProvideFoo() Foo {/* ... */}

func ProvideBar() Bar {/* ... */}

type FooBar struct {
    MyFoo Foo
    MyBar Bar
}

var Set = wire.NewSet(
    ProvideFoo,
    ProvideBar,
    wire.Struct(new(FooBar), "MyFoo", "MyBar"))

func NewFooBar() FooBar {
    wire.Build(Set)
    return FooBar{}
}

// wire_gen.go
func NewFooBar() FooBar {
	foo := ProvideFoo()
	bar := ProvideBar()
	fooBar := FooBar{
		MyFoo: foo,
		MyBar: bar,
	}
	return fooBar
}


```

wire.Struct 中后两个参数是结构体中需要初始化的字段，与 NewSet 的前两个参数一一对应。  
也可以使用 "*" 来指明需要生成所有字段。  
在结构体定义时的字段后面添加  `wire:"-"` 来避免被初始化。
