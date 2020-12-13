学习笔记

# 工程项目结构
> 路径即 package
> package 即架构设计

- 工程样例（解决gprc和http的兼容）
https://github.com/go-kratos/service-layout/

## 标准项目布局
- Standard Go Project Layout
https://github.com/golang-standards/project-layout/blob/master/README_zh.md
内含多个

- I'll take pkg over internal
https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/

## 基础库布局
- package oriented design
https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

## 服务应用项目布局
- 老师讲的这种实现是自己的思考成果，并不完全是 DDD

- DTO -> DO  mapstruct
  
- grpc服务API 到 http 路由
https://github.com/go-kratos/kratos/blob/v2/transport/http/service.go

## Lifecycle
- hook 的使用
https://github.com/go-kratos/kratos/blob/v2/app.go

- Wire
https://blog.golang.org/wire

# API 设计
## API primitive fields
- 指针解决零值-空值问题
https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto

## API error
### kratos 的实践
- 错误的判断
https://github.com/go-kratos/kratos/blob/v2/examples/kratos-demo/api/kratos/demo/errors/errors_errors.pb.go

- 生成错误实例
https://github.com/go-kratos/kratos/blob/v2/examples/kratos-demo/errors/errors.go

- 跨grpc传递错误
https://github.com/go-kratos/kratos/blob/v2/examples/kratos-demo/api/kratos/demo/v1/greeter_grpc.pb.go

- Google 的设计
https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto

## API 设计
- Google API guide
https://www.bookstack.cn/read/API-design-guide/API-design-guide-04-%E6%A0%87%E5%87%86%E6%96%B9%E6%B3%95.md
https://github.com/googleapis/googleapis

# 配置管理
## 分类
- 动态配置 
https://pkg.go.dev/expvar

## 函数式配置
- Self-referential functions and the design of options -- Rob Pike
https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html

- Functional options for friendly APIs -- Dave Cheney
https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

- 从配置文件到初始化参数
```protobuf
syntax = "proto3";
import "google/protobuf/duration.proto";
package config.redis.v1;
// redis config.
message redis {
  string network = 1;
  string address = 2;
  int32 database = 3;
  string password = 4;
  google.protobuf.Duration read_timeout = 5;
}
```

```go
func ApplyYAML(s *redis.Config, yml string) error {
  js, err := yaml.YAMLToJSON([]byte(yml))
  if err != nil {
    return err
  }
  return ApplyJSON(s, string(js))
}

// Options apply config to options.
func Options(c *redis.Config) []redis.Options {
  return []redis.Options{
    redis.DialDatabase(c.Database),
    redis.DialPassword(c.Password),
    redis.DialReadTimeout(c.ReadTimeout),
  }
}

func main() {
  // load config file from yaml.
  c := new(redis.Config)
  _ = ApplyYAML(c, loadConfig())
  r, _ := redis.Dial(c.Network, c.Address, Options(c)...)
}
```

# 测试
## 单元测试
- /api
    比较适合进行集成测试，直接测试 API，使用 API 测试框架(例如: yapi)，维护大量业务测试 case。
- /data
    docker compose 把底层基础设施真实模拟，因此可以去掉 infra 的抽象层。
- /biz
    依赖  repo、rpc client，利用 gomock 模拟 interface 的实现，来进行业务单元测试。
- /service
    依赖 biz 的实现，构建 biz 的实现类传入，进行单元测试。

## 集成测试
直接测API，用 docker-compose 在本地拉起 MySQL 或 redis