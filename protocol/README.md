# goplay
game server protocol

## Usage

## Protocol
Google Protobuf version 3.0

## Document

```
buy   code 3000+ 购买,商城
chat  code 2000+ 聊天,广播,公告
login code 1000+ 登录,注册
room  code 4000+ 游戏逻辑(niuniu,free,kong,classic,paohuzi),进入,离开,踢出,操作,投票,结算,记录,状态,刮奖
user  code 1020+ 全局配置,玩家数据,基础数据更新,绑定代理,银行,心跳,破产,抽奖活动,宝箱活动,vip,邮件,排位赛数据
vo    公共

dan      code 1020+ 排位赛
mail     code 1020+ 邮件
prize    code 1020+ 抽奖,宝箱活动
classic  code 4000+ 经典场
free     code 4000+ 自由场
niu      code 4000+ 牛牛
phz      code 4000+ 跑胡子
lottery  code 4000+ 全民刮奖

pub   公共

内部通信消息
protos
logger
data
```
