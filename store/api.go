package store

// Todo defines the structure for an API product
// swagger:model
type Todo struct {
	// the id for the product
	//
	// required: false
	// min: 1
	Id int // Unique identifier for the product
	// the description for this todo
	//
	// required: false
	// max length: 10000
	Description string
	// the completed for this todo
	//
	// required: false
	Completed bool
}
