package sisapi

import (
	"fmt"
	"os"
	"testing"
)

func TestPersonModule(t *testing.T) {

	api := NewSisApi(requireEnv("SIS_URL", t), requireEnv("SIS_USER", t), requireEnv("SIS_PASSWORD", t))
	api.SetLogger(&PrintSisLogger{})

	err := api.Authenticate()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	students, err := api.SearchPeople("kim smith")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if len(students) == 0 {
		fmt.Printf("people not returned\n")
		return
	}
	for _, student := range students {
		fmt.Printf("%s %s %s %s\n", student.StudentNumber, student.FirstName, student.LastName, student.Sex)
	}

	fmt.Printf("Found %d students\n", len(students))

}

func requireEnv(name string, t *testing.T) string {
	value := os.Getenv(name)
	if value == "" {
		t.Fatalf(fmt.Sprintf("Environment variable required: %s\n", name))
	}
	return value
}
