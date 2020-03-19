package dao

// swagger:response event
// Event Json request payload is as follows,
//{
//  "id": "1",
//  "title": "some title",
//  "description":  "some description"
//}
type Event struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
