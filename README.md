# fake115-go
115网盘助手Go版本，导出的目录什么样，导入就什么样😲，达到跟雷达功能一致的效果，且没有目录大小限制。



目前成功导出、导入多个TB级别的包，大的包甚至高达200TB，70万个文件😄。



## Installation



#### Clone



```bash
git clone https://github.com/gawwo/fake115-go
cd fake115-go
```



#### Get the dependencies



```bash
go get ./...
```



#### Build

```bash
go build -o fake115 .
```



## Getting Started



#### Prepare



从115浏览器中获取自己登陆后的cookie。



可以在程序目录下创建一个`cookies.txt`的文件存放cookie，也可以在使用时，添加`-c`参数设置cookie。



#### Export



- cid是指115文件夹的id，F12的开发者工具中查看network能找到它。



Usage:

```bash
fake115 -c "cookiexxxxx" <your export cid>
```



#### Import



Usage:

```bash
fake115 -c "cookiexxxxx" <your import cid> <import json file path>
```

