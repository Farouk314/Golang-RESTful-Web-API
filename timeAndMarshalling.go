package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// func main() {
// 	TimeFormats()
// }

//TimeFormats example
func TimeFormats() {

	certificate := Certificate{
		ID:        "1",
		Title:     "A certificate title",
		CreatedAt: time.Now(),
		OwnerID:   "A",
		Year:      2018,
		Note:      "Note",
		Transfer: &Transfer{
			To:     "To",
			Status: "Status",
		}}

	fmt.Println("Before marshalling:")
	fmt.Printf("%+v", certificate)
	fmt.Println("======")
	b, _ := json.Marshal(certificate)
	fmt.Println("Marshaled:")
	fmt.Printf("%+v", string(b))
	fmt.Println("")
	fmt.Println("Unmarshaled:")
	json.Unmarshal(b, certificate)
	fmt.Printf("%+v \n", certificate)
	bs, _ := json.Marshal(certificate)
	fmt.Println("Marshaled again:")
	fmt.Println(string(bs))

	//time
	p := fmt.Println
	p("Time Parsing layouts must reference time: Mon Jan 2 15:04:05 MST 2006")
	p("===================")
	t := time.Now()
	p("=================================")
	p("Printing time.Now().Format(time.RFC3339):")
	p(t.Format(time.RFC3339))
	p("=================================")

	p("=================================")
	t1, e := time.Parse(
		time.RFC3339,
		"2019-02-24T03:14:15+00:00")
	p("Printing time.Parse(time.RFC3339,2019-02-24T03:14:15+00:00)")
	if e != nil {
		p(e.Error())
	}
	p(t1)
	p("=================================")

	p("=================================")
	p("Printing t.Format(15:04)")
	p(t.Format("15:04"))
	p("Printing t.Format(Mon Jan _2 15:04:05 2006)")
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p("Printing t.Format(2006-01-02T15:04:05.999999-07:00")
	p(t.Format("2006-01-02T15:04:05.999999-07:00"))
	p("=================================")
	p("=================================")
	p("Setting form:=3:04PM")
	form := "3 04 PM"
	p("Printing time.Parse(form, 8 41 PM")
	t2, e := time.Parse(form, "8 41 PM")
	p(t2)
	p("=================================")

	// p("=================================")
	// fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
	// 	t.Year(), t.Month(), t.Day(),
	// 	t.Hour(), t.Minute(), t.Second())
	// p("Setting ansic and printing time.Parse(ansic, 841PM)")
	// ansic := "Mon Jan _2 15:04:05 2006"
	// _, e = time.Parse(ansic, "8:41PM")
	// p(e)
	// p("=================================")
}
