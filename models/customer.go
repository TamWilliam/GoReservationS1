package models

type Customer struct {
	ID        int
	Role      int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type TemplateData struct {
	StringData map[string]string
	IntData    map[string]int
}
