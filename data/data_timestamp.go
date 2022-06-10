package data

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

const timestampKey = "last-modified"

func (s *Store) GetRankTimestamp() (uint64, error) {
	var ts uint64

	err := s.storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tft-leagues"))
		val := b.Get([]byte(timestampKey))

		t, err := btoi(val)
		if err != nil {
			return err
		}

		ts = t
		return nil
	})

	if err != nil {
		return 0, err
	}

	return ts, nil
}

func upsertRankTimestamp(tx *bolt.Tx) error {
	b := tx.Bucket([]byte("tft-leagues"))
	ts, err := itob(uint64(time.Now().Unix()))
	if err != nil {
		return err
	}

	return b.Put([]byte(timestampKey), ts)
}

func (s *Store) UpdateRankTimestamp() error {
	return s.storage.Update(func(tx *bolt.Tx) error {
		return upsertRankTimestamp(tx)
	})
}
