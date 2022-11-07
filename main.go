package main

import (
	"net/http"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"

	// Forcibly load the `init()` function for bit hints
	_ "github.com/consensys/gnark/std/math/bits"
)

type HashCircuit struct {
	Key  frontend.Variable `json:"key" gnark:"key,public"`
	X    frontend.Variable `json:"x" gnark:"x,public"`
	Hash frontend.Variable `json:"hash"`
}

// Define declares the circuit"s constraints
func (circuit *HashCircuit) Define(api frontend.API) error {
	return nil
}

func main() {
	res, resErr := http.DefaultClient.Get(os.Args[0])
	if resErr != nil {
		panic("response error")
	}
	pkRes, pkResErr := http.DefaultClient.Get(os.Args[1])
	if pkResErr != nil {
		panic("response error")
	}
	vkRes, vkResErr := http.DefaultClient.Get(os.Args[2])
	if vkResErr != nil {
		panic("response error")
	}
	r1cs := groth16.NewCS(ecc.BN254)
	r1cs.ReadFrom(res.Body)
	pk := groth16.NewProvingKey(ecc.BN254)
	pk.ReadFrom(pkRes.Body)
	vk := groth16.NewVerifyingKey(ecc.BN254)
	vk.ReadFrom(vkRes.Body)
	// witness
	assignment := &HashCircuit{
		Key:  0,
		X:    1764,
		Hash: "15893827533473716138720882070731822975159228540693753428689375377280130954696",
	}
	// assignment := &HashCircuit{}
	// jsonErr := json.Unmarshal([]byte(os.Args[3]), assignment)
	// if jsonErr != nil {
	// 	fmt.Println(jsonErr.Error())
	// 	panic("failed to unmarshal json")
	// }
	// fmt.Printf("%+v", assignment)
	witness, _ := frontend.NewWitness(assignment, ecc.BN254)
	publicWitness, _ := witness.Public()
	proof, err1 := groth16.Prove(r1cs, pk, witness)
	if err1 != nil {
		panic("fail prove")
	}
	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic("Not verified")
	}
}
