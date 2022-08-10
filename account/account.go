package account

import (
	"account-merging/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ReadInAccounts will read the file from the provided filename and unmarshal the accounts json.
func ReadInAccounts(filename string) ([]model.Account, error) {
	if filename = strings.TrimSpace(filename); filename == "" {
		return nil, errors.New("cannot read file with an empty filename")
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s'; %w", filename, err)
	}

	accounts := []model.Account{}

	if err := json.NewDecoder(bytes.NewReader(content)).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounts json; %w", err)
	}

	return accounts, nil
}

// MergeAccounts will merge the provided accounts based on the relationships between each account's emails.
func MergeAccounts(accounts []model.Account) []model.Person {
	// will contain the relationships between all emails
	emailMapping := make(map[string][]string)
	// will contain the relationships of one email to the accounts it is present in
	// since the accounts don't have any unique identifiers, the indices of the provided slice will be used
	emailToAccountIdx := make(map[string][]int)

	for i, account := range accounts {
		if len(account.Emails) == 0 {
			continue
		}

		// track the account where the root email is present
		rootEmail := account.Emails[0]
		emailToAccountIdx[rootEmail] = append(emailToAccountIdx[rootEmail], i)

		// build the mapping of a emails in the accounts
		for j := 1; j < len(account.Emails); j++ {
			email := account.Emails[j]
			// all subsequent emails will have a relationship to the root email and vice-versa
			emailMapping[rootEmail] = append(emailMapping[rootEmail], email)
			emailMapping[email] = append(emailMapping[email], rootEmail)
			// track the account where the current email is present
			emailToAccountIdx[email] = append(emailToAccountIdx[email], i)
		}
	}

	// now, build email groups from the existing email mappings
	// each group will be representative of a single person's account
	seenEmails := make(map[string]bool)
	emailGroups := make([][]string, 0)

	for rootEmail := range emailToAccountIdx {
		if _, seen := seenEmails[rootEmail]; seen {
			continue // root email has already been added to a group
		}

		// each linked email needs to be traversed, this will hold those linked emails
		nextEmailQueue := make([]string, 0)
		// the list of one person's emails
		emailGroup := make([]string, 0)

		nextEmailQueue = append(nextEmailQueue, rootEmail)

		// check every email that has a relationship to the rootEmail and add it to the emailGroup if it hasn't been seen yet
		for len(nextEmailQueue) > 0 {
			nextEmail := nextEmailQueue[len(nextEmailQueue)-1]
			nextEmailQueue = nextEmailQueue[:len(nextEmailQueue)-1]

			emailGroup = append(emailGroup, nextEmail)

			for _, linkedEmail := range emailMapping[nextEmail] {
				if _, seen := seenEmails[linkedEmail]; !seen {
					nextEmailQueue = append(nextEmailQueue, linkedEmail)
					seenEmails[linkedEmail] = true
				}
			}
		}

		emailGroups = append(emailGroups, emailGroup)
	}

	// construct the person data using the previously created email groups
	seenAccounts := make(map[int]bool)
	persons := make([]model.Person, len(emailGroups))

	for i, eg := range emailGroups {
		p := model.Person{}

		for _, email := range eg {
			// Get the accounts related to this email
			accountIndices := emailToAccountIdx[email]

			for j, idx := range accountIndices {
				if _, seen := seenAccounts[idx]; !seen {
					seenAccounts[idx] = true
					// the name is arbitrary, so just picking the first one
					if j == 0 {
						p.Name = accounts[idx].Name
					}

					p.AddApplicationIfNotExists(accounts[idx].Application)
					p.AddEmailsIfNotExists(accounts[idx].Emails)
				}
			}
		}

		sort.Strings(p.Emails)
		sort.Ints(p.Applications)

		persons[i] = p
	}

	return persons
}
