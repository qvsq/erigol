package message

// Type is message type
type Type string

// Message item types
const (
	AddItemType     Type = "add"
	RemoveItemType  Type = "remove"
	GetItemType     Type = "get"
	GetAllItemsType Type = "list"
)

// Message consist of type and item
type Message struct {
	Type Type
	Body []byte
}
