package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const Name = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {

		md, _ := sci.Address.Metadata.(map[string]any)
		weightVal, _ := md["weight"]
		weight := weightVal.(float64)

		conns = append(conns, &weightConn{
			SubConn:       sc,
			weight:        int(weight),
			currentWeight: int(weight),
		})
	}

	return &Picker{}
}

type Picker struct {
	conns []*weightConn // 把服务节点存下来
	lock  sync.Mutex    //该接口会被并发调用
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {

	// 1. 计算出所有节点的总权重totalWeight，累加所有originWeight
	// 2. 计算当前服务节点的currentWeight+originWeight
	// 3. 挑选currentWeight最大的节点
	// 4. 使用currentWeight-totalWeight
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.conns) <= 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var total int
	var maxCC *weightConn

	for _, c := range p.conns {
		total += c.weight
		c.currentWeight += c.weight
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
	}

	//注意: 并不是请求响应回调回来才去调整，而是发出前就去调整，默认假设请求一直能成功发送
	maxCC.currentWeight -= total

	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		// 回调, 是否请求成功
		Done: func(info balancer.DoneInfo) {

		},
	}, nil

}

type weightConn struct {
	balancer.SubConn
	weight        int
	currentWeight int
}
