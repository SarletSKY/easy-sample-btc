package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
)

type Transaction struct {
	From      interface{} `json:"from"`      // 发起人
	To        interface{} `json:"to"`        // 接收者
	Amount    int         `json:"amount"`    // 数量
	Signature []byte      `json:"signature"` // 签名
}

func NewTransaction(from interface{}, to interface{}, amount int) Transaction {
	return Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (t *Transaction) completeHash() string {
	// 将Transaction进行json格式化
	// 2 使用x509将私钥进行序列化
	derFrom := crypto.FromECDSA(t.From.(*ecdsa.PrivateKey))
	derTo := crypto.FromECDSAPub(t.To.(*ecdsa.PublicKey))
	hashBytes := sha256.Sum256([]byte(bytes.NewBuffer(derFrom).String() + bytes.NewBuffer(derTo).String() + strconv.Itoa(t.Amount)))
	return hex.EncodeToString(hashBytes[:])
}

// 私钥签名
func (t *Transaction) sign(privateKey *ecdsa.PrivateKey) {
	var err error
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, t.HashToBytes())
	if err != nil {
		panic(fmt.Sprintf("签名异常:%v", err))
	}
	params := privateKey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	t.Signature = make([]byte, curveOrderByteSize*2)
	copy(t.Signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(t.Signature[curveOrderByteSize*2-len(sBytes):], sBytes)
}

// 公钥验证
func (t *Transaction) isValid() bool {
	// 当from是为nil时，是不需要验证的，因为是矿工发放奖励
	if t.From == "" {
		return true
	}
	curveOrderByteSize := t.From.(*ecdsa.PrivateKey).PublicKey.Curve.Params().P.BitLen() / 8
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(t.Signature[:curveOrderByteSize])
	s.SetBytes(t.Signature[curveOrderByteSize:])
	return ecdsa.Verify(&t.From.(*ecdsa.PrivateKey).PublicKey, t.HashToBytes(), r, s)
}

// 将hash(string)转换成[]byte
func (t *Transaction) HashToBytes() []byte {
	hashString := t.completeHash()
	hash, err := hex.DecodeString(hashString)
	if err != nil {
		panic(fmt.Sprintf("hash数据格式转换失败:%v", err))
	}
	return hash
}
