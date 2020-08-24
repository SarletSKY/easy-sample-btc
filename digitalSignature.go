package main

import (
	"crypto/ecdsa"
	//"crypto/elliptic"
	//"crypto/rand"
	//"encoding/hex"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func getKey() (*ecdsa.PrivateKey, error) {
	prk, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return prk, err
	}
	return prk, nil
}

func eccSign(data []byte, prk *ecdsa.PrivateKey) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, prk, data)
	if err != nil {
		return nil, err
	}
	params := prk.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)
	return signature, nil
}

func eccVerify(data, signature []byte, puk *ecdsa.PublicKey) bool {
	curveOrderByteSize := puk.Curve.Params().P.BitLen() / 8
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])
	return ecdsa.Verify(puk, data, r, s)
}

/*func main() {
	data := "00007f1f64109f1df066db39cdcfd7bb2343a02bf8a3054399f4772ed640300e"
	prk, err := getKey()
	puk := prk.PublicKey
	if err != nil {
		panic(err)
	}
	//bdata := []byte(data)
	bpuk := crypto.FromECDSAPub(&puk)
	lpuk, err := crypto.UnmarshalPubkey(bpuk)

	//测试是否转换成功
	eccData, err := eccSign([]byte(data), prk)
	if err != nil {
		panic(err)
	}
	fmt.Println(eccVerify([]byte(data), eccData, lpuk))
	//  fmt.Println(data)
}*/
