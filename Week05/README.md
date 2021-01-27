学习笔记

https://pkg.go.dev/golang.org/x/time/rate

# 限流
- RollingCounter  
https://github.com/go-kratos/kratos/tree/master/pkg/stat/metric
https://github.com/Netflix/Hystrix/blob/master/hystrix-core/src/main/java/com/netflix/hystrix/util/HystrixRollingNumber.java

- 限流算法  
github.com/go-kratos/kratos/pkg/ratelimit/bbr

- Consistent Hashing with Bounded Loads

# 熔断
- backoff and jitter
https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

# 降级

# 负载均衡
the choice-of-2
https://github.com/go-kratos/kratos/blob/6e49fe9ac61ee9a939d5fcd9a2df0e1d6ae250f0/pkg/net/rpc/warden/balancer/p2c/p2c.go#L78