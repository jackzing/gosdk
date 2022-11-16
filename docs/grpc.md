- [第六章 GRPC 接口](#第六章-grpc-接口)
	- [6.1 初始化](#61-初始化)
		- [6.1.1 绑定节点](#611-绑定节点)
		- [6.1.2 默认路径初始化](#612-默认路径初始化)
		- [6.1.3 自定义配置路径初始化](#613-自定义配置路径初始化)
	- [6.2 合约服务](#62-合约服务)
		- [6.2.1 部署](#621-部署)
		- [6.2.2 调用](#622-调用)
		- [6.2.3 管理](#623-管理)
		- [6.2.4 投票管理](#624-投票管理)
	- [6.3 交易服务](#63-交易服务)
		- [6.3.1 发送交易](#631-发送交易)
	- [6.4 DID服务](#64-did服务)
		- [6.4.1 发送交易](#641-发送交易)
		- [](#)
	- [6.5 MQ服务](#65-mq服务)
		- [6.5.1 初始化mq客户端](#651-初始化mq客户端)
		- [6.5.2 注册队列](#652-注册队列)
		- [6.5.3 解注册](#653-解注册)
		- [6.5.3 获取所有的队列名](#653-获取所有的队列名)
		- [6.5.4 消费队列](#654-消费队列)
		- [6.5.5 停止消费](#655-停止消费)

# 第六章 GRPC 接口


## 6.1 初始化

```go
// 初始化参数
type GRPCOption interface {
	apply(option *grpcOption)
}
```

### 6.1.1 绑定节点

```go
func BindNodes(s ...int) GRPCOption // 绑定节点
```

| 参数  |                                  |
| --- | -------------------------------- |
| 参数1 | 节点编号，从0开始，对应配置文件中grpc下ports设定的端口 |



### 6.1.2 默认路径初始化

```go
func NewGRPC(opt ...GRPCOption) *GRPC
```

| 参数   |        |
| ---- | ------ |
| 参数1  | 参数列表   |
| 返回值  |        |
| 返回值1 | GRPC实例 |

```go
var opt []rpc.GRPCOption
opt = append(opt, rpc.BindNodes(0, 1))
gp := rpc.NewGRPC(opt...)

// 不带参数
gp := rpc.NewGRPC()
```

### 6.1.3 自定义配置路径初始化

```go
func NewGRPCWithConfPath(path string, opts ...GRPCOption) *GRPC
```

| 参数   |                |
| ---- | -------------- |
| 参数1  | 配置目录 config 目录 |
| 参数1  | 参数列表           |
| 返回值  |                |
| 返回值1 | GRPC实例         |

```go
var opt []rpc.GRPCOption
opt = append(opt, rpc.BindNodes(0, 1))
gp := rpc.NewGRPCWithConfPath("../config", opt...)

// 不带参数
gp := rpc.NewGRPCWithConfPath("../config")
```

## 6.2 合约服务

**** 请求结束，手动close开启的流，否则会出现context cancle的警告**

### 6.2.1 部署

```go
func (c *ContractGrpc) DeployContract(trans *Transaction) (string, StdError)
func (c *ContractGrpc) DeployContractReturnReceipt(trans *Transaction) (*TxReceipt, StdError)
```

| 说明                          |          |
| --------------------------- | -------- |
| DeployContract              | 返回交易hash |
| DeployContractReturnReceipt | 返回交易回执   |

```go
gp := rpc.NewGRPC()
tg, err := gp.NewContractGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
if err != nil {
	return
}
defer tg.Close()
guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
pri := new(gm.SM2PrivateKey)
pri.FromBytes(common.FromHex(guomiPri), 0)
guomiKey := &account.SM2Key{
	&gm.SM2PrivateKey{
		K:         pri.K,
		PublicKey: pri.CalculatePublicKey().PublicKey,
	},
}
transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
transaction.Sign(guomiKey)


ans, err := tg.DeployContract(transaction)
// or
ans, err := tg.DeployContractReturnReceipt(transaction)

```



### 6.2.2 调用

```go
func (c *ContractGrpc) InvokeContract(trans *Transaction) (string, StdError)
func (c *ContractGrpc) InvokeContractReturnReceipt(trans *Transaction) (*TxReceipt, StdError)
```

| 说明                          |          |
| --------------------------- | -------- |
| InvokeContract              | 返回交易hash |
| InvokeContractReturnReceipt | 返回交易回执   |

```go
gp := rpc.NewGRPC()
tg, err := gp.NewContractGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
if err != nil {
	return
}
defer tg.Close()
guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
pri := new(gm.SM2PrivateKey)
pri.FromBytes(common.FromHex(guomiPri), 0)
guomiKey := &account.SM2Key{
	&gm.SM2PrivateKey{
		K:         pri.K,
		PublicKey: pri.CalculatePublicKey().PublicKey,
	},
}
transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
transaction.Sign(guomiKey)

ans, err := tg.DeployContractReturnReceipt(transaction)

ABI, err := abi.JSON(strings.NewReader(abiContract))
if err != nil {
	return
}
packed, err := ABI.Pack("getSum")
if err != nil {
	return
}
transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(addr.ContractAddress, packed)
transaction.Sign(guomiKey)
ans, err := tg.InvokeContract(transaction)
// or 
ans, err := tg.InvokeContractReturnReceipt(transaction)

if err != nil {
	return
}
```

### 6.2.3 管理

```go
func (c *ContractGrpc) MaintainContract(trans *Transaction) (string, StdError)
func (c *ContractGrpc) MaintainContractReturnReceipt(trans *Transaction) (*TxReceipt, StdError)
```

| 说明                            |          |
| ----------------------------- | -------- |
| MaintainContract              | 返回交易hash |
| MaintainContractReturnReceipt | 返回交易回执   |

```go
func demo() {
	g := NewGRPC()
	tg, err := g.NewContractGrpc(ClientOption{
		StreamNumber: 1,
	})
	defer tg.Close()
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	addrTransaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	addrTransaction.Sign(guomiKey)
	addr, err := tg.DeployContractReturnReceipt(addrTransaction)

	// freeze contract
	transactionFreeze := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Maintain(2, addr.ContractAddress, "")
	transactionFreeze.Sign(guomiKey)

	ans, err := tg.MaintainContract(transactionFreeze)

	// unfreeze
	transactionUnFreeze := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Maintain(3, addr.ContractAddress, "")
	transactionUnFreeze.Sign(guomiKey)

	ans2, err := tg.MaintainContractReturnReceipt(transactionUnFreeze)
}
```



### 6.2.4 投票管理

```go
func (c *ContractGrpc) ManageContractByVote(trans *Transaction) (string, StdError)
func (c *ContractGrpc) ManageContractByVoteReturnReceipt(trans *Transaction) (*TxReceipt, StdError)  
```

| 说明                                |          |
| --------------------------------- | -------- |
| ManageContractByVote              | 返回交易hash |
| ManageContractByVoteReturnReceipt | 返回交易回执   |

```go
func createProposalSuccessByGrpc(normalAccountKey account.Key, operation ...bvm.ContractOperation) {
	g := rpc.NewGRPC()
	cli, err := g.NewContractGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
	defer cli.Close()
	contractOpt := bvm.NewProposalCreateOperationForContract(operation...)
	payload := bvm.EncodeOperation(contractOpt)
	tx := rpc.NewTransaction(normalAccountKey.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(rpc.BVM)
	tx.Sign(normalAccountKey)
	re, err := cli.ManageContractByVote(tx)
	// or
	re, err := cli.ManageContractByVoteReturnReceipt(tx)
		
}
```

## 6.3 交易服务

**** 请求结束，手动close开启的流，否则会出现context cancle的警告**

### 6.3.1 发送交易

```go
func (t *TransactionGrpc) SendTransaction(trans *Transaction) (string, StdError)
func (t *TransactionGrpc) SendTransactionReturnReceipt(trans *Transaction) (*TxReceipt, StdError)
```

| 说明                           |          |
| ---------------------------- | -------- |
| SendTransaction              | 返回交易hash |
| SendTransactionReturnReceipt | 返回交易回执   |

```go
gp := rpc.NewGRPC()
tg, err := gp.NewTransactionGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
if err != nil {
	return
}
defer tg.Close()
guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
pri := new(gm.SM2PrivateKey)
pri.FromBytes(common.FromHex(guomiPri), 0)
guomiKey := &account.SM2Key{
	&gm.SM2PrivateKey{
		K:         pri.K,
		PublicKey: pri.CalculatePublicKey().PublicKey,
	},
}
transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
transaction.Sign(guomiKey)
ans, err := tg.SendTransaction(transaction)
// or
ans, err := tg.SendTransactionReturnReceipt(transaction)
```

## 6.4 DID服务

**** 请求结束，手动close开启的流，否则会出现context cancle的警告**

### 6.4.1 发送交易

```go
func (d *DidGrpc) SendDIDTransaction(trans *Transaction) (string, StdError)
func (d *DidGrpc) SendDIDTransactionReturnReceipt(trans *Transaction) (*TxReceipt, StdError)
```

| 说明                           |          |
| ---------------------------- | -------- |
| SendTransaction              | 返回交易hash |
| SendTransactionReturnReceipt | 返回交易回执   |

```go
// 初始化DID服务
func iniDID() {
	accountJson := `{"address":"0xfbca6a7e9e29728773b270d3f00153c75d04e1ad","version":"4.0","algo":"0x13","publicKey":"049c330d0aea3d9c73063db339b4a1a84d1c3197980d1fb9585347ceeb40a5d262166ee1e1cb0c29fd9b2ef0e4f7a7dfb1be6c5e759bf411c520a616863ee046a4","privateKey":"5f0a3ea6c1d3eb7733c3170f2271c10c1206bc49b6b2c7e550c9947cb8f098e3"}`
	key, _ := account.GenKeyFromAccountJson(accountJson, "")
	opt := bvm.NewDIDSetChainIDOperation("chainID_01")
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	_, err := rpc.SignAndInvokeContract(tx, key)
}

func demo() {
	iniDID()
	g := rpc.NewGRPC()
	tg, err := g.NewDidGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
	defer tg.Close()
	var accountJson string
	password := "hyper"
	accountJson, _ = account.NewAccountSm2(password)
	key, _ := account.GenKeyFromAccountJson(accountJson, password)
	suffix := common.RandomString(10)
	didKey := account.NewDIDAccount(key.(account.Key), "chainID_01", suffix)
	puKey, _ := rpc.GenDIDPublicKeyFromDIDKey(didKey)
	document := rpc.NewDIDDocument(didKey.GetAddress(), puKey, nil)
	transaction := rpc.NewTransaction(didKey.GetAddress()).Register(document)
	transaction.Sign(didKey)
	ans, err := tg.SendDIDTransaction(transaction)
	// or
	ans, err := tg.SendDIDTransactionReturnReceipt(transaction)
}
```

### 

## 6.5 MQ服务

使用grpc的MQ服务，需要将节点的`mq.broker.type` 配置为`grpc` ，且开启节点的MQ服务。

```javascript
[mq.broker] #代表启动哪种类型的mq
type = "grpc"
#type = "kafka"
#type = "rabbit"
```

### 6.5.1 初始化mq客户端

```go
func (g *GRPC) NewGrpcMQ() (*grpcMQ, error)
```

应用示例

```go
var opt []rpc.GRPCOption
opt = append(opt, rpc.BindNodes(0, 1))
g := rpc.NewGRPC(opt...)

que, err := g.NewGrpcMQ()
```

****注意**

**mq需要指定节点，否则无法正常使用**



### 6.5.2 注册队列

```go
func (g *grpcMQ) Register(meta *RegisterMeta) (*QueueRegister, StdError)
```

****注意 队列根据名字来区分，开发者不要重复注册同个名字的队列**

```go
func RegisterDemo() {
   var opt []rpc.GRPCOption
   opt = append(opt, rpc.BindNodes(0, 1))
   g := rpc.NewGRPC(opt...)

   que, err := g.NewGrpcMQ()

   if err != nil {
      return
   }
   ans, err := que.Register(&RegisterMeta{
      RoutingKeys: []routingKey{"MQBlock", "MQLog"},
      QueueName:   "test2",
      From:        "",
      Signature:   "",
      IsVerbose:   true,
      FromBlock:   "",
      ToBlock:     "",
      Addresses:   nil,
      Topics:      nil,
      Delay:       false,
   })
   if err != nil {
      return
   }
   // ...
}
```



### 6.5.3 解注册

```go
func (g *grpcMQ) UnRegister(meta *UnRegisterMeta) (*QueueUnRegister, StdError)
```

**解注册要保证队列确实存在，且没有客户端正在消费此队列**

应用实例

```go
func UnRegisterDemo() {
   var opt []rpc.GRPCOption
   opt = append(opt, rpc.BindNodes(0, 1))
   g := rpc.NewGRPC(opt...)

   que, err := g.NewGrpcMQ()
   if err != nil {
      return
   }
   ans, err := que.UnRegister(&UnRegisterMeta{
      QueueName: "test2",
      From:      "",
      Signature: "",
   })
   if err != nil {
      return
   }
   // ...
}
```

### 6.5.3 获取所有的队列名

```go
func (g *grpcMQ) GetAllQueueNames() ([]string, StdError)
```

| 返回值  |        |
| ---- | ------ |
| 返回值1 | 队列名的数组 |
| 返回值2 | err    |

应用实例

```go
func GetAllQueueNamesDemo() {
   var opt []rpc.GRPCOption
   opt = append(opt, rpc.BindNodes(0, 1))
   g := rpc.NewGRPC(opt...)

   que, err := g.NewGrpcMQ()
   if err != nil {
      return
   }
   ans, err := que.GetAllQueueNames()
   if err != nil {
      return
   }
   // ...
}
```

### 6.5.4 消费队列

```go
type ConsumeParams struct {
	QueueName string `json:"queueName"`
}

func (g *grpcMQ) Consume(meta *ConsumeParams) (api.GrpcApiMQ_ConsumeClient, StdError)
```

*****注意**

**1.消费队列的接口，我们仅提供流的返回，具体的recv逻辑开发者应当自己去处理**

**2.消费队列前要确保队列已经注册，且没有其他客户端在消费，如果有其他客户端在消费此队列可以主动调用stopcomsume来终止其他客户端的消费，因此开发者要保持清醒**

使用示例

```go
func Demo() {
   var opt []rpc.GRPCOption
   opt = append(opt, rpc.BindNodes(0, 1))
   g := rpc.NewGRPC(opt...)

   que, err := g.NewGrpcMQ()
   if err != nil {
	  log.Print(err)
      return
   }
   stream, err := que.Consume(&ConsumeParams{
      QueueName: "test2",
   })
   if err != nil {
	  log.Print(err)
      return
   }
   for {
      res, err := stream.Recv()
      if err != nil {
         break
      }
      // do your own
	  // 在这里添加自己的处理逻辑 
   }
}
```



### 6.5.5 停止消费

```go
type StopConsumeParams struct {
	QueueName string `json:"queueName"`
}

func (g *grpcMQ) StopConsume(meta *StopConsumeParams) (bool, StdError)
```

**调用停止消费接口，需要确保队列已经存在且有客户端正在消费队列，否则取消不成功**

应用实例

```go
func StopConsumeDemo() {
   var opt []rpc.GRPCOption
   opt = append(opt, rpc.BindNodes(0, 1))
   g := rpc.NewGRPC(opt...)

   que, err := g.NewGrpcMQ()
   if err != nil {
      return
   }
   ans, err := que.StopConsume(&StopConsumeParams{
      QueueName: "test2",
   })
   if err != nil {
	  return
   }
   // ...
}
```

