package main

import (
	"VSRT-Lang/internal"
	"VSRT-Lang/internal/database/postgres/migrations"
	"VSRT-Lang/internal/database/postgres/session_repository"
	"VSRT-Lang/internal/session"
	"database/sql"
	"fmt"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	//translator := net_translator.Translator{
	//	Client: http.Client{
	//		Transport:     nil,
	//		CheckRedirect: nil,
	//		Jar:           nil,
	//		Timeout:       0,
	//	},
	//	Backend: "http://localhost:8080/api/translator",
	//}

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=vsrt_lang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err := migrations.Run(db); err != nil {
		panic(err)
	}
	repository := session_repository.New(db)

	testSession := session.Session{
		User:       "",
		Name:       "test",
		Records:    nil,
		Translator: internal.MockTranslator{},
	}

	err = repository.Save(&testSession)
	if err != nil {
		panic(err)
	}

	record := testSession.SaveRecord(
		"simple sentence to show work of app",
		"simple sentence",
	)

	printRecord(record)

	err = repository.Save(&testSession)

	taken, err := repository.Take(0)
	if err != nil {
		panic(err)
	}

	fmt.Print("\n ========== FROM DB ========== \n\n")

	for _, takeRecord := range taken.Records {
		printRecord(takeRecord)
	}

}

func printRecord(takeRecord session.Record) {
	fmt.Printf("Phrase:\t\t%s\n", takeRecord.Phrase)
	fmt.Printf("Translations:\t%s\n", takeRecord.Translations)
	fmt.Printf("BaseForm:\t%s\n", takeRecord.BaseForm)
	fmt.Printf("Synonyms:\t%s\n", takeRecord.Synonyms)
	fmt.Printf("Antonyms:\t%s\n", takeRecord.Antonyms)
	fmt.Printf("Contexts:\t%s\n", takeRecord.Contexts)
}
