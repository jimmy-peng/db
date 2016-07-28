package raftwal

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/gyuho/db/raftwal/raftwalpb"
)

func TestRecordEncode(t *testing.T) {
	var (
		data = []byte("Hello World!")
		buf  = new(bytes.Buffer)
		enc  = newEncoder(buf, 0)

		tp = int64(0xABCDE)
	)
	enc.encode(&raftwalpb.Record{
		Type: tp,
		Data: data,
	})
	enc.flush()

	dec := newDecoder(ioutil.NopCloser(buf))
	rec := &raftwalpb.Record{}
	if err := dec.decode(rec); err != nil {
		t.Fatal(err)
	}
	if rec.Type != tp {
		t.Fatalf("type expected %x, got %x", tp, rec.Type)
	}
	if !bytes.Equal(rec.Data, data) {
		t.Fatalf("data expected %q, got %q", data, rec.Data)
	}
}