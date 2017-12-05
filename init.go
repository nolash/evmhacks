package main

import (
	"bytes"
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
	// bytecode is generated from this assembly code
	// (push value 10 (dec, int32, big endian) to offset 20 (dec)
	// store to memory
	// retrieve 1 byte at offset 20 + 32 - 1
	// return value (10)
	//	main:
	//		push 10
	//		push 20
	//		mstore
	//		push 1
	//		push 51
	//		return
	bytecode = common.FromHex("0x5b600a60145260016033f3")
	balance  = big.NewInt(int64(math.Pow(2, 7)))
)

func main() {

	// create key and derive address
	privkey, err := crypto.GenerateKey()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	addr := crypto.PubkeyToAddress(privkey.PublicKey)

	// set up backend
	auth := bind.NewKeyedTransactor(privkey)
	alloc := make(core.GenesisAlloc, 1)
	alloc[auth.From] = core.GenesisAccount{
		PrivateKey: crypto.FromECDSA(privkey),
		Balance:    balance,
	}
	sim := backends.NewSimulatedBackend(alloc)

	// create the evm interpreter
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
	evm := vm.NewEVM(ctx, sim.State(), sim.Config(), vm.Config{})
	ipr := vm.NewInterpreter(evm, vm.Config{})

	// set up and run contract
	ct := vm.NewContract(vm.AccountRef(addr), vm.AccountRef(addr), big.NewInt(0), 2000000)
	ct.SetCallCode(&addr, crypto.Keccak256Hash(bytecode), bytecode)
	r, err := ipr.Run(0, ct, []byte{})
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	// check result
	if !bytes.Equal(r, []byte{0x0a}) {
		fmt.Fprintf(os.Stderr, "expected [0x0a], got %v", r)
	}
}
