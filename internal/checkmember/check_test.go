package checkmember

import (
	"testing"

	"github.com/google/uuid"
)

func TestMemberCacheLifecycle(t *testing.T) {
	serverID := uuid.New()
	userID := uuid.New()
	cache := NewMemberCache(make(map[uuid.UUID]map[uuid.UUID]struct{}))

	if cache.IsMemberCache(serverID, userID) {
		t.Fatal("member should not exist before it is added")
	}

	cache.Add(serverID, userID)
	if !cache.IsMemberCache(serverID, userID) {
		t.Fatal("member should exist after it is added")
	}

	cache.RemoveMember(serverID, userID)
	if cache.IsMemberCache(serverID, userID) {
		t.Fatal("member should not exist after it is removed")
	}

	cache.Add(serverID, userID)
	cache.RemoveServer(serverID)
	if cache.IsMemberCache(serverID, userID) {
		t.Fatal("members should not exist after the server is removed")
	}
}
