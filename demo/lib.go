package demo

// Tomate if the resource subject.
type Tomate struct {
	Name string
}

// GetID ...
func (t Tomate) GetID() string {
	return t.Name
}

//go:generate lister vegetables_gen.go Tomate:Tomates
//go:generate jsoner json_vegetables_gen.go *Tomates:JSONTomates

// Controller of some resources.
type Controller struct {
}

// GetByID ...
func (t Controller) GetByID(id int) Tomate {
	return Tomate{}
}

// UpdateByID ...
func (t Controller) UpdateByID(id int, reqBody Tomate) Tomate {
	return Tomate{}
}

//go:generate jsoner json_controller_gen.go *Controller:JSONController
