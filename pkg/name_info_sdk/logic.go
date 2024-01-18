package name_info_sdk

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	genderHost      = "https://api.genderize.io"
	ageHost         = "https://api.agify.io"
	nationalityHost = "https://api.nationalize.io"
)

type RequestInfo int

const (
	Age RequestInfo = iota
	Gender
	Nationality
)

var (
	errCodeNotOK = errors.New("Something went wrong on the foreign API side")
)

type INameInfo interface {
	GetGenderInfoByName(name string) (*LikelyGender, error)
	GetAgeInfoByName(name string) (*LikelyAge, error)
	GetLikelyNationalityInfoByName(name string) (*LikelyNationality, error)
}

type NameInfo struct {
	apiKey string
	client *http.Client
}

func NewNameInfo(apiKey string) INameInfo {
	return &NameInfo{
		apiKey: apiKey,
		client: http.DefaultClient,
	}
}

func (n *NameInfo) buildURL(reqInfo RequestInfo, name string) url.URL {

	var parsedUrl *url.URL
	switch reqInfo {
	case Age:
		parsedUrl, _ = url.ParseRequestURI(ageHost)
	case Gender:
		parsedUrl, _ = url.ParseRequestURI(genderHost)
	case Nationality:
		parsedUrl, _ = url.ParseRequestURI(nationalityHost)
	}

	q := parsedUrl.Query()
	if n.apiKey != "" {
		q.Set("apikey", n.apiKey)
	}
	q.Set("name", name)
	parsedUrl.RawQuery = q.Encode()

	return *parsedUrl
}

func (n *NameInfo) DoHttpRequest(url string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, errCodeNotOK
	}

	return resp, nil
}

func (n *NameInfo) GetGenderInfoByName(name string) (*LikelyGender, error) {

	url := n.buildURL(Gender, name)

	resp, err := n.DoHttpRequest(url.String())
	if err != nil {
		return nil, err
	}

	var genderInfo GenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&genderInfo); err != nil {
		return nil, err
	}

	return &LikelyGender{
		Name:   name,
		Gender: genderInfo.Gender,
	}, nil
}

func (n *NameInfo) GetAgeInfoByName(name string) (*LikelyAge, error) {

	url := n.buildURL(Age, name)

	resp, err := n.DoHttpRequest(url.String())
	if err != nil {
		return nil, err
	}

	var ageInfo AgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ageInfo); err != nil {
		return nil, err
	}

	return &LikelyAge{
		Name: name,
		Age:  ageInfo.Age,
	}, nil
}

func (n *NameInfo) GetLikelyNationalityInfoByName(name string) (*LikelyNationality, error) {

	url := n.buildURL(Nationality, name)

	resp, err := n.DoHttpRequest(url.String())
	if err != nil {
		return nil, err
	}

	var nationalityInfo NationalityResponse
	if err := json.NewDecoder(resp.Body).Decode(&nationalityInfo); err != nil {
		return nil, err
	}

	if len(nationalityInfo.Countries) == 0 {
		return &LikelyNationality{
			Name:        name,
			Nationality: "",
		}, nil
	}

	var likelyNationality string
	var highestProb float64 = .0
	for _, v := range nationalityInfo.Countries {
		if v.Probability >= highestProb {
			highestProb = v.Probability
			likelyNationality = v.CountryId
		}
	}

	return &LikelyNationality{
		Name:        name,
		Nationality: likelyNationality,
	}, nil
}
