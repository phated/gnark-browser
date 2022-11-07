package circuits

import (
	"math"

	"github.com/consensys/gnark/frontend"
)

type InitCircuit struct {
	Radius        frontend.Variable `gnark:"r,public"`
	PlanetHashKey frontend.Variable `gnark:"PLANETHASH_KEY,public"`
	// SpaceTypeKey  frontend.Variable `gnark:"SPACETYPE_KEY,public"`
	// Scale         frontend.Variable `gnark:"scale,public"` // must be power of 2 at most 16384 so that DENOMINATOR works
	// MirrorX       frontend.Variable `gnark:"xMirror,public"`
	// MirrorY       frontend.Variable `gnark:"yMirror,public"`
	X          frontend.Variable `gnark:"x,secret"`
	Y          frontend.Variable `gnark:"y,secret"`
	LocationId frontend.Variable
	// Perlin     frontend.Variable `gnark:"perl,public"`
}

// Define declares the circuit"s constraints
func (circuit *InitCircuit) Define(api frontend.API) error {
	/* check abs(x), abs(y) <= 2^31 */
	api.AssertIsLessOrEqual(api.Add(circuit.X, 1<<31), math.MaxUint32)
	api.AssertIsLessOrEqual(api.Add(circuit.Y, 1<<31), math.MaxUint32)

	/* check x^2 + y^2 < r^2 */
	xSq := api.Mul(circuit.X, circuit.X)
	ySq := api.Mul(circuit.Y, circuit.Y)
	rSq := api.Mul(circuit.Radius, circuit.Radius)
	api.AssertIsEqual(api.Cmp(api.Add(xSq, ySq), rSq), -1)

	/* check x^2 + y^2 > 0.98 * r^2 */
	/* equivalently 100 * (x^2 + y^2) > 98 * r^2 */
	api.AssertIsEqual(api.Cmp(api.Mul(100, api.Add(xSq, ySq)), api.Mul(98, rSq)), 1)

	/* check MiMCSponge(x,y) = pub */
	/*
	   220 = 2 * ceil(log_5 p), as specified by mimc paper, where
	   p = 21888242871839275222246405745257275088548364400416034343698204186575808495617
	*/
	mimc := NewMiMCSponge(api, 1, 220, circuit.PlanetHashKey)

	inputs := []frontend.Variable{circuit.X}
	out := mimc.hash(inputs)

	// Gnark doesn't seem to have a way to do outputs! WTF!
	circuit.LocationId = out[0]

	return nil
}
