package model

// NOTE: the instructions said the Application field was a string value, but the accounts.json has the values as ints so I just went with that
type Account struct {
	Application int      `json:"application"`
	Emails      []string `json:"emails"`
	Name        string   `json:"name"`
}
