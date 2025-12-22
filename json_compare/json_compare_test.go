package json_compare

import (
	"encoding/json"
	"testing"
)

type Payload struct {
	ID    int                    `json:"id"`
	Name  string                 `json:"name"`
	Email string                 `json:"email"`
	Meta  map[string]interface{} `json:"meta"`
}

var p = Payload{
	ID:    42,
	Name:  "Gopher",
	Email: "gopher@example.com",
	Meta:  map[string]interface{}{"tags": []string{"go", "bench"}, "active": true},
}

func BenchmarkStdlibMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(p)
	}
}

func BenchmarkStdlibUnmarshal(b *testing.B) {
	data, _ := json.Marshal(p)
	var out Payload
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(data, &out)
	}
}
