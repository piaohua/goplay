# goplay

* this sample is a game server

## Installation
```
cd $GOPATH/src
go clone github.com/piaohua/goplay
```

## Usage:
```
cd $GOPATH/bin/ctrl
./ctrl

./ctrl build login
./login-bin -log_dir=./logs

./ctrl build dbms
./dbms-bin -log_dir=./logs

./ctrl build gate
./gate-bin -log_dir=./logs

./ctrl build hall
./hall-bin -log_dir=./logs
```

## Document
```
包引用规范
import (
    系统包

    自定义包

    第三方包
)

src/goplay/
protocol   协议文件目录 Google Protobuf version 3.0
pb         生成协议文件目录, packet unpack 文件
data       数据结构定义, 数据库连接操作 mongodb
tool       生成协议工具
game       逻辑处理 (niu, kong...),多个子目录

core       中心服务
gate       网关, 协议转发, 多个, 处理客户端连接, scoket packet unpack
hall       大厅服务,单个,处理服务注册
robot      机器人, 模拟客户端, 请求顺序login-gate
login      http请求登录,返回网关信息,单个,定时向中心分发器获取可用网关

src/
admin      后台 web 服务, base on beego
```

## package & program
```
protocol (proto)
    协议文件目录
    Google Protobuf version 3.0

pb  (package)
    生成协议文件目录
    packet unpack 文件
    rpacket runpack 机器人协议文件
    工具自动生成文件,无需手动修改

data (package)
    数据库操作
    数据结构
    参数定义

dbms (program)
    玩家数据缓存
    玩家数据中心
    logger日志中心
    唯一id管理
    房间列表管理
    房间基础信息
    后台配置数据加载
    邮箱管理
    投注活动

login (program)
    http请求节点信息
    登录节点分配
    支付回调请求

gate (program)
    websocket连接
    处理消息转发
    消息打包解包
    处理业务逻辑
    响应请求结果

hall (program)
    处理服务注册
    处理请求转发
    处理网关信息

game (package)
    处理业务逻辑
```

## TODO
* 配置文件动态加载,读取配置服务独立
* data数据库操作服务独立
* crontab服务
* admin修改包依赖,通信格式json改为pb,通信方式改为grpc
* dockerfile
* data数据库mgo操作优化
* dbms数据管理系统,玩家数据中心,优化拆分
* logging,mail,bets等服务拆分
* 版本控制
* game逻辑处理区分dbms,gate操作
* 返回消息按大小拆分发送,客户端处理粘包获取

## 启动顺序
    hall
    login
    dbms
    gate

## 停服顺序
    login
    gate
    dbms
    hall
