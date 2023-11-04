# 命令行工具 (sqlgen)

## 语法

使用以下语法从终端窗口运行 `sqlgen` 命令：

```bash
sqlgen [命令] [选项]
```

其中，命令 和 选项 分别是：

- 命令：指定您要执行的操作。
- 选项：指定可选标志。

## 快速参考

- 在当前目录创建配置文件 `sqlgen.yml`。

```bash
sqlgen init
```

- 检查 `sqlgen` 版本。

```bash
sqlgen version
```

- 为特定位置生成必要的代码，支持通配符。

```bash
# sqlgen generate [源文件]
sqlgen generate ./examples/*.go

# 或
sqlgen -c config.yml # 这会加载 `config.yml`

# 或
sqlgen # 这会加载 `sqlgen.yml`

```
