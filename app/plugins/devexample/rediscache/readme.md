RedisCache use step:

- copy directory to other dir which as same level with wp-go and remove .dev suffix
- ```shell
  go mod tidy && go build -gcflags all="-N -l" --race -buildmode=plugin -o redisCache.so main.go && cp ./redisCache.so ../wp-go/plugins/
  ```
- wp-go config.yaml adds redisCache plugin