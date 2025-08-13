# aegis

cli file encryption tool using aes256-gcm and chacha20-poly1305 algorithms. implements minio/sio dare format with pbkdf2 key derivation.

## build

using just:
```bash
just build
```

using go:
```bash
go build -o aegis cmd/aegis/main.go
```

## test

```bash
go test -bench=.
```

## usage

encrypt a file:
```bash
aegis encrypt input.txt -p password
aegis encrypt input.txt output.enc -p password
```

decrypt a file:
```bash
aegis decrypt input.txt.enc -p password  
aegis decrypt input.enc output.txt -p password
```