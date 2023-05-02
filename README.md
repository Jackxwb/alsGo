# alsGo - Golang Implementation of Another Looking-glass Server

Another Looking-glass Server 的Golang后端实现，[原项目地址](https://github.com/wikihost-opensource/als)

项目初衷是：摆脱Docker的依赖也能运行的现代化 Looking-glass 服务

> 警告：目前不适合高并发使用，且似乎存在内存泄漏的情况

## 使用方法

1. 下载并解压对应系统的编译文件，放入独立的文件夹中，该文件夹为主目录，拷贝项目中的public文件夹到主目录中
2. 项目依赖外部环境的 ping、iperf3、[speedtest-CLI(Ookla®)](https://www.speedtest.net/zh-Hans/apps/cli) ，请确保系统内可以正常调用，其中 iperf3 和 speedTest 可以直接把程序放在本程序的app文件夹里面(主目录\app)
3. 给编译后主程序添加可执行权限，直接运行即可

```cmd
目录结构(以Linux为例)

主目录
│  config.yml   (配置文件，不修改默认参数的情况下可以不需要)
├─app   (可选，放到全局环境也可以)
│  ├─iperf3
│  └─speedtest
├─public
│  ├─speedtest-static   (Web页面测速文件)
│  └─static
│      │  speedtest_worker.js   (前端测速脚本，项目来自https://github.com/librespeed/speedtest)
│      └─dist   (前端VUE打包文件)
│          │  favicon.ico
│          │  index.html
│          └─assets
│                  index.262fc5ee.js
│                  index.de51c7d3.css
│                  iPerf3.47688fd2.js
│                  SpeedtestDotNet.57c11e21.js
│                  vue3-apexcharts.common.ecc4956e.js
│
└─als   (Golang 后端编译文件)
```

## 硬件需求

- 理论上可以运行Golang的都可以，提供Linux amd64 和 Window 编译文件，其他版本可以自行下载项目后编译
- 内存需求：启动内存在3M左右，似乎有内存泄漏的地方，程序有定时任务，每5分钟主动回收一次内存，无请求状态下+主动回收，内存一般＜10MB。网页测速，每次会申请50MB左右的内存填充随机数据，内存占用是瞬间的，但是高并发下如果把系统内存吃完了会导致程序异常退出
- 磁盘需求：默认配置会在 `speedtest-static` 目录下创建静态测速文件：`1MB.test` `10MB.test` `100MB.test` 共占用111MB，可在配置文件中修改

## 配置

### 在主目录下创建config.yml配置文件

完整配置文件

```json
listenaddr: 0.0.0.0
port: 4000
testfiles: [1048576, 10485760, 104857600]
baseinfo:
  bandwidth: false
  display_traffic: true
  display_speedtest: true
  utilities_ping: true
  utilities_traceroute: true
  utilities_iperf3: true
  utilities_iperf3_port_min: 30000
  utilities_iperf3_port_max: 31000
  utilities_speedtestdotnet: true
  utilities_fakeshell: false
app:
  deftimeout: 90
  pingresulttemplate: ""
  pingresultnotlinktemplate: ""
  pingresultnotpingtemplate: ""
  iperffixedport: 0
netset:
  networkadapter:
    onlyshow: []
    dontshow: ["VMware Network Adapter"]
  hideexternalip: false
```


| 键         | 子树                      | 默认值                         | 说明                                                                                                                                                                                                                                                                                                                     |
| ------------ | --------------------------- | -------------------------------- |------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| listenaddr |                           | 0.0.0.0                        | 监听的IP地址                                                                                                                                                                                                                                                                                                                |
| port       |                           | 4000                           | 监听端口                                                                                                                                                                                                                                                                                                                   |
| testfiles  |                           | [1048576, 10485760, 104857600] | 静态测速文件大小，int数组类型，单位是B。比如：1048576B=1MB。<br/>原项目还有一个1G的测速文件，可自行添加 `,1073741824` 这段到数组里面来实现                                                                                                                                                                                                                               |
| baseinfo   |                           |                                | 前端页面配置项目，影响web显示效果                                                                                                                                                                                                                                                                                                     |
|            | bandwidth                 | false                          | 是否以带宽单位显示。如："Mbps"、"Gbps" 等。默认使用 "MB"、"GB" 等显示                                                                                                                                                                                                                                                                         |
|            | display_traffic           | true                           | 是否显示`服务器流量图`模块                                                                                                                                                                                                                                                                                                         |
|            | display_speedtest         | true                           | 是否显示`服务器网络测速`模块                                                                                                                                                                                                                                                                                                        |
|            | utilities_ping            | true                           | `网络工具`模块是否显示`ping`测试                                                                                                                                                                                                                                                                                                   |
|            | utilities_iperf3          | true                           | `网络工具`模块是否显示`iperf3`测速                                                                                                                                                                                                                                                                                                 |
|            | utilities_iperf3_port_min | 30000                          | iperf3服务随机端口最小值                                                                                                                                                                                                                                                                                                        |
|            | utilities_iperf3_port_man | 31000                          | iperf3服务随机端口最大值                                                                                                                                                                                                                                                                                                        |
|            | utilities_speedtestdotnet | true                           | `网络工具`模块是否显示`Speedteset.net`(speedtest-CLI测速)<br/>在服务器上运行由Ookla®提供的测速工具                                                                                                                                                                                                                                                |
|            | utilities_fakeshell       | false                          | `网络工具`模块是否显示`终端`。默认不显示。<br/>Golang版本不提供终端实现，也没有实现的计划，需要使用该功能请自行实现或使用原版Docker版本                                                                                                                                                                                                                                         |
| app        |                           |                                | 外部程序一些设置。对应Web页面中`网络工具`里面的功能实现相对应的设置                                                                                                                                                                                                                                                                                   |
|            | deftimeout                | 90                             | 运行外部程序的默认超时时间，超时后自动停止程序所在线程。单位为 秒 ，原项目为60秒                                                                                                                                                                                                                                                                             |
|            | pingresulttemplate        |                                | ping正常结果的正则模板，方便适配ping输出带其他语言的情况。如果模板配置出现问题，则不会返回数据到前端页面<br/>参考配置：<br/> ping结果 `64 bytes from xx.xx.xx.xx (xx.xx.xx.xx): icmp_seq=1 ttl=45 time=97.6 ms`<br/>ping模板 `from (?P<ip>.*?) \(.*\): icmp_seq=(?P<seq>.*?) ttl=(?P<ttl>.*) time=(?P<time>.*?) ms`<br/>(模板里的正则需要分组；程序默认有适配中文Window的ping输出，使用英文版Window的请自行配置模板) |
|            | pingresultnotlinktemplate |                                | ping超时的正则模板。程序参考原项目指令，Linux的默认ping命令是`ping 目标 -O -c 10` <br/>参考配置：<br/>ping结果 `no answer yet for icmp_seq=1`<br/>ping模板 `no answer yet for icmp_seq=`                                                                                                                                                                  |
|            | pingresultnotpingtemplate |                                | ping 无法解析目标时的模板。常在ping一个无法连通的域名时出现<br/>参考配置：<br/>ping结果 `bad address`<br/>ping模板 `bad address`                                                                                                                                                                                                                         |
|            | iperffixedport            | 0                              | 是否强制固定iperf3服务的端口，方便开放防火墙。默认为0，即不启用固定，使用随机端口。<br/>若启用该配置，上面设置的随机端口范围将无效                                                                                                                                                                                                                                                |
| netset     |                           |                                | 网络相关的设定                                                                                                                                                                                                                                                                                                                |
|            | networkadapter            |                                | 网络适配器相关的设定。影响Web页面上`服务器流量图`显示的内容                                                                                                                                                                                                                                                                                       |
|            | networkadapter>onlyshow   |                                | 只显示选定的网络适配器。优先级最高默认为空，即不启用                                                                                                                                                                                                                                                                                             |
|            | networkadapter>dontshow   |                                | 不显示的网络适配器。出现对应关键字的网络适配器均不记录数据<br/>配置参考:<br/>需要隐藏的适配器名称`VMware Network Adapter VMnet1`、`VMware Network Adapter VMnet8`<br/>对应配置`["VMware Network Adapter"]`                                                                                                                                                             |
|            | hideexternalip            | false                          | 隐藏自身外网IP，实验性功能。程序不去获取外网Ip，只展示内网IP到网页上。（程序在获取外网ipv4失败后也会尝试获取内网ip）                                                                                                                                                                                                                                                       |

## 目标/计划

- [ ]  HTML显示路由跟踪节点

## 自行编译、运行帮助
### Golang 后端编译
1. 下载项目代码 git clone https://github.com/Jackxwb/alsGo.git
2. 进入 `backend\go` 目录，此目录为Golang项目的主目录，在进行后端开发时，可直接使用IDE打开本目录
3. 打开终端，使用 `go build` 编译项目，编译不同系统请自行添加编译指令。本项目默认使用下面代码来生成来生成Linux-amd64版本可执行程序
```cmd
SET GOOS=linux
SET GOARCH=amd64

go build -o als
```
项目下的 `als` 既是Linux的后端可执行文件
### 前端编译
4. 进入 `ui` 目录，此目录为原项目的VUE前端
5. 打开终端，使用 `npm run install` 和 `npm run build` 来编译VUE项目
6. 打包程序结束后，将dist里面的静态文件复制到 `backend\go\public\static`下面
7. 检查 `backend\go\public\static` 目录下是否存在 `speedtest_worker.js` 文件，该文件来自 [librespeed/speedtest](https://github.com/librespeed/speedtest) 这个项目，若文件不存在可以去对应项目里查找

### 外部环境检查
8. 检查系统环境，ping命令、iperf3命令、speedtest命令能否执行。Window环境建议将iperf3和speedtest程序放入 `backend\go\app` 文件夹里面，Linux用户需要注意speedtest程序的版本，本项目只兼容Ookla®提供的[speedtest-CLI](https://www.speedtest.net/zh-Hans/apps/cli) ，不兼容使用apt安装的speedtest版本
9. 检查系统磁盘可用空间，第一次运行根据配置自动生成静态测速文件(默认生成 100M + 10M +1M 文件在 `backend\go\public\speedtest-static` 目录下)
10. 回到 `backend\go` 目录，运行 `als` 或 `als.exe` 即可启动项目

## 鸣谢

https://github.com/librespeed/speedtest

https://github.com/wikihost-opensource/als

[speedtest-CLI(Ookla®)](https://www.speedtest.net/zh-Hans/apps/cli)

## License

Code is licensed under MIT Public License.

* If you wish to support my efforts, keep the "Powered by LookingGlass" link intact.
