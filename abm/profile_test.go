package abm

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestSchema(t *testing.T) {
	is := is.New(t)
	j := `[{"id":0,"name":"Not selected"},{"id":1,"name":"Man"},{"id":2,"name":"Woman"}]  `

	var s schema
	is.NoErr(json.Unmarshal([]byte(j), &s))
	is.Equal(s[1].ID, "1")
	is.Equal(s[1].Value, "Man")
}
