package rpc

import (
	"context"

	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/jackzing/gosdk/grpc/pool"
)

type VersionGrpc struct {
	client                               protos.GrpcApiVersionClient
	setSupportedVersionPool              *pool.StreamPool
	setSupportedVersionReturnReceiptPool *pool.StreamPool
	num                                  int
	grpc                                 *GRPC
}

func (g *GRPC) NewVersionGrpc(opt ClientOption) (*VersionGrpc, error) {
	_, err := g.CheckClientOption(opt)
	if err != nil {
		return nil, err
	}
	client := protos.NewGrpcApiVersionClient(g.conn)
	return &VersionGrpc{
		client: client,
		num:    opt.StreamNumber,
		grpc:   g,
	}, nil
}

func (v *VersionGrpc) SetSupportedVersion() (string, StdError) {
	p, err := v.getSetSupportedVersionPool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)

	return v.grpc.sendAndRecvReturnString(stream, nil)
}

func (v *VersionGrpc) getSetSupportedVersionPool() (*pool.StreamPool, error) {
	if v.setSupportedVersionPool == nil {
		k, err := pool.NewStreamWithContext(v.grpc.config.MaxStreamLifetime(), v.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := v.client.SetSupportedVersion(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		v.setSupportedVersionPool = k
	}
	return v.setSupportedVersionPool, nil
}

func (v *VersionGrpc) SetSupportedVersionReturnReceipt() (*TxReceipt, StdError) {
	p, err := v.getSetSupportedVersionReturnReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)

	return v.grpc.sendAndRecv(stream, nil)
}

func (v *VersionGrpc) getSetSupportedVersionReturnReceiptPool() (*pool.StreamPool, error) {
	if v.setSupportedVersionReturnReceiptPool == nil {
		k, err := pool.NewStreamWithContext(v.grpc.config.MaxStreamLifetime(), v.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := v.client.SetSupportedVersionReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		v.setSupportedVersionReturnReceiptPool = k
	}
	return v.setSupportedVersionReturnReceiptPool, nil
}

func (v *VersionGrpc) Close() error {
	if v.setSupportedVersionPool != nil {
		err := v.setSupportedVersionPool.Close()
		if err != nil {
			return err
		}
	}
	if v.setSupportedVersionReturnReceiptPool != nil {
		err := v.setSupportedVersionReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
