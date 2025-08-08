## wp-go

This is a WordPress frontend written in Golang with Gin framework and sqlx db driver. Now had completed the displaying
of posts' listing, detail, comment, rss2 and so forth features and two themes. The lowest version of Golang requirement
is 1.23 and the newer, the better.

#### Primary features

- Substantially implemented the caching of whole site
- Use Golang plugin package to hook code
- Hot reloading configuration and plugin
- Traffic limit
- Rss2

#### Start up

rename config.example.yaml to config.yaml

```
go run app/cmd/main.go [-c configpath] [-p port]
```

#### Hot update configuration, change theme, load plugin , clear cache

- `kill -SIGUSR1 PID` update configuration and clear cache
- `kill -SIGUSR2 PID` clear cache

It will come into effect at next request.

#### Dynamically load plugin mechanism

You can hook totally 90% internal codes as or use the Golang plugin. So you can hook innate middlewares, routes, cache
driver and so on registered functions or add new themes same as php, and it would not stifle the performance.

[see example](app/plugins/devexample)

#### background settings support

You can set follow configurations in php WordPress background.

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

                                                        
