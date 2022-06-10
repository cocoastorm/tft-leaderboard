package data

import (
  bolt "go.etcd.io/bbolt"
  "ktn-x.com/tft-leaderboard/tft"
)

type Store struct {
	storage  *bolt.DB
	dataPath string
}

type Contestant struct {
	SequenceId uint64        `json:"sequence_id"`
	Summoner   *tft.Summoner `json:"summoner"`
}

func OpenDB(dataPath string) (*Store, error) {
	var err error

	db, err := bolt.Open(dataPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	store := &Store{
		storage:  db,
		dataPath: dataPath,
	}

	// create bbolt buckets
	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			[]byte("contestants"),
			[]byte("tft-leagues"),
		}

		for _, name := range buckets {
			_, err := tx.CreateBucketIfNotExists(name)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return store, err
}

func (s *Store) Close() error {
	return s.storage.Close()
}

