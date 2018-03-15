# SalamanderGoMonitor
基于go获得Linux系统负载（通过Linux ps命令）， web接口

# API

## GET /info
## 参数
|参数名|必选|类型|说明|
|--- |--- |--- |--- |
|sort|是|string|按什么排序，可为mem（内存），cpu|
|num|是|int|限制行数|

## 返回示例
```
查询成功
{
    "errcode": 0,
    "errmsg": "success",
    "res": "USER        PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND\nroot         12  0.0  0.3 157948  7932 pts/0    Ss+  11:50   0:00 curl -s localhost:8888/info?sort=mem\u0026num=10\nroot          1  0.4  0.2  48064  5852 ?        Ssl  11:50   0:00 ./app\nroot         19  0.0  0.1  36632  2788 ?        R    11:50   0:00 ps -aux --sort -pmem\nroot         18  0.0  0.0   4208   708 ?        S    11:50   0:00 head -10\n"
}
查询错误
{
    "errno":1,
    "text":"sort不能为空"
}
```
## 返回参数说明
|参数名|类型|说明|
|--- |--- |--- |
|errcode|int|错误码，0：成功；1：参数不正确 2:系统发生错误|
|errmsg|string|错误信息|
|res|string|查询结果，就是ps命令里的内容|
