/*
 *
 * 　　┏┓　　　┏┓+ +
 * 　┏┛┻━━━┛┻┓ + +
 * 　┃　　　　　　　┃
 * 　┃　　　━　　　┃ ++ + + +
 *  ████━████ ┃+
 * 　┃　　　　　　　┃ +
 * 　┃　　　┻　　　┃
 * 　┃　　　　　　　┃ + +
 * 　┗━┓　　　┏━┛
 * 　　　┃　　　┃
 * 　　　┃　　　┃ + + + +
 * 　　　┃　　　┃
 * 　　　┃　　　┃ +  神兽保佑
 * 　　　┃　　　┃    代码无bug
 * 　　　┃　　　┃　　+
 * 　　　┃　 　　┗━━━┓ + +
 * 　　　┃ 　　　　　　　┣┓
 * 　　　┃ 　　　　　　　┏┛
 * 　　　┗┓┓┏━┳┓┏┛ + + + +
 * 　　　　┃┫┫　┃┫┫
 * 　　　　┗┻┛　┗┻┛+ + + +
 *
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-05-14 02:48:25
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-07-02 06:14:15
 * @FilePath: /grid/home/rz1998/workspace/github.com/rz1998/invest/trade/tradeBasic/utilVol.go
 * @Description:
 *
 */

package tradeBasic

import (
	"math"

	investBasic "github.com/rz1998/invest-basic"
	"github.com/zeromicro/go-zero/core/logx"
)

func GetMinVol(uniqueCode string) int64 {
	code, exchangeCD := investBasic.GetSecInfo(uniqueCode)
	// 默认1
	var min int64 = 1
	switch exchangeCD {
	case investBasic.SSE, investBasic.SZSE, investBasic.BSE:
		// 沪深市场
		str2 := code[:2]
		str3 := code[:3]
		if str3 == "131" || str2 == "11" || str2 == "12" {
			// 深市逆回购、所有转债，10
			min = 10
		} else if str2 == "68" {
			// 科创板200
			min = 200
		} else {
			// 其他100
			min = 100
		}
	default:
		logx.Errorf("getVolMin unhandled exchangeCD %v\n", exchangeCD)
	}
	return min
}

// 根据代码返回最小下单单位
func GetTickVol(uniqueCode string) int64 {
	code, exchangeCD := investBasic.GetSecInfo(uniqueCode)
	// 默认1
	var tick int64 = 1
	switch exchangeCD {
	case investBasic.SSE, investBasic.SZSE, investBasic.BSE:
		// 沪深市场
		str2 := code[:2]
		str3 := code[:3]
		if str3 == "131" || str2 == "11" || str2 == "12" {
			// 深市逆回购、所有转债，10
			tick = 10
		} else if str2 == "68" {
			// 科创板200
			tick = 1
		} else {
			// 其他100
			tick = 100
		}
	default:
		logx.Errorf("getVolTick unhandled exchangeCD %v\n", exchangeCD)
	}
	return tick
}

// CalLimitVol 根据目标比例计算目标报单数量
func CalLimitVol(rateTar float64, vol, volMin, volTick int64) int64 {
	if volTick <= 0 {
		return 0
	}
	volTar := float64(vol) * rateTar
	if volTar < float64(volMin)/2 {
		return 0
	} else if volTar < float64(volMin) && volTar >= float64(volMin)/2 {
		return volMin
	}
	return volMin + int64(math.Round((volTar-float64(volMin))/float64(volTick)))*volTick
}

// CalTarVol 根据目标比例计算目标报单数量（根据代码自动判断最小下单单位）
func CalTarVol(uniqueCode string, rateTar float64, vol int64) int64 {
	return CalLimitVol(rateTar, vol, GetMinVol(uniqueCode), GetTickVol(uniqueCode))
}

type TypeRound uint

const (
	Round TypeRound = iota
	Ceil
	Floor
)

// CalTarVolByVal 根据目标金额，确定目标数量
func CalTarVolByVal(fees PFees, uniqueCode string, val float64, price int64, typeRound TypeRound) int64 {
	// 理论最大数量
	volTheory := val / float64(price) * 10000
	// 数量最小值
	volMin := GetMinVol(uniqueCode)
	if volTheory < float64(volMin)/2 {
		return 0
	} else if volTheory <= float64(volMin) && volTheory >= float64(volMin)/2 {
		return volMin
	}
	// 数量步长
	volTick := GetTickVol(uniqueCode)
	// 调整到一跳的倍数
	var vol int64
	switch typeRound {
	case Round:
		// 四舍五入到目标报单数量上
		vol = volMin + int64(math.Round((volTheory-float64(volMin))/float64(volTick)))*volTick
	case Ceil:
		// 向上取整到目标报单数量上
		vol = volMin + int64(math.Ceil((volTheory-float64(volMin))/float64(volTick)))*volTick
	case Floor:
		// 向下取整到目标报单数量上
		vol = volMin + int64(math.Floor((volTheory-float64(volMin))/float64(volTick)))*volTick
		// 向下取整需要考虑手续费充足与否
		for vol > 0 && (float64(vol)*float64(price)/10000)+CalFee(fees, vol, price) > val {
			vol -= volTick
		}
	default:
		logx.Errorf("CalTarVolByVal unhandled typeRound %v\n", typeRound)
	}
	return vol
}
