## 设计思想

在context做修改，添加一个middlewares数组。

中间件应用在group的上，而不是直接作用在engine上。

## 具体实践

### 修改旧的

1、router的handle函数，需要把在字典树里查到的handler函数放到middleware后面

2、修改ServeHTTP函数，先查找对应的一个或多个group，然后把group内的中间件都append到context里

### 增加新的

1、

## 一些理解

1、为什么ServeHTTP函数，先查找对应的一个或多个group，然后把group内的中间件都append到context里？
（因为每个请求都可能会经由不同的一组中间件处理，即每个request来的时候根据path再生成对应需要触发的那些中间件。


