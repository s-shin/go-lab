package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	type User struct {
		Name string          `json:"name,omitempty"`
		Age  *int            `json:"age,omitempty"`
		Data json.RawMessage `json:"data,omitempty"`
	}

	{
		alice := &User{
			Name: "alice",
			// https://stackoverflow.com/a/30716481
			Age:  func(x int) *int { return &x }(100),
			Data: nil,
		}

		bytes, err := json.Marshal(alice)
		require.NoError(t, err)
		require.JSONEq(t, `{ "name": "alice", "age": 100 }`, string(bytes))
	}

	{
		dataBytes, err := json.Marshal("secret")
		require.NoError(t, err)
		require.Equal(t, `"secret"`, string(dataBytes))

		anonymous := &User{
			Name: "",
			Age:  nil,
			Data: dataBytes,
		}

		bytes, err := json.Marshal(anonymous)
		require.NoError(t, err)
		require.JSONEq(t, `{ "data": "secret" }`, string(bytes))
	}

	{
		u := &User{}
		err := json.Unmarshal([]byte(`{ "age": 50, "data": ["array"] }`), u)
		require.NoError(t, err)
		require.Equal(t, 50, *u.Age)

		var arr []string
		err = json.Unmarshal(u.Data, &arr)
		require.NoError(t, err)
		require.Len(t, arr, 1)
		require.Equal(t, "array", arr[0])
	}

	{
		u := &User{}
		err := json.Unmarshal([]byte(`{ "name": "alice" }`), &u)
		require.NoError(t, err)
		require.Equal(t, "alice", u.Name)
		require.Nil(t, u.Age)
		require.Nil(t, u.Data)
	}
}
