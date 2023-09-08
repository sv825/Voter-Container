package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestVoterAPI(t *testing.T) {
	// Create a new VoterList instance backed by a Redis cache
	vl := NewVoterList()
	vl.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Create a new http.ServeMux and register the Voter API handlers
	mux := http.NewServeMux()
	RegisterHandlers(mux, vl)

	// Create a new test server using the http.ServeMux
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("GetVoters", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/voters")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var voters []Voter
		err = json.NewDecoder(resp.Body).Decode(&voters)
		if err != nil {
			t.Fatal(err)
		}

		if len(voters) != 0 {
			t.Errorf("expected 0 voters; got %d", len(voters))
		}
	})

	t.Run("AddVoter", func(t *testing.T) {
		voter := Voter{
			VoterID:   1,
			FirstName: "John",
			LastName:  "Doe",
			VoteHistory: []voterPoll{
				voterPoll{
					PollID:   1,
					VoteDate: time.Now(),
				},
			},
		}
	})