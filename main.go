package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// 创建链
	chain := NewChain()
	fmt.Println(chain)
	/*	// 创建区块
		block := NewBlock("", "转账十元")
		fmt.Println(block)
		// 将区块加入到链中
		chain.addBlockToChain(block)
		//再次进行验证
		fmt.Println(chain.chain)
		block2 := NewBlock("", "转账十个十元")

		chain.addBlockToChain(block2)
		fmt.Println(chain.chain)
		// 验证区块是否是合法的
		fmt.Println(chain.verificationChain())*/

	/*	// 创建交易
		t1 := NewTransaction("addr1", "addr2", 20)
		t2 := NewTransaction("addr2", "addr1", 10)
		// 将交易体到链上
		chain.addTransaction(t1)
		chain.addTransaction(t2)
		fmt.Println(fmt.Sprintf("%+v\n", chain))
		chain.mineTransactionPool("add3")
		fmt.Println(fmt.Sprintf("%+v\n", chain))*/

	// 生成两个用户
	privateKeySender, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	publicKeySender := privateKeySender.PublicKey
	privateKeyReceiver, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	publicKeyReceiver := privateKeyReceiver.PublicKey
	t1 := NewTransaction(privateKeySender, &publicKeyReceiver, 20)
	t2 := NewTransaction(privateKeyReceiver, &publicKeySender, 10)
	t1.sign(privateKeySender)
	fmt.Println(t1.isValid())
	//t1.Amount = 10
	chain.addTransaction(t1)
	chain.mineTransactionPool(&publicKeyReceiver)
	fmt.Printf("%+v\n", chain.chain[0].transaction)
	fmt.Printf("%+v\n", chain.chain[1].transaction)
	fmt.Println(t2)
}
