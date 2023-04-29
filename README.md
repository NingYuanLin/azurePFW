# Azure Pass Fire Wall
将该程序部署在a机器上，当a机器无法连通azure上的虚拟机时，自动切换虚拟机的ip
## 使用方法
### 1. 获取API
在azure网页右上角的cloud shell中，选择bash，然后执行
```
sub_id=$(az account list --query [].id -o tsv) && az ad sp create-for-rbac --role contributor --scopes /subscriptions/$sub_id
```
输出
```
{
  "appId": "xxx",
  "displayName": "xxx",
  "password": "xxx",
  "tenant": "xxx"
}
```
### 2. 安装
```
 go install github.com/NingYuanLin/azurePFW@latest
```
### 3. 生成配置文件
```
azurePFW config --create
```
```
Please input azure client id: ${获取API时获取到的appId}
Please input azure client secret: ${获取API时获取到的password}
Please input azure tenant id: ${获取API时获取到的tenant}
Please input subscription id: ${在订阅页获取订阅ID}
Please input resource group name: ${资源组的名字}
Please input the location of the resource (such as "Japan East"): ${地域名}
Please input the network interface name: ${网络接口的名字，可通过虚拟机=>网络=>网络接口}
Please input the ip configuration name (it will be detected automatically by default): ${可不填，也可网络=>ip配置=>名称}
```
### 4. 运行
```
azurePFW start
```
> 请注意打开虚拟机防火墙的icmp访问权限
## TODO:
* ipv6