# GoFurry

[English](./README.md)

GoFurry 是一个公益性质的兽圈网站，同时也是一个开源仓库。本仓库收录了站点前台、数据接口、采集服务和运维后台等多个服务的源码。

整个仓库按服务拆分组织。各服务可以独立开发、独立部署，同时共享同一套整体的数据结构设计和基础设施风格。

## 仓库内容

- `gofurry-nav-frontend`：导航站前端，基于 Vue
- `gofurry-nav-backend`：导航站后端接口，基于 Go
- `gofurry-nav-collector`：导航相关数据采集服务
- `gofurry-game-backend`：游戏相关后端接口
- `gofurry-game-collector`：游戏相关数据采集服务
- `gofurry-admin`：日常运维后台，前端嵌入二进制，便于部署
- `experimental`：实验性代码，不参与生产打包
- `tools`：辅助脚本和本地工具，不参与生产打包

## 技术栈

- Go
- Fiber
- PostgreSQL
- Redis
- Vue
- Tailwind CSS

## 打包方式

根目录提供了 `build.bat`，用于打包 Linux `amd64` 生产环境产物，输出到根目录的 `build/`。

打包全部服务：

```bat
build.bat all
```

单独打包某个服务：

```bat
build.bat gofurry-nav-backend
build.bat gofurry-nav-collector
build.bat gofurry-nav-frontend
build.bat gofurry-game-backend
build.bat gofurry-game-collector
build.bat gofurry-admin
```

说明：

- Go 二进制会使用面向生产环境的体积压缩参数进行构建
- `gofurry-admin` 会先构建前端，再把前端资源嵌入最终二进制
- `experimental` 和 `tools` 默认不参与打包

## 本地开发

各服务应在各自目录内独立开发和启动。

通常流程如下：

1. 进入目标服务目录
2. 安装该服务所需依赖
3. 准备本地配置以及 PostgreSQL / Redis 连接
4. 使用该服务自己的启动方式运行

前端服务通常使用：

- `npm install` 或 `npm ci`
- `npm run dev`

Go 服务通常使用：

- `go run . serve`

## 部署说明

生产环境应使用部署者自行准备的私有配置文件。

本仓库不会附带生产环境敏感信息，不应提交以下内容：

- 生产环境 PostgreSQL 地址和账号密码
- Redis 密码
- JWT 密钥
- TLS 私钥和证书私有文件
- 生产环境 `server.yaml` 或其他私有配置文件

根目录的 `.gitignore` 已经尽量覆盖常见敏感文件，但 `.gitignore` 只能防止后续继续提交，不能清除已经进入 Git 历史的内容。如果密钥曾经上传过，应该立即轮换。

## 项目说明

GoFurry 以公益项目的方式持续维护。本仓库开源的目的，是让站点的实现方式更透明，也方便后续维护、协作和扩展。

整个代码库采用多服务结构，而不是把所有能力强行堆进一个运行时。这样做更利于部署，也更适合不同模块按各自节奏演进。

## 参与贡献

欢迎提交 Issue 和 Pull Request。

提交代码时建议：

- 尽量把改动限制在对应服务内
- 不要提交本地或生产环境敏感配置
- 除非有明确必要，否则不要随意打破现有的服务边界

## License

见 [LICENSE](./LICENSE)。
