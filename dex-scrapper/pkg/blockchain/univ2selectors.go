package blockchain

import "encoding/hex"

type UniV2Selector struct {
	UniV2Methods map[string]string
}

func NewUniV2Selector() *UniV2Selector {

	var u UniV2Selector

	u.UniV2Methods = map[string]string{
		"fb3bdb41": "swapETHForExactTokens",
		"7ff36ab5": "swapExactETHForTokens",
		"b6f9de95": "swapExactETHForTokensSupportingFeeOnTransferTokens", //fee on transfer
		"18cbafe5": "swapExactTokensForETH",
		"791ac947": "swapExactTokensForETHSupportingFeeOnTransferTokens", //fee on transfer
		"38ed1739": "swapExactTokensForTokens",
		"5c11d795": "swapExactTokensForTokensSupportingFeeOnTransferTokens", //fee on transfer
		"4a25d94a": "swapTokensForExactETH",
		"8803dbee": "swapTokensForExactTokens",
	}

	return &u
}

func (u *UniV2Selector) IsUniV2(txInput []byte) bool {
	_, exists := u.UniV2Methods[hex.EncodeToString(txInput[0:4])]
	return exists
}
