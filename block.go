package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

//设置块对象
// 哈希值：data与上一个哈希进行哈希
type Block struct {
	preHash         string      // 前一个哈希
	hash            string      // 当前哈希 	// 时间戳在当前哈希中设置
	nonce           int         // 随机数
	transaction     interface{} // 将交易内容写成transaction
	TransactionTime int64       `json:"transaction_time"` // 交易时间
}

// block的构造函数
func NewBlock(preHash string, transactions interface{}) Block {
	b := Block{
		preHash:         preHash,
		nonce:           1,
		transaction:     transactions,
		TransactionTime: time.Now().Unix(),
	}
	b.hash = b.completeHash()
	return b
}

//  获取挖矿等级需求来进行判断hash是否满足
func getMineDiff(diff int) (str string) {
	for i := 0; i < diff; i++ {
		str = str + "0"
	}
	return
}

// 开始挖矿[链的挖矿难度]
func (b *Block) mine(diff int) {
	strDiff := getMineDiff(diff)

	// 在开始挖矿之前验证区块的交易合法性
	if !b.validateTransactions() {
		panic("发现异常交易，停止挖矿")
	}
	for {
		b.hash = b.completeHash()
		if b.hash[:diff] == strDiff {
			log.Println("挖矿结束", b.hash)
			break
		}
		b.nonce++
	}
}

// 交易的合法性
func (b *Block) validateTransactions() bool {
	// 所有交易进行判断合法性
	for _, transaction := range b.transaction.([]Transaction) {
		if !transaction.isValid() {
			return false
		}
	}
	return true
}

// 进行哈希
func (b *Block) completeHash() string {
	// 将Transaction进行json格式化
	transactionBytes, _ := json.Marshal(b.transaction)
	hashBytes := sha256.Sum256(
		[]byte(string(transactionBytes) +
			b.preHash +
			strconv.Itoa(b.nonce) +
			strconv.Itoa(int(b.TransactionTime)),
		))
	return hex.EncodeToString(hashBytes[:])
}
