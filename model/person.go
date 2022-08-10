package model

import "strings"

// NOTE: the instructions said the Applications field was a string array value, but the accounts.json has the values as ints so I just went with that
type Person struct {
	Applications []int    `json:"applications"`
	Emails       []string `json:"emails"`
	Name         string   `json:"name"`
}

func (p *Person) AddEmailIfNotExists(newEmail string) {
	for _, email := range p.Emails {
		if strings.EqualFold(email, newEmail) {
			return
		}
	}

	p.Emails = append(p.Emails, newEmail)
}

func (p *Person) AddEmailsIfNotExists(newEmails []string) {
	for _, email := range newEmails {
		p.AddEmailIfNotExists(email)
	}
}

func (p *Person) AddApplicationIfNotExists(newApplication int) {
	for _, application := range p.Applications {
		if newApplication == application {
			return
		}
	}

	p.Applications = append(p.Applications, newApplication)
}
