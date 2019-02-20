// API for querying and updating a moodle server
//
//        api := moodle.NewSisApi("https://moodle.example.com/moodle/", "user@example.com", "a0092ba9a9f5b45cdd2f01d049595bfe91")
//
//        // Search moodle courses
//        students, _ := api.GetStudents("kim smith")
//        for _, i := range students {
//                fmt.Printf("%s\n", i.GetDisplayName())
//        }
//
//
package sisapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// API Documentation
// https://docs.moodle.org/dev/Web_service_API_functions

type SisApi struct {
	base     string
	email    string
	password string
	token    string

	log   SisLogger
	fetch LookupUrl
}

func NewSisApi(base, email, password string) *SisApi {
	return &SisApi{
		base:     base,
		email:    email,
		password: password,
		token:    "",
		log:      &NilSisLogger{},
		fetch:    &DefaultLookupUrl{},
	}
}

type Person struct {
	Uuid               string `json:",omitempty"`
	StudentNumber      string `json:",omitempty"`
	FirstName          string `json:",omitempty"`
	PreferredFirstName string `json:",omitempty"`
	MiddleName         string `json:",omitempty"`
	LastName           string `json:",omitempty"`
	Sex                string `json:",omitempty"`
	Title              string `json:",omitempty"`
	Error              string `json:",omitempty"`
	ErrorDetails       string `json:",omitempty"`
}

type SisLogger interface {
	Debug(message string, items ...interface{}) error
}

type NilSisLogger struct {
}

func (ml *NilSisLogger) Debug(message string, items ...interface{}) error {
	return nil
}

func (m *SisApi) SetLogger(l SisLogger) {
	m.log = l
}

func (m *SisApi) Authenticate() error {
	url := fmt.Sprintf("%sapi/authenticate?email=%s&password=%s", m.base, url.QueryEscape(m.email), url.QueryEscape(m.password))
	m.log.Debug("Fetch: %s", url)
	body, _, _, err := m.fetch.GetUrl(url)

	if err != nil {
		return err
	}

	type AuthResponse struct {
		Token string
		Error string
	}
	var result AuthResponse

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return errors.New("Server returned unexpected response. " + err.Error())
	}

	if result.Error != "" {
		return errors.New(result.Error)
	}
	m.token = result.Token

	return nil
}

func (m *SisApi) GetPerson(uuid string) (*Person, error) {
	url := fmt.Sprintf("%sapi/person?token=%s&uuid=%s", m.base, m.token, url.QueryEscape(uuid))
	m.log.Debug("Fetch: %s", url)
	body, _, _, err := m.fetch.GetUrl(url)

	if err != nil {
		return nil, err
	}

	var result Person

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, errors.New("Server returned unexpected response. " + err.Error())
	}

	if result.Error != "" {
		return nil, errors.New(result.Error)
	}

	return &result, nil
}

func (m *SisApi) SearchPeople(q string) ([]Person, error) {
	url := fmt.Sprintf("%sapi/student.search?token=%s&q=%s", m.base, m.token, url.QueryEscape(q))
	m.log.Debug("Fetch: %s", url)
	body, _, _, err := m.fetch.GetUrl(url)

	if err != nil {
		return nil, err
	}

	var results []Person

	if err := json.Unmarshal([]byte(body), &results); err != nil {
		return nil, errors.New("Server returned unexpected response. " + err.Error())
	}

	return results[:], nil
}

// Get moodle account matching by email address.
func (m *SisApi) GetPersonByEmail(email string) (*Person, error) {
	url := fmt.Sprintf("%sapi/person?token=%s&email=%s", m.base, m.token, url.QueryEscape(email))
	m.log.Debug("Fetch: %s", url)
	body, _, _, err := m.fetch.GetUrl(url)

	if err != nil {
		return nil, err
	}

	var result Person
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, errors.New("Server returned unexpected response. " + err.Error())
	}

	if result.Error != "" {
		return nil, errors.New(result.Error)
	}

	if result.Uuid == "" {
		return nil, nil
	}

	return &result, nil
}

func (m *SisApi) RemovePersonFromGroup(personUuid, group string, intakeYear int, intakeSemester string) error {
	url := fmt.Sprintf("%sapi/group.remove?token=%s&person=%s&group=%s&year=%d&semester=%s", m.base, m.token, personUuid, url.QueryEscape(group), intakeYear, intakeSemester)
	m.log.Debug("Fetch: %s", url)

	body, _, _, err := m.fetch.GetUrl(url)
	if err != nil {
		return err
	}

	type Response struct {
		Success string
		Error   string
	}

	var result Response
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return errors.New("Server returned unexpected response. " + err.Error())
	}

	if result.Success != "" {
		return nil
	}

	return errors.New(result.Error)
}

func (m *SisApi) AddPersonToGroup(personUuid string, group string, intakeYear int, intakeSemester string) error {
	url := fmt.Sprintf("%sapi/group.add?token=%s&person=%s&group=%s&year=%d&semester=%s", m.base, m.token, personUuid, url.QueryEscape(group), intakeYear, intakeSemester)
	m.log.Debug("Fetch: %s", url)

	body, _, _, err := m.fetch.GetUrl(url)
	if err != nil {
		return err
	}

	type Response struct {
		Success string
		Error   string
	}

	var result Response
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return errors.New("Server returned unexpected response. " + err.Error())
	}

	if result.Success != "" {
		return nil
	}

	return errors.New(result.Error)
}

func (m *SisApi) SetUrlFetcher(fetch LookupUrl) {
	m.fetch = fetch
}

type PrintSisLogger struct {
}

func (ml *PrintSisLogger) Debug(message string, items ...interface{}) error {
	fmt.Printf(message+"\n", items...)
	return nil
}
