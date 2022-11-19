## 设计思想

利用trie树构建动态路由，并支持部分模板匹配。

## 具体实践

### 修改旧的

1、把router单独拆出一个文件出来，以便添加对动态路由的支持。

2、router里加上一个trie根节点的成员
### 增加新的

1、提供不同API来支持不同的content-type。

2、注意细节，需要先setHeader再调用WriteHeader，如果反过来则后面setHeader的不会生效。

## 一些理解

1、为什么router要单独拆出来呢？

因为可以发现，在engine.Get和engine.Post进行构造路由的时候，只会用到e.router.XXX，不会用到e的其他属性，所以我们可以把e.router单独作为一个文件封装起来，

2、在强调一下context什么时候可以用。所以在注册路由Get、Post、和addRouter的方法里面肯定都是不能出现context的。

3、全局函数很好做单元测试，因为入参都是可构造的，没有接受者即不需要考虑环境上下文信息。

