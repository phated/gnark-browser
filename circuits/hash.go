package circuits

import (
	"math"

	"github.com/consensys/gnark/frontend"
)

type HashCircuit struct {
	Key  frontend.Variable `json:"key" gnark:"key,public"`
	X    frontend.Variable `json:"x" gnark:"x,public"`
	Hash frontend.Variable `json:"hash"`
}

// Define declares the circuit"s constraints
func (circuit *HashCircuit) Define(api frontend.API) error {
	api.AssertIsLessOrEqual(circuit.X, math.MaxInt32)

	mimc := NewMiMCSponge(api, 1, 220, circuit.Key)

	inputs := []frontend.Variable{circuit.X}
	out := mimc.hash(inputs)
	api.AssertIsEqual(out[0], circuit.Hash)

	return nil
}
