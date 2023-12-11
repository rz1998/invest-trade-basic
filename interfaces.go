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
 * @Author: rz1998 rz1998@126.com
 * @Date: 2023-06-27 14:34:57
 * @LastEditors: rz1998 rz1998@126.com
 * @LastEditTime: 2023-09-17 15:59:27
 * @FilePath: /websocket/mnt/raid0/onedrive/zy/coding/workspace/github.com/rz1998/invest/trade/tradeBasic/interfaces.go
 * @Description:
 *
 */

package trade

import (
	"reflect"

	"github.com/rz1998/invest-basic/types/investBasic"
	"github.com/rz1998/invest-trade-basic/types/tradeBasic"
)

func NewSpiMD(structSpiMD interface{}) ISpiMD {
	t := reflect.TypeOf(structSpiMD)
	if t.Kind() == reflect.Ptr {
		//指针类型获取真正type需要调用Elem
		t = t.Elem()
	}
	newSpi := reflect.New(t).Interface()
	return newSpi.(ISpiMD)
}

func NewSpiTrader(structSpiTrader interface{}) ISpiTrader {
	t := reflect.TypeOf(structSpiTrader)
	if t.Kind() == reflect.Ptr {
		//指针类型获取真正type需要调用Elem
		t = t.Elem()
	}
	newSpi := reflect.New(t).Interface()
	return newSpi.(ISpiTrader)
}

type IApiMD interface {
	SetSpi(spi *ISpiMD)
	GetSpi() *ISpiMD
	Login(infoAc *tradeBasic.PInfoAc)
	Logout()
	Sub(uniqueCodes []string)
	Unsub(uniqueCodes []string)
	GetInfoAc() *tradeBasic.PInfoAc
	IsLogin() bool
	SetIsLogin(isLogin bool)
}

type IApiTrader interface {
	// Login 登录
	Login(infoAc *tradeBasic.PInfoAc)
	// Logout 登出
	Logout()
	// ReqOrder 报单请求
	ReqOrder(reqOrder *tradeBasic.PReqOrder)
	// ReqOrderBatch 批量报单请求
	ReqOrderBatch(reqOrders []*tradeBasic.PReqOrder)
	// ReqOrderAction 订单操作请求
	ReqOrderAction(reqOrderAction *tradeBasic.PReqOrderAction)
	// QryAcFund 查询资金
	QryAcFund()
	// QryAcLiability 查询负债
	QryAcLiability()
	// QryAcPos 查询持仓
	QryAcPos()
	// QryOrder 查询委托
	QryOrder(orderSys *tradeBasic.SOrderSys)
	// QryTrade 查询成交
	QryTrade()
	// SetSpi 设置回报监听
	SetSpi(spi *ISpiTrader)
	// GetSpi 获取回报监听
	GetSpi() *ISpiTrader
	// GetInfoSession 获取会话信息
	GetInfoSession() *tradeBasic.SInfoSessionTrader
	GenerateUniqueOrder(orderSys *tradeBasic.SOrderSys) string
}

type ISpiMD interface {
	// OnRtnMD 行情回报
	OnRtnMD(tick *investBasic.SMDTick)
	// OnFrontDisconnect 前置断开
	OnFrontDisconnect(reason int)
	// OnRspSub 订阅行情返回数据
	OnRspSub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool)
	// OnRspUnsub 取消订阅行情返回数据
	OnRspUnsub(code string, msgError *tradeBasic.SMsgError, idRequest int, isLast bool)
}

type ISpiTrader interface {
	// OnRtnOrder 报单回报
	OnRtnOrder(infoOrder *tradeBasic.SOrderInfo, isLast bool)
	// OnRtnTrade 成交回报
	OnRtnTrade(infoTrader *tradeBasic.STradeInfo, isLast bool)
	// OnRtnAcPos 账户持仓查询
	OnRtnAcPos(acPos *tradeBasic.SAcPos, isLast bool)
	// OnRtnAcFund 账户资金查询
	OnRtnAcFund(acFund *tradeBasic.SAcFund)
	// OnFrontDisconnect 前置断开回报
	OnFrontDisconnect(reason int)
	// OnErrRtnOrderInsert 报单错误回报
	OnErrRtnOrderInsert(reqOrder *tradeBasic.PReqOrder, msgError *tradeBasic.SMsgError)
	// OnErrRtnOrderAction 报单操作错误回报
	OnErrRtnOrderAction(reqOrderAction *tradeBasic.PReqOrderAction, msgError *tradeBasic.SMsgError)
}
