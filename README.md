
# これはなに？
Scratch3とCHaserサーバーを仲介するプロキシープログラムです。

Scratch3（ブラウザ）の仕様上、直接CHaserサーバーとTCP/IPのSocket通信を行うことはできません。そのためこのプログラムがクライアントとサーバーのやりとりを仲介します。

```mermaid
graph TD;
  id1(Chrome) --> id2(Scratch3(Extension)) -- (WebSocket) --> id3(go-scratch-chaser-proxy)
  id3 -- (socket) --> id4(CHaser Server)
  id4 --> id3
  id3 --> id2
```

# 使い方

```
go mod tidy
go build
./chrome-chaser-proxy
```
