package thunder_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/omeid/thunder"
	"github.com/omeid/thunder/codec/json"
	"github.com/omeid/thunder/codec/string"
)

type BasicStruct struct {
	Complex string
	Type    int
}

type Basics []BasicStruct

func (b Basics) Len() int           { return len(b) }
func (b Basics) Less(i, j int) bool { return b[i].Type < b[j].Type }
func (b Basics) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

var test_items = map[string]BasicStruct{
	"test":   {"Much complex", 0},
	"test1":  {"Much complex", 1},
	"test2":  {"Much complex", 2},
	"test3":  {"Much complex", 3},
	"test4":  {"Much complex", 4},
	"test5":  {"Much complex", 5},
	"test6":  {"Much complex", 6},
	"test7":  {"Much complex", 7},
	"test8":  {"Much complex", 8},
	"test9":  {"Much complex", 9},
	"test10": {"Much complex", 10},
	"test11": {"Much complex", 11},
}

func TestBucketPutGet(t *testing.T) {

	db, err := thunder.Open("testdata/tmp_json.bd", 0600, thunder.Options{
		KeyCodec:   strings.Codec(),
		ValueCodec: json.Codec(),
	})

	if err != nil {
		t.Fatal(err)
	}

	key := "test"
	value := BasicStruct{
		"Much complex",
		0,
	}

	tx := db.Begin(true)

	tx.CreateBucketIfNotExists("test")
	if tx.Err() != nil {
		t.Fatal(err)
	}

	err = tx.Bucket("test").Put(key, value)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	var gvalue BasicStruct

	err = db.Begin(false).Bucket("test").Get(key, &gvalue)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, gvalue) {
		t.Fatalf("Misss match %v != %v", value, gvalue)
	}

	db.Close()

}

func TestBucketAllWhere(t *testing.T) {

	db, err := thunder.Open("testdata/tmp_json.bd", 0600, thunder.Options{
		KeyCodec:   strings.Codec(),
		ValueCodec: json.Codec(),
	})

	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin(true)

	tx.CreateBucketIfNotExists("test_All")
	if tx.Err() != nil {
		t.Fatal(err)
	}

	err = tx.Bucket("test_All").Insert(test_items)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	//ALL MAP
	items := map[string]BasicStruct{}

	tx = db.Begin(false)

	err = tx.Bucket("test_All").All(&items)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(test_items, items) {
		t.Fatalf("Misss match %v != %v", test_items, items)
	}

	var values_slice Basics

	for _, v := range test_items {
		values_slice = append(values_slice, v)
	}

	sort.Sort(values_slice)
	//ALL Slice
	items_slice := Basics{}

	tx = db.Begin(false)
	err = tx.Bucket("test_All").All(&items_slice)
	if err != nil {
		t.Fatal(err)
	}

	sort.Sort(items_slice)
	if !reflect.DeepEqual(items_slice, values_slice) {
		t.Fatalf("Misss match %v != %v", items_slice, values_slice)
	}

	test_odd := map[string]BasicStruct{}

	for key, value := range test_items {
		if value.Type%2 != 0 {
			test_odd[key] = value
		}
	}

	items = map[string]BasicStruct{}

	err = db.Begin(false).Bucket("test_All").Where(func(v interface{}) bool {
		t, ok := v.(BasicStruct)
		if !ok {
			return false
		}

		return t.Type%2 != 0
	},
		&items)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(test_odd, items) {
		t.Fatalf("Misss match %v != %v", test_odd, items)
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}

	db.Close()

}
