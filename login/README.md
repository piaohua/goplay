# login

* http请求登录
* 返回网关信息
* 定时向中心分发器获取可用网关

## Installation

```
./ctrl build login linux
```

## Usage:

```
./login -log_dir=./logs > /dev/null 2&>1 &

./ctrl start login
./ctrl stop login
```
