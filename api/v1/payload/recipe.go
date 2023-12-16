package payload

// Recipe is the payload for the recipe entity
type Recipe struct {
	Name        *string      `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
}
