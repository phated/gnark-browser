package main

import (
	"fmt"
	"gnark-example/circuits"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	fmt.Println("Compiling Circuit")
	var circuit circuits.HashCircuit
	r1cs, compileErr := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if compileErr != nil {
		panic("compile error")
	}
	r1csFile, openErr := os.Create("./web/hash.r1cs")
	if openErr != nil {
		panic("open error")
	}
	_, writeErr := r1cs.WriteTo(r1csFile)
	if writeErr != nil {
		panic("write error")
	}
	fmt.Println("Generating Proving Key and Verifying Key")
	pk, vk, setupErr := groth16.Setup(r1cs)
	if setupErr != nil {
		panic("setup error")
	}
	pkFile, pkOpenErr := os.Create("./web/hash.pkey")
	if pkOpenErr != nil {
		panic("pk open error")
	}
	vkFile, vkOpenErr := os.Create("./web/hash.vkey")
	if vkOpenErr != nil {
		panic("vk open error")
	}
	_, pkWriteErr := pk.WriteTo(pkFile)
	if pkWriteErr != nil {
		panic("pk write error")
	}
	_, vkWriteErr := vk.WriteTo(vkFile)
	if vkWriteErr != nil {
		panic("vk write error")
	}
	assignment := &circuits.HashCircuit{
		Key:  0,
		X:    1764,
		Hash: "15893827533473716138720882070731822975159228540693753428689375377280130954696",
	}
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
	fmt.Println("Generating Solidity")
	verifier, verifierOpenErr := os.Create("./web/verifier.sol")
	if verifierOpenErr != nil {
		panic("verifier error")
	}
	verifierWriteErr := vk.ExportSolidity(verifier)
	if verifierWriteErr != nil {
		panic("verifier write error")
	}
}
