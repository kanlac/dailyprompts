# 码位与字节

byte 类型就是 uint8 的别名，rune 就是 int32 的别名。它们本质上都是整型！uint8 范围是 0-255，溢出了又会从 0 开始。

它们均可以以十进制、二进制、字符的形式打印出来：

```go
func printBytes() {
	var b byte
	for i, max := 0, 1000; i < max; i++ {
		fmt.Printf("i: %d, b in character: %c, b in decimal: %d, b in binary: %b, b in hexadecimal: %x\n", i, b, b, b, b)
		if b == 255 {
			fmt.Println("byte overflow!")
			break
		}
		b++
	}
}

func printRunes() {
	var r rune
	for i, max := 0, 1000; i < max; i++ {
		fmt.Printf("i: %d, r in character: %c, r in decimal: %d, r in binary: %b, r in hexadecimal: %x\n", i, r, r, r, r)
		r++
	}
}

/*  输出
byte:
...
i: 251, b in character: û, b in decimal: 251, b in binary: 11111011, b in hexadecimal: fb
i: 252, b in character: ü, b in decimal: 252, b in binary: 11111100, b in hexadecimal: fc
i: 253, b in character: ý, b in decimal: 253, b in binary: 11111101, b in hexadecimal: fd
i: 254, b in character: þ, b in decimal: 254, b in binary: 11111110, b in hexadecimal: fe
i: 255, b in character: ÿ, b in decimal: 255, b in binary: 11111111, b in hexadecimal: ff
byte overflow!

rune:
...
i: 996, r in character: Ϥ, r in decimal: 996, r in binary: 1111100100, r in hexadecimal: 3e4
i: 997, r in character: ϥ, r in decimal: 997, r in binary: 1111100101, r in hexadecimal: 3e5
i: 998, r in character: Ϧ, r in decimal: 998, r in binary: 1111100110, r in hexadecimal: 3e6
i: 999, r in character: ϧ, r in decimal: 999, r in binary: 1111100111, r in hexadecimal: 3e7
/*
```

for-range 获取码位，下标索引获取字节

```go
s := "hi你好👋"
fmt.Println("runes: ")
for _, r := range s {
	fmt.Printf("%v ", r)
}

fmt.Println("\nbytes: ")
for i := 0; i < len(s); i++ {
	fmt.Printf("%v ", s[i])
}
fmt.Println("\n\nlen(s): ", len(s))

/* output
runes: 
104 105 20320 22909 128075 
bytes: 
104 105 228 189 160 229 165 189 240 159 145 139 

len(s):  12
*/
```