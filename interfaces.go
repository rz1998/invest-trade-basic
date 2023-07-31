package trade

import (
	"github.com/rz1998/invest-basic/types/investBasic"
	"github.com/rz1998/invest-trade-basic/types/tradeBasic"
)

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
