# 如何通过 set 命令提高脚本健壮性？

```bash
# TLDR;
# 在脚本的 shebang 后加上这一行，提高脚本的健壮性
set -eu

# 当脚本遇到错误（命令返回非零值）时立即退出
# 'e' stands for 'errexit'
# 建议在任何脚本中使用
set -e

# 在执行命令之前，将命令打印到标准错误输出。这有助于调试脚本
set -x
```

### 如果脚本中涉及管道操作
    
建议在需要处理包含管道（`|`）操作的 shell 脚本时使用 `set -o pipefail`。`set -o pipefail` 会改变管道命令的返回状态：如果管道中的任何命令失败（返回非零值），整个管道的返回值将为失败。这有助于在脚本中捕获和处理管道操作中的错误。

默认情况下，管道的返回值是最后一个命令的返回值。这可能导致在管道操作中的错误被忽略，因为后续命令可能仍然成功执行。

例如，考虑以下脚本：

```bash
#!/bin/bash
set -e

grep "some_pattern" input_file.txt | sort > output_file.txt
echo "Operation successful"
```

在这个例子中，如果 `grep` 命令未找到匹配的模式，它将返回非零值。然而，由于管道操作，`sort` 命令仍然会执行，管道将返回 `sort` 的返回值。因此，脚本不会因为 `grep` 命令的失败而退出，而是继续执行 `echo` 命令。

要避免这种情况并捕获管道中的错误，可以使用 `set -o pipefail`，如下所示：

```bash
#!/bin/bash
set -e
set -o pipefail

grep "some_pattern" input_file.txt | sort > output_file.txt
echo "Operation successful"
```

现在，如果 `grep` 命令失败，管道将返回非零值，因为 `set -e`，脚本将立即退出，不会执行 `echo` 命令。这样，你可以更容易地发现和处理脚本中的错误。

总之，在处理包含管道操作的脚本时，建议使用 `set -o pipefail`，以确保捕获和处理管道中的错误。
    
### 如果脚本中涉及环境变量
    
    ```bash
    # 当使用未设置的变量时，脚本会退出。这有助于捕获因未设置变量而导致的错误
    # 建议在有用到环境变量的脚本中使用
    # 字母 "u" 代表 "nounset"。"nounset" 的意思是 "未设置的变量"，表示当脚本尝试使用未设置的变量时，将触发错误并立即退出脚本
    set -u
    ```
    

### End
要查看当前 shell 环境中的所有变量和设置，可以使用不带任何选项的 **`set`** 命令。这将显示一个包含当前环境变量、函数和设置的列表。

注意：不同的 shell 可能支持不同的选项和行为。这里的示例主要针对 Bourne Again Shell（Bash），但许多其他 shell（如 sh、zsh、ksh 等）也支持类似的选项。