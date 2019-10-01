CF Worker代理
===

[TOC]

另一个Goagent

# 如何使用
## 编译
运行目录下的**build.sh**, 之后程序会自动构建并创建binary.tar.gz 包, 可以从里面找到对应平台的可执行文件

## 部署
### cf-worker部署
将代码目录cf_worker中的index.js 复制到cloudflare worker编辑器中, 
并修改其中的CHECK_CODE 变量, 这个是用来校验合法来源的, 避免别人用你的服务, 
需要注意的是, CHECK_CODE字段与下面要讲到的json配置中的code要一致, 不然你自己也连不上。

### 本地部署(代理前端)

这里是基础的配置项. 下面我们每个项都来讲下用法
```json=
{
  "remote": [
    {"host": "这里填写worker的域名, 不需要加https:// 前缀, 在这个数组里面可以配置多个map", "code": "这里填写校验码, 需要跟js里面的一致"}
  ],
  "default_code": "默认的校验码, 如果remote项里面没有填写校验码, 则默认使用这个",
  "bind_host": "这里填写监听地址"
}
```

* remote 这个字段配置的是代理地址列表
* * host worker域名, 不要加 **https://前缀** 也不要加 **/** 后缀, 就单纯域名
* * code 这里涩校验码, 要跟index.js中的CHECK_CODE的值一致 
* default_code 忽略它吧, 手动配置remote项中的code值即可。
* bind_host 本地监听地址。

下面给一个配置好了的样例(需要自己替换里面的host)
```json=
{
  "remote": [
    {"host": "areyouok.dirtycat.workers.dev", "code": "hahaha"},
    {"host": "areyouok2.dirtycat.workers.dev", "code": "hahaha"}
  ],
  "default_code": "hello world",
  "bind_host": ":8080"
}
```

## 使用 
1. 自签证书配置
部署完cf-worker和代理前端后, 还需要导入本地证书, 这个代理前端本质就是一个中间人代理, 劫持你的https请求, 将其转发给worker进行处理。
自签证书可以使用目录下的 **create_cert.sh** 进行生成, 之后将其导入到chrome中。
导入路径:
三个点->设置->高级->隐私设置和安全性->管理证书->授权中心->导入生成的ca文件。

2. 浏览器代理配置
浏览器插件可以使用[SwitchyOmega](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?hl=zh-CN)
插件配置的配置代理地址为你配置中的 **bind_host**(例如上面填写的地址为 **:8080**, 那么就配置地址为**localhost**, 端口为 **8080**), 代理类型为**http**, 之后就阔以开心的上网啦。



# 原理
```flow=

                                                       +--------+
                                                       |        |
                                                       |        |
                                                       |  GFW   |
          +-------------+                              |        |                         +-----------+
          |             +----------------------------->+  X  X  |                         |           |
          |   Browser   |           BLOCK              |   XX   |                         |  Website  |
          |             +<-----------------------------+   XX   |                         |           |
          +-------------+                              |  X  X  |                         +-----------+
                                                       |        |
                                                       |        |
                                                       |        |
                                                       |        |
                                                       |        |
          +-------------+    +-----------------+       |        |      +------------+     +-----------+
          |             +--->+                 +---------------------->+            +---->+           |
          |   Browser   |    |  ProxyFrontend  |       |  PASS  |      | Cloudflare |     |  Website  |
          |             +<---+                 +<----------------------+            +<----+           |
          +-------------+    +-----------------+       |        |      +------------+     +-----------+
                                                       |        |
                                                       +--------+


```

# 优点
不需要后端服务器(自己的)

# 缺陷(待完善)
1. 部分网站打开有问题, 会出现404。
2. ws不支持。
3. 对端只能是http/https服务器, 也就是只能代理http请求, 无法代理所有的流量(受限cloudflare worker的fetch函数)。 
4. 其他。

# 其他
不熟悉前端, 前端代码是复制别人的proxy代码并改了一些东西, 可能某些地方改得有问题, 哈哈。
Cloudflare的worker支持每天10w次免费请求, 如果还不够用, 就创建多个帐号吧。