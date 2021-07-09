package entity

// Country is a country definition struct
type Country struct {
	ID   int    `json:"id"`
	ISO2 string `json:"iso2"`
	Name string `json:"name"`
}
