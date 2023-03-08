# Azure Pass Fire Wall
将该程序部署在a机器上，当a机器无法连通azure上的虚拟机时，自动切换虚拟机的ip
## 使用方法
### 1. 获取API
请参考：https://www.ddml.net/thread-11319.htm
### 2. 安装
```
 go install github.com/NingYuanLin/azurePFW@latest
```
### 3. 生成配置文件
```
azurePFW config --create
```
### 4. 运行
```
azurePFW start
```
> 请注意打开虚拟机防火墙的icmp访问权限
## TODO:
* ipv6