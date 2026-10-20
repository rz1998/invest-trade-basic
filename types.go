package tradeBasic

import (
	"fmt"

	investBasic "github.com/rz1998/invest-basic"
)

type EFlagOffset uint

const (
	EFlagOffset_NONE EFlagOffset = iota
	// Open FlagOffset_Open 开仓
	Open
	// Close FlagOffset_Close 平仓
	Close
	// ForceClose FlagOffset_ForceClose 强平
	ForceClose
	// CloseToday FlagOffset_CloseToday 平今
	CloseToday
	// CloseYesterday FlagOffset_CloseYesterday 平昨
	CloseYesterday
	// ForceOff 强减
	ForceOff
	// LocalForceClose 本地强平
	LocalForceClose
)

type ETypeOrderAction uint

const (
	ETypeOrderAction_NONE ETypeOrderAction = iota
	// Delete 撤单
	Delete
	// Modify 修改订单
	Modify
)

type ESourceInfo uint

const (
	ESourceInfo_NONE ESourceInfo = iota
	// RETURN 来自交易所回报
	RETURN
	// QUERY 来自查询
	QUERY
)

type EStatusOrderSubmit uint

const (
	EStatusOrderSubmit_NONE EStatusOrderSubmit = iota
	// InsertSubmitted 已经提交
	InsertSubmitted
	// CancelSubmitted 撤单已经提交
	CancelSubmitted
	// ModifySubmitted 修改已经提交
	ModifySubmitted
	// Accepted 已经接受
	Accepted
	// InsertRejected 报单已经被拒绝
	InsertRejected
	// CancelRejected 撤单已经被拒绝
	CancelRejected
	// ModifyRejected 改单已经被拒绝
	ModifyRejected
)

type EStatusOrder uint

const (
	EStatusOrder_NONE EStatusOrder = iota
	// Unknown 未知
	Unknown
	// NotTraded 未成交
	NotTraded
	// PartialTraded 部分成交
	PartialTraded
	// AllTraded 全部成交/已成
	AllTraded
	// Canceled 全撤
	Canceled
	// PartialCanceled 部分撤单
	PartialCanceled
)

type ETypeTrade uint

const (
	ETypeTrade_NONE ETypeTrade = iota
	// Common 普通成交
	Common
	// OptionsExecution 期权执行
	OptionsExecution
	// OTC OTC成交
	OTC
	// EFPDerived 期转现衍生成交
	EFPDerived
	// CombinationDerived 组合衍生成交
	CombinationDerived
	// SplitCombination 组合持仓拆分为单一持仓,初始化不应包含该类型的持仓
	SplitCombination
	// BlockTrade 大宗成交
	BlockTrade
)

type ESourcePrice uint

const (
	ESourcePrice_NONE ESourcePrice = iota
	// LastPrice 前成交价
	LastPrice
	// Buy 买委托价
	Buy
	// Sell 卖委托价
	Sell
	// OTC 场外成交价
	ESourcePrice_OTC
)

type SMsgError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// PApiBroker 中介api接口信息
type PApiBroker struct {
	// 接口名称
	ApiName string `json:"apiName"`
	// 经济商代码
	BrokerId string `json:"brokerId"`
	// 交易前置URL
	UrlFrontTrader []string `json:"urlFrontTrader"`
	// 行情前置URL
	UrlFrontMD []string `json:"urlFrontMD"`
}

// PAuthTrader 穿透式监管信息
type PAuthTrader struct {
	BrokerId        string `json:"brokerId"`
	UsrID           string `json:"usrID"`
	AppID           string `json:"appID"`
	AuthCode        string `json:"authCode"`
	UserProductInfo string `json:"userProductInfo"`
}

// PInfoAc 用户账户信息
type PInfoAc struct {
	ApiName string `json:"apiName"`
	// 经纪公司代码
	BrokerId string `json:"brokerId"`
	// 投资者编号
	InvestorId string `json:"investorId"`
	// 账户账号
	UserId string `json:"userId"`
	// 账户名
	UserName string `json:"userName"`
	// 密码
	Psw string `json:"psw"`
}

// UniqueCode 生成账户唯一识别
func (infoAc *PInfoAc) UniqueCode() string {
	return fmt.Sprintf("%s|%s|%s", infoAc.BrokerId, infoAc.InvestorId, infoAc.UserId)
}

type SAcFund struct {
	// 可取资金
	ValWithdraw float64 `json:"valWithdraw"`
	// 可用资金
	ValAvailable float64 `json:"valAvailable"`
	// 冻结资金
	ValFrozen float64 `json:"valFrozen"`
	// 保证金占用
	Margin float64 `json:"margin"`
	// 总权益：可用+持仓（支持某些接口）
	ValBalance float64 `json:"valBalance"`
	// 账户识别
	IdAccount string `json:"idAccount,optional"`
}

// SAcPos 账户持仓
type SAcPos struct {
	// 持仓证券
	UniqueCode string `json:"uniqueCode"`
	// 持仓方向
	TradeDir investBasic.EDirTrade `json:"tradeDir"`
	// 成交均价
	PriceOpen int64 `json:"priceOpen"`
	// 结算价
	PriceSettle int64 `json:"priceSettle"`
	// 持仓总量
	VolTotal int64 `json:"volTotal"`
	// 昨仓量
	VolYd int64 `json:"volYd"`
	// 今仓量
	VolTd int64 `json:"volTd"`
	// 总冻结量
	VolFrozenTotal int64 `json:"volFrozenTotal"`
	// 昨仓冻结数量
	VolFrozenYd int64 `json:"volFrozenYd"`
	// 今仓冻结数量
	VolFrozenTd int64 `json:"volFrozenTd"`
	// VolCanUse 是柜台给的权威可卖数量 (QMT 的 m_nCanUseVolume).
	//
	// 为什么是指针而不是 int64: 「字段存在且为 0」和「字段不存在」是两件
	// 完全不同的事. 指针的 nil = 字段不存在, 走下面的 Go 公式
	// (VolYd - VolFrozenYd); 非 nil = 柜台的权威值, 哪怕是 0 也要照用.
	//
	// 把 0 当成「没有数据」而退回公式, 是 2026-10-16 al-013 当日 63 笔
	// 拒单事故的直接成因: 第一笔卖单受理后柜台在几秒内把可卖量从 900
	// 翻到 0, 旧逻辑据此推出「仍可卖 900 股」, 同一笔持仓被重复下单
	// 12 次, 11 次被 251005(可用数量不足) 拒, 而且撤单窗口也撤不掉
	// (价格不在钓鱼窗口内). best_qmt.py 的 v29 用 can_use==0 即 0
	// 修掉了它, 这里是同一件事的 Go 侧表达.
	//
	// 指针同时保持了 JSON 兼容: 老文件没有这个键, 反序列化出来就是 nil,
	// 行为与改动前完全一致 —— 对其它使用方零回归.
	VolCanUse *int64 `json:"volCanUse,omitempty"`
	// 占用保证金
	Margin float64 `json:"margin"`
	// 账户识别
	IdAccount string `json:"idAccount,optional"`
}

// PReqOrder 报单请求
type PReqOrder struct {
	// 报单时间
	Timestamp int64 `json:"timestamp"`
	// 报单来源识别
	OrderRef string `json:"orderRef"`
	// 交易方向
	Dir investBasic.EDirTrade `json:"dir"`
	// 交易类型
	FlagOffset EFlagOffset `json:"flagOffset"`
	// 标的
	UniqueCode string `json:"uniqueCode"`
	// 价格
	Price int64 `json:"price"`
	// 数量
	Vol int64 `json:"vol"`
}

// SInfoSessionTrader ctp会话里的信息缓存
type SInfoSessionTrader struct {
	IdFront           int    `json:"idFront"`
	IdSession         int    `json:"idSession"`
	MaxOrderRef       string `json:"maxOrderRef"`
	MaxOrderActionRef int    `json:"maxOrderActionRef"`
}

// SOrderSys 报单系统信息
type SOrderSys struct {
	// 申报时间
	Timestamp int64 `json:"timestamp"`
	// 请求编号
	IdRequest int `json:"idRequest"`
	// 前置编号
	IdFront int `json:"idFront"`
	// 会话编号
	IdSession int `json:"idSession"`
	// 报单引用编号
	OrderRef string `json:"orderRef"`
	// 本地报单编号
	IdOrderLocal string `json:"idOrderLocal"`
	// 报单编号
	IdOrderSys string `json:"idOrderSys"`
	// 交易所代码
	IdExchange string `json:"idExchange"`
	// 来源类型（查询/回报）
	SourceInfo ESourceInfo `json:"sourceInfo"`
}

func (s *SOrderSys) GenerateIDLocal() string {
	return fmt.Sprintf("%d|%d|%s", s.IdFront, s.IdSession, s.OrderRef)
}

func (s *SOrderSys) GenerateIDServer() string {
	return fmt.Sprintf("%s|%s", s.IdExchange, s.IdOrderSys)
}

// PReqOrderAction 撤单请求
type PReqOrderAction struct {
	// 报单操作引用
	OrderActionRef string `json:"orderActionRef"`
	// 报单唯一识别
	OrderSys *SOrderSys `json:"orderSys"`
	// 操作标志
	OrderAction ETypeOrderAction `json:"orderAction"`
}

// SOrderStatus 订单状态信息
type SOrderStatus struct {
	// 最后修改时间
	Timestamp int64 `json:"timestamp"`
	// 撤销时间
	TimeCancel int64 `json:"timeCancel"`
	// 订单提交状态
	StatusOrderSubmit EStatusOrderSubmit `json:"statusOrderSubmit"`
	// 订单状态
	StatusOrder EStatusOrder `json:"statusOrder"`
	// 状态信息
	StatusMsg string `json:"statusMsg"`
	// 已成交数量
	VolTraded int64 `json:"volTraded"`
	// 总数量
	VolTotal int64 `json:"volTotal"`
}

// SOrderInfo 报单信息
type SOrderInfo struct {
	OrderSys    *SOrderSys    `json:"orderSys"`
	ReqOrder    *PReqOrder    `json:"reqOrder"`
	OrderStatus *SOrderStatus `json:"orderStatus"`
	// 账户识别
	IdAccount string `json:"idAccount,optional"`
}

func (infoOrder SOrderInfo) String() string {
	return fmt.Sprintf("%+v %+v %+v", *infoOrder.OrderSys, *infoOrder.ReqOrder, *infoOrder.OrderStatus)
}

// STradeStatus 成交状态信息
type STradeStatus struct {
	// 时间戳
	Timestamp int64 `json:"timestamp"`
	// 成交编号
	IdTrade string `json:"idTrade"`
	// 成交类型
	TypeTrade ETypeTrade `json:"typeTrade"`
	// 成交价来源
	SourcePrice ESourcePrice `json:"sourcePrice"`
	// 价格
	Price int64 `json:"price"`
	// 数量
	Vol int64 `json:"vol"`
	// 成交额
	Val float64 `json:"val"`
	// 保证金占用
	Margin float64 `json:"margin"`
	// 交易费用信息
	Fees float64 `json:"fees"`
	// 成交信息来源
	TradeSource ESourceInfo `json:"tradeSource"`
}

// STradeInfo 成交信息
type STradeInfo struct {
	// 交易日
	Date string `json:"date"`
	// 下单请求信息
	ReqOrder *PReqOrder `json:"reqOrder"`
	// 交易系统信息
	OrderSys *SOrderSys `json:"orderSys"`
	// 成交状态信息
	TradeStatus *STradeStatus `json:"tradeStatus"`
	// 账户识别
	IdAccount string `json:"idAccount,optional"`
}

func (infoTrade STradeInfo) String() string {
	return fmt.Sprintf("%s %+v %+v %+v", infoTrade.Date, *infoTrade.ReqOrder, *infoTrade.OrderSys, *infoTrade.TradeStatus)
}

// STransferInfo 转账信息
type STransferInfo struct {
	Datetime  string  `json:"datetime"`
	IdRecord  string  `json:"idRecord"`
	IdAccount string  `json:"idAccount"`
	Amount    float64 `json:"amount"`
}
