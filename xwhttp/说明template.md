## 设计思想

在context做修改，添加engine。
在engine层面增加全局的静态文件，没有实现group级别的

## 具体实践

### 修改旧的

1、修改c.HTML函数，增加html模板渲染。

### 增加新的

1、

## 一些理解

1、为什么ServeHTTP函数，先查找对应的一个或多个group，然后把group内的中间件都append到context里？
（因为每个请求都可能会经由不同的一组中间件处理，即每个request来的时候根据path再生成对应需要触发的那些中间件。


