package tradeBasic

import (
	"math"
	"sort"
	"sync"
	"time"

	investBasic "github.com/rz1998/invest-basic"
	"github.com/rz1998/invest-stra-signal/types/signalStra"
	"github.com/zeromicro/go-zero/core/logx"
)

// 根据最大单笔委托数量上限，等量拆分报单
func SepOrderBySameVol(maxVol int64, reqOrder *PReqOrder) []*PReqOrder {
	if reqOrder == nil {
		return nil
	}
	var reqOrders []*PReqOrder
	vol := reqOrder.Vol
	for vol > maxVol {
		reqOrders = append(reqOrders, &PReqOrder{
			Timestamp:  reqOrder.Timestamp,
			Dir:        reqOrder.Dir,
			FlagOffset: reqOrder.FlagOffset,
			UniqueCode: reqOrder.UniqueCode,
			Price:      reqOrder.Price,
			Vol:        maxVol,
		})
		vol -= maxVol
	}
	if vol > 0 {
		reqOrders = append(reqOrders, &PReqOrder{
			Timestamp:  reqOrder.Timestamp,
			Dir:        reqOrder.Dir,
			FlagOffset: reqOrder.FlagOffset,
			UniqueCode: reqOrder.UniqueCode,
			Price:      reqOrder.Price,
			Vol:        vol,
		})
	}
	return reqOrders
}

// CancelOrderBatch 根据报单方向批量撤单
/*
 * dir 需要撤单的报单方向
 * flagOffset 开平标记
 * mapOrderInfo ucOrder : infoOrder
 */
func CancelOrderBatch(dir investBasic.EDirTrade, flagOffset EFlagOffset, mapOrderInfo *sync.Map) []*PReqOrderAction {
	strMethod := "CancelOrderBatch"
	var rtns []*PReqOrderAction
	if mapOrderInfo == nil {
		logx.Infof("%s stopped by %s\n", strMethod, "no mapOrderInfo")
		return rtns
	}
	mapOrderInfo.Range(func(key, value any) bool {
		orderInfo := value.(*SOrderInfo)
		// 筛选目标报单
		if orderInfo.ReqOrder != nil {
			if dir != investBasic.EMPTY &&
				dir != orderInfo.ReqOrder.Dir {
				return true
			}
			if flagOffset != EFlagOffset_NONE &&
				flagOffset != orderInfo.ReqOrder.FlagOffset {
				return true
			}
		}
		// 生成撤单信息
		rtns = append(rtns, &PReqOrderAction{
			OrderSys:    orderInfo.OrderSys,
			OrderAction: Delete,
		})
		return true
	})
	return rtns
}

// 尝试给信号分配资金，计算最终未利用金额
func TestSignalBuy(pFees PFees, signals []*signalStra.SSignalStra, valAvailable, valBalance float64) (reqOrders []*PReqOrder, valBalanceLeft float64) {
	valBalanceLeft = valBalance
	logx.Infof("TestSignalBuy valAvailable %f, valBalance %f", valAvailable, valBalance)
	if len(signals) == 0 {
		return reqOrders, valBalanceLeft
	}
	if valAvailable <= 0 || valBalance <= 0 {
		return reqOrders, valBalanceLeft
	}
	size := len(signals)
	rateLeft := 1.0
	timestamp := time.Now().UnixMilli()
	for i, signal := range signals {
		logx.Debugf("TestSignalBuy signal %+v\n", *signal)
		if signal.TradeDir != investBasic.LONG {
			continue
		}
		var vol int64
		var rateVal float64
		var valTar float64
		if signal.Vol > 0 {
			// 信号里有数量的，以信号为准
			vol = signal.Vol
		} else {
			// 信号里没有数量的，以priority设置的仓位系数为准
			// 仓位系数
			rateVal = signal.Priority
			if signal.Priority <= 0 {
				rateVal = 1.0 / float64(size)
			}
			// 可用资金
			valTar = valBalanceLeft / rateLeft * rateVal
			if i == size-1 {
				// 最后一个向下取整
				vol = CalTarVolByVal(pFees, signal.UniqueCode, valTar, signal.Price, Floor)
			} else {
				// 其他在剔除费用后四舍五入
				valTar = CalValLeft(pFees, valTar)
				vol = CalTarVolByVal(pFees, signal.UniqueCode, valTar, signal.Price, Round)
			}
		}
		valTar = float64(signal.Price) / 10000.0 * float64(vol)
		if vol > 0 {
			// 生成下单文件
			reqOrder := &PReqOrder{
				Timestamp:  timestamp,
				Dir:        investBasic.LONG,
				FlagOffset: Open,
				UniqueCode: signal.UniqueCode,
				Price:      signal.Price,
				Vol:        vol,
			}
			// 根据上限拆分
			reqOrders = append(reqOrders, SepOrderBySameVol(1000000, reqOrder)...)
			// 更新缓存
			rateLeft -= rateVal
			valBalanceLeft -= valTar
		}
	}
	logx.Infof("TestSignalBuy result: valBalanceLeft %f,valAvailable %f\n", valBalanceLeft, valAvailable)
	if valAvailable > 0 && valBalance-valBalanceLeft > valAvailable {
		// 如果理论使用的金额高于总可用，用可用重新算一遍
		return TestSignalBuy(pFees, signals, valAvailable, valBalanceLeft+valAvailable)
	}
	return reqOrders, valBalanceLeft
}

// AllocateSignalBuy 根据股价顺序倒叙分配资金，选择资金利用率最高的进行下单
func AllocateSignalBuy(pFees PFees, signalBuys []*signalStra.SSignalStra, valAvailable, valBalance float64) []*PReqOrder {
	// 顺序排列价格
	sort.Slice(signalBuys, func(i, j int) bool {
		return signalBuys[i].Price < signalBuys[j].Price
	})
	reqsAsc, leftAsc := TestSignalBuy(pFees, signalBuys, valAvailable, valBalance)
	// 倒序排列价格
	sort.Slice(signalBuys, func(i, j int) bool {
		return signalBuys[i].Price > signalBuys[j].Price
	})
	reqsDesc, leftDesc := TestSignalBuy(pFees, signalBuys, valAvailable, valBalance)
	if leftAsc > leftDesc {
		return reqsDesc
	} else {
		return reqsAsc
	}
}

// AllocateSignalSell 根据信号和持仓，生成卖出报单
func AllocateSignalSell(mapPos *sync.Map, signalSells []*signalStra.SSignalStra) []*PReqOrder {
	strMethod := "AllocateSignalSell"
	var reqOrders []*PReqOrder
	if len(signalSells) == 0 {
		return reqOrders
	}
	timestamp := time.Now().UnixMilli()
	for _, signal := range signalSells {
		if signal.TradeDir != investBasic.SHORT {
			continue
		}
		if acPos, hasPos := GetAcPos(mapPos, signal.UniqueCode, investBasic.LONG); hasPos && GetVolSellableFromAcPos(acPos) > 0 {
			// 有持仓的情况下，根据信号计算卖出数量
			volSellable := GetVolSellableFromAcPos(acPos)
			var vol int64
			if signal.Vol > 0 {
				// 信号中有数量时，以信号中的数量为准
				if signal.Vol > volSellable {
					vol = volSellable
				} else {
					vol = signal.Vol
				}
			} else {
				if signal.Priority <= 0 {
					// 没有设置仓位则全清
					vol = volSellable
				} else {
					// 有设置则根据设置的比例进行四舍五入
					vol = int64(math.Min(float64(CalTarVol(signal.UniqueCode, signal.Priority, acPos.VolTotal)), float64(volSellable)))
				}
			}
			// 生成报单
			reqOrder := &PReqOrder{
				Timestamp:  timestamp,
				Dir:        investBasic.SHORT,
				FlagOffset: Close,
				UniqueCode: signal.UniqueCode,
				Price:      signal.Price,
				Vol:        vol,
			}
			// 根据上限拆分
			if reqOrder.Price > 0 {
				reqOrders = append(reqOrders, SepOrderBySameVol(1000000, reqOrder)...)
			}
		} else {
			logx.Infof("%s stopped by no acPos %s", strMethod, signal.UniqueCode)
		}
	}
	return reqOrders
}
