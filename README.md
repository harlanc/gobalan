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

[中文文档](https://github.com/harlanc/gobalan/blob/master/README_CN.md)

## Introduction

Gobalan is a TCP load balancer that supports high network throughput and also support a special load balancing algorithm based on the machine performance.It is a per-connection load balancer.


## Concepts

Next are some concepts for gobalan:

### Worker Service

**Workere Service** is used for representing tcp services need to be load balanced.


### Worker

**Worker** runs on a machine for **Worker Service**'s health check and monitoring the load of the machine(e.g. CPU usage rate,Memory usage rate and network IO usage rate).It reports the above information to **Master** at some specified time interval.

### Master

**Master** is used for collecting **Worker Service**s' health information and machine load information.It also provides service for client to pick a proper worker node using the configured load balance algorithm.

## Load Balance Mechanism

![](http://qiniu.harlanc.vip/2.9.2020_3:20:19.png)

Look at the above picture,it describes how gobalan works:

- Firstly, start **Master** node and all the **Workers** nodes,then **workers** will establish connections with **Master**, then send relevant information including health check information and machine load information to **Master**.

    **Master** will register the load balance algorithms for the next picking.
  
- Secondly,when a client needs to request **Worker Service**,it will firstly request the **Master** to pick a proper worker node using the configured load balance algorithm.

- Thirdly,after the client gets the **Worker Service**'s ip and port,the client will establish a connection with the selected worker node and then proceed with the service request. 

### Load balance algorithms

Now gobalan supports two load balance algorithms:

- RoundRobin It is a tranditional load balance algorithm.
- OptimalPerformance
Now the machine performance scoring formula is:

        score = cpu usage rate + memory usage rate + networkIO read rate + networkIO write rate
The lower score and the better performance.

You can even implement your own load balance algorithm based on the gobalan framework.

## Configurations

We need to do some configurations for Worker and Master,the configuration file config.ini is under config folder.

### Master Configurations

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

The Master section is used for master configuration:

- **IsMaster** If you use the current machine as a gobalan master then set it to **true** or else set it to **false**.
- **Port** The port is the master service port for collecting workers' information and picking worker node for clients. 
- **LBAlgorithm** Use this option to specify load balance algorithm.Here we use the abbreviation of algorithms, that is OP or RR.

### Worker Configurations

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
The worker section is used for worker configuration:

- **IsWorker** If you use the current machine as a gobalan worker then set it to **true** or else set it to **false**.
- **MasterIP** Specify master IP that the current worker will connect to.
- **MasterPort** Specify the master port that the current worker will connect to.
- **ServicePort** The ServicePort is the Port of Worker Service which is run on worker machine.
- **HeartbeatInterval** Specify the heartbeat time interval.

The section Worker.LoadReport is the child section of Worker:

- **LoadReportInterval** Specify the machine load report time interval.
- **MaxNetworkBandwidth** Provide the current maximum network bandwidth.
- **NetworkAdapterName** Specify the network adapter used by current machine.

## How to Use 

The Steps are as follows:

- You should install gobalan on the machine running [Worker Service](https://github.com/harlanc/gobalan#worker-service) as a worker and also install it on another machine as a master.
- Do some configuratinos accrording to [Configurations](https://github.com/harlanc/gobalan#configurations).
- Start both the worker and master service.
- You need to do some extra development work to request master and parse the result for getting the real service IP and port.


## Dependencies

- [grpc-go](https://github.com/grpc/grpc-go) gobalan uses grpc-go as network communication framework.
- [statgo](https://github.com/akhenakh/statgo) gobalan uses statgo to read the machine load in real time,bacause it is a [libstatgrab](http://www.i-scream.org/libstatgrab/) binding for Golang, so you should install the library on your work OS.









