- [第五章 RPC 接口](#第五章-rpc-接口)
	- [5.1 初始化](#51-初始化)
		- [5.1.1 默认路径初始化RPC结构体](#511-默认路径初始化rpc结构体)
		- [5.1.2 带路径初始化RPC结构体](#512-带路径初始化rpc结构体)
		- [5.1.3 生成绑定节点的RPC实例](#513-生成绑定节点的rpc实例)
		- [5.1.4 创建默认RPC结构体](#514-创建默认rpc结构体)
		- [5.1.5 设置RPC实例的namespace](#515-设置rpc实例的namespace)
		- [5.1.6 设置重发次数](#516-设置重发次数)
		- [5.1.7 设置第一次轮询时间间隔](#517-设置第一次轮询时间间隔)
		- [5.1.8 设置第一次轮询次数](#518-设置第一次轮询次数)
		- [5.1.9 设置第二次轮询时间间隔](#519-设置第二次轮询时间间隔)
		- [5.1.10 设置第二次轮询次数](#5110-设置第二次轮询次数)
		- [5.1.11 设置发送重连请求的时间间隔](#5111-设置发送重连请求的时间间隔)
		- [5.1.12 开启https请求](#5112-开启https请求)
		- [5.1.13 增加节点](#5113-增加节点)
		- [5.1.14 开启Tcert认证](#5114-开启tcert认证)
		- [5.1.15 关闭连接](#5115-关闭连接)
	- [5.2 合约服务](#52-合约服务)
		- [5.2.1 编译合约](#521-编译合约)
		- [5.2.2 部署合约](#522-部署合约)
		- [5.2.3 调用合约方法](#523-调用合约方法)
		- [5.2.4 通过投票管理合约](#524-通过投票管理合约)
		- [5.2.5 获取合约字节码](#525-获取合约字节码)
		- [5.2.6 获取指定账户合约数](#526-获取指定账户合约数)
		- [5.2.7 管理合约](#527-管理合约)
		- [5.2.8 获取合约状态](#528-获取合约状态)
		- [5.2.9 通过合约名获取合约状态](#529-通过合约名获取合约状态)
		- [5.2.10 获取合约部署者](#5210-获取合约部署者)
		- [5.2.11 通过合约名获取合约部署者](#5211-通过合约名获取合约部署者)
		- [5.2.12 获取合约部署时间](#5212-获取合约部署时间)
		- [5.2.13 通过合约名获取部署时间](#5213-通过合约名获取部署时间)
		- [5.2.14 查询用户已部署合约列表](#5214-查询用户已部署合约列表)
		- [5.2.15 根据extraId查询交易](#5215-根据extraid查询交易)
		- [5.2.16 调用合约方法(同时返回交易信息)](#5216-调用合约方法同时返回交易信息)
	- [5.3 交易服务](#53-交易服务)
		- [5.3.1 获取交易的TxVersion](#531-获取交易的txversion)
		- [5.3.2 获取区块号区间内的交易信息](#532-获取区块号区间内的交易信息)
		- [5.3.3 通过交易hash获取交易](#533-通过交易hash获取交易)
		- [5.3.4 通过区块hash和交易序号获取交易信息](#534-通过区块hash和交易序号获取交易信息)
		- [5.3.5 通过区块号和交易序号查询交易](#535-通过区块号和交易序号查询交易)
		- [5.3.6 通过区块号区间获取交易平均处理时间](#536-通过区块号区间获取交易平均处理时间)
		- [5.3.7 通过区块hash获取区块上交易数](#537-通过区块hash获取区块上交易数)
		- [5.3.8 通过区块number获取区块上交易数](#538-通过区块number获取区块上交易数)
		- [5.3.9 获取链上所有交易数量](#539-获取链上所有交易数量)
		- [5.3.10 根据时间范围分页查询交易信息](#5310-根据时间范围分页查询交易信息)
		- [5.3.11 根据时间以及合约地址分页查询交易](#5311-根据时间以及合约地址分页查询交易)
		- [5.3.12 根据时间以及合约名分页查询交易](#5312-根据时间以及合约名分页查询交易)
		- [5.3.13 通过交易hash获取交易回执(带轮询)](#5313-通过交易hash获取交易回执带轮询)
		- [5.3.14 通过交易hash获取交易回执(不带轮询)](#5314-通过交易hash获取交易回执不带轮询)
		- [5.3.15 同步发送交易](#5315-同步发送交易)
		- [5.3.16 查询指定区块范围内的非法交易数量](#5316-查询指定区块范围内的非法交易数量)
		- [5.3.17 根据区块号查询区块内的非法交易列表](#5317-根据区块号查询区块内的非法交易列表)
		- [5.3.18 根据区块哈希查询区块内的非法交易列表](#5318-根据区块哈希查询区块内的非法交易列表)
		- [5.3.19 获取链上的非法交易数](#5319-获取链上的非法交易数)
		- [5.3.20 通过交易hash获取产生了checkpoint之后的交易回执](#5320-通过交易hash获取产生了checkpoint之后的交易回执)
		- [5.3.21 获取节点设置的交易Gas价格](#5321-获取节点设置的交易gas价格)
		- [5.3.22 设置当前rpc的gasPrice为节点的gasPrice](#5322-设置当前rpc的gasprice为节点的gasprice)
	- [5.4 节点服务](#54-节点服务)
		- [5.4.1 新建节点](#541-新建节点)
		- [](#)
		- [5.4.2 获取区块链节点信息](#542-获取区块链节点信息)
		- [5.4.3 获取随机节点hash](#543-获取随机节点hash)
		- [5.4.4 从指定节点获取hash](#544-从指定节点获取hash)
		- [5.4.5 删除NVP节点](#545-删除nvp节点)
		- [5.4.6 NVP断开与VP节点的链接](#546-nvp断开与vp节点的链接)
		- [5.4.7 获取节点状态信息](#547-获取节点状态信息)
		- [5.4.8 替换节点证书](#548-替换节点证书)
	- [5.5 区块服务](#55-区块服务)
		- [5.5.1 获取最后一个区块的信息](#551-获取最后一个区块的信息)
		- [5.5.2 分页获取区块信息](#552-分页获取区块信息)
		- [5.5.3 通过区块hash获取区块信息](#553-通过区块hash获取区块信息)
		- [5.5.4 通过区块号获取区块信息](#554-通过区块号获取区块信息)
		- [5.5.5 计算区间内区块平均生成时间](#555-计算区间内区块平均生成时间)
		- [5.5.6 查询区间内区块生成的速度以及TPS](#556-查询区间内区块生成的速度以及tps)
		- [5.5.7 查询创世区块号](#557-查询创世区块号)
		- [5.5.8 查询区块高度](#558-查询区块高度)
	- [5.6 账户服务](#56-账户服务)
		- [5.6.1 获取账户角色](#561-获取账户角色)
		- [5.6.2 根据角色查询账户](#562-根据角色查询账户)
		- [5.6.3 获取账户状态](#563-获取账户状态)
		- [5.6.4 获取指定账户的证明路径](#564-获取指定账户的证明路径)
	- [5.7 MQ服务](#57-mq服务)
		- [5.7.1 获取mq客户端](#571-获取mq客户端)
		- [5.7.2 注册MQ channel](#572-注册mq-channel)
		- [5.7.3 注销MQ channel](#573-注销mq-channel)
		- [5.7.4 获取指定节点所有的queue名](#574-获取指定节点所有的queue名)
		- [5.7.5 与broker建立连接](#575-与broker建立连接)
		- [5.7.6 获取节点当前exchanger名](#576-获取节点当前exchanger名)
		- [5.7.7 删除exchanger](#577-删除exchanger)
		- [5.7.8 添加监听](#578-添加监听)
	- [5.8 Kvsql服务](#58-kvsql服务)
	- [5.9 webSocket服务](#59-websocket服务)
		- [5.9.1 获取websocket客户端](#591-获取websocket客户端)
		- [5.9.2 订阅提案事件](#592-订阅提案事件)
		- [5.9.3 订阅事件](#593-订阅事件)
		- [5.9.4 取消订阅](#594-取消订阅)
		- [5.9.5 获取节点所有的订阅信息](#595-获取节点所有的订阅信息)
		- [5.9.6 关闭节点的websocket连接](#596-关闭节点的websocket连接)
	- [5.10 归档服务](#510-归档服务)
		- [5.10.1 列出所有快照信息](#5101-列出所有快照信息)
		- [5.10.2 数据直接归档](#5102-数据直接归档)
		- [5.10.3 检查数据归档是否完成](#5103-检查数据归档是否完成)
		- [5.10.4 查询数据归档状态](#5104-查询数据归档状态)
		- [5.10.5 查询最近一次归档的状态](#5105-查询最近一次归档的状态)
		- [5.10.6 在某个已经存在的区块高度归档](#5106-在某个已经存在的区块高度归档)
	- [5.11 配置服务](#511-配置服务)
		- [5.11.1 查询提案](#5111-查询提案)
		- [5.11.2 查询配置信息](#5112-查询配置信息)
		- [5.11.3 查询Hosts](#5113-查询hosts)
		- [5.11.4 查询VSet](#5114-查询vset)
		- [5.11.5 查询链上角色权重信息](#5115-查询链上角色权重信息)
		- [5.11.6 检查链上角色是否存在](#5116-检查链上角色是否存在)
		- [5.11.7 根据合约命名查询合约地址](#5117-根据合约命名查询合约地址)
		- [5.11.8 根据合约地址获取合约名](#5118-根据合约地址获取合约名)
		- [5.11.9 获取所有<合约地址, 合约名>的映射](#5119-获取所有合约地址-合约名的映射)
		- [5.11.10 获取创世信息](#51110-获取创世信息)
	- [5.12 权限服务](#512-权限服务)
		- [5.12.1 为账户增加节点级角色](#5121-为账户增加节点级角色)
		- [5.12.2 为账户删除节点级角色](#5122-为账户删除节点级角色)
		- [5.12.3 查询账户的节点级角色](#5123-查询账户的节点级角色)
		- [5.12.4 查询节点级某角色的账户列表](#5124-查询节点级某角色的账户列表)
		- [5.12.5 查询所有的节点级角色](#5125-查询所有的节点级角色)
		- [5.12.6 设置节点级接口权限管理规则列表](#5126-设置节点级接口权限管理规则列表)
		- [5.12.7 查询节点级接口权限管理规则列表](#5127-查询节点级接口权限管理规则列表)
	- [5.13 DID服务](#513-did服务)
		- [5.13.1 发送DID交易](#5131-发送did交易)
		- [5.13.2 查询chainID](#5132-查询chainid)
		- [5.13.3 查询DID文档](#5133-查询did文档)
		- [5.13.4 查询凭证基础信息](#5134-查询凭证基础信息)
		- [5.13.5 检查凭证是否有效](#5135-检查凭证是否有效)
		- [5.13.6 检查凭证是否吊销](#5136-检查凭证是否吊销)
		- [5.13.7 获取节点chainId](#5137-获取节点chainid)
		- [5.13.8 设置本地chainId为节点chainId](#5138-设置本地chainid为节点chainid)
		- [5.13.9 获取DID账户的公钥](#5139-获取did账户的公钥)
		- [5.13.10 新建DIDDocument](#51310-新建diddocument)
		- [5.13.11 新建DID凭证](#51311-新建did凭证)
	- [5.14 文件服务](#514-文件服务)
		- [5.14.1 文件上传](#5141-文件上传)
		- [5.14.2 文件下载](#5142-文件下载)
		- [5.14.3 文件信息更新](#5143-文件信息更新)
		- [5.14.4 推送文件](#5144-推送文件)
		- [5.14.5 通过extraId获取文件信息FileExtra](#5145-通过extraid获取文件信息fileextra)
		- [5.14.6  通过filter获取文件信息FileExtra](#5146--通过filter获取文件信息fileextra)
		- [5.14.7 通过交易哈希获取文件信息FileExtra](#5147-通过交易哈希获取文件信息fileextra)
	- [5.15 跨分区服务](#515-跨分区服务)
		- [5.15.1 部署跨分区合约](#5151-部署跨分区合约)
		- [5.15.2 调用跨分区方法](#5152-调用跨分区方法)
	- [5.16 交易证明](#516-交易证明)
        - [5.16.1 通过交易hash获取证明路径](#5161-通过交易hash获取证明路径)
        - [5.16.2 获取状态数据的证明路径](#5162-获取状态数据的证明路径)
        - [5.16.3 验证状态数据的证明路径是否正确](#5163-验证状态数据的证明路径是否正确)
	- [5.17 版本管理服务](#517-版本管理服务)
		- [5.17.1 二进制版本上链](#5171-二进制版本上链)
		- [5.17.2 查询链运行版本信息](#5172-查询链运行版本信息)
		- [5.17.3 查询节点支持的版本信息](#5173-查询节点支持的版本信息)
		- [5.17.4 查询Hyperchain版本对应的链级细分版本](#5174-查询Hyperchain版本对应的链级细分版本)

# 第五章 RPC 接口



## 5.1 初始化

### 5.1.1 默认路径初始化RPC结构体

```go
func NewRPC() *RPC
```

| 说明   | 默认路径为初始化RPC处的上层文件夹下的conf文件，即相对路径../conf/，若配置文件夹在别处，则需使用传路径的方式 |
| ---- | ------------------------------------------------------------- |
| 返回值  |                                                               |
| 返回值1 | RPC结构实例                                                       |

应用实例

```go
rp := rpc.NewRPC()
```



### 5.1.2 带路径初始化RPC结构体

```go
func NewRPCWithPath(confRootPath string) *RPC
```

| 说明   | 传入conf配置文件夹路径来获取RPC结构体 |
| ---- | ---------------------- |
| 参数   |                        |
| 参数1  | 配置文件目录的路径, 详情参见1.1部分   |
| 返回值  |                        |
| 返回值1 | RPC结构体实例               |

应用实例

```go
rpc.NewRPCWithPath("../conf")
```



### 5.1.3 生成绑定节点的RPC实例

```go
func (r *RPC) BindNodes(nodeIndexes ...int) (*RPC, error)
```

| 说明   | 指定RPC绑定的节点, 默认RPC节点是配置文件中指定的全部节点                                                   |
| ---- | ---------------------------------------------------------------------------------- |
| 参数   |                                                                                    |
| 参数1  | 节点编号。如nodeIndex为1,2，即对应配置文件中jsonRPC.nodes中的1号和2号节点，返回的RPC对象将只会在1号节点和2号节点间负载均衡发送请求。 |
| 返回值  |                                                                                    |
| 返回值1 | 如果未指定节点号，则为原来的RPC对象；否则为一个新的RPC对象，与原有对象互不影响                                         |

应用实例

```go
h1 := rpc.NewRPC()
fmt.Printf("%p \n", h1)
h2, _ := h1.BindNodes(1,2) // h2 与 h1 不同
h3, _ := h1.BindNodes() // h3与h1为同一个对象
h4, _ := h1.BindNodes(1,2,3,4) // // h3 与 h1 不同
```



### 5.1.4 创建默认RPC结构体

```go
func DefaultRPC(nodes ...*Node) *RPC
```

| 说明   | 生成一个带有默认参数的RPC结构体         |
| ---- | ------------------------- |
| 参数   |                           |
| 参数1  | 节点列表。可以通过NewNode接口创建，见5.4 |
| 返回值  |                           |
| 返回值1 | RPC实例                     |

应用实例

```go
rp := rpc.DefaultRPC(NewNode("localhost", "8081", "11001"))
```



### 5.1.5 设置RPC实例的namespace

```go
func (rpc *RPC) Namespace(ns string) *RPC
```

应用实例

```go
rp := rpc.NewRPC()
rp.Namespace("global")
```



### 5.1.6 设置重发次数

重发次数主要应用在轮询的场景

```go
func (rpc *RPC) ResendTimes(resTime int64) *RPC
```

应用实例

```go
rp := rpc.NewRPC()
rp.ResendTimes(3)
```

### 5.1.7 设置第一次轮询时间间隔

单位ms

```go
func (rpc *RPC) FirstPollInterval(fpi int64) *RPC
```

```go
rp := rpc.NewRPC()
rp.FirstPollInterval(100)
```

### 5.1.8 设置第一次轮询次数

```go
func (rpc *RPC) FirstPollTime(fpt int64) *RPC 
```

```go
rp := rpc.NewRPC()
rp.FirstPollTime(10)
```



### 5.1.9 设置第二次轮询时间间隔

单位ms

```go
func (rpc *RPC) SecondPollInterval(fpi int64) *RPC
```

```go
rp := rpc.NewRPC()
rp.SecondPollInterval(1000)
```

### 5.1.10 设置第二次轮询次数

```go
func (rpc *RPC) SecondPollTime(fpt int64) *RPC 
```

```go
rp := rpc.NewRPC()
rp.SecondPollTime(10)
```

### 5.1.11 设置发送重连请求的时间间隔

单位ms

```go
func (rpc *RPC) ReConnTime(rct int64) *RPC
```

```go
rp := rpc.NewRPC()
rp.ReConnTime(10000)
```



### 5.1.12 开启https请求

```go
func (rpc *RPC) Https(tlscaPath, tlspeerCertPath, tlspeerPrivPath string) *RPC
```

| 说明   | 使用该方法后将开启tls认证，并使用https发送请求 |
| ---- | --------------------------- |
| 参数   |                             |
| 参数1  | tls ca证书的路径                 |
| 参数2  | tls cert证书的路径               |
| 参数3  | tls 私钥文件的路径                 |
| 返回值  |                             |
| 返回值1 | RPC实例                       |

```go
rp := rpc.NewRPC().Https("../conf/certs/tls/tlsca.ca", "../conf/certs/tls/tls_peer.cert", "../conf/certs/tls/tls_peer.priv")
```



### 5.1.13 增加节点

```go
func (rpc *RPC) AddNode(url, rpcPort, wsPort string) *RPC
```

| 参数  |                 |
| --- | --------------- |
| 参数1 | 节点url           |
| 参数2 | 节点rpc服务端口       |
| 参数3 | 节点websocket服务端口 |

```go
rp := rpc.NewRPC()
rp.AddNode("127.0.0.1", "8081", "10001")
```

### 5.1.14 开启Tcert认证

```go
func (rpc *RPC) Tcert(cfca bool, sdkcertPath, sdkcertPrivPath, uniquePubPath, uniquePrivPath string) *RPC
```

| 说明   | 使用该方法后将开启tcert认证 |
| ---- | ---------------- |
| 参数   |                  |
| 参数1  | 是否使用ca           |
| 参数2  | sdk cert证书的路径    |
| 参数3  | sdk 私钥文件的路径      |
| 参数4  | unique.pub证书的路径  |
| 参数5  | unique.priv证书的路径 |
| 返回值  |                  |
| 返回值1 | RPC实例            |

应用实例

```go
rp := rpc.NewRPC().Tcert(true, "../conf/certs/sdkcert.cert", "../conf/certs/sdkcert.priv", "../conf/certs/unique.pub", "../conf/certs/unique.priv")
```



### 5.1.15 关闭连接

```go
func (rpc *RPC) Close() 
```

```go
rp := rpc.NewRPC()
// ...
rp.Close()
```

## 5.2 合约服务

### 5.2.1 编译合约

```go
func (rpc *RPC) CompileContract(code string) (*CompileResult, StdError)
```

| 说明   | 编译合约，Solidity合约支持在线编译，目前只支持solidity 0.5.0版本以下的合约,不推荐使用远程编译，因为远程编译支持的合约版本较低，不支持最新版本solidity合约和多合约编译。 |
| ---- | --------------------------------------------------------------------------------------------------- |
| 参数   |                                                                                                     |
| 参数1  | 合约源码                                                                                                |
| 返回值  |                                                                                                     |
| 返回值1 | 编译结果                                                                                                |
| 返回值2 | error                                                                                               |

应用实例

```go
rpcAPI := rpc.NewRPCWithPath("../../conf")
path := "../../conf/contract/Accumulator.sol"
contract, _ := common.ReadFileAsString(path)
cr, err := rpcAPI.CompileContract(contract)
```

### 5.2.2 部署合约

```go
func (rpc *RPC) SignAndDeployContract(transaction *Transaction, key interface{}) (*TxReceipt, StdError)
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 交易体      |
| 参数2  | 交易发起方key |
| 返回值  |          |
| 返回值1 | 交易回执     |
| 返回值2 | error    |

应用实例

```go
hrpc := rpc.NewRPCWithPath("../conf")
js, err := account.NewAccountSm2("12345678")
gmAcc, err := account.GenKeyFromAccountJson(js, "12345678")
newAddress := gmAcc.(*account.SM2Key).GetAddress()
tx := rpc.NewTransaction(newAddress.Hex()).Deploy(hexutil.Encode([]byte("KVSQL")))
hrpc.SignAndDeployContract(tx, gmAcc)
```

### 5.2.3 调用合约方法

> （若需要同时返回交易信息，可查看5.2.16）

```go
func (rpc *RPC) SignAndInvokeContract(transaction *Transaction, key interface{}) (*TxReceipt, StdError)
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 交易体      |
| 参数2  | 交易发起方key |
| 返回值  |          |
| 返回值1 | 交易回执     |
| 返回值2 | error    |

应用实例

```go
rp := rpc.NewRPC()
accountJson := `{"address":"0xfbca6a7e9e29728773b270d3f00153c75d04e1ad","version":"4.0","algo":"0x13","publicKey":"049c330d0aea3d9c73063db339b4a1a84d1c3197980d1fb9585347ceeb40a5d262166ee1e1cb0c29fd9b2ef0e4f7a7dfb1be6c5e759bf411c520a616863ee046a4","privateKey":"5f0a3ea6c1d3eb7733c3170f2271c10c1206bc49b6b2c7e550c9947cb8f098e3"}`
key, _ := account.GenKeyFromAccountJson(accountJson, "")
opt := bvm.NewDIDSetChainIDOperation("chainID_01")
payload := bvm.EncodeOperation(opt)
// 发送交易
tx := rpc.NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
rp.SignAndInvokeContract(tx, key)
```

### 5.2.4 通过投票管理合约

```go
func (rpc *RPC) SignAndManageContractByVote(transaction *Transaction, key interface{}) (*TxReceipt, StdError) 
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 交易体      |
| 参数2  | 交易发起方key |
| 返回值  |          |
| 返回值1 | 交易回执     |
| 返回值2 | error    |

应用实例

```go
rp := rpc.NewRPC()
source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
bin := `6060604052341561000f57600080fd5b5b6104c78061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635b6beeb914610049578063e15fe02314610120575b600080fd5b341561005457600080fd5b6100a4600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061023a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e55780820151818401525b6020810190506100c9565b50505050905090810190601f1680156101125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561012b57600080fd5b6101be600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061034f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ff5780820151818401525b6020810190506101e3565b50505050905090810190601f16801561022c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102426103e2565b6000826040518082805190602001908083835b60208310151561027b57805182525b602082019150602081019050602083039250610255565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103425780601f1061031757610100808354040283529160200191610342565b820191906000526020600020905b81548152906001019060200180831161032557829003601f168201915b505050505090505b919050565b6103576103e2565b816000846040518082805190602001908083835b60208310151561039157805182525b60208201915060208101905060208303925061036b565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906103d79291906103f6565b508190505b92915050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061043757805160ff1916838001178555610465565b82800160010185558215610465579182015b82811115610464578251825591602001919060010190610449565b5b5090506104729190610476565b5090565b61049891905b8082111561049457600081600090555060010161047c565b5090565b905600a165627a7a723058208ac1d22e128cf8381d7ac66b4c438a6a906ccf5ee583c3a9e46d4cdf7b3f94580029`
ope := bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(bin), "evm", nil)
contractOpt := bvm.NewProposalCreateOperationForContract(ope)
payload := bvm.EncodeOperation(contractOpt)
tx := rpc.NewTransaction(privateKey.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
re, err := rp.SignAndManageContractByVote(tx)
```



### 5.2.5 获取合约字节码

```go
func (rpc *RPC) GetCode(contractAddress string) (string, StdError) 
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 合约地址  |
| 返回值  |       |
| 返回值1 | 合约字节码 |
| 返回值2 | error |

```go
rp := rpc.NewRPC()
// contractAddr 合约地址
rp.GetCode(contractAddr)
```

### 5.2.6 获取指定账户合约数

```go
func (rpc *RPC) GetContractCountByAddr(accountAddress string) (uint64, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 账户地址  |
| 返回值  |       |
| 返回值1 | 合约数量  |
| 返回值2 | error |

应用示例

```go
rp := rpc.NewRPC()
pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
newAddress := h[12:]
count, err := rp.GetContractCountByAddr(common.BytesToAddress(newAddress).Hex())
```

### 5.2.7 管理合约

```go
func (rpc *RPC) SignAndMaintainContract(transaction *Transaction, key interface{}) (*TxReceipt, StdError)
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 交易体      |
| 参数2  | 交易发起方key |
| 返回值  |          |
| 返回值1 | 交易回执     |
| 返回值2 | error    |

应用实例

```go
rp := rpc.NewRPC()
transactionUpdate := rp.NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(1, contractAddress, compileUpdate.Bin[0])
receiptUpdate, err := rp.SignAndMaintainContract(transactionUpdate, key)
```

### 5.2.8 获取合约状态

```go
func (rpc *RPC) GetContractStatus(contractAddress string) (string, StdError)
```

| 参数   |                                                         |
| ---- | ------------------------------------------------------- |
| 参数1  | 合约地址                                                    |
| 返回值  |                                                         |
| 返回值1 | 合约状态，返回值为normal表示正常，frozen表示冻结，non-contract表示非合约即普通转账交易 |
| 返回值2 | error                                                   |

应用实例

```go
rp := rpc.NewRPC()
// contractAddr 合约地址
rp.GetContractStatus(contractAddr)
```



### 5.2.9 通过合约名获取合约状态

```go
func (rpc *RPC) GetContractStatusByName(contractName string) (string, StdError) 
```

| 参数   |                                                         |
| ---- | ------------------------------------------------------- |
| 参数1  | 合约名                                                     |
| 返回值  |                                                         |
| 返回值1 | 合约状态，返回值为normal表示正常，frozen表示冻结，non-contract表示非合约即普通转账交易 |
| 返回值2 | error                                                   |



应用实例

```go
rp := rpc.NewRPC()
// contractName 合约名
rp.GetContractStatusByName(contractName)
```



### 5.2.10 获取合约部署者

```go
func (rpc *RPC) GetCreator(contractAddress string) (string, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 合约地址  |
| 返回值  |       |
| 返回值1 | 合约部署者 |
| 返回值2 | error |

应用实例

```go
rp := rpc.NewRPC()
// contractAddr 合约地址
rp.GetCreator(contractAddr)
```



### 5.2.11 通过合约名获取合约部署者

```go
func (rpc *RPC) GetCreatorByName(contractName string) (string, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 合约名   |
| 返回值  |       |
| 返回值1 | 合约部署者 |
| 返回值2 | error |

应用实例

```go
rp := rpc.NewRPC()
// contractName 合约名
rp.GetCreatorByName(contractName)
```



### 5.2.12 获取合约部署时间

```go
func (rpc *RPC) GetCreateTime(contractAddress string) (string, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 合约地址   |
| 返回值  |        |
| 返回值1 | 合约部署时间 |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
// contractAddr 合约地址
rp.GetCreateTime(contractAddr)
```

### 5.2.13 通过合约名获取部署时间

```text
func (rpc *RPC) GetCreateTimeByName(contractName string) (string, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 合约名    |
| 返回值  |        |
| 返回值1 | 合约部署时间 |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
// contractName 合约名
rp.GetCreateTimeByName(contractName)
```

### 5.2.14 查询用户已部署合约列表

```go
func (rpc *RPC) GetDeployedList(address string) ([]string, StdError)
```

| 参数   |           |
| ---- | --------- |
| 参数1  | 用户地址      |
| 返回值  |           |
| 返回值1 | 已部署合约地址列表 |
| 返回值2 | error     |

应用实例

```go
rp := rpc.NewRPC()
rp.GetDeployedList(gmKey.GetAddress())
```

### 5.2.15 根据extraId查询交易

```go
func (rpc *RPC) GetTransactionsByExtraID(extraId []interface{}, txTo string, detail bool, mode int, metadata *Metadata) (*PageResult, StdError)
```

| 参数  |                                                                                                                                                                                               |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 参数1 | extraId的列表                                                                                                                                                                                    |
| 参数2 | 交易接收方地址                                                                                                                                                                                       |
| 参数3 | 是否返回交易详情，true则返回TransactionResult， false返回交易摘要TransactionSummary                                                                                                                              |
| 参数4 | 表示本次查询请求的查询模式，目前有0、1、2三个值可选，默认为0。0 表示按序精确查询模式，即筛选出的的交易 extraId 数组的数值和顺序都与查询条件完全一致。1 表示非按序精确查询模式，即筛选出的交易 extraId 数组包含查询条件里指定的全部数值，顺序无要求。2 表示非按序匹配查询模式，即筛选出的交易 extraId 数组包含部分或全部查询条件指定的值，且顺序无要求 |
| 参数5 | 查询参数                                                                                                                                                                                          |



```go
rp := rpc.NewRPC()
rp.GetTransactionsByExtraID([]interface{}{extraId}, "", true, 0, nil)
```



### 5.2.16 调用合约方法(同时返回交易信息)

> 在调用合约方法的同时返回交易信息，也可通过写__轮询__的方式去查交易信息

```go
func (rpc *RPC) SignAndInvokeContractCombineReturns(transaction *Transaction, key interface{}) (*TxReceipt, *TransactionInfo, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 交易体   |
| 参数2  | 用户key |
| 返回值  |       |
| 返回值1 | 交易回执  |
| 返回值2 | 交易信息  |
| 返回值3 | err   |

```go
rp := rpc.NewRPC()
accountJSON, _ := account.NewAccountED25519("12345678")
ekey, err := account.GenKeyFromAccountJson(accountJSON, "12345678")
assert.Nil(t, err)
newAddress := ekey.(*account.ED25519Key).GetAddress()

transaction := rpc.NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
transaction.Sign(ekey)
txreceipt, info, err := rp.SignAndInvokeContractCombineReturns(transaction, ekey)
```



## 5.3 交易服务

### 5.3.1 获取交易的TxVersion

```go
func (rpc *RPC) GetTxVersion() (string, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
rp.GetTxVersion()
```

### 5.3.2 获取区块号区间内的交易信息

```go
func (rpc *RPC) GetTransactionsByBlkNumWithLimit(start, end uint64, metadata *Metadata) (*PageResult, StdError)
```

| 参数  |       |
| --- | ----- |
| 参数1 | 起始区块号 |
| 参数2 | 末尾区块号 |
| 参数3 | 查询参数  |

应用实例

```go
rp := rpc.NewRPC()
block, err := rp.GetLatestBlock()
metadata := &Metadata{
		PageSize: 1,
		Bookmark: &Bookmark{
			BlockNumber: 1,
			TxIndex:     0,
		},
		Backward: false,
	}
pageResult, err := rp.GetTransactionsByBlkNumWithLimit(block.Number-1, block.Number, metadata)
```

### 5.3.3 通过交易hash获取交易

```go
func (rpc *RPC) GetTransactionByHash(txHash string) (*TransactionInfo, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
rp.GetTransactionByHash("0x2a6776578")
```

### 5.3.4 通过区块hash和交易序号获取交易信息

```go
func (rpc *RPC) GetTxByBlkHashAndIdx(blkHash string, index uint64) (*TransactionInfo, StdError)
```

| 参数  |             |
| --- | ----------- |
| 参数1 | 区块hash      |
| 参数2 | 交易在当前区块内的序号 |

应用实例

```go
rp := rpc.NewRPC()
tx, err := rp.GetTransactionByHash(toHash)
rp.GetTxByBlkHashAndIdx(tx.BlockHash, tx.TxIndex)
```



### 5.3.5 通过区块号和交易序号查询交易

```go
func (rpc *RPC) GetTxByBlkNumAndIdx(blkNum, index uint64) (*TransactionInfo, StdError)
```

| 参数  |             |
| --- | ----------- |
| 参数1 | 区块号         |
| 参数2 | 交易在当前区块内的序号 |

应用实例

```go
rp := rpc.NewRPC()
tx, err := rp.GetTransactionByHash(toHash)
rp.GetTxByBlkNumAndIdx(tx.BlockNumber, tx.TxIndex)
```

### 5.3.6 通过区块号区间获取交易平均处理时间

```go
func (rpc *RPC) GetTxAvgTimeByBlockNumber(from, to uint64) (uint64, StdError) 
```

| 参数  |       |
| --- | ----- |
| 参数1 | 起始区块号 |
| 参数2 | 末尾区块号 |

应用实例

```go
rp := rpc.NewRPC()
block, _ := rp.GetLatestBlock()
ts, err := rp.GetTxAvgTimeByBlockNumber(block.Number-2, block.Number)
```



### 5.3.7 通过区块hash获取区块上交易数

```go
func (rpc *RPC) GetBlkTxCountByHash(blkHash string) (uint64, StdError)
```

| 参数  |        |
| --- | ------ |
| 参数1 | 区块hash |

应用实例

```go
rp := rpc.NewRPC()
block, err := rp.GetLatestBlock()
count, err := rp.GetBlkTxCountByHash(block.Hash)
```

### 5.3.8 通过区块number获取区块上交易数

```go
func (rpc *RPC) GetBlkTxCountByNumber(blkNum string) (uint64, StdError) 
```

| 参数  |          |
| --- | -------- |
| 参数1 | 区块number |

应用实例

```go
rp := rpc.NewRPC()
block, err := rp.GetLatestBlock()
hex := "0x" + strconv.FormatUint(block.Number, 16)
count, err := rp.GetBlkTxCountByNumber(hex)
```

### 5.3.9 获取链上所有交易数量

```go
func (rpc *RPC) GetTxCount() (*TransactionsCount, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
rp.GetTxCount()
```



### 5.3.10 根据时间范围分页查询交易信息

```go
func (rpc *RPC) GetTxByTimeWithLimit(start, end uint64, metadata *Metadata) (*PageTxs, StdError)
```

| 参数  |      |
| --- | ---- |
| 参数1 | 开始时间 |
| 参数2 | 截止时间 |
| 参数3 | 查询参数 |

应用实例

```go
rp := rpc.NewRPC()
rp.GetTxByTimeWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		})

```

### 5.3.11 根据时间以及合约地址分页查询交易

```go
func (rpc *RPC) GetTxByTimeAndContractAddrWithLimit(start, end uint64, metadata *Metadata, contractAddr string) (*PageTxs, StdError)
```

| 参数  |      |
| --- | ---- |
| 参数1 | 开始时间 |
| 参数2 | 截止时间 |
| 参数3 | 查询参数 |
| 参数4 | 合约地址 |

应用实例

```go
rp := rpc.NewRPC()
// contractAddr 合约地址
rp.GetTxByTimeAndContractAddrWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		}, contractAddr)

```





### 5.3.12 根据时间以及合约名分页查询交易



```go
func (rpc *RPC) GetTxByTimeAndContractNameWithLimit(start, end uint64, metadata *Metadata, contractName string) (*PageTxs, StdError)
```

| 参数  |      |
| --- | ---- |
| 参数1 | 开始时间 |
| 参数2 | 截止时间 |
| 参数3 | 查询参数 |
| 参数4 | 合约名  |

应用实例

```go
rp := rpc.NewRPC()
// contractName 合约名
rp.GetTxByTimeAndContractNameWithLimit(uint64(start), uint64(end), &Metadata{
			PageSize: 10,
			Backward: false,
		}, contractName)

```

### 5.3.13 通过交易hash获取交易回执(带轮询)

```go
func (rpc *RPC) GetTxReceiptByPolling(txHash string, isPrivateTx bool) (*TxReceipt, StdError, bool)
```

| 参数  |        |
| --- | ------ |
| 参数1 | 交易hash |
| 参数2 | 是否隐私交易 |

应用实例

```go
rp := rpc.NewRPC()
// txhash 交易hash
rp.GetTxReceiptByPolling(txhash, false)
```



### 5.3.14 通过交易hash获取交易回执(不带轮询)

```go
func (rpc *RPC) GetTxReceipt(txHash string, isPrivateTx bool) (*TxReceipt, StdError)
```

| 参数  |        |
| --- | ------ |
| 参数1 | 交易hash |
| 参数2 | 是否隐私交易 |

应用实例

```go
rp := rpc.NewRPC()
// txhash 交易hash
rp.GetTxReceiptByPolling(txhash, false)
```

### 5.3.15 同步发送交易

```text
func (rpc *RPC) SignAndSendTx(transaction *Transaction, key interface{}) (*TxReceipt, StdError)
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 交易体      |
| 参数2  | 交易发送方key |
| 返回值  |          |
| 返回值1 | 交易回执     |
| 返回值2 | error    |

应用实例

```go
rp := rpc.NewRPC()
accountJSON, _ := account.NewAccountED25519("12345678")
ekey, err := account.GenKeyFromAccountJson(accountJSON, "12345678")
newAddress := ekey.(*account.ED25519Key).GetAddress()
//address为交易接受方地址
transaction := rpc.NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
rp.SignAndSendTx(transaction, ekey)
```



### 5.3.16 查询指定区块范围内的非法交易数量

```go
func (rpc *RPC) GetInvalidTransactionsByBlkNumWithLimit(start, end uint64, metadata *Metadata) (*PageResult, StdError)
```

| 参数   |              |
| ---- | ------------ |
| 参数1  | 起始区块号        |
| 参数2  | 结束区块号        |
| 参数3  | 查询参数         |
| 返回值  |              |
| 返回值1 | 该区块号区间内的交易列表 |
| 返回值2 | error        |

应用实例

```go
rp := rpc.NewRPC()
latestBlock, err := rp.GetLatestBlock()
metadata := &rpc.Metadata{
	PageSize: 5,
}
pageResult, err := rp.GetInvalidTransactionsByBlkNumWithLimit(latestBlock-1, latestBlock, metadata)
```



### 5.3.17 根据区块号查询区块内的非法交易列表

```go
func (rpc *RPC) GetInvalidTransactionsByBlkNum(blkNum uint64) ([]TransactionInfo, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 区块号    |
| 返回值  |        |
| 返回值1 | 非法交易列表 |
| 返回值2 | error  |

应用示例

```go
rp := rpc.NewRPC()
txInfos, err := rp.GetInvalidTransactionsByBlkNum(5)
```

### 5.3.18 根据区块哈希查询区块内的非法交易列表

```go
func (rpc *RPC) GetInvalidTransactionsByBlkHash(hash string) ([]TransactionInfo, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 区块hash |
| 返回值  |        |
| 返回值1 | 非法交易列表 |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
latestBlock, err := rp.GetLatestBlock()
rp.GetInvalidTransactionsByBlkHash(latestBlock.Hash)
```

### 5.3.19 获取链上的非法交易数

```go
func (r *RPC) GetInvalidTxCount() (uint64, StdError)
```

| 返回值  |       |
| ---- | ----- |
| 返回值1 | 非法交易数 |
| 返回值2 | error |

应用示例

```go
rp := rpc.NewRPC()
count, err := rp.GetInvalidTxCount()
```

### 5.3.20 通过交易hash获取产生了checkpoint之后的交易回执

```go
func (rpc *RPC) GetTxConfirmedReceipt(txHash string) (*TxReceipt, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 交易hash |
| 返回值  |        |
| 返回值1 | 交易回执   |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
count, err := rp.GetTxConfirmedReceipt(txhash)
```



### 5.3.21 获取节点设置的交易Gas价格

```go
func (rpc *RPC) GetGasPrice() (int64, error)
```

| 返回值  |              |
| ---- | ------------ |
| 返回值1 | 节点设置的交易gas价格 |
| 返回值2 | error        |

应用实例

```go
rp := rpc.NewRPC()
price, err := rp.GetGasPrice()
```



### 5.3.22 设置当前rpc的gasPrice为节点的gasPrice

```go
func (rpc *RPC) SetGasPrice() error
```

应用实例

```go
rp := rpc.NewRPC()
err := rp.SetGasPrice()
```





## 5.4 节点服务

### 5.4.1 新建节点

```go
func NewNode(url string, rpcPort string, wsPort string) (node *Node)
```

| 参数   |                 |
| ---- | --------------- |
| 参数1  | 节点url           |
| 参数2  | 节点rpc服务端口       |
| 参数3  | 节点websocket服务端口 |
| 返回值  |                 |
| 返回值1 | 节点Node实例        |

应用实例

```go
rpc := rpc.DefaultRPC(NewNode("localhost", "8081", "11001"))
```

### 

### 5.4.2 获取区块链节点信息

```go
func (rpc *RPC) GetNodes() ([]NodeInfo, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
info, err := rp.GetNodes()
```



### 5.4.3 获取随机节点hash

```go
func (rpc *RPC) GetNodeHash() (string, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
info, err := rp.GetNodeHash()
```



### 5.4.4 从指定节点获取hash

```go
func (rpc *RPC) GetNodeHashByID(id int) (string, StdError)
```

| 参数  |                              |
| --- | ---------------------------- |
| 参数1 | 节点编号，从1开始，不要超过RPC结构体中设置的节点个数 |

应用实例

```go
rp := rpc.NewRPC()
rp.GetNodeHashByID(1)
```



### 5.4.5 删除NVP节点

```go
func (rpc *RPC) DeleteNodeNVP(hash string) (bool, StdError)
```

| 参数  |        |
| --- | ------ |
| 参数1 | 节点hash |

应用实例

```go
rp := rpc.NewRPC()
rp.DeleteNodeNVP("0x")
```



### 5.4.6 NVP断开与VP节点的链接

```go
func (rpc *RPC) DisconnectNodeVP(hash string) (bool, StdError)
```

| 参数  |        |
| --- | ------ |
| 参数1 | 节点hash |

应用实例

```go
rp := rpc.NewRPC()
rp.DisconnectNodeVP("0x")
```

### 5.4.7 获取节点状态信息

```go
func (rpc *RPC) GetNodeStates() ([]NodeStateInfo, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
info, err := rp.GetNodeStates()
```



### 5.4.8 替换节点证书

```go
func (rpc *RPC) ReplaceNodeCerts(hostname string) (string, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 节点名    |
| 返回值  |        |
| 返回值1 | 交易hash |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
rp.ReplaceNodeCerts("node1")
```



## 5.5 区块服务

```go
// Block is packaged result of Block
type Block struct {
	Version      string // block version
	Number       uint64
	Hash         string
	ParentHash   string
	WriteTime    uint64
	AvgTime      int64
	TxCounts     uint64
	MerkleRoot   string
	Transactions []TransactionInfo
}

```

### 5.5.1 获取最后一个区块的信息

```go
func (rpc *RPC) GetLatestBlock() (*Block, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
rp.GetLatestBlock()
```



### 5.5.2 分页获取区块信息

```go
type PageResult struct {
   HasMore bool        `json:"hasmore"` // 是否还有符合条件的区块或交易
   Data    interface{} `json:"data"`    // 查询到的区块信息
}

type Bookmark struct {
	BlockNumber uint64 `json:"blkNum"`
	TxIndex     int64  `json:"txIndex"`
}

type Metadata struct {
	PageSize int32     `json:"pagesize"`
	Bookmark *Bookmark `json:"bookmark"`

	// true means to search backward from the bookmark position,
	// otherwise to search forward from the bookmark position
	Backward bool `json:"backward"`
}

```

```go
func (rpc *RPC) GetBlocksWithLimit(from, to uint64, isPlain bool, metadata *Metadata) (*PageResult, StdError)
```

| 参数  |                      |
| --- | -------------------- |
| 参数1 | 查询起始区块号              |
| 参数2 | 查询最后区块号              |
| 参数3 | 是否在区块中去掉交易信息， true去除 |
| 参数4 | 查询参数                 |

应用实例

```go
rp := rpc.NewRPC()
latestBlock, err := rp.GetLatestBlock()
metadata := &rpc.Metadata{
	PageSize: 5,
}
pageResult, err := rp.GetBlocksWithLimit(latestBlock.Number-1, latestBlock.Number, true, metadata)
```



### 5.5.3 通过区块hash获取区块信息

```go
func (rpc *RPC) GetBlockByHash(blockHash string, isPlain bool) (*Block, StdError)
```

| 参数  |                     |
| --- | ------------------- |
| 参数1 | 区块hash              |
| 参数2 | 区块信息是否去除交易信息，true去除 |

应用实例

```go
rp := rpc.NewRPC()
latestBlock, err := rp.GetLatestBlock()
rp.GetBlockByHash(latestBlock.Hash, true)
```

### 5.5.4 通过区块号获取区块信息

```go
func (rpc *RPC) GetBlockByNumber(blockNum interface{}, isPlain bool) (*Block, StdError)
```

| 参数  |                             |
| --- | --------------------------- |
| 参数1 | 区块号，可使用"latest"获取latest区块信息 |
| 参数2 | 区块信息是否去除交易信息，true去除         |

应用示例

```go
rp := rpc.NewRPC()
latestBlock, err := rp.GetLatestBlock()
rp.GetBlockByNumber("latest", false)
block, err := rp.GetBlockByNumber(latestBlock.Number, true)
```



### 5.5.5 计算区间内区块平均生成时间

```go
func (rpc *RPC) GetAvgGenTimeByBlockNum(from, to uint64) (int64, StdError)
```

| 参数  |             |
| --- | ----------- |
| 参数1 | 起始区块号(区间起始) |
| 参数2 | 末尾区块号(区间末尾) |

应用实例

```go
rp := rpc.NewRPC() 
block, err := rp.GetLatestBlock()
avgTime, err := rp.GetAvgGenTimeByBlockNum(block.Number-2, block.Number)
```

### 5.5.6 查询区间内区块生成的速度以及TPS

```go
type TPSInfo struct {
   StartTime     string
   EndTime       string
   TotalBlockNum uint64
   BlocksPerSec  float64
   Tps           float64
}
```

```text
func (rpc *RPC) QueryTPS(startTime, endTime uint64) (*TPSInfo, StdError)
```

| 参数  |      |
| --- | ---- |
| 参数1 | 起始时间 |
| 参数2 | 终止时间 |

应用实例

```go
rp := rpc.NewRPC()
rp.QueryTPS(1, 1635410111)
```

### 5.5.7 查询创世区块号

```go
func (rpc *RPC) GetGenesisBlock() (string, StdError)
```

```go
rp := rpc.NewRPC()
blkNum, err := rp.GetGenesisBlock()
```



### 5.5.8 查询区块高度

即查询最新的区块号

```go
func (rpc *RPC) GetChainHeight() (string, StdError)
```

```go
rp := rpc.NewRPC()
blkNum, err := rp.GetChainHeight()
```



## 5.6 账户服务

### 5.6.1 获取账户角色

```go
func (rpc *RPC) GetRoles(account string) ([]string, StdError)
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 账户地址     |
| 返回值  |          |
| 返回值1 | 账户所有角色列表 |
| 返回值2 | error    |

应用实例

```go
pwd := "12347890"
rp := rpc.NewRPC()
newAccount, err := account.NewAccount(pwd)
newKey, err := account.GenKeyFromAccountJson(newAccount, pwd)
newKey.(*account.ECDSAKey).GetAddress()
roles, err := rp.GetRoles(newKey.(*account.ECDSAKey).GetAddress().Hex())
```

### 5.6.2 根据角色查询账户

```go
func (rpc *RPC) GetAccountsByRole(role string) ([]string, StdError) 
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 角色名    |
| 返回值  |        |
| 返回值1 | 账户地址列表 |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
rp.GetAccountsByRole("admin")
```

### 5.6.3 获取账户状态

```go
func (rpc *RPC) GetAccountStatus(address string) (string, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 账户地址   |
| 返回值  |        |
| 返回值1 | 账户状态描述 |
| 返回值2 | error  |

应用实例

```go
rp := rpc.NewRPC()
rp.GetAccountStatus("0xfbca6a7e9e29728773b270d3f00153c75d04e1ad")
```



### 5.6.4 获取指定账户的证明路径

```go
func (rpc *RPC) GetAccountProof(account string) (*AccountProofPath, StdError)
```

| 参数   |           |
| ---- | --------- |
| 参数1  | 账户地址      |
| 返回值  |           |
| 返回值1 | 指定账户的证明路径 |
| 返回值2 | error     |

应用示例

```go

proofPath, err := rpcAPI.GetAccountProof(key.GetAddress())
res := rpc.ValidateAccountProof(key.GetAddress(),proofPath)
```



## 5.7 MQ服务

### 5.7.1 获取mq客户端

```go
func (rpc *RPC) GetMqClient() *MqClient
```

```go
rp := rpc.NewRPC()
rp.GetMqClient()
```

### 5.7.2 注册MQ channel

```go
// RegisterMeta mq register
type RegisterMeta struct {
   //queue related
   RoutingKeys []routingKey `json:"routingKeys,omitempty"`
   QueueName   string       `json:"queueName,omitempty"`
   //self info
   From      string `json:"from,omitempty"`
   Signature string `json:"signature,omitempty"`
   // block accounts
   IsVerbose bool `json:"isVerbose"`
   // vm log criteria
   FromBlock string           `json:"fromBlock,omitempty"`
   ToBlock   string           `json:"toBlock,omitempty"`
   Addresses []common.Address `json:"addresses,omitempty"`
   Topics    [][]common.Hash  `json:"topics,omitempty"`
   Delay     bool             `json:"delay"`
}

// QueueRegister MQ register result
type QueueRegister struct {
	QueueName     string
	ExchangerName string
}
```

```go
func (mc *MqClient) Register(id uint, meta *RegisterMeta) (*QueueRegister, StdError)
```

| 参数   |                              |
| ---- | ---------------------------- |
| 参数1  | 节点id                         |
| 参数2  | 事件相关参数， 调用NewRegisterMeta来构造 |
| 返回值  |                              |
| 返回值1 | 注册成功的队列queue以及exchanger名称    |
| 返回值2 | error                        |

应用实例

```go
rp := rpc.NewRPC()
client := rp.GetMqClient()
var hash common.Hash
hash.SetString("123")
guomiKey, _ := gm.GenerateSM2Key()
pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
newAddress := h[12:]
queneName := fmt.Sprintf("testQueue%d", time.Now().Unix())
rm := rpc.NewRegisterMeta(common.BytesToAddress(newAddress).Hex(), queneName, MQBlock).SetTopics(1, hash)
rm.Sign(guomiKey)
regist, err := client.Register(1, rm)
```

### 5.7.3 注销MQ channel

```go
// UnRegisterMeta UnRegisterMeta
type UnRegisterMeta struct {
   From         string
   QueueName    string
   ExchangeName string
   Signature    string
}
// QueueUnRegister MQ unRegister result
type QueueUnRegister struct {
	Count   uint
	Success bool
	Error   error
}
```

```go
func (mc *MqClient) UnRegister(id uint, meta *UnRegisterMeta) (*QueueUnRegister, StdError)
```

| 参数  |                              |
| --- | ---------------------------- |
| 参数1 | 节点编号id，从1开始                  |
| 参数2 | 注销参数， 调用NewUnRegisterMeta来构造 |

应用实例

```go
rmq := rpc.NewUnRegisterMeta(address, queneName, regist.ExchangerName)
rmq.Sign(guomiKey)
client.UnRegister(1, rmq)
```

### 5.7.4 获取指定节点所有的queue名

```go
func (mc *MqClient) GetAllQueueNames(id uint) ([]string, StdError)
```

| 参数   |             |
| ---- | ----------- |
| 参数1  | 节点编号id，从1开始 |
| 返回值  |             |
| 返回值1 | queue名列表    |
| 返回值2 | error       |



应用实例

```go
rp := rpc.NewRPC()
client := rp.GetMqClient()
queues, err := client.GetAllQueueNames(1)
```

### 5.7.5 与broker建立连接

```go
func (mc *MqClient) InformNormal(id uint, brokerURL string) (bool, StdError)
```

| 参数   |                              |
| ---- | ---------------------------- |
| 参数1  | 节点编号id，从1开始                  |
| 参数2  | brokerURL， 为空则表示使用平台配置的默认URL |
| 返回值  |                              |
| 返回值1 | 是否正常建立连接                     |
| 返回值2 | error                        |

```go
success, err := client.InformNormal(1, "")
```

### 5.7.6 获取节点当前exchanger名

```go
func (mc *MqClient) GetExchangerName(id uint) (string, StdError)
```

| 参数   |             |
| ---- | ----------- |
| 参数1  | 节点编号id，从1开始 |
| 返回值  |             |
| 返回值1 | exchanger名  |
| 返回值2 | error       |

```go
client.GetExchangerName(1)
```

### 5.7.7 删除exchanger

```go
func (mc *MqClient) DeleteExchange(id uint, exchange string) (bool, StdError)
```

| 参数   |            |
| ---- | ---------- |
| 参数1  | 节点编号id     |
| 参数2  | exchange 名 |
| 返回值  |            |
| 返回值1 | 是否删除成功     |
| 返回值2 | error      |

```go
client.DeleteExchange(1, regist.ExchangerName)
```

### 5.7.8 添加监听

```go
// MqListener handle register
type MqListener interface {
   HandleDelivery(data []byte)
}
```

```go
func (mc *MqClient) Listen(queue, url string, autoAck bool, listener MqListener) StdError
```

| 参数  |              |
| --- | ------------ |
| 参数1 | queue名       |
| 参数2 | queue 的监听url |
| 参数3 | 是否自动确认       |
| 参数4 | 对获取到的数据处理    |

应用实例

```go
listener := new(MqListener)
err := client.Listen("hello", DefaultAmdpURL, true, *listener)
```

## 5.8 Kvsql服务





## 5.9 webSocket服务

```go
// WsEventHandler web socket event handler
// note: if you unsubscribe a event, the OnClose() will never be called
// even if the connection is closed
type WsEventHandler interface {
   // when subscribe success
   OnSubscribe()
   // when unsubscribe
   OnUnSubscribe()
   // when receive notification
   OnMessage([]byte)
   // when connection closed
   OnClose()
}

type TestEventHandler struct {
}

func (h *TestEventHandler) OnSubscribe() {
	fmt.Println("订阅成功！")
}

func (h *TestEventHandler) OnUnSubscribe() {
	fmt.Println("取消订阅成功！")
}

func (h *TestEventHandler) OnMessage(message []byte) {
	fmt.Printf("收到信息: %s\n", message)
}

func (h *TestEventHandler) OnClose() {
	fmt.Println("连接关闭回调调用！")
}
// --------------------- 
type EventFilter interface {
	Serialize() interface{}
	GetEventType() EventType
}

const (
	// BLOCKEVENT block
	BLOCKEVENT EventType = "block"
	// SYSTEMSTATUSEVENT systemStatus
	SYSTEMSTATUSEVENT EventType = "systemStatus"
	// LOGSEVENT logs
	LOGSEVENT EventType = "logs"
)

// BlockEventFilter block filter
type BlockEventFilter struct {
	eventType EventType
	BlockInfo bool
}
// SystemStatusFilter system status filter
type SystemStatusFilter struct {
	eventType         EventType
	Modules           []string `json:"modules,omitempty"`
	ModulesExclude    []string `json:"modules_exclude,omitempty"`
	Subtypes          []string `json:"subtypes,omitempty"`
	SubtypesExclude   []string `json:"subtypes_exclude,omitempty"`
	ErrorCodes        []string `json:"error_codes,omitempty"`
	ErrorCodesExclude []string `json:"error_codes_exclude,omitempty"`
}

// LogsFilter logs filter
type LogsFilter struct {
	eventType EventType
	FromBlock uint64           `json:"fromBlock,omitempty"`
	ToBlock   uint64           `json:"toBlock,omitempty"`
	Addresses []string         `json:"addresses,omitempty"`
	Topics    [4][]common.Hash `json:"topics,omitempty"`
}

```

### 5.9.1 获取websocket客户端

```go
func (rpc *RPC) GetWebSocketClient() *WebSocketClient
```

| 说明  | 获取WebSokcet客户端，所有的web socket的接口都是基于web socket客户端 |
| --- | ------------------------------------------------ |

```go
wsRPC := rpc.NewRPC()
wsCli := wsRPC.GetWebSocketClient()
```

### 5.9.2 订阅提案事件

```go
func (wscli *WebSocketClient) SubscribeForProposal(nodeIndex int, eventHandler WsEventHandler) (SubscriptionID, StdError)
```

| 参数   |                                                                                                                                                                              |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 参数1  | 节点序号，从1开始。代表节点编号，与hpc.toml配置文件中jsonRPC下的nodes配置一致，1号表示向1号节点订阅，2号表示向2号节点订阅，以此类推。                                                                                              |
| 参数2  | 用户自定义回调函数，需要实现WsEventHandler接口，当事件发生时会自动触发对应回调函数：在订阅成功时会触发OnSubscribe()，在取消订阅成功时会触发OnUnSubscribe()，在接收到事件推送时会触发OnMessage([]byte)，在连接关闭时会触发OnClose()（若已经取消订阅，那么将不会再触发关闭连接回调）。 |
| 返回值  |                                                                                                                                                                              |
| 返回值1 | 订阅id                                                                                                                                                                         |
| 返回值2 | error                                                                                                                                                                        |

```go
wsCli := wsRPC.GetWebSocketClient()
wsCli.SubscribeForProposal(1, &TestEventHandler{})
```



### 5.9.3 订阅事件

```go
func (wscli *WebSocketClient) Subscribe(nodeIndex int, filter EventFilter, eventHandler WsEventHandler) (SubscriptionID, StdError)
```

| 说明   | 用来向hyperchain中的某个节点订阅某个事件                                                                                                                                                    |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 参数   |                                                                                                                                                                              |
| 参数1  | 节点序号，从1开始。代表节点编号，与hpc.toml配置文件中jsonRPC下的nodes配置一致，1号表示向1号节点订阅，2号表示向2号节点订阅，以此类推。                                                                                              |
| 参数2  | 事件的过滤条件，要求实现EventFilter接口                                                                                                                                                    |
| 参数3  | 用户自定义回调函数，需要实现WsEventHandler接口，当事件发生时会自动触发对应回调函数：在订阅成功时会触发OnSubscribe()，在取消订阅成功时会触发OnUnSubscribe()，在接收到事件推送时会触发OnMessage([]byte)，在连接关闭时会触发OnClose()（若已经取消订阅，那么将不会再触发关闭连接回调）。 |
| 返回值  |                                                                                                                                                                              |
| 返回值1 | 订阅id                                                                                                                                                                         |
| 返回值2 | error                                                                                                                                                                        |

```go
bf := NewBlockEventFilter()
bf.BlockInfo = true
wsCli := wsRPC.GetWebSocketClient()
wid, err := wsCli.Subscribe(1, bf, &TestEventHandler{})
```

### 5.9.4 取消订阅

```go
func (wscli *WebSocketClient) UnSubscribe(id SubscriptionID) StdError 
```

| 参数  |      |
| --- | ---- |
| 参数1 | 订阅id |

```go
wsCli.UnSubscribe(wid)
```



### 5.9.5 获取节点所有的订阅信息

```go
func (wscli *WebSocketClient) GetAllSubscription(nodeIndex int) ([]Subscription, StdError)
```

| 参数  |            |
| --- | ---------- |
| 参数1 | 节点id， 从1开始 |

```go
infos, err := wsCli.GetAllSubscription(1)
```

### 5.9.6 关闭节点的websocket连接

```go
func (wscli *WebSocketClient) CloseConn(nodeIndex int) StdError
```

| 参数  |            |
| --- | ---------- |
| 参数1 | 节点id， 从1开始 |

```go
infos, err := wsCli.CloseConn(1)
```

## 5.10 归档服务

### 5.10.1 列出所有快照信息

```go
type Manifest struct {
    Height         uint64 `json:"height"`
	Genesis        uint64 `json:"genesis"`
	BlockHash      string `json:"hash"`
	FilterID       string `json:"filterId"`
	MerkleRoot     string `json:"merkleRoot"`
	Namespace      string `json:"Namespace"`
	TxCount        uint64 `json:"txCount"`
	InvalidTxCount uint64 `json:"invalidTxCount,omitEmpty"` // online invalid tx record number
	Status         uint   `json:"status"` // 0: old snapshot; 1: snapshot not finish; 2: snapshot finish; 3: recover; 4: sync
	DBVersion      string `json:"dbVersion"` // db version, "0.0": no metadb; "1.0": has metadb; "1.1": has didDB and credentialDB
	// use for hyperchain
	Date string `json:"date"`
}
// Manifests Manifests
type Manifests []Manifest
```

```go
func (rpc *RPC) ListSnapshot() (Manifests, StdError)
```

```go
rp := rpc.NewRPC()
rp.ListSnapshot()
```

### 5.10.2 数据直接归档

```go
func (rpc *RPC) ArchiveNoPredict(blockNumber uint64) (string, StdError)
```

| 参数   |         |
| ---- | ------- |
| 参数1  | 归档的区块高度 |
| 返回值  |         |
| 返回值1 | 快照id    |
| 返回值2 | error   |

应用实例

```go
rp.ArchiveNoPredict(20)
```

### 5.10.3 检查数据归档是否完成

```go
func (rpc *RPC) QueryArchiveExist(filterID string) (bool, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 快照id  |
| 返回值  |       |
| 返回值1 | 是否完成  |
| 返回值2 | error |

应用实例

```go
rp.QueryArchiveExist("0xcc2cc319fe2a5782ea15433206745a8f")
```

### 5.10.4 查询数据归档状态

```go
func (rpc *RPC) QueryArchive(filterID string) (string, StdError)
```

| 参数   |          |
| ---- | -------- |
| 参数1  | 快照id     |
| 返回值  |          |
| 返回值1 | 数据归档状态描述 |
| 返回值2 | error    |

应用实例

```go
rp.QueryArchive("0xcc2cc319fe2a5782ea15433206745a8f")
```

### 5.10.5 查询最近一次归档的状态

```go
// ArchiveResult used for return archive result, tell caller which step is processing
type ArchiveResult struct {
	FilterID string `json:"filterId"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
}

func (rpc *RPC) QueryLatestArchive() (*ArchiveResult, StdError) 
```

```go
rp.QueryLatestArchive()
```



### 5.10.6 在某个已经存在的区块高度归档

```go
func (rpc *RPC) MakeSnapshot4Flato(blockHeight interface{}) (string, StdError)
```

| 说明   | 在某已存在的blockHeight高度进行快照操作，要求入参区块高度小于或等于链上节点最新的checkpoint。(仅flato适用) |
| ---- | ------------------------------------------------------------------- |
| 参数   |                                                                     |
| 参数1  | 区块高度，小于或等于链上节点最新的checkpoint                                         |
| 返回值  |                                                                     |
| 返回值1 | 快照标号，用来查询是否已经归档                                                     |
| 返回值2 | err                                                                 |

```go
res, err := MakeSnapshot4Flato(1)
```





## 5.11 配置服务

### 5.11.1 查询提案

```go
// ProposalRaw ProposalRaw
type ProposalRaw struct {
   ID        uint64      `json:"id,omitempty"`
   Code      string      `json:"code,omitempty"`
   Timestamp int64       `json:"timestamp,omitempty"`
   Timeout   int64       `json:"timeout,omitempty"`
   Status    string      `json:"status,omitempty"`
   Assentor  []*VoteInfo `json:"assentor,omitempty"`
   Objector  []*VoteInfo `json:"objector,omitempty"`
   Threshold uint32      `json:"threshold,omitempty"`
   Score     uint32      `json:"score,omitempty"`
   Creator   string      `json:"creator,omitempty"`
   Version   string      `json:"version,omitempty"`
   Type      string      `json:"type,omitempty"`
   Completed string      `json:"completed,omitempty"`
   Cancel    string      `json:"cancel,omitempty"`
}
type VoteInfo struct {
   Addr   string `json:"addr,omitempty"`
   TxHash string `json:"txHash,omitempty"`
}

func (rpc *RPC) GetProposal() (*ProposalRaw, StdError)
```

应用实例

```go
rp.GetProposal()
```

### 5.11.2 查询配置信息

```go
func (rpc *RPC) GetConfig() (string, StdError) 
```

应用实例

```go
rp.GetConfig()
```



### 5.11.3 查询Hosts

```go
func (rpc *RPC) GetHosts(role string) (map[string][]byte, StdError) 
```

| 参数   |                         |
| ---- | ----------------------- |
| 参数1  | 角色名                     |
| 返回值  |                         |
| 返回值1 | key为hostname，value为节点公钥 |
| 返回值2 | error                   |

应用实例

```go
rp.GetHosts("vp")
```

### 5.11.4 查询VSet

```go
func (rpc *RPC) GetVSet() ([]string, StdError)
```

| 返回值  |        |
| ---- | ------ |
| 返回值1 | vset列表 |
| 返回值2 | error  |

应用实例

```go
rp.GetVSet()
```

### 5.11.5 查询链上角色权重信息

```go
func (rpc *RPC) GetAllRoles() (map[string]int, StdError)
```

| 返回值  |           |
| ---- | --------- |
| 返回值1 | 角色以及对应的权重 |
| 返回值2 | error     |

应用实例

```go
rp.GetAllRoles()
```

### 5.11.6 检查链上角色是否存在

```go
func (rpc *RPC) IsRoleExist(role string) (bool, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 角色名   |
| 返回值  |       |
| 返回值1 | 是否存在  |
| 返回值2 | error |

应用实例

```go
rp.IsRoleExist("admin")
```

### 5.11.7 根据合约命名查询合约地址

```go
func (rpc *RPC) GetAddressByName(name string) (string, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 合约名   |
| 返回值  |       |
| 返回值1 | 合约地址  |
| 返回值2 | error |

应用实例

```go
rp.GetAddressByName("HashContract")
```

### 5.11.8 根据合约地址获取合约名

```go
func (rpc *RPC) GetNameByAddress(address string) (string, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 合约地址  |
| 返回值  |       |
| 返回值1 | 合约名   |
| 返回值2 | error |

应用实例

```go
rp.GetNameByAddress("0x8485147cbf02dec93ee84f81824a3b60e355f5cd")
```

### 5.11.9 获取所有<合约地址, 合约名>的映射

```go
func (rpc *RPC) GetAllCNS() (map[string]string, StdError)
```

```go
rp.GetAllCNS()
```

### 5.11.10 获取创世信息

```go
func (rpc *RPC) GetGenesisInfo() (string, StdError)
```

```go
rp.GetGenesisInfo()
```

## 5.12 权限服务

### 5.12.1 为账户增加节点级角色

```go
func (rpc *RPC) AddRoleForNode(address string, roles ...string) StdError
```

| 参数  |      |
| --- | ---- |
| 参数1 | 账户地址 |
| 参数2 | 角色   |

应用实例

```go
// acAddr 账户地址
rp.AddRoleForNode(acAddr, "admin")
```

### 5.12.2 为账户删除节点级角色

```go
func (rpc *RPC) DeleteRoleFromNode(address string, roles ...string) StdError
```

| 参数  |      |
| --- | ---- |
| 参数1 | 账户地址 |
| 参数2 | 角色   |

应用实例

```go
// acAddr 账户地址
rp.DeleteRoleFromNode(acAddr, "admin")
```

### 5.12.3 查询账户的节点级角色

```go
func (rpc *RPC) GetRoleFromNode(address string) ([]string, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 账户地址  |
| 返回值  |       |
| 返回值1 | 角色列表  |
| 返回值2 | error |

应用实例

```go
rp.GetRoleFromNode(acAddr)
```

### 5.12.4 查询节点级某角色的账户列表

```go
func (rpc *RPC) GetAddressFromNode(role string) ([]string, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 角色名    |
| 返回值  |        |
| 返回值1 | 账户地址列表 |
| 返回值2 | error  |

应用实例

```go
rp.GetAddressFromNode("admin")
```

### 5.12.5 查询所有的节点级角色

```go
func (rpc *RPC) GetAllRolesFromNode() ([]string, StdError)
```

| 返回值  |       |
| ---- | ----- |
| 返回值1 | 角色名列表 |
| 返回值2 | error |

### 5.12.6 设置节点级接口权限管理规则列表

```go
// InspectorRule is the rule of api filter
type InspectorRule struct {
	// AllowAnyone determines whether the resources can be accessed freely by anyone
	AllowAnyone bool `json:"allow_anyone" mapstructure:"allow_anyone"`

	// AuthorizedRoles determine who can access the resource if the resources can not be accessed freely
	AuthorizedRoles []string `json:"authorized_roles" mapstructure:"authorized_roles"`

	// ForbiddenRoles determine who can not access the resources though he has the authorized roles
	ForbiddenRoles []string `json:"forbidden_roles" mapstructure:"forbidden_roles"`

	// ID is the identity sequence number for priority
	ID int `json:"id" mapstructure:"id"`

	// Name is the identity string for reading
	Name string `json:"name" mapstructure:"name"`

	// To is  the `to` address used to define resources of tx api
	Method []string `json:"method" mapstructure:"method"`
}

func (rpc *RPC) SetRulesInNode(rules []*InspectorRule) StdError 
```

```go
rp := rpc.NewRPC()
rule := &InspectorRule{
		Name:            "rule",
		ID:              1,
		AllowAnyone:     false,
		AuthorizedRoles: []string{"accountManager"},
		Method:          []string{"account_*"},
	}
err := rp.SetRulesInNode([]*InspectorRule{rule})

```

### 5.12.7 查询节点级接口权限管理规则列表

```go
func (rpc *RPC) GetRulesFromNode() ([]*InspectorRule, StdError)
```

| 返回值  |            |
| ---- | ---------- |
| 返回值1 | 接口权限管理规则列表 |
| 返回值2 | error      |



## 5.13 DID服务

### 5.13.1 发送DID交易

```go
func (rpc *RPC) SendDIDTransaction(transaction *Transaction, key interface{}) (*TxReceipt, StdError)
```

应用实例

```go
rp := rpc.NewRPC()
password := "password"
accountJson, _ = account.NewAccountSm2(password)
key, _ := account.GenKeyFromAccountJson(accountJson, password)
suffix := "didAddress_suffix"
//生成DID账户
didKey := account.NewDIDAccount(key.(account.Key), rp.GetChainID(), suffix)
//生成DID账户的公钥
puKey, _ := rpc.GenDIDPublicKeyFromDIDKey(didKey)
//生成DID账户对应的文档
admins := []string{didKey.GetAddress()}
document := rpc.NewDIDDocument(didKey.GetAddress(), puKey, admins)
//构建DID账户注册交易
tx := rpc.NewTransaction(didKey.GetAddress()).Register(document)
//发送DID交易
res, err := rp.SendDIDTransaction(tx, didKey)
```



### 5.13.2 查询chainID

```go
func (rpc *RPC) GetChainID() string {
    return rpc.chainID
}
```

### 5.13.3 查询DID文档

```go
type DIDDocument struct {
	DidAddress string                 `json:"didAddress,omitempty"`
	State      int                    `json:"state,omitempty"`
	PublicKey  *DIDPublicKey          `json:"publicKey,omitempty"`
	Admins     []string               `json:"admins,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

func (rpc *RPC) GetDIDDocument(didAddress string) (*DIDDocument, StdError)
```

| 参数   |             |
| ---- | ----------- |
| 参数1  | did账户地址     |
| 返回值  |             |
| 返回值1 | didDocument |
| 返回值2 | error       |

### 5.13.4 查询凭证基础信息

```go
func (rpc *RPC) GetCredentialPrimaryMessage(id string) (*DIDCredential, StdError)
```

| 参数   |         |
| ---- | ------- |
| 参数1  | 凭证id    |
| 返回值  |         |
| 返回值1 | did凭证信息 |
| 返回值2 | error   |

应用实例

```go
rp := rpc.NewRPC()
_, err = rp.GetCredentialPrimaryMessage(credId)
```

### 5.13.5 检查凭证是否有效

```go
func (rpc *RPC) CheckCredentialValid(id string) (bool, StdError)
```

| 参数   |           |
| ---- | --------- |
| 参数1  | 凭证id      |
| 返回值  |           |
| 返回值1 | did凭证是否有效 |
| 返回值2 | error     |

```go
//credID为待检查的凭证ID
valid, _ := rp.CheckCredentialValid(credID)
```



### 5.13.6 检查凭证是否吊销

```go
func (rpc *RPC) CheckCredentialAbandoned(id string) (bool, StdError) 
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 凭证id   |
| 返回值  |        |
| 返回值1 | 凭证是否吊销 |
| 返回值2 | error  |

```go
//credID为待检查的凭证ID
isAbandoned, _ := rp.CheckCredentialAbandoned(credID)
```

### 5.13.7 获取节点chainId

```go
func (rpc *RPC) GetNodeChainID() (string, StdError)
```

| 返回值  |         |
| ---- | ------- |
| 返回值1 | chainId |
| 返回值2 | error   |

```go
rp := rpc.NewRPC()
res, _ := rp.GetNodeChainID()
```



### 5.13.8 设置本地chainId为节点chainId

```go
func (rpc *RPC) SetLocalChainID() error
```

```go
rp.SetLocalChainID()
```

### 5.13.9 获取DID账户的公钥

```go
func GenDIDPublicKeyFromDIDKey(didKey *account.DIDKey) (*DIDPublicKey, error)
```

| 参数   |         |
| ---- | ------- |
| 参数1  | did账户   |
| 返回值  |         |
| 返回值1 | did账户公钥 |
| 返回值2 | error   |

```go
password := "hyper"
accountJson, _ := account.NewAccountSm2(password)
key, _ := account.GenDIDKeyFromAccountJson(accountJson, password)
suffix := common.RandomString(10)
//生成did账户
didKey := account.NewDIDAccount(key.(account.Key), rpc.chainID, suffix)
//生成公钥
puKey, _ := rpc.GenDIDPublicKeyFromDIDKey(didKey)
```

### 5.13.10 新建DIDDocument

```go
func NewDIDDocument(didAddress string, publicKey *DIDPublicKey, admins []string) *DIDDocument
```

| 参数   |           |
| ---- | --------- |
| 参数1  | DID账户的地址  |
| 参数2  | DID账户的公钥  |
| 参数3  | DID账户的管理员 |
| 返回值  |           |
| 返回值1 | DID账户的文档  |

```go
//生成DID账户
didKey := account.NewDIDAccount(key.(account.Key), rpc.chainID, suffix)
//生成账户公钥
puKey, _ := rpc.GenDIDPublicKeyFromDIDKey(didKey)
admins := []string{didKey.GetAddress()}
//生成DID账户文档
document := rpc.NewDIDDocument(didKey.GetAddress(), puKey, admins)
```

### 5.13.11 新建DID凭证



```go
func NewDIDCredential(ctype, issuer, holder, subject string, issuanceDate, expirationDate int64) *DIDCredential
```

| 参数   |           |
| ---- | --------- |
| 参数1  | DID凭证的类型  |
| 参数2  | DID凭证的签发者 |
| 参数3  | DID凭证的持有者 |
| 参数4  | 凭证信息      |
| 参数5  | 凭证签发时间    |
| 参数6  | 凭证过期时间    |
| 返回值  |           |
| 返回值1 |           |

```go
//issuer和holder都为DID账户
cred := rpc.NewDIDCredential("type", issuer.GetAddress(), holder.GetAddress(), "credential message", time.Now().UnixNano(), time.Now().UnixNano()+1e11)
```



## 5.14 文件服务

```go
type FileExtra struct {
   Hash            string   `json:"hash"`
   FileName        string   `json:"file_name"`
   FileSize        int64    `json:"file_size"`
   UpdateTime      string   `json:"update_time"`
   NodeList        []string `json:"node_list"`
   UserList        []string `json:"user_list,omitempty"`
   FileDescription string   `json:"file_description"`
}
```

### 5.14.1 文件上传

```go
func (rpc *RPC) FileUpload(filePath string, description string, userList []string, nodeIdList []int, pushNodes []int, accountJson string, password string) (string, StdError)
```

| 参数          |        |
| ----------- | ------ |
| filePath    | 文件路径   |
| description | 文件描述信息 |
| userList    | 用户白名单  |
| nodeIdList  | 节点白名单  |
| pushNodes   | 上传节点   |
| accountJson | 账户     |
| password    | 账户密码   |
| 返回值         |        |
| 返回值1        | 文件hash |
| 返回值2        | err    |

```go
filePath := filepath.Join(dirPath, "upload1.txt")
nodeIdList := []int{1, 2, 3}
// 设置文件的白名单
userList := []string{fmKey.GetAddress().Hex()}
txHash, err = rpc.FileUpload(filePath, "des", userList, nodeIdList, nodeIdList, accountJson, password)
```



### 5.14.2 文件下载

```go
func (rpc *RPC) FileDownload(tarPath, hash, owner string, nodeID int, accountJson string, password string) (string, StdError)
```

| 参数          |                                                                |
| ----------- | -------------------------------------------------------------- |
| tarPath     | tarPath有两种使用：1.传有效目录，会在给路径下以hash作为文件名保存文件；2.传有效文件路径,对该文件进行断点续传 |
| hash        | 文件hash                                                         |
| owner       | 文件持有者                                                          |
| nodeID      | 表示发送下载请求的节点ID                                                  |
| accountJson | 账户                                                             |
| password    | 账户密码                                                           |
| 返回值         |                                                                |
| 返回值1        | 下载路径                                                           |
| 返回值2        | err                                                            |

```go
fmTempKey, _  := account.GenKeyFromAccountJson(accountJson, password)
fmKey         := fmTempKey.(*account.SM2Key)
addr := fmKey.GetAddress()
txFrom := hex.EncodeToString(addr[:])
downloadPath, err := rp.FileDownload(dirPath, fileHash, txFrom, 1, accountJson, password)
```



### 5.14.3 文件信息更新

```go
func (rpc *RPC) FileUpdate(fileUpdateTX *Transaction) StdError
```

**** **文件信息放在交易体的extra字段**

| 参数           |     |
| ------------ | --- |
| fileUpdateTX | 交易体 |
| 返回值          |     |
| 返回值1         | err |

```go
fileExtra, err := rp.GetFileExtraByExtraId(fileHash)
assert.Nil(err)
newUserList := append(fileExtra.UserList, fmKey2.GetAddress().Hex())
fileExtra.UserList = newUserList
newFileExtraJson, err := fileExtra.ToJson()
assert.Nil(err)
addr := fmKey.GetAddress()
fileUpdateTx := rpc.NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0).Extra(newFileExtraJson)
fileUpdateTx.SetExtraIDString(fileHash)
assert.Nil(err)
fileUpdateTx.Sign(fmKey)
err = rp.FileUpdate(fileUpdateTx)
```

### 5.14.4 推送文件

```go
func (rpc *RPC) FilePush(hash string, pushNodes []int, accountJson, password string, nodeID int) (string, StdError)
```

| 参数          |         |
| ----------- | ------- |
| hash        | 文件hash  |
| pushNodes   | 推送到节点列表 |
| accountJson | 账户      |
| password    | 账户密码    |
| nodeID      | 发送请求的节点 |
| 返回值         |         |
| 返回值1        | 推送结果    |
| 返回值2        | err     |

```go
rp.FilePush(fileHash, []int{1, 2, 3}, accountJson, password, 1)
```



### 5.14.5 通过extraId获取文件信息FileExtra

```go
func (rpc *RPC) GetFileExtraByExtraId(extraId string) (*FileExtra, error)
```

| 参数      |         |
| ------- | ------- |
| extraId | extraId |
| 返回值     |         |
| 返回值1    | 文件信息    |
| 返回值2    | err     |



### 5.14.6  通过filter获取文件信息FileExtra

```go
func (rpc *RPC) GetFileExtraByFilter(from, extraId string) (*FileExtra, StdError)
```

应用示例

```go
addr := fmKey.GetAddress()
txFrom := hex.EncodeToString(addr[:])
fileExtra, err := rp.GetFileExtraByFilter(txFrom, fileHash)
```

### 5.14.7 通过交易哈希获取文件信息FileExtra

```go
func (rpc *RPC) GetFileExtraByTxHash(txHash string) (*FileExtra, StdError)
```

## 5.15 跨分区服务

跨分区服务需要配合BVM进行使用，详情见**4.2.9** 章节

```go
type CrossChainMethod string

const (
   // InvokeAnchorContract method used for normal cross chain request
   InvokeAnchorContract CrossChainMethod = "invokeAnchorContract"
   // InvokeTimeoutContract method used for timeout cross chain request
   InvokeTimeoutContract CrossChainMethod = "invokeTimeoutContract"
   // InvokeContract method used for cross chain request
   InvokeContract CrossChainMethod = "invokeContract"
)
```

### 5.15.1 部署跨分区合约

```go
func (rpc *RPC) SignAndDeployCrossChainContract(transaction *Transaction, key interface{}) (*TxReceipt, StdError)
```

| 参数   |       |
| ---- | ----- |
| 参数1  | 交易体   |
| 参数2  | 账户key |
| 返回值  |       |
| 返回值1 | 交易回执  |



### 5.15.2 调用跨分区方法

```go
func (rpc *RPC) SignAndInvokeCrossChainContract(transaction *Transaction, methodName CrossChainMethod, key interface{}) (*TxReceipt, StdError)
```

| 参数   |                       |
| ---- | --------------------- |
| 参数1  | 交易体                   |
| 参数2  | 方法名（CrossChainMethod） |
| 参数3  | 账户key                 |
| 返回值  |                       |
| 返回值1 | 交易回执                  |

```go
rp:=rpc.NewRPC()
key, _ := account.NewAccountFromAccountJSON(accountJsons[0], "")
operation1 := bvm.NewSystemAnchorOperation(bvm.RegisterAnchor, "node1", "ns2")
payload1 := bvm.EncodeOperation(operation1)
tx1 := rpc.NewTransaction(key.GetAddress().Hex()).Invoke(operation1.Address(), payload1).VMType(BVM)
rp.Namespace("global")
re1, err := rpc.SignAndInvokeCrossChainContract(tx1, InvokeAnchorContract, key)

```



## 5.16 交易证明

```go
// MerkleProofNode struct
type MerkleProofNode struct {
	Hash  []byte `json:"hash,omitempty"`
	Index int    `json:"index,omitempty"`
}

// MerkleProofPath struct
type MerkleProofPath []*MerkleProofNode

// TxProofPath represents the result returned by tx proof query.
type TxProofPath struct {
   TxProof types.MerkleProofPath `json:"txProof"`
}

// ProofParam contains ledger info and key info two parts param
type ProofParam struct {
	Meta *LedgerMetaParam `json:"meta"`
	Key  *KeyParam        `json:"key"`
}

// LedgerMetaParam is the ledger info related user-defined param
// 用来指定证明操作基于哪个账本来执行
type LedgerMetaParam struct {
	// 用于指定archiveReader上的数据目录
	SnapshotID string `json:"snapshotID"`
	// 在当前的路径下，具体要基于哪个区块的状态获取（验证）证明数据	
	SeqNo      uint64 `json:"seqNo"`
}

// 是指在合约中操作账本数据时所使用的逻辑条件，以此来生成用户想要验证的状态数据
type KeyParam struct {
	// 合约地址
	Address   types.Address `json:"address"`
	// 合约中需要验证的数据结构的变量名	
	FieldName string        `json:"fieldName"`
	// 通过哪些参数可在对应的合约变量下获取到状态数据，以字符串形式传入	
	Params    []string      `json:"params"`
	// 合约虚拟机类型，目前仅支持HVM
	VMType    string        `json:"vmType"`
}

// StateProof is the proof path for a ledger key
type StateProof struct {
    // 从statedb中获取的证明数据
	StatePath   types.ProofPath `json:"statePath"`
	// accountdb中获取的证明数据
	AccountPath types.ProofPath `json:"accountPath"`
}
// Inode struct
type Inode struct {
	Key   []byte `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
	Hash  []byte `json:"hash,omitempty"`
}

// Inodes struct
type Inodes []*Inode

// ProofNode struct
type ProofNode struct {
	// 是否为数据节点
	IsData bool   `json:"isData"`
	// byte数组，base64编码，表示当前节点的最小的key
	Key    []byte `json:"key,omitempty"`
	// 节点的哈希值，base64编码，表示当前节点所有数据的hash
	Hash   []byte `json:"hash,omitempty"`
	// 数组结构，每一个元素都是一个json结构：{key是base64编码的byte数组，hash是base64编码的byte数组}	
	Inodes Inodes `json:"inodes,omitempty"`
	// index表示的是在向下层寻找key的过程中，当前节点使用的是哪个位置的分支继续向下寻找的
	Index  int    `json:"index"`
}

// ProofPath struct
type ProofPath []*ProofNode


```

### 5.16.1 通过交易hash获取证明路径

```go
func (rpc *RPC) GetTxProof(txhash string) (*TxProofPath, StdError)
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 交易hash |
| 返回值  |        |
| 返回值1 | 证明路径   |
| 返回值2 | err    |

```go
block, _ := rp.GetLatestBlock()
info, err := rp.GetTxByBlkHashAndIdx(block.Hash, 0)
if err != nil {
   return
}
res, err2 := rp.GetTxProof(info.Hash)
if err2 != nil {
   return
}
```

### 5.16.2 获取状态数据的证明路径

```go
func (rpc *RPC) GetStateProof(proofParam *ProofParam) (*StateProof, StdError)
```

| 参数   |                                           |
| ---- | ----------------------------------------- |
| 参数1  | 证明用到的参数， 详见5.16部分struct结构                 |
| 返回值  |                                           |
| 返回值1 | 指定合约状态值的证明路径StateProof， 详见5.16部分的struct结构 |
| 返回值2 | err                                       |



```java
// 以HVM为例，需要查询的数据结构变量应为被表述了@StoreField的属性：

@StoreField
private Person p1; // -> fieldName:p1, params无

@StoreFiled
private HyperMap<String, String> map1;
map1.put("key1", "value1"); // -> fieldName:map1, params：key1

@StoreField
private HyperList<String> list1;
list1.add("value1"); // -> fieldName:list1, params：0(list的索引)

@StoreField
private HyperTable table1;
table1.put("col", "cf", "col", "value1"); // -> fieldName:table1, params:[col,cf,col]

@StoreFiled
private NestedMap<String, String> nestedMap1;
nestedMap1.put("key1", "value1"); // -> fieldName:nestedMap1, params：key1, 若为多层嵌套则需要传入多个key组成的数组，直至最后一层NestedMap

```

### 5.16.3 验证状态数据的证明路径是否正确

```go
func (rpc *RPC) ValidateStateProof(proofParam *ProofParam, stateProof *StateProof, merkleRoot string) (bool, StdError)
```

| 参数   |                                                         |
| ---- | ------------------------------------------------------- |
| 参数1  | 状态数据路径参数                                                |
| 参数2  | 查询获得的状态数据路径结构                                           |
| 参数3  | 验证路径区块对应的账本hash，即LedgerMetaParam的seqNo所对应的区块的merkleRoot |
| 返回值  |                                                         |
| 返回值1 | 证明路径是否合法                                                |
| 返回值2 | err                                                     |

```go
	id := "0x5b1a5bb7b10d15bc9d47701eed9c9349"
	seq := 2
	contractAddr := "0x6de31be7a30204189d70bd202340c6d9b395523e"
	merkleRoot := "0xaa2fd673656f4bada6ff6d8588498239eeb3202214a24005d6cf0138a9f30a79"
	proofParam := &ProofParam{
		Meta: &LedgerMetaParam{
			SnapshotID: id,
			SeqNo: uint64(seq),
		},
		Key:  &KeyParam{
			Address:   types.HexToAddress(contractAddr),
			FieldName: "hyperMap1",
			Params:    []string{"key1"},
			VMType:    "HVM",
		},
	}
	// get proof
	proof, err := rp.GetStateProof(proofParam)
	// validate proof
	ok, err := rp.ValidateStateProof(proofParam, proof, merkleRoot)


```

## 5.17 版本管理服务

### 5.17.1 二进制版本上链
- 接口描述：让节点构造一笔交易将节点二进制支持的版本信息上链。该交易实质上由节点账户发起，其他普通账户无法构造。用户再根据接口返回的交易哈希查询交易回执，来确认版本上链是否成功。版本上链的目的是将节点支持的版本信息（SupportedVersion）上链，由内置合约计算出共识网络里的所有节点共有的、大于当前运行版本的版本号有哪些，用于指导下一次系统升级。
- 接口定义：

```go
func (rpc *RPC) SetSupportedVersionByID(id int) (*TxReceipt, StdError) {}
```

| 参数   |                   解释                                 |
| ---- | ------------------------------------------------------- |
| 参数1  | 接收这次请求的节点id号，接收该请求的节点将进行版本上链          |
| 返回值1 | 交易回执                                                |
| 返回值2 | err                                                     |

应用实例：

```go
rpc, err := NewJsonRPC()
// 发送请求
receipt, err := rpc.SetSupportedVersionByID(1)
// 解析回执
if receipt != nil {
	bvmResult := bvm.Decode(receipt.Ret)
	if bvmResult.Success {
		// 表示版本上链成功
	}
}
```


### 5.17.2 查询链运行版本信息
- 接口描述：从节点本地账本里查询当前运行版本 RunningVersion 以及可用于下一次升级的版本 AvailableVersion。

在进行系统升级之前，通常需要先调用该接口来查询目前链可以升级到的目标版本有哪些。为了便于客户端理解与使用，本接口对 tx version、block version 等链级细分版本做了一层映射，映射到了一个大的版本号里，比如 ”2.9.0”，”2.10.0”。下面用一个JSON-RPC请求为例，结果表示Hyperchain ”2.10.0” 支持的最大交易版本为 4.1。通过版本映射，用户只需要关心Hyperchain版本号，无需关心该版本号下面的细分链级版本内容。

```bash
# Request
curl localhost:8081 --data '{"jsonrpc":"2.0","method":"version_getVersions","params":[],"id":1}'

# Response
# 结果表示，当前账本运行的版本为v2.9.0，可用于下一次升级的版本是v2.10.0
{
    "jsonrpc": "2.0",
    "namespace": "global",
    "id": 1,
    "code": 0,
    "message": "SUCCESS",
    "result": {
        "availableHyperchainVersions": {
            "2.10.0": {
                "tx_version": "4.1",
                "block_version": "4.0",
                "encode_version": "0"
            }
        },
        "runningHyperchainVersions": {
            "2.9.0": {
                "tx_version": "4.0",
                "block_version": "4.0",
                "encode_version": "0"
            }
        }
    }
}
```

*注意：该接口不会生成新的交易，是从节点本地账本里直接查询数据。如果查询节点账本落后，则本地账本不是链最新数据。如果要保证获取的是最新数据，请从BVM系统升级版本管理内置合约VersionContract GetLatestVersions()方法查询。*

- 接口定义：

```go
type VersionTag string

type RunningVersion map[VersionTag]string

type VersionResult struct {
	AvailableHyperchainVersion map[string]types.RunningVersion `json:"availableHyperchainVersions"`
	RunningHyperchainVersion   map[string]types.RunningVersion `json:"runningHyperchainVersions"`
}

func (rpc *RPC) GetVersions() (*VersionResult, StdError) {}
```

| 参数   |                   解释                                 |
| ---- | ------------------------------------------------------- |
| 返回值1 | 账本里记录的当前正在运行的Hyperchain版本以及可用于下一次系统升级的Hyperchain版本    |
| 返回值2 | err                                                     |

应用实例：

```go
rpc, err := NewJsonRPC()
// 发送请求
vr, err := rpc.GetVersions()
if err != nil {
	return err
}
// 获取当前链正在运行的版本信息
t.Logf("%#v", vr.RunningHyperchainVersion) // map[string]types.RunningVersion{"2.9.0":types.RunningVersion{"block_version":"4.0", "encode_version":"0", "tx_version":"4.0"}}
// 获取可用于下次升级的版本信息
t.Logf("%#v", vr.AvailableHyperchainVersion) // map[string]types.RunningVersion{"2.10.0":types.RunningVersion{"block_version":"4.0", "encode_version":"0", "tx_version":"4.1"}}
```


### 5.17.3 查询节点支持的版本信息
- 接口描述：从节点本地账本里查询指定节点可以支持的版本信息。比如可以支持的区块版本有哪些、可以支持的交易版本有哪些。

*注意：该接口不会生成新的交易，是从节点本地账本里直接查询数据。如果查询节点账本落后，则本地账本不是链最新数据。如果要保证获取的是最新数据，请从BVM系统升级版本管理内置合约VersionContract GetSupportedVersionByHostname()方法查询。*

- 接口定义：

```go
// VersionTag defines version classification.
type VersionTag string

// SupportedVersion defines which versions the node supports, the node
// may support multiple version for one tag.
type SupportedVersion map[VersionTag][]string


func (rpc *RPC) GetSupportedVersionByHostname(hostname string) (types.SupportedVersion, StdError) {}
```

| 参数   |                   解释                                 |
| ---- | ------------------------------------------------------- |
| 参数1 |    节点hostname
| 返回值1 |   节点支持的版本列表  |
| 返回值2 | err                                                     |

应用实例：

```go
rpc, err := NewJsonRPC()
// 发送请求，查询node1的版本支持列表
sv, gerr := rpc.GetSupportedVersionByHostname("node1")
t.Logf("%#v", sv) // types.SupportedVersion{"block_version":[]string{"1.5", "1.6", "1.7", "1.8", "2.0", "2.1", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9", "3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "4.0"}, "encode_version":[]string{"0"}, "tx_version":[]string{"2.0", "2.1", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9", "3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "4.0", "4.1"}}
```

### 5.17.4 查询Hyperchain版本对应的链级细分版本
- 接口描述：从接收请求的节点查询指定Hyperchain版本支持的链级细分版本的最大版本号。

*注意：该接口只是从接收请求的节点二进制里获取这些信息，并不是从账本中查询，结果取决于接收请求节点的二进制版本。*

- 接口定义：

```go
type VersionTag string

type RunningVersion map[VersionTag]string

func (rpc *RPC) GetHyperchainVersionFromBin(hyperchainVersion string) (types.RunningVersion, StdError) {}
```

| 参数   |                   解释                                 |
| ---- | ------------------------------------------------------- |
| 参数1 |    hyperchain版本
| 返回值1 |   该版本对应的链级最大版本号  |
| 返回值2 | err                                                     |

应用实例：

```go
rpc, err := NewJsonRPC()
rvs, err := rpc.GetHyperchainVersionFromBin("2.9.0")
t.Logf("%#v", rvs) // types.RunningVersion{"block_version":"4.0", "encode_version":"0", "tx_version":"4.0"}
```