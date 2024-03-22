## wp-go

[en readme](https://github.com/fthvgb1/wp-go/blob/master/readme_en.md)

一个go写的WordPress的前端，功能比较简单，只有列表页和详情页,rss2，主题只有twentyfifteen和twentyseventeen两套主题，插件的话只有一个简单的列表页的摘要生成和enlighter代码高亮。本身只用于展示文章及评论。要求go的版本在1.20以上，越新越好。。。

#### 特色功能

- 基本实现全站缓存，并且可防止缓存击穿
- 列表页也可以高亮语法格式化显示代码
- 简易插件扩展开发机制、配置后支持热加载更新
- 使用.so扩展主题、插件、路由等
- 丰富繁杂的配置，呃，配置是有点儿多，虽然大部分都是可选项。。。
- 添加评论或panic时发邮件通知，包涵栈调用和请求信息
- 简单的流量限制中间件，可以限制全瞬时最大请求数量
- 除配置文件外将所有静态资源都打包到执行文件中
- 支持密码查看，且cookie信息可被php版所验证
- 支持rss2订阅
- 热更新配置、切换主题、清空缓存
    - kill -SIGUSR1 PID 更新配置和清空缓存
    - kill -SIGUSR2 PID 清空缓存

#### 运行
```
go run app/cmd/main.go [-c configpath] [-p port]
```

#### 数据显示支持程度

| 页表  | 支持程度                                        |
|-----|---------------------------------------------|
| 列表页 | 首页/搜索/归档/分类/标签/作者 分页列表                      |
| 详情页 | 显示内容、评论并可以添加评论(转发的php处理，需要配置php版的添加评论的url)  |
| 侧边栏 | 支持旧版  近期文章、近期评论、规档、分类、其它操作  显示及设置, 支持新版  分类 |

#### 后台设置支持程度

- 仪表盘
    - 外观
        - 小工具
            - 搜索
            - 规档
            - 近期文章
            - 近期评论
            - 分类
            - 其它操作

- 设置-
    - 常规
        - 站点标题
        - 副标题
    - 阅读
        - 博客页面至多显示数量
        - Feed中显示最近数量
    - 讨论
      - 其他评论设置
          - `启用|禁止`评论嵌套，最多嵌套层数
          - 分页显示评论，每页显示评论条数，默认显示`最前/后`页
          - 在每个页面顶部显示 `新旧`评论

#### 主题支持程度

| twentyfifteen | twentyseventeen |
|---------------|-----------------|
| 站点身份          | 站点身份            |
| 颜色            | 颜色              |
| 页眉图片          | 页眉媒体            |
| 背景图片          | 额外css           |
| 额外css         |                 |

#### 插件机制

分为对列表页文章数据的修改的插件和对影响整个程序表现的插件

| 列表页文章数据插件           | 整个程序表现的插件                            |
|---------------------|--------------------------------------|
| digest  自动生成指定长度的摘要 | enlighter 代码高亮(需要在后台安装enlighterjs插件) |
|                     | hiddenLogin 隐藏登录入口                   |

#### 其它

用的gin框架和sqlx,在外面封装了层查询的方法。

#### 鸣谢

<a href="https://jb.gg/OpenSourceSupport"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." width="50%"></a>

