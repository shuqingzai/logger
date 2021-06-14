# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.0.6](https://github.com/shuqingzai/logger/compare/v0.0.5...v0.0.6) (2021-06-14)


### Features

* **file(文件驱动):** 1. 支持配置 file_map_size 设置文件句柄map的初始化数量,2. 支持配置 log_split_type 设置文件切割方式,3. 支持配置 log_split_size 设置文件切割大小, 3 支持 log_file_ext 设置日志文件后缀 .txt或.log ([7f194f9](https://github.com/shuqingzai/logger/commit/7f194f90c3157ddeb48e99813e4fe8d0697bc432))

### [0.0.5](https://github.com/shuqingzai/logger/compare/v0.0.4...v0.0.5) (2021-06-13)


### Features

* **file(文件驱动):** 支持配置切割文件方式: 1 支持配置按小时切割文件 2 支持配置按文件大小切割文件 ([6a636cd](https://github.com/shuqingzai/logger/commit/6a636cd38e20fbf5a51abbd5670e786f63ae17f8))
* **other(其他):** 文件驱动初始化统一入口 InitLogger() ([bcbd3b4](https://github.com/shuqingzai/logger/commit/bcbd3b4b20c36e95ff5493f88b521656ec37de47))

### [0.0.4](https://github.com/shuqingzai/logger/compare/v0.0.3...v0.0.4) (2021-06-13)

### [0.0.3](https://github.com/shuqingzai/logger/compare/v0.0.2...v0.0.3) (2021-06-13)

### 0.0.2 (2021-06-13)


### Features

* **file(文件驱动):** 新增检查文件后缀，如果不是 .log结尾的文件后缀，自动补上 ([3591562](https://github.com/shuqingzai/logger/commit/3591562cacaa69e99f6a4f38a2d764715ec709f5))
* **other(其他):** 首次commit ([2ede038](https://github.com/shuqingzai/logger/commit/2ede038f9c13dff71050bc00ad43f5ff7184f76c))
