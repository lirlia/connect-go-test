# connect-go-test

## gRPCリクエストにハッシュをつけてサーバで認証する

grpc-go ではできなかったが、クライアントでリクエストに対して生成したハッシュをサーバでも同じロジックで生成し、想定しない通信が飛んでこないように真正性を担保する仕組みをやってみた。

client

```sh
❯ go run ./cmd/client/main.go
2022/12/01 03:17:16 Hello, Jane!
```

server

```sh
❯ go run ./cmd/server/main.go
2022/12/01 03:17:41 request hash matched! expected: e8cad83aa92b16b3c39f4f597bba94d493f463652e0d1778c6c10f1169b2d3f7, actual: e8cad83aa92b16b3c39f4f597bba94d493f463652e0d1778c6c10f1169b2d3f7 
```
