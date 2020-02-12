# gobalan 

[![Build Status][1]][2] [![Go Report Card][3]][4] [![MIT licensed][5]][6] [![code coverage][7]][8] 
 
[1]: https://travis-ci.org/harlanc/gobalan.svg?branch=master
[2]: https://travis-ci.org/harlanc/gobalan
[3]: https://goreportcard.com/badge/github.com/harlanc/gobalan
[4]: https://goreportcard.com/report/github.com/harlanc/gobalan
[5]: https://img.shields.io/badge/license-MIT-blue.svg
[6]: LICENSE
[7]: https://codecov.io/gh/harlanc/gobalan/branch/master/graph/badge.svg
[8]: https://codecov.io/gh/harlanc/gobalan


[EN](https://github.com/harlanc/gobalan)


## 介绍

golalan是一个支持高网络吐吞量的TCP负载均衡器，并且实现了基于机器性能评分的负载均衡算法。gobalan是对TCP连接进行的负载均衡（而非grpc-go的每次调用都会进行负载均衡）。


## 概念

下面是gobalan的一些概念：

### Worker Service

**Worker Service**表示需要被负载均衡的tcp服务。

### Worker

**Worker**用来对**Worker Service**进行健康检查，并且监控**Worker Service**所在机器的负载情况（包括CPU使用率、内存使用率和带宽使用率）。这些信息会在指定的时间间隔内上报给**Master**。

### Master

**Master**用来收集**Worker**节点上**Worker Service**的健康信息和此节点的机器负载信息，维护健康节点列表；实现各种负载均衡算法，并为客户端提供接口，返回由负载均衡算法得出的最合适的**Worker**节点信息。

## 负载均衡机制

![](http://qiniu.harlanc.vip/2.9.2020_3:20:19.png)

上图描述了gobalan是如何工作的：

- 首先，配置好**Master**和**Worker**的配置文件之后，分别启动这些服务。**Worker**会和**Master**建立连接，上报**Worker Service**健康度信息和机器负载信息。
 Master会注册所有负载均衡算法。
 
- 其次，当一个客户端需要请求**Worker Service**的时候，首先请求**Master**,返回由当前负载均衡算法计算出来的最合适的**Worker**节点的信息。
- 最后，客户端解析出**Worker Service**的ip地址和端口号，直接同**Worker service**建立连接进行服务的请求。

### 负载均衡算法

目前gobalan支持两种负载均衡算法：

- 轮询 轮询就是传统的负载均衡算法。
- 最优机器性能 现在的机器性能评分公式为：

        分数 = cpu使用率 + 内存使用率 + 上行带宽使用率+ 下行带宽使用率

分数越低，代表机器性能越好，这个算法还需要进一步改进，需要为各个使用率加上权重。
你也可以在gobalan框架的基础上实现自己的负载均衡算法。

## 配置文件

在启动服务之前需要队**Master** 和**Worker**做一些配置，配置文件config.ini放在项目根目录的config文件夹下面。


### Master 配置

    ; Load Balance Master configurations.
    [Master]
    ; Use the current server as a master or not.
    IsMaster = true
    ; The port that will be connected by workers and clients.
    Port = 5388
    ; Load balancing algorithm
    ; [0] Roundrobin abbr RR
    ; [1] OptimalPerformance abbr OP
    LBAlgorithm = OP

在Master section中配置Master。

- IsMaster 用于表示是否把当前机器当做Master节点。
- Port 为提供Master服务（收集Worker健康信息和机器负载信息）的端口号。
- LBAlgorithm 用于指定负载均衡算法，这里使用简称，RR为轮询，OP为最优机器性能。

### Worker 配置

    ; Load Balance Worker configurations.
    [Worker]
    ; Use the current server as a worker or not.
    IsWorker = false
    ; The Master IP that current worker will connect to.
    MasterIP = 192.0.2.62
    ; The Master port that current worker will connect to.
    MasterPort = 5388 
    ; The service port monitored by worker. 
    ServicePort = 143
    ; The heartbeat sending interval (seconds) used by worker.
    HeartbeatInterval = 5
    ; Load report parameters about worker.
    [Worker.LoadReport]
    ; The worker server load report interval (seconds).
    LoadReportInterval = 60
    ; Max network bandwidth(Mb)
    MaxNetworkBandwidth = 100
    ; The adapter name monitored by worker
    NetworkAdapterName = eth0
    
在 Worker section中配置 Worker：

- IsWorker 用于表示是否把当前机器当做Worker节点。
- MasterIP 表示此Worker需要连接到的MasterIP。
- MasterPort 表示此Woker需要连接到的Master端口。
- ServicePort 此Worker需要监控的Worker Service的端口号。
- HeartbeatInterval Worker节点向Master发送心跳的时间间隔。

Worker.LoadReport是worker section的子section:

- LoadReportInterval 指定Worker向Master发送机器负载信息的时间间隔。
- MaxNetworkBandwidth 当前最大的网络带宽。
- NetworkAdapterName Worker监控的网卡名称。


## 如何使用

步骤如下：

- 在运行Worker service的所有机器上安装gobalan Worker，并且在另外一台机器上安装一个gobalan Master.
- 对Master和Worker做一些配置。
- 启动Worker和Master。
- 使用者需要做一些额外的开发工作：在客户端调用Master提供的接口来请求worker节点的信息。gobalan中提供了golang版本的客户端代码（[rpcpickclient.go](https://github.com/harlanc/gobalan/blob/master/balancer/rpcpickclient.go)）。

## 依赖

- [grpc-go](https://github.com/grpc/grpc-go) gobalan 使用了grpc-go作为网络通信框架。
- [statgo](https://github.com/akhenakh/statgo) gobalan使用statgo读取机器负载信息，因为statgo是对[libstatgrab](http://www.i-scream.org/libstatgrab/)进行的golang封装，因此需要在你的worker服务器上安装这个LIbrary。








