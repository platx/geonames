package value

type ContinentCode string

const (
	ContinentCodeAfrica       ContinentCode = "AF"
	ContinentCodeAsia         ContinentCode = "AS"
	ContinentCodeEurope       ContinentCode = "EU"
	ContinentCodeNorthAmerica ContinentCode = "NA"
	ContinentCodeOceania      ContinentCode = "OC"
	ContinentCodeSouthAmerica ContinentCode = "SA"
	ContinentCodeAntarctica   ContinentCode = "AN"
)

type Continent struct {
	Code ContinentCode
	Name string
}
