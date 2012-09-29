package knapsack

type Cargo interface {
	Weight(index int) (w int)
	Value(index int) (v int)
	Quantity() (q int)
}

type DefaultCargo struct {
	weight []int
	value  []int
}

func (self *DefaultCargo) Weight(index int) int {
	return self.weight[index]
}

func (self *DefaultCargo) Value(index int) int {
	return self.value[index]
}

func (self *DefaultCargo) Quantity() int {
	return len(self.weight)
}

func NewDefaultCargo(w, v []int) *DefaultCargo {
	l := func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}(len(w), len(v))
	return &DefaultCargo{w[:l], v[:l]}
}

func KnapValue(budget int, cargo Cargo) int {
	kap := make([]int, budget+1)
	kap[0] = 0
	q := cargo.Quantity()
	for w := 1; w < budget+1; w++ {
		kap[w] = func() int {
			max := 0
			for i := 0; i < q; i++ {
				if cargo.Weight(i) <= w {
					if kap[w-cargo.Weight(i)]+cargo.Value(i) > max {
						max = kap[w-cargo.Weight(i)] + cargo.Value(i)
					}
				}
			}
			return max
		}()
	}
	return kap[budget]
}

func KnapValueNoRep(budget int, cargo Cargo) int {
	kap := make([][]int, budget+1)
	q := cargo.Quantity() + 1
	for i := 0; i < budget+1; i++ {
		kap[i] = make([]int, q)
	}
	for i := 0; i < q; i++ {
		kap[0][i] = 0
	}
	for i := 0; i < budget+1; i++ {
		kap[i][0] = 0
	}
	for i := 1; i < q; i++ {
		for w := 1; w < budget+1; w++ {
			if cargo.Weight(i-1) > w {
				kap[w][i] = kap[w][i-1]
			} else {
				kap[w][i] = func(a, b int) int {
					if a > b {
						return a
					}
					return b
				}(kap[w-cargo.Weight(i-1)][i-1]+cargo.Value(i-1), kap[w][i-1])
			}
		}
	}
	return kap[budget][q-1]
}
