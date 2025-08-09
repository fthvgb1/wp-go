#### limitation

- only linux platform
- the golang version, build flags of plugin same as wp-go
- same plugin name only loaded once

#### build customized plugin steps

- set wp-go config.yaml `pluginPath`
- plugin package must be named `main`
- plugin export func signature `func Pluginname( *wp.Handle) `
- plugin's `go.mod` replace wp-go package => local wp-go path
- ```shell
  go build buildmode=plugin [other flags same as wp-go] -o pluginname.so && mv pluginname.so pluginPath
  ```

#### load plugin

- add `pluginname` into wp-go config.yaml `plugins` item
- reload wp-go configuration (`kill -SIGUSR1 wp-go pid`)