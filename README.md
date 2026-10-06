# ansi-screen

把带 ANSI 转义的字节流喂进一块固定大小的屏幕，最后输出屏幕文本。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/ansi --sample=cursor
go run ./cmd/ansi --sample=clamp
go run ./cmd/ansi --sample=clear
go run ./cmd/ansi --sample=work
```

## 口径

- **可打印字符**：写到光标位置，光标右移一格。
- **光标定位**：`ESC [ row ; col H`（`f` 同义），行列为 1 起。
- **边界**：光标位置钳制在屏幕范围内，越界不写坏别的行。
- **清屏**：`ESC [ 2 J` 把整块屏幕清成空格，光标回到左上角。
- **颜色**：`m` 这类只改属性的序列不产生字符。

## 不变量

- 屏幕始终是 `h` 行、每行 `w` 列。
- 转义序列不会变成可见字符。
- `scanned` 不随输入长度平方级放大：两千个字符写进 80×24 的屏幕时不超过 9000。

## 输出契约

每个场景打印一行 JSON，字段 `rows` 是逐行文本；
`--sample=work` 只打印 `{"cells": N, "scanned": N}`。
