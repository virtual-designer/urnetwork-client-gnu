package core

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"log"
	"github.com/urnetwork/connect"
)

type APINetwork struct {
	Jwt string `json:"by_jwt"`
	Name string `json:"name"`
}

type APIVerificationRequired struct {
	UserAuth string `json:"user_auth,omitempty"`
}

type APILoginError struct {
	Message string `json:"message"`
}

type APILoginWithPasswordResponse struct {
	Network *APINetwork `json:"network,omitempty"`
	VerificationRequired *APIVerificationRequired `json:"verification_required,omitempty"`
	Error *APILoginError `json:"error,omitempty"`
}

type APILocationGroupResult struct {
	LocationGroupId *connect.Id `json:"location_group_id"`
	Name            string      `json:"name"`
	ProviderCount   int         `json:"provider_count,omitempty"`
	Promoted        bool        `json:"promoted,omitempty"`
	MatchDistance   int         `json:"match_distance,omitempty"`
}

type APILocationResult struct {
	LocationId   *connect.Id  `json:"location_id"`
	LocationType string       `json:"location_type"`
	Name         string       `json:"name"`
	// FIXME
	City string `json:"city,omitempty"`
	// FIXME
	Region string `json:"region,omitempty"`
	// FIXME
	Country           string      `json:"country,omitempty"`
	CountryCode       string      `json:"country_code,omitempty"`
	CityLocationId    *connect.Id `json:"city_location_id,omitempty"`
	RegionLocationId  *connect.Id `json:"region_location_id,omitempty"`
	CountryLocationId *connect.Id `json:"country_location_id,omitempty"`
	ProviderCount     int         `json:"provider_count,omitempty"`
	MatchDistance     int         `json:"match_distance,omitempty"`
}

type APILocationDeviceResult struct {
	ClientId   *connect.Id `json:"client_id"`
	DeviceName string      `json:"device_name"`
}

type APIGetLocationsResponse struct {
	Specs []*connect.ProviderSpec `json:"specs"`
	// this includes groups that show up in the location results
	// all `ProviderCount` are from inside the location results
	// groups are suggestions that can be used to broaden the search
	Groups []*APILocationGroupResult `json:"groups"`
	// this includes all parent locations that show up in the location results
	// every `CityId`, `RegionId`, `CountryId` will have an entry
	Locations []*APILocationResult       `json:"locations"`
	Devices   []*APILocationDeviceResult `json:"devices"`
}

const API_BASE_URL = "https://api.bringyour.com"
var client = &http.Client {}

func makeRoute(path string) string {
	return API_BASE_URL + path
}

func AttemptLoginWithPassword(email string, password string) (*APILoginWithPasswordResponse, error) {
	url := makeRoute("/auth/login-with-password")
	data := map[string]string {
		"user_auth": email,
		"password": password,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var result APILoginWithPasswordResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetLocations(jwt string) (*APIGetLocationsResponse, error) {
	url := makeRoute("/network/provider-locations")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "URnetworkClientGNU/1.0.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var result APIGetLocationsResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
