# Golang http.Server安全退出：容易被误用的Shutdown()方法

这个仓库展示了一种正确的（希望是）和一种错误的安全退出`http.Server`的方法。

详情可以参考我写的这篇文章：https://nanmu.me/zh-cn/posts/2021/go-http-server-shudown-done-right/

## 运行示例

如果你有现成的Go环境，可以直接安装：

```bash
# 一种正确的安全退出姿势
go install github.com/nanmu42/Go-HTTP-Server-Gracefully-Shutdown-Done-Right/better-way@latest
# 一种不正确的安全退出姿势
go install github.com/nanmu42/Go-HTTP-Server-Gracefully-Shutdown-Done-Right/wrong-way@latest
```

然后就可以启动和访问了：

```bash
$ better-way 
2021/09/23 00:16:03 listening on port 3100...
```

```bash
$ wrong-way 
2021/09/23 00:19:00 listening on port 3000...
```

启动后可以用浏览器进行访问，在访问的同时退出服务试试看效果。如果手速不够快，可以`-h`看下选项。

Windows和Mac没有测试，由于Syscall的不同，可能需要改监听的信号名称才能编译。

## License

MIT
