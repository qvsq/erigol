package orderedmap_test

import (
	"testing"

	"bloxroute/internal/orderedmap"
)

func Test(t *testing.T) {
	m := orderedmap.New()

	m.Store("big power", "big responsibility")
	m.Store("why you", "so serious")
	m.Store("only one", "lord of the ring")
	m.Store("bond", "james bond")
	m.Delete("why you")
	m.Store("only one", "hour a night")

	expResult := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "big power",
			Value: "big responsibility",
		},
		{
			Key:   "only one",
			Value: "hour a night",
		},
		{
			Key:   "bond",
			Value: "james bond",
		},
	}

	for i, el := 0, m.Front(); el != nil; i, el = i+1, el.Next() {
		if expResult[i].Key != el.Key {
			t.Errorf("Wrong key, exp: '%s', got: '%s'", expResult[i].Key, el.Value)
		}

		if expResult[i].Value != el.Value {
			t.Errorf("Wrong value, exp: '%s', got: '%s'", expResult[i].Key, el.Value)
		}
	}
}
