package plugin

import (
	"strings"

	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/meshplus/crypto"
)

//Set 遍历algoAll，如果通过测试则调用success
//e             : 代表go语言的plugin
//index         :
//l             : 插件的Level方法
//functionType  : 目标位置
//algoAll       : 目标算法
func (mux *externalMux) Set(e *externalPlugin, index int, l crypto.Level,
	functionType common.Function, algoAll []int,
	test func(k uint64, f crypto.Level, name string) error,
	success func(*externalMux, crypto.Level, uint64), fail func(*externalMux, uint64)) {
	defer func() {
		if err := recover(); err != nil {
			mux.logger.Errorf("load plugin %v index %v panic :%v", e.GetName(), index, err)
		}
	}()
	supports, _ := l.GetLevel()
	supportsMap := slice2map(supports)
	//遍历本重function的全部algo
	for _, algo := range algoAll {
		k := common.GetKey(functionType, algo)
		//info记录了当前k上已经加载的插件的信息，例如当前最高的level等
		info := mux.mapping[k]
		if info == nil {
			info = &functionInfo{
				software: true,
			}
		}

		invalidErr := test(k, l, e.GetName())
		if invalidErr != nil {
			if !strings.HasSuffix(invalidErr.Error(), crypto.ErrNotSupport.Error()) {
				mux.logger.Warningf("crypto mux: externalMux plugin self-test fail: %v", invalidErr.Error())
			}
		}
		if _, support := supportsMap[algo]; support && invalidErr == nil {
			mux.logger.Noticef("crypto mux: externalMux function %v for %v from %v is loading...",
				common.Function(k>>32).String(), common.GetModeName(common.GetModeFromKey(k)), e.GetName())
			success(mux, l, k)
			info.index, info.from, info.software = index, e, false
		} else {
			fail(mux, k)
			info.software = true
		}
		mux.mapping[k] = info
	}
}

func slice2map(in []int) map[int]int {
	ret := make(map[int]int, len(in))
	for i, v := range in {
		ret[v] = i
	}
	return ret
}
