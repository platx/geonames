package value

type Cities string

const (
	// Cities500 all cities with a population > 500 or seats of adm div down to PPLA4 (ca 185.000).
	Cities500 = "cities500"
	// Cities1000 all cities with a population > 1000 or seats of adm div down to PPLA3 (ca 130.000).
	Cities1000 = "cities1000"
	// Cities5000 all cities with a population > 5000 or PPLA (ca 50.000).
	Cities5000 = "cities5000"
	// Cities15000 all cities with a population > 15000 or capitals (ca 25.000).
	Cities15000 = "cities15000"
)
