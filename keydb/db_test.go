// cryptctl - Copyright (c) 2017 SUSE Linux GmbH, Germany
// This source code is licensed under GPL version 3 that can be found in LICENSE file.
package keydb

import (
	"os"
	"reflect"
	"testing"
	"time"
)

const TEST_DIR = "/tmp/cryptctl-dbtest"

func TestRecordCRUD(t *testing.T) {
	defer os.RemoveAll(TEST_DIR)
	os.RemoveAll(TEST_DIR)
	db, err := OpenDB(TEST_DIR)
	if err != nil {
		t.Fatal(err)
	}
	// Insert two records
	aliveMsg := AliveMessage{
		Hostname:  "host1",
		IP:        "ip1",
		Timestamp: time.Now().Unix(),
	}
	rec1 := Record{
		UUID:             "1",
		Key:              []byte{0, 1, 2, 3},
		MountPoint:       "/tmp/2",
		MountOptions:     []string{"rw", "noatime"},
		MaxActive:        1,
		AliveIntervalSec: 1,
		AliveCount:       4,
		AliveMessages:    map[string][]AliveMessage{},
	}
	rec1Alive := rec1
	rec1Alive.LastRetrieval = aliveMsg
	rec1Alive.AliveMessages = map[string][]AliveMessage{aliveMsg.IP: []AliveMessage{aliveMsg}}
	rec2 := Record{
		UUID:             "2",
		Key:              []byte{0, 1, 2, 3},
		MountPoint:       "/tmp/2",
		MountOptions:     []string{"rw", "noatime"},
		MaxActive:        1,
		AliveIntervalSec: 1,
		AliveCount:       4,
		AliveMessages:    map[string][]AliveMessage{},
	}
	rec2Alive := rec2
	rec2Alive.LastRetrieval = aliveMsg
	rec2Alive.AliveMessages = map[string][]AliveMessage{aliveMsg.IP: []AliveMessage{aliveMsg}}
	if seq, err := db.Upsert(rec1); err != nil || seq != "1" {
		t.Fatal(err, seq)
	}
	if seq, err := db.Upsert(rec2); err != nil || seq != "2" {
		t.Fatal(err, seq)
	}
	// Match sequence number in my copy of records with their should-be ones
	rec1.ID = "1"
	rec1Alive.ID = "1"
	rec2.ID = "2"
	rec2Alive.ID = "2"
	// Select one record and then select both records
	if found, rejected, missing := db.Select(aliveMsg, true, "1", "doesnotexist"); !reflect.DeepEqual(found, map[string]Record{rec1.UUID: rec1Alive}) ||
		!reflect.DeepEqual(rejected, []string{}) ||
		!reflect.DeepEqual(missing, []string{"doesnotexist"}) {
		t.Fatal(found, rejected, missing)
	}
	if found, rejected, missing := db.Select(aliveMsg, true, "1", "doesnotexist", "2"); !reflect.DeepEqual(found, map[string]Record{rec2.UUID: rec2Alive}) ||
		!reflect.DeepEqual(rejected, []string{"1"}) ||
		!reflect.DeepEqual(missing, []string{"doesnotexist"}) {
		t.Fatal(found, rejected, missing)
	}
	if found, rejected, missing := db.Select(aliveMsg, false, "1", "doesnotexist", "2"); !reflect.DeepEqual(found, map[string]Record{rec1.UUID: rec1Alive, rec2.UUID: rec2Alive}) ||
		!reflect.DeepEqual(rejected, []string{}) ||
		!reflect.DeepEqual(missing, []string{"doesnotexist"}) {
		t.Fatal(found, rejected, missing)
	}
	// Update alive message on both records
	newAlive := AliveMessage{
		Hostname:  "host1",
		IP:        "ip1",
		Timestamp: time.Now().Unix(),
	}
	if rejected := db.UpdateAliveMessage(newAlive, "1", "2", "doesnotexist"); !reflect.DeepEqual(rejected, []string{"doesnotexist"}) {
		t.Fatal(rejected)
	}
	if len(db.RecordsByUUID["1"].AliveMessages["ip1"]) != 2 || len(db.RecordsByUUID["2"].AliveMessages["ip1"]) != 2 {
		t.Fatal(db.RecordsByUUID)
	}
	if len(db.RecordsByID["1"].AliveMessages["ip1"]) != 2 || len(db.RecordsByID["2"].AliveMessages["ip1"]) != 2 {
		t.Fatal(db.RecordsByUUID)
	}
	// Erase a record
	if err := db.Erase("doesnotexist"); err == nil {
		t.Fatal("did not error")
	}
	if err := db.Erase(rec1.UUID); err != nil {
		t.Fatal(err)
	}
	if found, rejected, missing := db.Select(aliveMsg, true, "1"); len(found) != 0 ||
		!reflect.DeepEqual(rejected, []string{}) ||
		!reflect.DeepEqual(missing, []string{"1"}) {
		t.Fatal(found, rejected, missing)
	}
	// Reload database and test query once more (2 is already retrieved and hence it shall be rejected)
	db, err = OpenDB(TEST_DIR)
	if err != nil {
		t.Fatal(err)
	}
	if found, rejected, missing := db.Select(aliveMsg, true, "1", "2"); len(found) != 0 ||
		!reflect.DeepEqual(rejected, []string{"2"}) ||
		!reflect.DeepEqual(missing, []string{"1"}) {
		t.Fatal(found, missing)
	}
}

func TestOpenDBOneRecord(t *testing.T) {
	defer os.RemoveAll(TEST_DIR)
	os.RemoveAll(TEST_DIR)
	db, err := OpenDB(TEST_DIR)
	if err != nil {
		t.Fatal(err)
	}
	rec := Record{
		UUID:         "a",
		Key:          []byte{1, 2, 3},
		MountPoint:   "/a",
		MountOptions: []string{},
		LastRetrieval: AliveMessage{
			Hostname:  "host1",
			IP:        "ip1",
			Timestamp: 3,
		},
		AliveMessages: make(map[string][]AliveMessage),
	}
	if seq, err := db.Upsert(rec); err != nil || seq != "1" {
		t.Fatal(err)
	}
	dbOneRecord, err := OpenDBOneRecord(TEST_DIR, "a")
	if err != nil {
		t.Fatal(err)
	}
	if len(dbOneRecord.RecordsByUUID) != 1 {
		t.Fatal(dbOneRecord.RecordsByUUID)
	}
	rec.ID = "1"
	if recA, found := dbOneRecord.GetByUUID("a"); !found || !reflect.DeepEqual(recA, rec) {
		t.Fatal(recA, found)
	}
	if recA, found := dbOneRecord.GetByID("1"); !found || !reflect.DeepEqual(recA, rec) {
		t.Fatal(recA, found)
	}
	if _, found := dbOneRecord.GetByUUID("doesnotexist"); found {
		t.Fatal("false positive")
	}
	if _, found := dbOneRecord.GetByID("78598123"); found {
		t.Fatal("false positive")
	}
}

func TestList(t *testing.T) {
	defer os.RemoveAll(TEST_DIR)
	db, err := OpenDB(TEST_DIR)
	if err != nil {
		t.Fatal(err)
	}
	// Insert three records and get them back in sorted order
	rec1 := Record{
		UUID:         "a",
		Key:          []byte{1, 2, 3},
		MountPoint:   "/a",
		MountOptions: []string{},
		LastRetrieval: AliveMessage{
			Hostname:  "host1",
			IP:        "ip1",
			Timestamp: 3,
		},
		AliveMessages: make(map[string][]AliveMessage),
	}
	rec1NoKey := rec1
	rec1NoKey.Key = nil
	rec2 := Record{
		UUID:         "b",
		Key:          []byte{1, 2, 3},
		MountPoint:   "/b",
		MountOptions: []string{},
		LastRetrieval: AliveMessage{
			Hostname:  "host1",
			IP:        "ip1",
			Timestamp: 1,
		},
		AliveMessages: make(map[string][]AliveMessage),
	}
	rec2NoKey := rec2
	rec2NoKey.Key = nil
	rec3 := Record{
		UUID:         "c",
		Key:          []byte{1, 2, 3},
		MountPoint:   "/c",
		MountOptions: []string{},
		LastRetrieval: AliveMessage{
			Hostname:  "host1",
			IP:        "ip1",
			Timestamp: 2,
		},
		AliveMessages: make(map[string][]AliveMessage),
	}
	rec3NoKey := rec3
	rec3NoKey.Key = nil
	if seq, err := db.Upsert(rec1); err != nil || seq != "1" {
		t.Fatal(err, seq)
	}
	if seq, err := db.Upsert(rec2); err != nil || seq != "2" {
		t.Fatal(err)
	}
	if seq, err := db.Upsert(rec3); err != nil || seq != "3" {
		t.Fatal(err)
	}
	rec1NoKey.ID = "1"
	rec2NoKey.ID = "2"
	rec3NoKey.ID = "3"
	recs := db.List()
	if !reflect.DeepEqual(recs[0], rec1NoKey) ||
		!reflect.DeepEqual(recs[1], rec3NoKey) ||
		!reflect.DeepEqual(recs[2], rec2NoKey) {
		t.Fatal(recs)
	}
}
