package domain

type Team struct {
	Name string
}

type TeamWithUsers struct {
	TeamName string
	Members  []User
}
