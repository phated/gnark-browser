package main

import (
	"fmt"
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
	r1csRes, r1csResErr := http.DefaultClient.Get(os.Args[0])
	if r1csResErr != nil {
		fmt.Println(r1csResErr.Error())
		panic("r1cs response error")
	}
	pkRes, pkResErr := http.DefaultClient.Get(os.Args[1])
	if pkResErr != nil {
		fmt.Println(pkResErr.Error())
		panic("pk response error")
	}
	vkRes, vkResErr := http.DefaultClient.Get(os.Args[2])
	if vkResErr != nil {
		fmt.Println(vkResErr.Error())
		panic("vk response error")
	}
	r1cs := groth16.NewCS(ecc.BN254)
	_, r1csReadErr := r1cs.ReadFrom(r1csRes.Body)
	if r1csReadErr != nil {
		fmt.Println(r1csReadErr.Error())
		panic("r1cs read error")
	}
	pk := groth16.NewProvingKey(ecc.BN254)
	_, pkReadErr := pk.ReadFrom(pkRes.Body)
	if pkReadErr != nil {
		fmt.Println(pkReadErr.Error())
		panic("pk read error")
	}
	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, vkReadErr := vk.ReadFrom(vkRes.Body)
	if vkReadErr != nil {
		fmt.Println(vkReadErr.Error())
		panic("vk read error")
	}
	// witness
	assignment := &HashCircuit{
		Key:  0,
		X:    1764,
		Hash: "15893827533473716138720882070731822975159228540693753428689375377280130954696",
	}
	// This crashes with a nil reference panic
	// var assignment HashCircuit
	// jsonErr := json.Unmarshal([]byte(os.Args[3]), &assignment)
	// if jsonErr != nil {
	// 	fmt.Println(jsonErr.Error())
	// 	panic("failed to unmarshal json")
	// }
	// fmt.Printf("%+v", assignment)
	witness, witnessErr := frontend.NewWitness(assignment, ecc.BN254)
	if witnessErr != nil {
		fmt.Println(witnessErr.Error())
		panic("witness error")
	}
	publicWitness, pubWitnessErr := witness.Public()
	if pubWitnessErr != nil {
		fmt.Println(pubWitnessErr.Error())
		panic("pubc witness error")
	}
	proof, proveErr := groth16.Prove(r1cs, pk, witness)
	if proveErr != nil {
		fmt.Println(proveErr.Error())
		panic("fail prove")
	}
	verifyErr := groth16.Verify(proof, vk, publicWitness)
	if verifyErr != nil {
		fmt.Println(verifyErr.Error())
		panic("Not verified")
	}
}
