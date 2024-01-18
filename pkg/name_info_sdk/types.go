package name_info_sdk

type GenderResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type NationalityResponse struct {
	Count     int       `json:"count"`
	Name      string    `json:"name"`
	Countries []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type LikelyGender struct {
	Name   string
	Gender string
}

type LikelyAge struct {
	Name string
	Age  int
}

type LikelyNationality struct {
	Name        string
	Nationality string
}
