// structs
package main

type FrFjson struct {
	Posts struct {
		Id          string   `json:"id"`
		Body        string   `json:"body"`
		PostedTo    []string `json:"postedTo"`
		CreatedAt   string   `json:"createdAt"`
		CreatedBy   string   `json:"createdBy"`
		Attachments []string `json:"attachments"`
		Likes       []string `json:"likes"`
		Comments    []string `json:"comments"`
	} `json:"posts"`
	Comments []struct {
		Body      string `json:"body"`
		UpdatedAt string `json:"updatedAt"`
		Likes     string `json:"likes"`
		CreatedBy string `json:"createdBy"`
	} `json:"comments"`
}

