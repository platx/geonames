package value

type AdminCode struct {
	// First fipscode (subject to change to iso code), see exceptions below, see file admin1Codes.txt
	// for display names of this code; max length is 20 characters.
	First string `url:"adminCode1"`
	// Second level administrative division, a county in the US, see file admin2Codes.txt; max length is 80 characters.
	Second string `url:"adminCode2"`
	// Third level administrative division, max length is 20 characters.
	Third string `url:"adminCode3"`
	// Fourth level administrative division, max length is 20 characters.
	Fourth string `url:"adminCode4"`
	// Fifth level administrative division, max length is 20 characters.
	Fifth string `url:"adminCode5"`
}

type AdminDivision struct {
	Code string
	Name string
}

type AdminDivisions struct {
	First  AdminDivision
	Second AdminDivision
	Third  AdminDivision
	Fourth AdminDivision
	Fifth  AdminDivision
}
