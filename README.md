## wp-go
a simple front of WordPress build with golang.

一个go写的WordPress的前端，功能比较简单，只有列表页和详情页,rss2，主题只有一个twentyfifteen主题，插件的话只有一个简单的列表页的摘要生成和enlighter代码高亮。本身只用于展示文章，不支持restful api调用，添加评论走的转发请求到php的WordPress。因为大量用了泛型功能，所以要求go的版本在1.18以上。

#### 特色功能

- 缓存配置
- 添加评论或panic时发邮件通知，包涵栈调用和请求信息
- 简单的流量限制中间件
- 除配置文件外将所有静态资源都打包到执行文件中
- 支持密码查看，且cookie信息可被php版所验证
- 支持rss2订阅

#### 其它
用的gin框架和sqlx,在外面封装了层查询的方法。后台可以设置的比较少，大部分设置还没打通。