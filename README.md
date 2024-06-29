# Telegram 付费进群机器人

一个方便部署的Telegram付费机器人，用户支付成功后会发送邀请链接给用户进入群组

### 特性
    ①使用易支付通用接口
    ②使用sqlite作为数据库，闭着眼睛部署
    ③机器人监控群组进群状态，如果进群的人未支付，则踢出

### DEMO


https://github.com/yuimoi/tg_pay_gate/assets/119736684/9ec4f9db-d81a-4968-afb2-d797a6a28728



### 使用说明
自行编译或者下载右侧Release编译好的文件，直接运行即可，运行文件的同一目录下要附带.env文件夹并在文件夹中写好配置文件。

机器人拉进群组，并设置为管理员，群组可以设置为私有，以防直接进入

### 支付说明
使用通用易支付接口，[Mopay](mopay.vip)或自行寻找易支付进行接入 

因为需要接收支付成功的回调信息，所以需要对外开放http请求，运行时附带参数可修改http运行端口 `--port 8086`

### nginx反代参考配置
    location / {
        proxy_pass http://127.0.0.1:8086;
    }


### 配置说明
配置文件保存在`.env`文件夹中，修改即可
#### config.json
    "tg_bot_token": 机器人的token，前往BotFather获取，并设置为群组的管理员，群组
    "group_id": 群组id，可以在web端点击群组后在浏览器地址栏中可以看到（群组id前面是带负号的）
    "price": 进群费用（使用双引号括起来）
    "host": 绑定的域名，用于易支付发起订单时，拼接回调地址
    
    "proxy": 代理，一般不用开


#### epay_config.json
    "pid": 易支付的pid，数字用双引号括起来
    "key": 易支付的key
    "url": 易支付发起订单的url，有些易支付后台显示的url不以submit.php结尾，可能要自己加上
    "pay_type": 易支付的支付类型
    
    "notify_url": 保持默认


Tg: [@nulllllllll](https://t.me/nulllllllll)
