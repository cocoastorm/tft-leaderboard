package data

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"ktn-x.com/tft-leaderboard/tft"
)

type Store struct {
	storage *bolt.DB
	dataPath string
}

type Contestant struct {
	SequenceId uint64 `json:"sequence_id"`
	Summoner *tft.Summoner `json:"summoner"`
}

func OpenDB(dataPath string) (*Store, error) {
	var err error

	db, err := bolt.Open(dataPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	store := &Store{
		storage: db,
		dataPath: dataPath,
	}

	// create bbolt buckets
	err = db.Update(func(tx *bolt.Tx) error {
		var err error

		_, err = tx.CreateBucketIfNotExists([]byte("contestants"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("tft-leagues"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("discord"))
		if err != nil {
			return err
		}

		return nil
	})

	return store, err
}

func (s *Store) Close() error {
	return s.storage.Close()
}

func (s *Store) UpdateContestant(c *Contestant) error {
	return s.storage.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("contestants"))

		id, err := b.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to get next contestant sequence: %s", err)
		}

		c.SequenceId = id
		encoded, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to encode contestant into json: %s", err)
		}

		iid, err := itob(id)
		if err != nil {
			return fmt.Errorf("failed to binary sequence id: %s", err)
		}

		return b.Put(iid, encoded)
	})
}

func (s *Store) UpdateContestants(contestants []*Contestant) error {
	return s.storage.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("contestants"))

		for _, c := range contestants {
			if c == nil {
				continue
			}

			id, err := b.NextSequence()
			if err != nil {
				return fmt.Errorf("failed to get next contestant sequence: %s", err)
			}

			c.SequenceId = id
			encoded, err := json.Marshal(c)
			if err != nil {
				return fmt.Errorf("failed to encode contestant into json: %s", err)
			}

			iid, err := itob(c.SequenceId)
			if err != nil {
				return fmt.Errorf("failed to binary sequence id: %s", err)
			}

			// todo: handle error here
			b.Put(iid, encoded)
		}

		return nil
	})
}

func (s *Store) ListContestants() ([]*Contestant, error) {
	collection := make([]*Contestant, 0, 7)
	err := s.storage.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte("contestants"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			contestant := &Contestant{}
			err := json.Unmarshal(v, contestant)

			if err != nil {
				// todo: error aggregator per contestant
				fmt.Println(err)
				continue
			}

			collection = append(collection, contestant)
		}

		return nil
	})

	return collection, err
}

func (s *Store) ListContestantRanks() (tft.RankResults, error) {
	collection := make(tft.RankResults, 0, 7)
	err := s.storage.View(func(tx *bolt.Tx) error {
		contestantBucket := tx.Bucket([]byte("contestants"))
		contestantCursor := contestantBucket.Cursor()

		leagueBucket := tx.Bucket([]byte("tft-leagues"))

		for k, rawContestant := contestantCursor.First(); k != nil; k, rawContestant = contestantCursor.Next() {
			contestant := &Contestant{}
			err := json.Unmarshal(rawContestant, contestant)

			if err != nil {
				// todo: error aggregator per contestant
				fmt.Println(err)
				continue
			}

			summoner := contestant.Summoner

			pair := &tft.TftPair{
				Summoner: summoner,
				Rank: &tft.TftLeague{},
			}
			
			if rawLeague := leagueBucket.Get([]byte(summoner.Id)); rawLeague != nil {
				err := json.Unmarshal(rawLeague, pair.Rank)
				if err != nil {
					// todo: error aggregator per tft league
					err = fmt.Errorf("failed to decode tft league: %s", err)
				}
			} else {
				// explicitly set rank to nil here to avoid empty json object
				pair.Rank = nil

				// todo: error aggregator per tft league
				err = fmt.Errorf("no tft rank for %s [%s]", summoner.Name, summoner.Id)
				fmt.Println(err)
			}

			collection = append(collection, pair.Transform())
		}

		return nil
	});

	return collection, err
}

func (s *Store) UpdateContestantRanks(items []*tft.TftPair) error {
	return s.storage.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tft-leagues"))

		for _, pair := range items {
			summoner := pair.Summoner

			// rank can be empty, if so skip it
			if pair.Rank == nil {
				err := fmt.Errorf("skipping %s [%s]", summoner.Name, summoner.Id)
				fmt.Println(err)
				continue
			}

			encoded, err := json.Marshal(pair.Rank)
			if err != nil {
				err = fmt.Errorf("failed to update rank for %s [%s]: %s", summoner.Name, summoner.Id, err)
				fmt.Println(err)
				continue
			}

			b.Put([]byte(summoner.Id), encoded)
		}

		return nil
	})
}
