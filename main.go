package main

import (
	"fmt"
	"net/url"
	"time"
)

// Link struct
type Link struct {
	ID        int       `json:"id,string,omitempty"`
	Title     string    `json:"title,string"`
	URL       url.URL   `json:"url,string"`
	CreatedAt time.Time `json:"createdat,string"`
}

// Main
func main() {
	a := App{}
	a.Initialise()
	fmt.Println("Running...")
	a.Run(":8000")

	// InitInMemoryData()

	// s := &http.Server{
	// 	Addr:    ":8000",
	// 	Handler: InitCors(InitHandlers),
	// }

	// log.Fatal(s.ListenAndServe())
	// time.Now().Date()

}

// UnmarshalJSON receiver func
// func (c *Certificate) UnmarshalJSON(j []byte) error {
// 	var rawStrings map[string]string
// 	fmt.Println("UnmarshalJSON Custom called..")
// 	err := json.Unmarshal(j, &rawStrings)
// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range rawStrings {
// 		fmt.Println("At key: " + k + "and value: " + v)
// 		if strings.ToLower(k) == "id" {
// 			fmt.Println("id")
// 			c.ID = v
// 		}
// 		if strings.ToLower(k) == "title" {
// 			fmt.Println("title")
// 			c.Title = v
// 		}
// 		if strings.ToLower(k) == "createdAt" {
// 			fmt.Println("createdAt")
// 			t, err := time.Parse(time.RFC3339, v)
// 			if err != nil {
// 				return err
// 			}
// 			c.CreatedAt = t
// 		}
// 		if strings.ToLower(k) == "ownerId" {
// 			fmt.Println("ownerId")
// 			c.OwnerID = v
// 		}
// 		if strings.ToLower(k) == "year" {
// 			fmt.Println("year")
// 			c.Year, err = strconv.Atoi(v)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if strings.ToLower(k) == "note" {
// 			fmt.Println("note")
// 			c.Note = "I CHANGED IT"
// 		}
// 		// if strings.ToLower(k) == "transfer" {
// 		// 	c.Transfer.To = v //
// 		// 	c.Transfer.Status = v
// 		// }

// 	}

// 	return nil
// }

// // MarshalJSON Certificate
// func (c Certificate) MarshalJSON() ([]byte, error) {
// 	basicCertificate := struct {
// 		ID        string `json:"id"`
// 		Title     string `json:"title"`
// 		CreatedAt string `json:"createdAt"`
// 		OwnerID   string `json:"ownerId"`
// 		Year      int    `json:"year"`
// 		Note      string `json:"note"`
// 		Transfer  string `json:"transfer"`
// 	}{
// 		ID:        c.ID,
// 		Title:     c.Title,
// 		CreatedAt: c.CreatedAt.Format(time.RFC3339),
// 		OwnerID:   c.OwnerID,
// 		Year:      c.Year,
// 		Note:      c.Note,
// 		Transfer:  "To",
// 	}

// 	return json.Marshal(basicCertificate)
// }

// // Main applicatiton entry point
// func main() {
// 	// a := gruckful.App{}
// 	// a.Initialise()
// 	// a.Run(":8000")

// 	var link Link

// 	basicLink := Link{
// 		ID        int    `json:"id"`
// 		Title     string `json:"title"`
// 		URL       string `json:"url"`
// 		CreatedAt string `json:"createdat"`
// 	}{
// 		ID:        1,
// 		Title:     "Title",
// 		URL:       "http://localhost:8000",
// 		CreatedAt: time.Now().Format(time.RFC3339),
// 	}

// 	bs, _ := link.MarshalJSON()
// 	basicLink.UnmarshalJSON(bs)

// 	// fakePostManData := struct {
// 	// 	ID        int    `json:"id"`
// 	// 	Title     string `json:"title"`
// 	// 	URL       string `json:"url"`
// 	// 	CreatedAt string `json:"createdat"`
// 	// }{
// 	// 	ID:        1,
// 	// 	Title:     "Title",
// 	// 	URL:       "http://localhost:800",
// 	// 	CreatedAt: time.Now().Format(time.ANSIC),
// 	// }

// 	// bytesss, _ = json.Marshal(fakePostManData)
// 	// link.UnmarshalJSON(bytesss)

// 	// Convert postman data to Link struct
// 	// b, err := json.Marshal(fakePostManData)
// 	// if err != nil {
// 	// 	fmt.Println("errororororor")
// 	// }
// 	// link.UnmarshalJSON(b) // Need this be
// 	// fmt.Println(link.ID)

// 	// Convet Link struct to json data
// 	// jbs, _ := link.MarshalJSON(),
// 	// var rawStrings map[string]string
// 	// json.Unmarshal(jbs, &rawStrings)
// 	// fmt.Println("maps:", rawStrings)
// }
