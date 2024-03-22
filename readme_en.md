## wp-go

A WordPress front-end written in Go, with relatively simple functions,  only the list page and detail page, rss2, the theme only has two sets of twentyfifteen and twentyseventeen themes, and the plug-in only has a  simple summary generation of the list page and highlighter code  highlighting. It is only used to display articles and comments. The go version is required to be above 1.20, the newer the better. . .

#### 

#### Special feature

- Basically realize the whole site cache and prevent cache breakdown
- The list page can also highlight the syntax and format the code
- Simple plug-in extension development mechanism , support hot loading update after configuration
- can build .so  extend  theme , plug-in ,route and so on
- Rich and complicated configurations, uh, there are a lot of configurations, although most of them are optional. . .
- Send an email notification when adding a comment or panic, including stack calls and request information
- Simple traffic limit middleware, which can limit the maximum number of requests in an instant
- Package all static resources into the execution file except the configuration file
- Support password viewing, and the cookie information can be verified by the php version
- Support rss2 subscription
- Hot update configuration, switch theme, clear cache
    - kill -SIGUSR1 PID update configuration and clear cache
    - kill -SIGUSR2 PID clear cache


#### start up
```
go run app/cmd/main.go [-c configpath] [-p port]
```

#### The data show the degree of support

| page table   | Support                                                                                                                                                                         |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| List         | Home/Search/Archive/Categories/Tags/Author Paginated List                                                                                                                       |
| Details page | Display content, comments and add comments (for forwarding php processing, you  need to configure the url for adding comments in the php version)                               |
| Sidebar      | Support the old version of recent articles, recent comments, regulations,  categories, and other operation display and settings, and support the  new version of classification |

#### background settings support

- dashboard
    - Appearance
        - Widgets
            - search
            - categories
            - recent posts
            - recent comments
            - archives
            - meta
- settings
    - General Settings
        - site title
        - subtitle
    - reading
        - Blog pages show at most
        - Syndication feeds show the most recent
    - discussion
        - Other comment settings 
          - `enable|disable` threaded (nested) comments levels deep
          - Break comments into pages with top level comments per page number and the `last|first` page displayed by default
          - Comments should be displayed with the `newer|older`comments

#### Theme support

| twentyfifteen    | twentyseventeen |
|------------------|-----------------|
| site identity    | site identity   |
| colors           | color           |
| header image     | Header Media    |
| Background image | additional css  |
| additional css   |                 |

#### Plug-in mechanism

It is divided into plug-ins that modify the data of list pages and plug-ins that affect the performance of the entire program

| List page article data plugin                                  | Plugins for whole program performance                                                     |
|----------------------------------------------------------------|-------------------------------------------------------------------------------------------|
| digest Automatically generate a digest of the specified length | Enlighter code highlighting (enlighterjs plug-in needs to be installed in the background) |
|                                                                | hiddenLogin hidden login entry                                                            |

#### Others

The gin framework and sqlx used encapsulate the layer query method outside.

#### Thanks

[![jetbrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png)](https://jb.gg/OpenSourceSupport)