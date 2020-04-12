package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
)

const (
	testCert    = ""
	testPrivKey = ""
)

func main() {
	cert, err := ioutil.ReadFile(testCert)
	if err != nil {
		fmt.Errorf("can't to read cert, error: %s", err)
	}

	privKey, err := ioutil.ReadFile(testPrivKey)
	if err != nil {
		fmt.Errorf("can't to read privkey, error: %s", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	err = wallet.Put("admin", gateway.NewX509Identity("Org1MSP", string(cert), string(privKey)))
	if err != nil {
		fmt.Errorf(err.Error())
	}

	gw, err := gateway.Connect(gateway.WithConfig(config.FromFile("./../connection-org1.yaml")),
		gateway.WithIdentity(wallet, "admin"),
		gateway.WithDiscovery(true),
		gateway.WithCommitHandler(gateway.DefaultCommitHandlers.OrgAny))
	if err != nil {
		panic(err)
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		panic(err)
	}

	contract := network.GetContract("mycc")

	result, err := contract.SubmitTransaction("Create", "b", "50")
	if err != nil {
		panic(err)
	}

	result, err = contract.EvaluateTransaction("Read", "b")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))

	transient := make(map[string][]byte)
	transient["someTransient"] = []byte("mytransientdata")
	txn, err := contract.CreateTransaction("Update", gateway.WithTransient(transient))
	if err != nil {
		panic(err)
	}
	result, err = txn.Submit("b", "100")
	if err != nil {
		panic(err)
	}
}
