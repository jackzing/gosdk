package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/go-hpc-common/utils"
	"github.com/hyperchain/go-hpc-msp/plugin"
	"github.com/hyperchain/go-hpc-msp/tls"
	"github.com/hyperchain/gosdk/common"
	"github.com/hyperchain/gosdk/grpc/pool"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type GRPC struct {
	tcm       *TCertManager
	config    *pool.Config
	namespace string
	im        *inspectorManager
	conn      *grpc.ClientConn
	gopts     grpcOption
}

type GRPCOption interface {
	apply(option *grpcOption)
}

// grpc options, append at here if need
type grpcOption struct {
	bindNode []int
}

type funcDialOption struct {
	f func(*grpcOption)
}

func (fdo *funcDialOption) apply(do *grpcOption) {
	fdo.f(do)
}

func newFuncDialOption(f func(*grpcOption)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

// BindNodes binding grpc nodes, begin at 0
func BindNodes(s ...int) GRPCOption {
	return newFuncDialOption(func(o *grpcOption) {
		o.bindNode = s
	})
}

type ClientOption struct {
	StreamNumber int
}

var (
	grpcLogger = common.GetLogger("grpc")
)

func NewGRPC(opt ...GRPCOption) *GRPC {
	return NewGRPCWithConfPath(common.DefaultConfRootPath, opt...)
}

func NewGRPCWithConfPath(path string, opts ...GRPCOption) *GRPC {
	cf := pool.NewConfigWithPath(path)
	tcm := NewTCertManager(cf.Viper(), path)
	// TODO change after support GetTxVersion
	rpc := NewRPCWithPath(path)
	gg := &GRPC{
		tcm:       tcm,
		config:    cf,
		namespace: cf.Namespace(),
		im:        rpc.im,
	}
	for _, opt := range opts {
		opt.apply(&gg.gopts)
	}
	conn, err := gg.newGrpcConn()
	if err != nil {
		panic(err)
	}
	gg.conn = conn
	return gg
}

func (g *GRPC) newGrpcConn() (*grpc.ClientConn, error) {
	var (
		conn *grpc.ClientConn
	)
	var opt []grpc.DialOption
	if g.config.IsTls() {
		clientcreds, err := tls.NewClientTLSFromFile(g.config.TlscaPath(), g.config.TlsDomain(), plugin.GetSoftwareEngine())
		if err != nil {
			grpcLogger.Errorf("NewClientTLSFromFile err %v", err)
			return nil, err
		}
		opt = append(opt, grpc.WithTransportCredentials(clientcreds))
	} else {
		opt = append(opt, grpc.WithInsecure())
	}
	opt = append(opt, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(g.config.MaxRecvMsgSize()), grpc.MaxCallSendMsgSize(g.config.MaxSendMsgSize())))

	ctx, cel := context.WithTimeout(context.Background(), g.config.DailTimeout())
	defer cel()
	var ips []string
	totalNodes := len(g.config.Targets())
	for i := 0; i < totalNodes; i++ {
		ips = append(ips, g.config.GetDailStringByIndex(i))
	}

	if len(g.gopts.bindNode) > 0 {
		ips = nil
		for _, v := range g.gopts.bindNode {
			if v < 0 || v >= totalNodes {
				grpcLogger.Errorf("error bind node. should [0, %d], get %d, check it", totalNodes-1, v)
				return nil, errors.New(fmt.Sprintf("error bind node. should [0, %d], get %d, check it", totalNodes-1, v))
			}
			ips = append(ips, g.config.GetDailStringByIndex(v))
		}
	}

	resolver.Register(&grpcResolverBuilder{addrs: ips})
	opt = append(opt, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", grpcScheme, grpcServiceName), opt...)
	if err != nil {
		grpcLogger.Errorf("DialContext err %v", err)
		return nil, err
	}
	return conn, err
}

func (g *GRPC) CheckClientOption(copt ClientOption) (bool, error) {
	if copt.StreamNumber <= 0 {
		return false, errors.New("num value is error, should > 0")
	}
	return true, nil
}

const (
	grpcScheme      = "grpc"
	grpcServiceName = "w.hyperchain.com"
)

type grpcResolverBuilder struct {
	addrs []string
}

func (e *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &grpcResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			grpcServiceName: e.addrs,
		},
	}
	r.start()
	return r, nil
}
func (*grpcResolverBuilder) Scheme() string { return grpcScheme }

type grpcResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *grpcResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*grpcResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*grpcResolver) Close()                                  {}

func convertTxToSendTxArgsProto(transaction *Transaction) *protos.SendTxArgs {
	fromBytes := []byte(transaction.from)
	toBytes := []byte(transaction.to)
	var pt *protos.ParticipantMsg
	if transaction.participant != nil {
		pt.Initiator = transaction.participant.Initiator
		for _, v := range transaction.participant.Withholding {
			pt.Withholding = append(pt.Withholding, v)
		}

	}

	return &protos.SendTxArgs{
		From:                string(fromBytes),
		To:                  string(toBytes),
		Value:               transaction.value,
		Payload:             transaction.payload,
		Signature:           transaction.signature,
		Timestamp:           transaction.timestamp,
		Simulate:            transaction.simulate,
		Nonce:               transaction.nonce,
		Extra:               transaction.extra,
		VmType:              transaction.vmType,
		Opcode:              int32(transaction.opcode),
		SnapshotID:          "",
		ExtraIDInt64Array:   transaction.extraIdInt64,
		ExtraIDStringArray:  transaction.extraIdString,
		CName:               transaction.cName,
		ExpirationTimestamp: transaction.expirationTimestamp,
		GasPrice:            transaction.gasPrice,
		GasLimit:            transaction.gasLimit,
		Participant:         pt,
	}
}

func convertReceiptResultProtoToTxReceipt(receipt *protos.ReceiptResult) *TxReceipt {
	return &TxReceipt{
		Version:         receipt.Version,
		TxHash:          receipt.TxHash,
		VMType:          receipt.VMType,
		ContractAddress: receipt.ContractAddress,
		ContractName:    receipt.ContractName,
		Ret:             receipt.Ret,
		Log:             convertLogTransProtoToTxLogs(receipt.Log),
		Valid:           receipt.Valid,
		ErrorMsg:        receipt.ErrorMsg,
		PrivTxHash:      receipt.TxHash,
	}
}

func convertLogTransProtoToTxLogs(logTrans []*protos.LogTrans) []TxLog {
	txLogs := make([]TxLog, len(logTrans))
	for _, logTran := range logTrans {
		txLog := TxLog{
			Address:     logTran.Address,
			Topics:      logTran.Topics,
			Data:        logTran.Data,
			BlockNumber: logTran.BlockNumber,
			TxHash:      logTran.TxHash,
			TxIndex:     logTran.TxIndex,
			Index:       logTran.Index,
		}
		txLogs = append(txLogs, txLog)
	}
	return txLogs
}

func (g *GRPC) prepareCommonReq(sendTxArgsProto *protos.SendTxArgs) (*protos.CommonReq, StdError) {
	grpcLogger.Debugf("[PrepareCommonReq] sendArgs %+v", sendTxArgsProto)
	commonReq := &protos.CommonReq{
		Namespace: g.namespace,
		Auth:      &protos.Auth{},
	}
	if sendTxArgsProto != nil {
		if sendTxArgsProto.Simulate {
			return nil, NewSystemError(errors.New("暂不支持simulate接口"))
		}
		params, err := proto.Marshal(sendTxArgsProto)
		if err != nil {
			grpcLogger.Errorf("marshal error %v", err)
			return nil, NewSystemError(err)
		}
		commonReq.Params = params
	}

	if g.im.enable {
		now := time.Now().UnixNano()
		var addr string
		if len(g.im.didAddress) != 0 {
			if utils.IsHex(g.im.didAddress) {
				addr = g.im.didAddress
			} else {
				addr = utils.ToHex([]byte(g.im.didAddress))
			}
		} else {
			addr = g.im.key.GetAddress().Hex()
		}
		sig, err := SignWithDID(g.im.key, authNeedHash(&Authentication{
			Address:   addr,
			Timestamp: now,
		}), false, false, len(g.im.didAddress) != 0)
		if err != nil {
			grpcLogger.Errorf("sign error %v", err)
			return nil, NewSystemError(err)
		}
		commonReq.Auth = &protos.Auth{
			Address:   addr,
			Timestamp: now,
			Signature: sig,
		}
	}
	marshed, err := proto.Marshal(commonReq)
	if err != nil {
		grpcLogger.Errorf("marshal error %v", err)
		return nil, NewSystemError(err)
	}
	if g.tcm != nil {
		sin, err := g.tcm.GetSDKCert().Sign(marshed)
		if err != nil {
			grpcLogger.Errorf("sign error %v", err)
			return nil, NewSystemError(err)
		}
		commonReq.TCert = g.tcm.GetECert()
		commonReq.Signature = common.Bytes2Hex(sin)
	}
	return commonReq, nil
}

func (g *GRPC) prepareMqCommReq(meta interface{}) (*protos.CommonReq, error) {
	marsh, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	commonReq := &protos.CommonReq{
		Namespace: g.namespace,
		Auth:      &protos.Auth{},
		Params:    marsh,
	}
	if g.im.enable {
		now := time.Now().UnixNano()
		var addr string
		if len(g.im.didAddress) != 0 {
			if utils.IsHex(g.im.didAddress) {
				addr = g.im.didAddress
			} else {
				addr = utils.ToHex([]byte(g.im.didAddress))
			}
		} else {
			addr = g.im.key.GetAddress().Hex()
		}
		sig, err := SignWithDID(g.im.key, authNeedHash(&Authentication{
			Address:   addr,
			Timestamp: now,
		}), false, false, len(g.im.didAddress) != 0)
		if err != nil {
			grpcLogger.Errorf("sign error %v", err)
			return nil, NewSystemError(err)
		}
		commonReq.Auth = &protos.Auth{
			Address:   addr,
			Timestamp: now,
			Signature: sig,
		}
	}
	marshed, err := proto.Marshal(commonReq)
	if err != nil {
		grpcLogger.Errorf("marshal error %v", err)
		return nil, NewSystemError(err)
	}
	if g.tcm != nil {
		sin, err := g.tcm.GetSDKCert().Sign(marshed)
		if err != nil {
			grpcLogger.Errorf("sign error %v", err)
			return nil, NewSystemError(err)
		}
		commonReq.TCert = g.tcm.GetECert()
		commonReq.Signature = common.Bytes2Hex(sin)
	}
	return commonReq, nil
}

func (g *GRPC) sendAndRecvReturnString(stream *pool.IdleStream, sendTxArgsProto *protos.SendTxArgs) (string, StdError) {
	if stream == nil {
		return "", NewSystemError(errors.New("system is busy"))
	}
	commonReq, err1 := g.prepareCommonReq(sendTxArgsProto)
	if err1 != nil {
		grpcLogger.Errorf("prepareCommonReq err %v", err1)
		return "", err1
	}
	grpcLogger.Debugf("[REQUEST] %+v", commonReq)
	err := stream.GetStream().Send(commonReq)
	if err != nil {
		grpcLogger.Errorf("Send err %v", err)
		return "", NewSystemError(err)
	}
	var ans *protos.CommonRes
	ans, err = stream.GetStream().Recv()
	if err != nil {
		grpcLogger.Errorf("Recv err %v", err)
		return "", NewSystemError(err)
	}
	grpcLogger.Debugf("[RESPONSE] %s", formatCommonRes(ans))
	if ans.Code != SuccessCode {
		grpcLogger.Errorf("response not success code: %d, codeDesc: %s", ans.Code, ans.CodeDesc)
		return "", NewServerError(int(ans.Code), ans.CodeDesc)
	}
	return common.BytesToHash(ans.Result).Hex(), nil
}

func (g *GRPC) sendAndRecv(stream *pool.IdleStream, sendTxArgsProto *protos.SendTxArgs) (*TxReceipt, StdError) {
	if stream == nil {
		return nil, NewSystemError(errors.New("system is busy"))
	}
	commonReq, err1 := g.prepareCommonReq(sendTxArgsProto)
	if err1 != nil {
		grpcLogger.Errorf("prepareCommonReq err %v", err1)
		return nil, err1
	}
	grpcLogger.Debugf("[REQUEST] %+v", commonReq)
	err := stream.GetStream().Send(commonReq)
	if err != nil {
		grpcLogger.Errorf("Send err %v", err)
		return nil, NewSystemError(err)
	}

	var ans *protos.CommonRes
	ans, err = stream.GetStream().Recv()
	if err != nil {
		grpcLogger.Errorf("Recv err %v", err)
		return nil, NewSystemError(err)
	}
	var ret = new(protos.ReceiptResult)
	err = proto.Unmarshal(ans.Result, ret)
	if err != nil {
		grpcLogger.Errorf("Unmarshal err %v", err)
		return nil, NewSystemError(err)
	}
	grpcLogger.Debugf("[RESPONSE] %s", formatReceiptCommonRes(ans))
	if ans.Code != SuccessCode {
		grpcLogger.Errorf("response not success code: %d, codeDesc: %s", ans.Code, ans.CodeDesc)
		return nil, NewServerError(int(ans.Code), ans.CodeDesc)
	}
	return convertReceiptResultProtoToTxReceipt(ret), nil
}

func (g *GRPC) Close() error {
	return g.conn.Close()
}

func formatReceiptCommonRes(ans *protos.CommonRes) string {
	var ret = new(protos.ReceiptResult)
	proto.Unmarshal(ans.Result, ret)
	return fmt.Sprintf(`{"namespace": %s, "code": %d, "code_desc": %s, "result": %s}`, ans.Namespace, ans.Code, ans.CodeDesc, ret)
}

func formatCommonRes(ans *protos.CommonRes) string {
	return fmt.Sprintf(`{"namespace": %s, "code": %d, "code_desc": %s, "result": %s}`, ans.Namespace, ans.Code, ans.CodeDesc, common.BytesToHash(ans.Result).Hex())
}
