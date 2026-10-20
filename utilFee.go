/*
 *
 *    ┏┓　　　┏┓
 *  ┏┛┻━━━┛┻┓
 *  ┃　　　　　　　┃
 *  ┃　　　━　　　┃
 *  ┃　＞　　　＜　┃
 *  ┃　　　　　　　┃
 *  ┃...　⌒　...　┃
 *  ┃　　　　　　　┃
 *  ┗━┓　　　┏━┛
 *      ┃　　　┃
 *      ┃　　　┃
 *      ┃　　　┃
 *      ┃　　　┃  神兽保佑
 *      ┃　　　┃  代码无bug
 *      ┃　　　┃
 *      ┃　　　┗━━━┓
 *      ┃　　　　　　　┣┓
 *      ┃　　　　　　　┏┛
 *      ┗┓┓┏━┳┓┏┛
 *        ┃┫┫　┃┫┫
 *        ┗┻┛　┗┻┛
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-05-14 02:50:17
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-05-14 07:51:25
 * @FilePath: /tradeBasic/utilFee.go
 * @Description: 费用相关的计算
 *
 */

package tradeBasic

// 费用收取方式
type PFees struct {
	// 按量收取的部分（计算时，FeeByVol*vol)
	FeeByVol float64
	// 按笔收取的部分（计算时，每笔直接加上FeeByNum)
	FeeByNum float64
	// 按交易额收取的部分（计算时，FeeByVal*val)
	FeeByVal float64
	// 最低收费金额，当最终计算出的费用少于此金额时，以此为准
	FeeMin float64
}

// CalValLeastByFee 根据最低佣金额及佣金率，计算单次下单最小额度
/*
 * rateFeeByVal 按金额计算的佣金率
 * feeMin  最低绝对佣金额
 */
func CalValLeastByFee(rateFeeByVal, feeMin float64) float64 {
	return feeMin / rateFeeByVal
}

// 计算佣金
func CalFee(fees PFees, price, vol int64) float64 {
	fee := fees.FeeByVol*float64(vol) + fees.FeeByNum + fees.FeeByVal*float64(price)/10000
	if fee < fees.FeeMin {
		fee = fees.FeeMin
	}
	return fee
}

// 计算股票资金需求
func CalTarAvailable(fees PFees, rateTax float64, price, vol int64) float64 {
	// 计算交易金额
	val := float64(price) / 10000 * float64(vol)
	// 计算收费金额
	fee := CalFee(fees, price, vol)
	return val + fee
}

// 估算剔除手续费之后实际可用的金额
func CalValLeft(fees PFees, val float64) float64 {
	return (val - fees.FeeByNum) / (1 + fees.FeeByVal)
}
