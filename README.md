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
pb         生成协议文件目录, pack unpack 文件
data       数据结构定义, 数据库连接操作 mongodb
logger     日志服务,结构可以在protocol协议中注册

test       测试
tool       生成协议工具

gate       网关, 协议转发, 多个, 处理客户端连接, scoket pack unpack, 版本控制
core       中心分发器,单个,网关,登录服,逻辑服注册,中心服, 只负责服务的注册和分配处理

room       房间服务,remote获取通信
game       逻辑服 (niu, kong...),多个,多个子目录,只是桌上操作,拓展多少个牌桌
hall       中心服,大厅,单个,玩家数据状态缓存,后台配置数据加载,remote获取房间服务
robot      机器人, 模拟客户端, 请求顺序login-gate
login      http请求登录,返回网关信息,单个,定时向中心分发器获取可用网关

src/
admin      后台 web 服务
```

## package & program
>protocol</br>
>>协议文件目录</br>
>>Google Protobuf version 3.0</br>

>pb</br>
>>生成协议文件目录</br>
>>pack unpack 文件</br>
>>机器人协议文件</br>
>>工具自动生成文件,无需手动修改</br>

>data</br>
>>1个端口</br>
>>数据库操作</br>
>>room唯一id,列表管理,房间基础信息</br>
>>logger日志写入记录,pid</br>
>>玩家数据中心</br>

>login</br>
>>2端口,1个rpc,1个http</br>
>>加密返回网关信息</br>

>gate</br>
>>2个端口,1个rpc,1个websocket</br>
>>处理消息转发</br>
>>消息打包解包</br>
>>各个游戏状态管理</br>

>hall</br>
>>1个端口</br>
>>处理服务注册,连接数</br>
>>玩家数据缓存</br>
>>网关信息</br>
>>处理请求消息</br>

>game</br>
>>1个端口</br>
>>处理业务逻辑</br>
>>暂时包括排位,排行</br>
>>TODO 支付回调</br>

>web
>>2端口,1个rpc,1个http</br>

## TODO
* 配置文件动态加载,读取配置服务独立
* data数据库操作服务独立
* crontab服务
* admin修改包依赖, 通信格式json改为pb, 通信方式改为grpc
* dockerfile
* data数据库mgo操作优化
* dbms数据管理系统,玩家数据中心,优化拆分
* logging

client -(http)- login -(rpc)- hall

client -(http)- web -(rpc)- hall

gate -(rpc)- game
gate -(rpc)- hall

game -(rpc)- room
game -(rpc)- logger

//
hall |--- web    |---  client
     |--- login  |---  client
     |--- gate   |---  client
                 |---  game  |--- logger
                 |---        |--- room

## 启动顺序
    hall
    login
    dbms
    gate
