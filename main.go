package main

import (
	"account-merging/account"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

var (
	filename string
	writeTo  bool
)

const (
	resultJsonFilename = "./result.json"
)

func init() {
	flag.StringVar(&filename, "file", "./data/accounts.json", "the filename of the JSON file containing un-merged accounts (defaults to './data/accounts.json')")
	flag.BoolVar(&writeTo, "out", false, "write merge result to file ('./results.json') ")
	flag.Parse()
}

func main() {
	log.Info().Msg("starting application")

	log.Info().Msgf("filename: %s", filename)
	// parse accounts json file
	accounts, err := account.ReadInAccounts(filename)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read in accounts json")
	}

	// merge the accounts
	persons := account.MergeAccounts(accounts)

	log.Info().Int("count", len(persons)).Msgf("accounts successfully merged")

	res, err := json.MarshalIndent(persons, "", "    ")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal indent merged accounts")
	}

	if writeTo {
		file, err := os.Create(resultJsonFilename)
		if err != nil {
			log.Error().Err(err).Msgf("failed to create file %s", writeTo)
		}

		defer file.Close()

		if _, err := file.Write(res); err != nil {
			log.Error().Err(err).Msgf("failed to write to file %s", writeTo)
		}
	}

	fmt.Printf("%s\n", res)
}
