# Golang API for SIS

This is a simple library that wraps the SIS Json  Web Service API. 
It is part of a project that is currently in use in a live production environment, and
is under active development.

No warranty is given or implied. Use at your own risk.


## Example usage

	// Setup
	api := sis.NewSisApi("https://site.example.com/", "Fi04DbaAa9f5b45cdd2f01d5f1", l)

	// Search students
	students, _ := api.SearchStudents("John Smith")
	for _, i := range students {
		fmt.Printf("%s\n", i.StudentNumber)
	}


