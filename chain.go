package main

import (
	"crypto/ecdsa"
	"log"
)

// 链的结构体
type Chain struct {
	chain           []Block
	diff            int           // 挖矿等级
	TransactionPool []Transaction // 矿池
	MinerReward     int           // 矿池奖励
}

// 初始化链
func NewChain() *Chain {
	return &Chain{
		chain:           []Block{firstBlock()},
		TransactionPool: make([]Transaction, 0),
		MinerReward:     50,
		diff:            4,
	}
}

// 挖矿,先是生成奖励，然后开始挖矿生成区块
func (c *Chain) mineTransactionPool(minerAddress *ecdsa.PublicKey) {
	// 挖矿是一笔交易
	t1 := NewTransaction("", minerAddress, c.MinerReward)
	// 将这笔挖矿交易加入到矿池上
	c.TransactionPool = append(c.TransactionPool, t1)
	// 开始挖矿, 生成区块
	newBlock := NewBlock(c.getLastBlock().hash, c.TransactionPool)
	newBlock.mine(c.diff)
	// 将挖到矿加入到链上
	c.TransactionPool = []Transaction{}
	c.chain = append(c.chain, newBlock)
}

// 设置祖先区块
func firstBlock() Block {
	return NewBlock("", "我是祖先")
}

// 将区块加入到链中
func (c *Chain) addBlockToChain(block Block) {
	//将区块的上一个哈希进行拼接
	// 先获取链上的最后一个区块，然后hash值重新计算
	block.preHash = c.getLastBlock().hash
	// 在哈希加入之前要满足条件才允许加入, 矿工进行挖矿
	block.mine(c.diff)
	c.chain = append(c.chain, block)
}

// 链上的最后一个区块
func (c *Chain) getLastBlock() Block {
	return c.chain[len(c.chain)-1]
}

// 验证区块的合法性
func (c *Chain) verificationChain() bool {
	// 如果只有一个初始区块，则返回true
	if len(c.chain) == 1 {
		if c.chain[0].hash != c.chain[0].completeHash() {
			return false
		}
		return true
	}
	for i := 0; i < len(c.chain)-1; i++ {
		if !c.chain[i].validateTransactions() {
			log.Println("非法交易")
			return false
		}
		// 前一个区块的hash是否是后一个区块的前一个哈希
		if c.chain[i].hash != c.chain[i].completeHash() {
			log.Println("数据被篡改")
			return false
		}
		if c.chain[i].hash != c.chain[i+1].preHash {
			log.Println("前后区块断裂")
			return false
		}
	}
	return true
}

// 添加交易到区块上
func (c *Chain) addTransaction(transaction Transaction) {
	// 再添加到区块上的交易之前，要验证合法性
	if !(transaction.isValid()) {
		panic("invalid transaction")
	}
	log.Println("valid transaction")
	c.TransactionPool = append(c.TransactionPool, transaction)
}
