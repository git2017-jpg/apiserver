apiserver-go一款基于Go构建的微服务框架，可以快速构建API服务进行业务开发，遵循SOLID设计原则

## 目录结构
```bash
├── Makefile                     # 项目管理文件
├── README.md                    # 项目说明文件
├── cmd                          # 项目启动入口文件
├── configs                      # 配置文件统一存放目录
├── internal                     # 业务目录
│   ├── cache                    # 基于业务封装的cache
│   ├── handler                  # http 接口
│   ├── model                    # 数据库 model
│   ├── repository               # 数据访问层, eg: 数据库,redis, 第三发api等
│   ├── ecode                    # 业务自定义错误码
│   ├── mock                     # mock 文件，用于单元测试
│   ├── routers                  # http服务和业务路由
│   ├── tasks                    # 任务定义和处理，包含即时、延迟和定时任务
│   └── service                  # 业务逻辑层
├── pkg                          # 项目公共库目录
├── test                         # 单元测试依赖的配置文件，主要是供docker使用的一些环境配置文件
└── scripts                      # 存放用于执行各种构建，安装，分析等操作的脚本
```