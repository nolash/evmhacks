package main

import (
	"fmt"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

var (
	balance = big.NewInt(int64(math.Pow(2, 7)))
)

func main() {

	// create key and derive address
	privkey, err := crypto.GenerateKey()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// set up database
	//	dbdir, err := ioutil.TempDir("", "sqlite-evmhack-")
	//	if err != nil {
	//		log.Error(err.Error())
	//		os.Exit(1)
	//	}
	//	defer os.RemoveAll(dbdir)
	//	db, err := ethdb.NewLDBDatabase(dbdir, 20, 4)
	//	if err != nil {
	//		log.Error(err.Error())
	//		os.Exit(1)
	//	}

	// set up backend
	//genesis := GenesisBlockForTesting(db, addr, &balance)
	auth := bind.NewKeyedTransactor(privkey)
	alloc := make(core.GenesisAlloc, 1)
	alloc[auth.From] = core.GenesisAccount{
		PrivateKey: crypto.FromECDSA(privkey),
		Balance:    balance,
	}
	sim := backends.NewSimulatedBackend(alloc)

	ctx := vm.Context{
		CanTransfer: func(state vm.StateDB, addr common.Address, amount *big.Int) bool {
			return false
		},
		Transfer: func(state vm.StateDB, laddr common.Address, raddr common.Address, amount *big.Int) {
			return
		},
		GetHash: func(uint64) common.Hash {
			return common.StringToHash("foo")
		},
	}

	e := vm.NewEVM(ctx, sim.State(), sim.Config(), vm.Config{})

	fmt.Printf("%v\n", e)
}
