package trade

import (
	"context"
	"fmt"
	"github.com/rz1998/invest-basic/types/investBasic"
	"github.com/rz1998/invest-trade-basic/types/tradeBasic"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"sync/atomic"
	"time"
)

type SpiMDPrint struct {
}

// OnRtnMD 行情回报
func (spi *SpiMDPrint) OnRtnMD(tick *investBasic.SMDTick) {
	logx.Infof("OnRtnMD %+v", *tick)
}

// OnFrontDisconnect 前置断开
func (spi *SpiMDPrint) OnFrontDisconnect(reason int) {
	logx.Infof("OnFrontDisconnect %d", reason)
}

// OnRspSub 订阅行情返回数据
func (spi *SpiMDPrint) OnRspSub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspSub %s %+v %d %t", code, msgError, idRequest, isLast)
}

// OnRspUnsub 取消订阅行情返回数据
func (spi *SpiMDPrint) OnRspUnsub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspUnsub %s %+v %d %t", code, msgError, idRequest, isLast)
}

type SpiMDCache struct {
	// uniqueCode : SMDTick
	MapMD *sync.Map
}

// OnRtnMD 行情回报
func (spi *SpiMDCache) OnRtnMD(tick *investBasic.SMDTick) {
	if tick == nil {
		return
	}
	spi.MapMD.Store(tick.UniqueCode, tick)
}

// OnFrontDisconnect 前置断开
func (spi *SpiMDCache) OnFrontDisconnect(reason int) {
	logx.Infof("OnFrontDisconnect %d", reason)
}

// OnRspSub 订阅行情返回数据
func (spi *SpiMDCache) OnRspSub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspSub %s %+v %d %t", code, msgError, idRequest, isLast)
}

// OnRspUnsub 取消订阅行情返回数据
func (spi *SpiMDCache) OnRspUnsub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspUnsub %s %+v %d %t", code, msgError, idRequest, isLast)
}

// SpiMDMinute 行情整合为分钟行情
type SpiMDMinute struct {
	// uniqueCode : SMDTick
	MapMD       *sync.Map
	MapMDMinute *sync.Map
}

// OnRtnMD 行情回报
func (spi *SpiMDMinute) OnRtnMD(md *investBasic.SMDTick) {
	if md == nil {
		return
	}
	mdLastAny, ok := spi.MapMD.Load(md.UniqueCode)
	if ok {
		mdLast := mdLastAny.(*investBasic.SMDTick)
		// 时间只取最新（qos<2可能会时间错乱）
		if md.Timestamp <= mdLast.Timestamp {
			return
		}
		timeLast := time.UnixMilli(mdLast.Timestamp)
		timeThis := time.UnixMilli(md.Timestamp)
		if timeLast.Hour() != timeThis.Hour() ||
			timeLast.Minute() != timeThis.Minute() {
			// 最新行情的小时或者分钟不等于上一行情（分钟改变），更新分钟行情
			var mds []*investBasic.SMDTick
			mdsAny, ok := spi.MapMDMinute.Load(md.UniqueCode)
			if ok {
				mds = mdsAny.([]*investBasic.SMDTick)
			}
			mdLast.Timestamp = mdLast.Timestamp / 60000 * 60000
			mds = append(mds, mdLast)
			spi.MapMDMinute.Store(md.UniqueCode, mds)
		}
	}
	spi.MapMD.Store(md.UniqueCode, md)
}

// OnFrontDisconnect 前置断开
func (spi *SpiMDMinute) OnFrontDisconnect(reason int) {
	logx.Infof("OnFrontDisconnect %d", reason)
}

// OnRspSub 订阅行情返回数据
func (spi *SpiMDMinute) OnRspSub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspSub %s %+v %d %t", code, msgError, idRequest, isLast)
}

// OnRspUnsub 取消订阅行情返回数据
func (spi *SpiMDMinute) OnRspUnsub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspUnsub %s %+v %d %t", code, msgError, idRequest, isLast)
}

type SpiMDOuter struct {
	FuncOnRtnMD func(tick *investBasic.SMDTick)
}

func (spi *SpiMDOuter) OnRtnMD(tick *investBasic.SMDTick) {
	spi.FuncOnRtnMD(tick)
}

func (spi *SpiMDOuter) OnFrontDisconnect(reason int) {
	logx.Infof("OnFrontDisconnect %d", reason)
}

func (spi *SpiMDOuter) OnRspSub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspSub %s %+v %d %t", code, msgError, idRequest, isLast)
}

func (spi *SpiMDOuter) OnRspUnsub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspUnsub %s %+v %d %t", code, msgError, idRequest, isLast)
}

type SpiTraderPrint struct {
}

// OnRtnOrder 报单回报
func (spi *SpiTraderPrint) OnRtnOrder(infoOrder *tradeBasic.SOrderInfo, isLast bool) {
	if infoOrder == nil {
		return
	}
	logx.Infof("%s %+v %t", "OnRtnOrder", *infoOrder, isLast)
}

// OnRtnTrade 成交回报
func (spi *SpiTraderPrint) OnRtnTrade(infoTrader *tradeBasic.STradeInfo, isLast bool) {
	if infoTrader == nil {
		return
	}
	logx.Infof("%s %+v %t", "OnRtnTrade", *infoTrader, isLast)
}

// OnRtnAcPos 账户持仓查询
func (spi *SpiTraderPrint) OnRtnAcPos(acPos *tradeBasic.SAcPos, isLast bool) {
	if acPos == nil {
		return
	}
	logx.Infof("%s %+v %t", "OnRtnAcPos", *acPos, isLast)

}

// OnRtnAcFund 账户资金查询
func (spi *SpiTraderPrint) OnRtnAcFund(acFund *tradeBasic.SAcFund) {
	if acFund == nil {
		return
	}
	logx.Infof("%s %+v", "OnRtnAcFund", *acFund)
}

// OnFrontDisconnect 前置断开回报
func (spi *SpiTraderPrint) OnFrontDisconnect(reason int) {
	logx.Infof("%s %d", "OnFrontDisconnect", reason)
}

// OnErrRtnOrderInsert 报单错误回报
func (spi *SpiTraderPrint) OnErrRtnOrderInsert(reqOrder *tradeBasic.PReqOrder, msgError *tradeBasic.SMsgError) {
	if reqOrder == nil || msgError == nil {
		return
	}
	logx.Infof("%s %+v %+v", "OnErrRtnOrderInsert", *reqOrder, *msgError)

}

// OnErrRtnOrderAction 报单操作错误回报
func (spi *SpiTraderPrint) OnErrRtnOrderAction(reqOrderAction *tradeBasic.PReqOrderAction, msgError *tradeBasic.SMsgError) {
	if reqOrderAction == nil || msgError == nil {
		return
	}
	logx.Infof("%s %+v %+v", "OnErrRtnOrderAction", *reqOrderAction, *msgError)
}

type SpiTraderBenchmark struct {
}

// OnRtnOrder 报单回报
func (spi *SpiTraderBenchmark) OnRtnOrder(infoOrder *tradeBasic.SOrderInfo, isLast bool) {
	if infoOrder == nil {
		return
	}
	logx.Infof("%s %+v %t", "OnRtnOrder", *infoOrder, isLast)
}

// OnRtnTrade 成交回报
func (spi *SpiTraderBenchmark) OnRtnTrade(infoTrader *tradeBasic.STradeInfo, isLast bool) {
	if infoTrader == nil {
		return
	}
	logx.Infof("%s %+v %t", "OnRtnTrade", *infoTrader, isLast)
}

// OnRtnAcPos 账户持仓查询
func (spi *SpiTraderBenchmark) OnRtnAcPos(acPos *tradeBasic.SAcPos, isLast bool) {
	if acPos == nil {
		return
	}
	logx.Infof("%s %+v %t", "OnRtnAcPos", *acPos, isLast)

}

// OnRtnAcFund 账户资金查询
func (spi *SpiTraderBenchmark) OnRtnAcFund(acFund *tradeBasic.SAcFund) {
	if acFund == nil {
		return
	}
	logx.Infof("%s %+v", "OnRtnAcFund", *acFund)
}

// OnFrontDisconnect 前置断开回报
func (spi *SpiTraderBenchmark) OnFrontDisconnect(reason int) {
	logx.Infof("%s %d", "OnFrontDisconnect", reason)
}

// OnErrRtnOrderInsert 报单错误回报
func (spi *SpiTraderBenchmark) OnErrRtnOrderInsert(reqOrder *tradeBasic.PReqOrder, msgError *tradeBasic.SMsgError) {
	if reqOrder == nil || msgError == nil {
		return
	}
	logx.Infof("%s %+v %+v", "OnErrRtnOrderInsert", *reqOrder, *msgError)

}

// OnErrRtnOrderAction 报单操作错误回报
func (spi *SpiTraderBenchmark) OnErrRtnOrderAction(reqOrderAction *tradeBasic.PReqOrderAction, msgError *tradeBasic.SMsgError) {
	if reqOrderAction == nil || msgError == nil {
		return
	}
	logx.Infof("%s %+v %+v", "OnErrRtnOrderAction", *reqOrderAction, *msgError)
}

type SpiMDBenchmark struct {
	timestamp int64
}

// OnRtnMD 行情回报
func (spi *SpiMDBenchmark) OnRtnMD(tick *investBasic.SMDTick) {
	timeLast := atomic.LoadInt64(&spi.timestamp)
	timeThis := time.Now().UnixMilli()
	fmt.Println(float64(timeThis-timeLast) / 1000.0)
	atomic.StoreInt64(&spi.timestamp, timeThis)
}

// OnFrontDisconnect 前置断开
func (spi *SpiMDBenchmark) OnFrontDisconnect(reason int) {
	logx.Infof("OnFrontDisconnect %d", reason)
}

// OnRspSub 订阅行情返回数据
func (spi *SpiMDBenchmark) OnRspSub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspSub %s %+v %d %t", code, msgError, idRequest, isLast)
}

// OnRspUnsub 取消订阅行情返回数据
func (spi *SpiMDBenchmark) OnRspUnsub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool) {
	logx.Infof("OnRspUnsub %s %+v %d %t", code, msgError, idRequest, isLast)
}

// SpiTraderCache 缓存所有信息
type SpiTraderCache struct {
	GenerateUniqueOrder func(orderSys *tradeBasic.SOrderSys) string
	AcFund              *tradeBasic.SAcFund
	// ucOrder : SOrderInfo
	MapOrder *sync.Map
	// uniqueCode : EDirTrade : SAcPos
	MapPos    *sync.Map
	CancelPos context.CancelFunc
}

func (spi *SpiTraderCache) GetOrderInfo(ucOrder string) (*tradeBasic.SOrderInfo, bool) {
	orderAny, ok := spi.MapOrder.Load(ucOrder)
	if ok {
		return orderAny.(*tradeBasic.SOrderInfo), ok
	}
	return nil, ok
}

func (spi *SpiTraderCache) GetAcPos(uniqueCode string, dir investBasic.EDirTrade) (*tradeBasic.SAcPos, bool) {
	mapAny, ok := spi.MapPos.Load(uniqueCode)
	if ok {
		mapDir := mapAny.(*sync.Map)
		posAny, hasPos := mapDir.Load(dir)
		if hasPos {
			return posAny.(*tradeBasic.SAcPos), true
		}
	}
	return nil, false
}

// OnRtnOrder 报单回报
func (spi *SpiTraderCache) OnRtnOrder(infoOrder *tradeBasic.SOrderInfo, isLast bool) {
	if infoOrder == nil || infoOrder.OrderSys == nil {
		return
	}
	spi.MapOrder.Store(spi.GenerateUniqueOrder(infoOrder.OrderSys), infoOrder)
}

// OnRtnTrade 成交回报
func (spi *SpiTraderCache) OnRtnTrade(infoTrader *tradeBasic.STradeInfo, isLast bool) {

}

// OnRtnAcPos 账户持仓查询
func (spi *SpiTraderCache) OnRtnAcPos(acPos *tradeBasic.SAcPos, isLast bool) {
	if isLast {
		if spi.CancelPos != nil {
			spi.CancelPos()
		}
	}
	if acPos == nil {
		return
	}
	var mapDir *sync.Map
	mapAny, ok := spi.MapPos.Load(acPos.UniqueCode)
	if ok {
		mapDir = mapAny.(*sync.Map)
	} else {
		mapDir = &sync.Map{}
		spi.MapPos.Store(acPos.UniqueCode, mapDir)
	}
	mapDir.Store(acPos.TradeDir, acPos)
}

// OnRtnAcFund 账户资金查询
func (spi *SpiTraderCache) OnRtnAcFund(acFund *tradeBasic.SAcFund) {
	spi.AcFund = acFund
}

// OnFrontDisconnect 前置断开回报
func (spi *SpiTraderCache) OnFrontDisconnect(reason int) {

}

// OnErrRtnOrderInsert 报单错误回报
func (spi *SpiTraderCache) OnErrRtnOrderInsert(reqOrder *tradeBasic.PReqOrder, msgError *tradeBasic.SMsgError) {

}

// OnErrRtnOrderAction 报单操作错误回报
func (spi *SpiTraderCache) OnErrRtnOrderAction(reqOrderAction *tradeBasic.PReqOrderAction, msgError *tradeBasic.SMsgError) {

}

// SpiTraderOuter 用于暴露接口
type SpiTraderOuter struct {
	// OnRtnOrder 报单回报
	FuncOnRtnOrder func(infoOrder *tradeBasic.SOrderInfo, isLast bool)
	// OnRtnTrade 成交回报
	FuncOnRtnTrade func(infoTrader *tradeBasic.STradeInfo, isLast bool)
	// OnRtnAcPos 账户持仓查询
	FuncOnRtnAcPos func(acPos *tradeBasic.SAcPos, isLast bool)
	// OnRtnAcFund 账户资金查询
	FuncOnRtnAcFund func(acFund *tradeBasic.SAcFund)
	// OnFrontDisconnect 前置断开回报
	FuncOnFrontDisconnect func(reason int)
	// OnErrRtnOrderInsert 报单错误回报
	FuncOnErrRtnOrderInsert func(reqOrder *tradeBasic.PReqOrder, msgError *tradeBasic.SMsgError)
	// OnErrRtnOrderAction 报单操作错误回报
	FuncOnErrRtnOrderAction func(reqOrderAction *tradeBasic.PReqOrderAction, msgError *tradeBasic.SMsgError)
}

// OnRtnOrder 报单回报
func (spi *SpiTraderOuter) OnRtnOrder(infoOrder *tradeBasic.SOrderInfo, isLast bool) {
	spi.FuncOnRtnOrder(infoOrder, isLast)
}

// OnRtnTrade 成交回报
func (spi *SpiTraderOuter) OnRtnTrade(infoTrader *tradeBasic.STradeInfo, isLast bool) {
	spi.FuncOnRtnTrade(infoTrader, isLast)
}

// OnRtnAcPos 账户持仓查询
func (spi *SpiTraderOuter) OnRtnAcPos(acPos *tradeBasic.SAcPos, isLast bool) {
	spi.FuncOnRtnAcPos(acPos, isLast)
}

// OnRtnAcFund 账户资金查询
func (spi *SpiTraderOuter) OnRtnAcFund(acFund *tradeBasic.SAcFund) {
	spi.FuncOnRtnAcFund(acFund)
}

// OnFrontDisconnect 前置断开回报
func (spi *SpiTraderOuter) OnFrontDisconnect(reason int) {
	spi.FuncOnFrontDisconnect(reason)
}

// OnErrRtnOrderInsert 报单错误回报
func (spi *SpiTraderOuter) OnErrRtnOrderInsert(reqOrder *tradeBasic.PReqOrder, msgError *tradeBasic.SMsgError) {
	spi.FuncOnErrRtnOrderInsert(reqOrder, msgError)
}

// OnErrRtnOrderAction 报单操作错误回报
func (spi *SpiTraderOuter) OnErrRtnOrderAction(reqOrderAction *tradeBasic.PReqOrderAction, msgError *tradeBasic.SMsgError) {
	spi.FuncOnErrRtnOrderAction(reqOrderAction, msgError)
}
