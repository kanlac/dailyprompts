# 如何创建软链接和硬链接，它们是什么区别？

软链接（soft link）相当于快捷方式，它有不同的文件的索引节点 inode，如果原文件删除或被重命名，它就无效了；硬链接（hard link）相当于给同一份文件不同的名字，因此它们的 inode 相同，但删除一个文件，不会影响另一个。

创建 soft link：

`$ ln -s hello.txt shortcut.txt`

创建 hard link：

`$ ln hello.txt another.txt`