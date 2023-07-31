package trade

import (
	"fmt"
	"github.com/rz1998/invest-basic/types/investBasic"
	"github.com/rz1998/invest-trade-basic/types/tradeBasic"
	"math"
	"strings"
	"sync"
)

// GetSecInfo 唯一码转换为代码和交易所
func GetSecInfo(uniqueCode string) (code string, exchangeCD investBasic.ExchangeCD) {
	contents := strings.Split(uniqueCode, ".")
	if len(contents) > 1 {
		return contents[0], investBasic.ExchangeCD(contents[1])
	} else {
		return contents[0], ""
	}
}

// CalLimitPrice 计算涨跌停价(股票四舍五入)
func CalLimitPrice(rateStop float64, preClose, tick int64) int64 {
	if tick <= 0 {
		return 0
	}
	return int64(math.Round(float64(preClose)/float64(tick)*(1.0+rateStop))) * tick
}

// CalTarPrice 根据证券代码计算目标涨跌幅度价格 priceLimit 是涨跌停价，rateTar>0用涨停价，rateTar<0用跌停价
func CalTarPrice(uniqueCode string, rateTar float64, price, priceLimit int64) int64 {
	code, exchangeCD := GetSecInfo(uniqueCode)
	var tick int64
	switch exchangeCD {
	case investBasic.SSE, investBasic.SZSE:
		// 沪深市场
		str1 := code[:1]
		str2 := code[:2]
		if str1 == "0" || str1 == "3" || str1 == "6" {
			// 股票，100一跳
			tick = 100
		} else if exchangeCD == investBasic.SSE && str2 == "11" {
			// 上交所可转债，100一跳
			tick = 100
		} else {
			tick = 10
		}
	default:
		fmt.Printf("CalTarPrice unhandled exchangeCD %v\n", exchangeCD)
	}
	if priceLimit < 0 {
		// 忽略限制
		return CalLimitPrice(rateTar, price, tick)
	} else {
		if rateTar > 0 {
			if price > priceLimit {
				fmt.Printf("CalTarPrice error, price over limitUpper! rataTar %d price %d priceLimit %d\n", rateTar, price, priceLimit)
				return price
			} else {
				return int64(math.Min(float64(priceLimit), float64(CalLimitPrice(rateTar, price, tick))))
			}
		} else {
			if price < priceLimit {
				fmt.Printf("CalTarPrice error, price under limitLower! rataTar %d price %d priceLimit %d\n", rateTar, price, priceLimit)
				return price
			} else {
				return int64(math.Max(float64(priceLimit), float64(CalLimitPrice(rateTar, price, tick))))
			}
		}
	}
}

// CalLimitVol 根据目标比例计算目标报单数量
func CalLimitVol(rateTar float64, vol, volTick int64) int64 {
	if volTick <= 0 {
		return 0
	}
	return int64(math.Round(float64(vol)/float64(volTick)*(rateTar))) * volTick
}

// CalValLeastByFee 根据最低佣金额及佣金率，计算单次下单最小额度
/*
 * rateFeeBrokerage 佣金率
 * floorFeeBrokerage  最低绝对佣金额
 */
func CalValLeastByFee(rateFeeBrokerage, floorFeeBrokerage float64) float64 {
	return floorFeeBrokerage / rateFeeBrokerage
}

// 根据代码返回最小下单单位
func getVolTick(uniqueCode string) int64 {
	code, exchangeCD := GetSecInfo(uniqueCode)
	// 默认1
	var tick int64 = 1
	switch exchangeCD {
	case investBasic.SSE, investBasic.SZSE:
		// 沪深市场
		str2 := code[:2]
		str3 := code[:3]
		if str3 == "131" || str2 == "11" || str2 == "12" {
			// 深市逆回购、所有转债，10
			tick = 10
		} else {
			// 其他100
			tick = 100
		}
	default:
		fmt.Printf("getVolTick unhandled exchangeCD %v\n", exchangeCD)
	}
	return tick
}

// CalTarVol 根据目标比例计算目标报单数量（根据代码自动判断最小下单单位）
func CalTarVol(uniqueCode string, rateTar float64, vol int64) int64 {
	return CalLimitVol(rateTar, vol, getVolTick(uniqueCode))
}

// CalTarVolByVal 根据目标金额，确定目标数量
func CalTarVolByVal(uniqueCode string, val float64, price int64) int64 {
	//
	volTheory := val / float64(price) * 10000
	// 向下取证到目标报单数量上
	volTick := getVolTick(uniqueCode)
	return int64(math.Floor(volTheory/float64(volTick))) * volTick
}

// GetVolSellableFromAcPos 根据持仓计算可卖数量
func GetVolSellableFromAcPos(acPos *tradeBasic.SAcPos) (vol int64) {
	if acPos == nil {
		return vol
	}
	_, exchangeCD := GetSecInfo(acPos.UniqueCode)
	switch exchangeCD {
	case investBasic.SSE, investBasic.SZSE:
		// 股票昨仓
		vol = acPos.VolYd - acPos.VolFrozenYd
	default:
		// 其他全部
		vol = acPos.VolTotal - acPos.VolFrozenTotal
	}
	return vol
}

// 根据订单的开平标致和买卖方向，选择受影响的持仓方向
/*
 * @param dirTrade 目标交易方向
 * @param flagOffset 开平标致
 * @return 所需持仓的方向
 */
func chooseDirPos(dirTrade investBasic.EDirTrade, flagOffset tradeBasic.EFlagOffset) investBasic.EDirTrade {
	var dir investBasic.EDirTrade
	switch flagOffset {
	case tradeBasic.Open:
		dir = dirTrade
	case tradeBasic.Close, tradeBasic.CloseToday, tradeBasic.CloseYesterday,
		tradeBasic.ForceClose, tradeBasic.LocalForceClose, tradeBasic.ForceOff:
		if dirTrade == investBasic.LONG {
			dir = investBasic.SHORT
		} else {
			dir = investBasic.LONG
		}
	}
	return dir
}

// UpdateAcPosByOrderInfo 用报单回报更新持仓
/*
 * mapAcPos dirTrade : uniqueCode : acPos
 */
func UpdateAcPosByOrderInfo(mapAcPos *sync.Map, infoOrder *tradeBasic.SOrderInfo) *tradeBasic.SAcPos {
	if mapAcPos == nil {
		fmt.Printf("%s stopped by %s\n", "UpdateAcPosByOrderInfo", "no map acPos")
		return nil
	}
	if infoOrder == nil || infoOrder.ReqOrder == nil || infoOrder.OrderStatus == nil {
		fmt.Printf("%s stopped by %s\n", "UpdateAcPosByOrderInfo", "no infoOrder")
		return nil
	}
	// 获取对应持仓
	uniqueCode := infoOrder.ReqOrder.UniqueCode
	dirPos := chooseDirPos(infoOrder.ReqOrder.Dir, infoOrder.ReqOrder.FlagOffset)
	mapDirVal, hasMapDir := mapAcPos.Load(dirPos)
	var mapDir *sync.Map
	if hasMapDir {
		mapDir = mapDirVal.(*sync.Map)
	} else {
		mapDir = &sync.Map{}
		mapAcPos.Store(dirPos, mapDir)
	}
	posVal, hasPos := mapDir.Load(uniqueCode)
	var acPos *tradeBasic.SAcPos
	if hasPos {
		acPos = posVal.(*tradeBasic.SAcPos)
	} else {
		acPos = &tradeBasic.SAcPos{
			UniqueCode: uniqueCode,
		}
		mapDir.Store(uniqueCode, acPos)
	}
	// 分情况处理
	switch infoOrder.OrderStatus.StatusOrder {
	case tradeBasic.AllTraded:
		handleTrade(acPos, infoOrder.ReqOrder.FlagOffset, infoOrder.OrderStatus.VolTotal)
	case tradeBasic.PartialCanceled:
		handleTrade(acPos, infoOrder.ReqOrder.FlagOffset, infoOrder.OrderStatus.VolTraded)
	case tradeBasic.Canceled:
		handleCancel(acPos, infoOrder.ReqOrder.FlagOffset, infoOrder.OrderStatus.VolTotal)
		//case PartialTraded:
		//	fmt.Printf("%s stopped by %s %+v\n", "UpdateAcPosByOrderInfo", "unhandled status", infoOrder)
	}
	return acPos
}

// 处理
func handleCancel(acPos *tradeBasic.SAcPos, flagOffset tradeBasic.EFlagOffset, volCancel int64) {
	fmt.Printf("%s old acPos %+v\n", "handleCancel", acPos)
	// 计算今昨数量
	var frozenYd int64 = 0
	var frozenTd int64 = 0
	switch flagOffset {
	case tradeBasic.CloseToday:
		frozenTd = volCancel
	case tradeBasic.CloseYesterday:
		frozenYd = volCancel
	case tradeBasic.Close, tradeBasic.ForceClose, tradeBasic.ForceOff, tradeBasic.LocalForceClose:
		// 优先平昨，其次平今
		if volCancel <= acPos.VolFrozenYd {
			frozenYd = volCancel
		} else {
			frozenYd = acPos.VolFrozenYd
		}
		frozenTd = volCancel - frozenYd
	}
	// 更新持仓信息
	FrozenTd := acPos.VolFrozenTd - frozenTd
	if FrozenTd > 0 {
		acPos.VolFrozenTd = FrozenTd
	} else {
		acPos.VolFrozenTd = 0
	}
	FrozenYd := acPos.VolFrozenYd - frozenYd
	if FrozenYd > 0 {
		acPos.VolFrozenYd = FrozenYd
	} else {
		acPos.VolFrozenYd = 0
	}
	acPos.VolFrozenTotal = acPos.VolFrozenYd + acPos.VolFrozenTd
	fmt.Printf("%s new acPos %+v\n", "handleCancel", acPos)
}

func handleTrade(acPos *tradeBasic.SAcPos, flagOffset tradeBasic.EFlagOffset, volTrade int64) {
	fmt.Printf("%s old acPos %+v\n", "handleTrade", acPos)
	// 计算今昨数量
	var volYd int64 = 0
	var volTd int64 = 0
	var frozenYd int64 = 0
	var frozenTd int64 = 0
	switch flagOffset {
	case tradeBasic.Open:
		volTd = volTrade
	case tradeBasic.CloseToday:
		volTd = -volTrade
		frozenTd = volTd
	case tradeBasic.CloseYesterday:
		volYd = -volTrade
		frozenYd = volYd
	case tradeBasic.Close, tradeBasic.ForceClose, tradeBasic.ForceOff, tradeBasic.LocalForceClose:
		// 优先平昨，其次平今
		if volTrade <= acPos.VolYd {
			volYd = -volTrade
		} else {
			volYd = -acPos.VolYd
		}
		volTd = -(volTrade + volYd)
		frozenYd = volYd
		frozenTd = volTd
	}
	// 更新持仓信息
	acPos.VolYd = acPos.VolYd + volYd
	FrozenTd := acPos.VolFrozenTd + frozenTd
	if FrozenTd > 0 {
		acPos.VolFrozenTd = FrozenTd
	} else {
		acPos.VolFrozenTd = 0
	}
	FrozenYd := acPos.VolFrozenYd + frozenYd
	if FrozenYd > 0 {
		acPos.VolFrozenYd = FrozenYd
	} else {
		acPos.VolFrozenYd = 0
	}
	acPos.VolTd = acPos.VolTd + volTd
	acPos.VolTotal = acPos.VolTotal + volYd + volTd
	fmt.Printf("%s new acPos %+v\n", "handleTrade", acPos)
}

// CancelOrderBatch 根据报单方向批量撤单
/*
 * dir 需要撤单的报单方向
 * flagOffset 开平标记
 * mapOrderInfo ucOrder : infoOrder
 */
func CancelOrderBatch(dir investBasic.EDirTrade, flagOffset tradeBasic.EFlagOffset, mapOrderInfo *sync.Map) []*tradeBasic.PReqOrderAction {
	strMethod := "CancelOrderBatch"
	var rtns []*tradeBasic.PReqOrderAction
	if mapOrderInfo == nil {
		fmt.Printf("%s stopped by %s\n", strMethod, "no mapOrderInfo")
		return rtns
	}
	mapOrderInfo.Range(func(key, value any) bool {
		orderInfo := value.(*tradeBasic.SOrderInfo)
		// 筛选目标报单
		if orderInfo.ReqOrder != nil {
			if dir != investBasic.EMPTY &&
				dir != orderInfo.ReqOrder.Dir {
				return true
			}
			if flagOffset != tradeBasic.EFlagOffset_NONE &&
				flagOffset != orderInfo.ReqOrder.FlagOffset {
				return true
			}
		}
		// 生成撤单信息
		rtns = append(rtns, &tradeBasic.PReqOrderAction{
			OrderSys:    orderInfo.OrderSys,
			OrderAction: tradeBasic.Delete,
		})
		return true
	})
	return rtns
}
