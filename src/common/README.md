common
========

# 简介
----------

go基础库。可通过init_common.sh安装所有外部依赖库。

# 配置指南

工程结构如下：

```bash
/GOPATH--------------自己随便配，线上目录是~/difoil/
  |
  /data--------------数据目录
  /logs--------------日志目录
  /bin---------------可执行文件目录
  /src---------------源码目录
    |
    /github.com------不解释
    /vidmate.com-----nemo的代码统一放这里
      |
      /common--------公共模块(http://gitlab.alibaba-inc.com/nemo/common)
      /difoil--------VidMate相关(http://gitlab.alibaba-inc.com/nemo/difoil)
      /ligo----------UGC相关(http://gitlab.alibaba-inc.com/nemo/ligo)
        |
        /cover_black_detect---黑屏自动下线服务（Owner海洪）
        /fake_maker-----------机器人服务（Owner海洪）
        /ligo_comment---------评论服务（Owner良强）
        /ligo_front-----------UGC推荐接入层（Owner理恩）
        /ligo_zipper----------UGC推荐拉链层（Owner理恩）
        /ligo_green-----------UGC去黄服务（Owner为涛）
```

需要先配路径、再拉代码，如下例：

```bash
cd GOPATH
mkdir src && cd src
mkdir vidmate.com && cd vidmate.com
git clone http://gitlab.alibaba-inc.com/nemo/common
git clone http://gitlab.alibaba-inc.com/nemo/difoil
git clone http://gitlab.alibaba-inc.com/nemo/ligo
```