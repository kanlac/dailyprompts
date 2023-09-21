# 计算文件大小

## 表面大小与实际大小有什么区别？为什么会造成这种差异？

**实际大小（Actual Size）** 是文件或目录实际占用的磁盘空间，考虑到了文件系统的块大小。即使一个只有 1 字节的文件也会占用一个完整的磁盘块（块大小因文件系统而异，通常为 4KB 或更大）。

**表面大小（Apparent Size）** 是文件或目录中实际数据的大小，不考虑文件系统的块大小或任何其他底层存储细节，通常比实际大小要小。**传输文件时，传输的是文件的表面大小而不是实际大小**。

新建一个空文本文件，du 查看是 0，写入一个字符，du 查看是 4K（即便用 `du --block-size=1 a.txt` 查看也是 4096），而用 `ll` 命令或者 `du --apparent-size` 查看是 2B。

## `du` 的默认单位是？

严格来说，`du` 命令默认的单位不是 KB，而是块的数量！只不过系统默认的块大小是 1024B，所以可能 du 刚好和 KB 一样大！注意块的数量最小是 1，没有小数。

以下四个命令是等价的：

- `du foo.txt`
- `du -k foo.txt`
- `du --block-size=1K foo.txt`
- `du --block-size=1024 foo.txt`

## 用什么命令可以以字节单位显示一个文件/目录的实际大小？

不是 `-b` 哦，和 `-k` `-m` 不一样，`-b` 看的表面大小，要注意区分。答案是：`du --block-size=1`。

## 使用 Go 计算 *nix 系统的文件大小

在 Go 中，用 `os.FileInfo.Size()` 获取的大小是逻辑大小，不是真正磁盘使用的空间。真正的磁盘大小是 `Blocks` 字段和块大小（通常是 512 字节）的乘积，可以使用 `syscall` 包获取：

```go
// get the actual disk size of a single file
func getFileDiskSize(path string) (int64, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return 0, errors.Wrapf(err, "cannot get the file stat of %s", path)
	}
	// In Unix, the block size is usually 512 bytes
	// Multiply the block count by 512 to get actual disk size
	return stat.Blocks * 512, nil
}

// get the actual disk size of a directory recursively
func getDirDiskSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Sys().(*syscall.Stat_t).Blocks * 512
		return nil
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}
```

在文件系统中，块（block）是用于存储数据的基础单位。操作系统通常以块为单位从硬盘读取或写入数据。块大小一般是硬盘和文件系统类型的属性，常见的大小包括 512 字节、1KB、4KB 等。

通过实践比较 `du --block-size=1024` 和 `du`（默认不指定）的输出发现是一致的，由此可以说明 du 的默认块大小是 1024。然而 Go 程序中以 512B 的块大小计算出的结果是最近似的。**看来文件系统的块大小和 du 的块大小似乎不是一回事**。

实际观察发现，`du` 计算出来的大小会要比 Go 的 `syscall` 包算出来的大一些，也许是因为前者除了文件的实际数据块，还可能包括额外的信息，如元数据、间接块等。

经过协商，leader 表示就直接用 du 好了，一个原因是它会比较准确，另一个是省的研究统计的信息为什么不一样的问题。

新版的用 du 实现的函数：

```go
func getDiskSize(path string) (int64, error) {
	var cmd *exec.Cmd

	// Check if path is a file or directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot get the file stat of %s", path)
	}

	if fileInfo.IsDir() {
		// It's a directory
		cmd = exec.Command("du", "-s", "--block-size=1", path)
	} else {
		// It's a file
		cmd = exec.Command("du", "--block-size=1", path)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return 0, errors.Wrapf(err, "cannot get the disk size of %s", path)
	}

	s := strings.Fields(out.String())[0] // First field is the size in bytes
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot parse disk size: %s", s)
	}
	return size, nil
}
```