package pricing

type Price struct {
	Id       int
	Amount   int64 					// in microgbp per minute
	Software string
	Tier     int
}
