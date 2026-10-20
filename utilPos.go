/*
 *  ┌───┐   ┌───┬───┬───┬───┐ ┌───┬───┬───┬───┐ ┌───┬───┬───┬───┐ ┌───┬───┬───┐
 *  │Esc│   │ F1│ F2│ F3│ F4│ │ F5│ F6│ F7│ F8│ │ F9│F10│F11│F12│ │P/S│S L│P/B│  ┌┐    ┌┐    ┌┐
 *  └───┘   └───┴───┴───┴───┘ └───┴───┴───┴───┘ └───┴───┴───┴───┘ └───┴───┴───┘  └┘    └┘    └┘
 *  ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───────┐ ┌───┬───┬───┐ ┌───┬───┬───┬───┐
 *  │~ `│! 1│@ 2│# 3│$ 4│% 5│^ 6│& 7│* 8│( 9│) 0│_ -│+ =│ BacSp │ │Ins│Hom│PUp│ │N L│ / │ * │ - │
 *  ├───┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─────┤ ├───┼───┼───┤ ├───┼───┼───┼───┤
 *  │ Tab │ Q │ W │ E │ R │ T │ Y │ U │ I │ O │ P │{ [│} ]│ | \ │ │Del│End│PDn│ │ 7 │ 8 │ 9 │   │
 *  ├─────┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴─────┤ └───┴───┴───┘ ├───┼───┼───┤ + │
 *  │ Caps │ A │ S │ D │ F │ G │ H │ J │ K │ L │: ;│" '│ Enter  │               │ 4 │ 5 │ 6 │   │
 *  ├──────┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴────────┤     ┌───┐     ├───┼───┼───┼───┤
 *  │ Shift  │ Z │ X │ C │ V │ B │ N │ M │< ,│> .│? /│  Shift   │     │ ↑ │     │ 1 │ 2 │ 3 │   │
 *  ├─────┬──┴─┬─┴──┬┴───┴───┴───┴───┴───┴──┬┴───┼───┴┬────┬────┤ ┌───┼───┼───┐ ├───┴───┼───┤ E││
 *  │ Ctrl│    │Alt │         Space         │ Alt│    │    │Ctrl│ │ ← │ ↓ │ → │ │   0   │ . │←─┘│
 *  └─────┴────┴────┴───────────────────────┴────┴────┴────┴────┘ └───┴───┴───┘ └───────┴───┴───┘
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-05-14 06:41:38
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-09-20 10:24:43
 * @FilePath: /tradeNew/home/rz1998/workspace/github.com/rz1998/invest/trade/tradeBasic/utilPos.go
 * @Description:
 *
 */

package tradeBasic

import (
	"sync"

	investBasic "github.com/rz1998/invest-basic"
	"github.com/zeromicro/go-zero/core/logx"
)

// 根据订单的开平标致和买卖方向，选择受影响的持仓方向
/*
 * @param dirTrade 目标交易方向
 * @param flagOffset 开平标致
 * @return 所需持仓的方向
 */
func chooseDirPos(dirTrade investBasic.EDirTrade, flagOffset EFlagOffset) investBasic.EDirTrade {
	var dir investBasic.EDirTrade
	switch flagOffset {
	case Open:
		dir = dirTrade
	case Close, CloseToday, CloseYesterday,
		ForceClose, LocalForceClose, ForceOff:
		if dirTrade == investBasic.LONG {
			dir = investBasic.SHORT
		} else {
			dir = investBasic.LONG
		}
	}
	return dir
}

// 找到指定品种制定方向的持仓
// mapPos uniqueCode : dir : acPos
func GetAcPos(mapPos *sync.Map, uniqueCode string, dir investBasic.EDirTrade) (*SAcPos, bool) {
	mapAny, ok := mapPos.Load(uniqueCode)
	if ok {
		mapDir := mapAny.(*sync.Map)
		posAny, hasPos := mapDir.Load(dir)
		if hasPos {
			return posAny.(*SAcPos), true
		}
	}
	return nil, false
}

// GetVolSellableFromAcPos 根据持仓计算可卖数量
func GetVolSellableFromAcPos(acPos *SAcPos) (vol int64) {
	if acPos == nil {
		return vol
	}
	_, exchangeCD := investBasic.GetSecInfo(acPos.UniqueCode)
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

// UpdateAcPosByOrderInfo 用报单回报更新持仓
/*
 * mapAcPos dirTrade : uniqueCode : acPos
 */
func UpdateAcPosByOrderInfo(mapAcPos *sync.Map, infoOrder *SOrderInfo) *SAcPos {
	if mapAcPos == nil {
		logx.Infof("%s stopped by %s\n", "UpdateAcPosByOrderInfo", "no map acPos")
		return nil
	}
	if infoOrder == nil || infoOrder.ReqOrder == nil || infoOrder.OrderStatus == nil {
		logx.Infof("%s stopped by %s\n", "UpdateAcPosByOrderInfo", "no infoOrder")
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
	var acPos *SAcPos
	if hasPos {
		acPos = posVal.(*SAcPos)
	} else {
		acPos = &SAcPos{
			UniqueCode: uniqueCode,
		}
		mapDir.Store(uniqueCode, acPos)
	}
	// 分情况处理
	switch infoOrder.OrderStatus.StatusOrder {
	case AllTraded:
		handleTrade(acPos, infoOrder.ReqOrder.FlagOffset, infoOrder.OrderStatus.VolTotal)
	case PartialCanceled:
		handleTrade(acPos, infoOrder.ReqOrder.FlagOffset, infoOrder.OrderStatus.VolTraded)
	case Canceled:
		handleCancel(acPos, infoOrder.ReqOrder.FlagOffset, infoOrder.OrderStatus.VolTotal)
	}
	return acPos
}

// 处理
func handleCancel(acPos *SAcPos, flagOffset EFlagOffset, volCancel int64) {
	logx.Infof("%s old acPos %+v\n", "handleCancel", acPos)
	// 计算今昨数量
	var frozenYd int64 = 0
	var frozenTd int64 = 0
	switch flagOffset {
	case CloseToday:
		frozenTd = volCancel
	case CloseYesterday:
		frozenYd = volCancel
	case Close, ForceClose, ForceOff, LocalForceClose:
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
	logx.Infof("%s new acPos %+v\n", "handleCancel", acPos)
}

func handleTrade(acPos *SAcPos, flagOffset EFlagOffset, volTrade int64) {
	logx.Infof("%s old acPos %+v\n", "handleTrade", acPos)
	// 计算今昨数量
	var volYd int64 = 0
	var volTd int64 = 0
	var frozenYd int64 = 0
	var frozenTd int64 = 0
	switch flagOffset {
	case Open:
		volTd = volTrade
	case CloseToday:
		volTd = -volTrade
		frozenTd = volTd
	case CloseYesterday:
		volYd = -volTrade
		frozenYd = volYd
	case Close, ForceClose, ForceOff, LocalForceClose:
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
	logx.Infof("%s new acPos %+v\n", "handleTrade", acPos)
}
