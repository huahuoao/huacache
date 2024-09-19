# huacache

## 简述
基于go语言实现的高性能缓存，专注于memcache缓存，在部分场景能够替代memcache提供更高性能的缓存服务
## 概念介绍

- **Client**：客户端，用于连接到 Huacache 节点进行 CRUD 操作。

- **Group**：缓存组，用户管理一组缓存。不同组之间可以存放相同的 key，互不影响。进行 CRUD 操作之前需要绑定组使用。

- **Key**：缓存 key，用于唯一标识缓存数据。

- **Value**：缓存值，存储缓存数据，可以是字符串或 Go 的结构体对象。

## 特性
- 基于Reactor模式异步事件驱动架构、BlueBell自定义高性能异步通信协议，高性能高并发低内存占用
- 原生支持集群模式，采用主从节点，实现Ketama-go一致性hash协议，无感负载均衡，分片存储。
- 基于Docker容器化构建部署，天然支持云原生环境
- 打包大小仅10MB，开箱即用
- 由字节跳动sonic提供序列化支持 https://github.com/bytedance/sonic
- 由gnet提供网络库支持 https://github.com/panjf2000/gnet
## 性能测试

## 快速入门
### Docker
```shell
git clone https://github.com/huahuoao/huacache
```
```shell
cd huacache
```
```shell
docker build -t huacache .
```
```shell
docker run -itd -p 9000:9000 huacache
```
### 源码编译
提供golang环境即可编译运行，要求go版本>=1.23.0
### Golang客户端
请移步 https://github.com/huahuoao/huacache-go
### Java客户端（待开发）
