package entity

// Item represents an item in an order.
type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number" gorm:"foreignKey:TrackNumber"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
