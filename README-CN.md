# Sub-Store-Manager-Cli

[English](./README.md) | 简体中文

用于管理 [Sub-Store](https://github.com/sub-store-org/Sub-Store) 后端服务的命令行工具。由于该工具基于 [Docker](https://www.docker.com/)，您需要先安装 Docker。


# 安装

您可以通过运行以下命令来安装这个工具：

```bash
curl -sSL https://sub-store-org.github.io/resource/ssm/install.sh | bash
```

或者您可以下载 [Release 文件](https://github.com/DesnLee/Sub-Store-Manager-Cli/releases)，并手动配置环境变量和程序名称。


# 使用方法

如果您正在使用脚本安装则可以直接执行以下命令。如果您正在进行手动安装，请确保每个命令中的程序名称匹配你的可执行文件名。

### new

创建一个新的 Sub-Store Docker 容器并运行，如果镜像不存在将自动构建。

```bash
ssm new
```

该命令支持以下 `flag`：

- `--interface` 或 `-i` ：当您传递 `-i` 标志时，将创建前端容器而不是后端容器。默认行为是创建一个后端容器。
 
- `--name` 或 `-n` ：容器的唯一名称，默认名称为 `ssm-backend`，如果您传递了 `-i` 标标志来创建一个前端容器，则默认名称为 `ssm-frontend`。此名称将用于管理持久化数据，只要不手动删除该名称的持久化数据，或者使用 `ssm delete -c` 标志执行删除操作，无论是如何删除/重建容器，只要使用此名称都可以访问该数据。

- `--version` 或 `-v` ：一个 [Sub-Store Release](https://github.com/sub-store-org/Sub-Store/releases) 的版本字符串，默认获取最新版本。如果您传递了 `-i` 标志来创建前端容器，则 `-v` 标志将被忽略，它总是使用最新版本的前端。

- `--port` 或 `-p` ：指定端口映射，默认为 `3000`，且必须可用，如果您传递了 `-i` 标标志来创建一个前端容器，则默认端口为 `80`。如果你想使用域名访问服务，则需要使用反向代理工具（如 Nginx 或 Caddy）手动代理该端口。


### update

更新一个 Sub-Store Docker 容器，确保镜像已经存在且正在运行。

```bash
ssm update
```

该命令支持以下 `flag`：

- `--name` 或 `-n` ：一个正在运行的容器名称，默认名称为 `ssm-backend`。

- `--version` 或 `-v` ：一个 [Sub-Store Release](https://github.com/sub-store-org/Sub-Store/releases) 的版本字符串，默认获取最新版本。如果您更新目标为前端容器，则 `-v` 标志将被忽略，它总是使用最新版本的前端。


### start

启动一个未在运行的 Sub-Store Docker 容器，默认名称为 `ssm-backend`。

> 基本等价于 `docker start <name>`.

```bash
ssm new <name>
```


### stop

停止一个正在运行的 Sub-Store Docker 容器，默认名称为 `ssm-backend`。

> 基本等价于 `docker stop <name>`.

```bash
ssm stop <name>
```


### delete

删除一个 Sub-Store Docker 容器，默认名称为 `ssm-backend`。

> 基本等价于 `docker rm <name>`.

```bash
ssm delete <name>
```

该命令支持以下 `flag`：
 
- `--clear` 或 `-c` : 同时删除容器的持久化数据。如果删除的为前端镜像，则 `-c` 标志将被忽略，因为前端镜像没有持久化数据。


### list

列出所有 Sub-Store Docker 容器。

> 基本等价于 `docker ps -a` 并过滤以 `ssm` 镜像启动的容器。

```bash
ssm ls
```


### version

查看当前 Sub-Store-Manager-Cli 的版本。

```bash
ssm version
```


# 卸载

如果您使用脚本安装则可以直接执行以下命令。如果您使用手动安装，请手动移除您的可执行文件。

```bash
rm -rf /usr/local/bin/ssm
```

如果您想同时删除持久化数据，可以执行以下命令：

```bash
rm -rf ~/.ssm
```


# License
GPL-3.0 License
