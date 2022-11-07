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
		fmt.Println(compileErr.Error())
		panic("compile error")
	}
	r1csFile, openErr := os.Create("./web/hash.r1cs")
	if openErr != nil {
		fmt.Println(openErr.Error())
		panic("r1cs open error")
	}
	_, writeErr := r1cs.WriteTo(r1csFile)
	if writeErr != nil {
		fmt.Println(writeErr.Error())
		panic("r1cs write error")
	}
	fmt.Println("Generating Proving Key and Verifying Key")
	pk, vk, setupErr := groth16.Setup(r1cs)
	if setupErr != nil {
		fmt.Println(setupErr.Error())
		panic("groth16 setup error")
	}
	pkFile, pkOpenErr := os.Create("./web/hash.pkey")
	if pkOpenErr != nil {
		fmt.Println(pkOpenErr.Error())
		panic("pk open error")
	}
	vkFile, vkOpenErr := os.Create("./web/hash.vkey")
	if vkOpenErr != nil {
		fmt.Println(vkOpenErr.Error())
		panic("vk open error")
	}
	_, pkWriteErr := pk.WriteTo(pkFile)
	if pkWriteErr != nil {
		fmt.Println(pkWriteErr.Error())
		panic("pk write error")
	}
	_, vkWriteErr := vk.WriteTo(vkFile)
	if vkWriteErr != nil {
		fmt.Println(vkWriteErr.Error())
		panic("vk write error")
	}
	assignment := &circuits.HashCircuit{
		Key:  0,
		X:    1764,
		Hash: "15893827533473716138720882070731822975159228540693753428689375377280130954696",
	}
	witness, witnessErr := frontend.NewWitness(assignment, ecc.BN254)
	if witnessErr != nil {
		fmt.Println(witnessErr.Error())
		panic("witness error")
	}
	publicWitness, pubWitnessErr := witness.Public()
	if pubWitnessErr != nil {
		fmt.Println(pubWitnessErr.Error())
		panic("pub witness error")
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
	fmt.Println("Generating Solidity")
	verifier, verifierOpenErr := os.Create("./web/verifier.sol")
	if verifierOpenErr != nil {
		fmt.Println(verifierOpenErr.Error())
		panic("verifier open error")
	}
	verifierWriteErr := vk.ExportSolidity(verifier)
	if verifierWriteErr != nil {
		fmt.Println(verifierWriteErr.Error())
		panic("verifier write error")
	}
}
