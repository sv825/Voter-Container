package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

type VoterList struct {
	rdb *redis.Client
}

// constructor for VoterList struct
func NewVoterList() *VoterList {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &VoterList{
		rdb: rdb,
	}
}

// Get all voter resources including all voter history for each voter (note we will discuss the concept of "paging" later, for now you can ignore)
func (vl *VoterList) GetVoters() ([]Voter, error) {
	keys, err := vl.rdb.Keys(ctx, "voter:*").Result()
	if err != nil {
		return nil, err
	}
	voters := make([]Voter, 0, len(keys))
	for _, key := range keys {
		voter, err := vl.GetVoter(key)
		if err != nil {
			return nil, err
		}
		voters = append(voters, voter)
	}
	return voters, nil
}

// Get a single voter resource with voterID=:id including their entire voting history.
func (vl *VoterList) GetVoter(id string) (Voter, error) {
	val, err := vl.rdb.Get(ctx, id).Result()
	if err != nil {
		return Voter{}, err
	}
	var v Voter
	err = json.Unmarshal([]byte(val), &v)
	if err != nil {
		return Voter{}, err
	}
	return v, nil
}

// POST version adds one to the "database"
func (vl *VoterList) AddVoter(v Voter) error {
	key := fmt.Sprintf("voter:%d", v.VoterID)
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = vl.rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Gets the JUST the voter history for the voter with VoterID = :id
func (vl *VoterList) GetVoterPolls(id string) ([]voterPoll, error) {
	voter, err := vl.GetVoter(id)
	if err != nil {
		return nil, err
	}
	return voter.VoteHistory, nil
}

// Gets JUST the single voter poll data with PollID = :id and VoterID = :id.
func (vl *VoterList) GetVoterPoll(voterID string, pollID uint) (voterPoll, error) {
	voterPolls, err := vl.GetVoterPolls(voterID)
	if err != nil {
		return voterPoll{}, err
	}
	for _, vp := range voterPolls {
		if vp.PollID == pollID {
			return vp, nil
		}
	}
	return voterPoll{}, fmt.Errorf("poll not found")
}

// POST version adds one to the "database"
func (vl *VoterList) AddVoterPoll(voterID string, vp voterPoll) error {
	voter, err := vl.GetVoter(voterID)
	if err != nil {
		return err
	}
	voter.VoteHistory = append(voter.VoteHistory, vp)
	err = vl.AddVoter(voter)
	if err != nil {
		return err
	}
	return nil
}

// Returns a "health" record indicating that the voter API is functioning properly and some metadata about the API.
func (vl *VoterList) HealthCheck() string {
	return "API is functioning properly"
}
