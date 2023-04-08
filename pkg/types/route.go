package types

type Node struct {
	Pair *Pair
	
	//this means 
	// if path[0].token1 == path[1].token0 => inverse[0] = false
	// if path[0].token1 == path[1].token1 => inverse[0] = true
	// inverse length is len(path) -1, because it's the same as number of hops
	IsInverse bool
}
type Route struct {
	Path      []*Node	
}
