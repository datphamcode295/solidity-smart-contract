package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	todo "github.com/datphamcode295/solidity-smart-contract/gen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	b, err := ioutil.ReadFile("wallet/UTC--2022-06-12T14-25-14.179365454Z--834bec369b146648ec0e4170705d3a2725ab71be")
	if err != nil {
		log.Fatal(err)
	}

	key, err := keystore.DecryptKey(b, "password")
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial("https://kovan.infura.io/v3/2ee533b24d4d4bbd97523ecf4113e8bf")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	add := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), add)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000)
	auth.Nonce = big.NewInt(int64(nonce))

	a, tx, _, err := todo.DeployStore(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----------------------------------")
	fmt.Println(a.Hex())
	fmt.Println(tx.Hash().Hex())
	fmt.Println("-----------------------------------")
	// 0x934A955652F6A2858F48f655aE501A569258E922 -- contract address
	// 0x9132e0fb7ce59a70f58bf489223f0df7555b8a0ec0c8bb8d9be9d638c6c9be4b -- transaction address
}
