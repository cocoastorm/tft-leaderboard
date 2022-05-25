package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

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

const timestampKey = "last-modified"

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
			[]byte("discord"),
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
	err := s.storage.View(func(tx *bolt.Tx) error {
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
				Rank:     make(map[string]*tft.TftLeague),
			}

			prefix := []byte(summoner.Id)
			c := leagueBucket.Cursor()

			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				rank := &tft.TftLeague{}

				err := json.Unmarshal(v, rank)
				if err != nil {
					// todo: error aggregator per tft league
					err = fmt.Errorf("failed to decode tft league: %s", err)
					fmt.Println(err)
				}

				pair.Rank[rank.QueueType] = rank
			}

			if len(pair.Rank) != 0 {
				collection = append(collection, pair.Transform())
				return nil
			}

			// fallback behavior, reference summoner id directly
			rawLeague := leagueBucket.Get([]byte(summoner.Id))

			if rawLeague != nil {
				mainRank := &tft.TftLeague{}

				err := json.Unmarshal(rawLeague, mainRank)
				if err != nil {
					// todo: error aggregator per tft league
					err = fmt.Errorf("failed to decode tft league: %s", err)
					fmt.Println(err)
				} else {
					pair.Rank[tft.RankQueueType] = mainRank
				}
			} else {
				// explicitly set rank to nil here to avoid empty json object
				pair.Rank = nil
			}

			collection = append(collection, pair.Transform())
		}

		return nil
	})

	return collection, err
}

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

func upsertRankTimestamp(b *bolt.Bucket) error {
	ts, err := itob(uint64(time.Now().Unix()))
	if err != nil {
		return err
	}

	return b.Put([]byte(timestampKey), ts)
}

func (s *Store) UpdateRankTimestamp() error {
	return s.storage.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tft-leagues"))
		return upsertRankTimestamp(b)
	})
}

// upsertRankLeague adds a rank league by its queue type.
// it concatenates the summoner's id and queue type as its key.
// <summonerId>.<queueType>
func upsertRankLeague(b *bolt.Bucket, summonerId []byte, item *tft.TftLeague) error {
	keyPieces := [][]byte{summonerId, []byte(item.QueueType)}
	key := bytes.Join(keyPieces, []byte("."))

	encoded, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return b.Put(key, encoded)
}

func (s *Store) UpdateContestantRanks(items []*tft.TftPair) error {
	return s.storage.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tft-leagues"))

		for _, pair := range items {
			summoner := pair.Summoner

			// rank can be empty, if so skip it
			if len(pair.Rank) == 0 {
				err := fmt.Errorf("rank list empty, skipping %s [%s]", summoner.Name, summoner.Id)
				fmt.Println(err)
				continue
			}

			// backwards compatibility: store main ranked league with summoner id as the key
			normal, ok := pair.Rank[tft.RankQueueType]
			if ok {
				encoded, err := json.Marshal(normal)
				if err != nil {
					err = fmt.Errorf("failed to update normal rank for %s [%s]: %s", summoner.Name, summoner.Id, err)
					fmt.Println(err)
				}

				err = b.Put([]byte(summoner.Id), encoded)
				if err != nil {
					return err
				}
			}

			// store all rank types
			for queueType, rank := range pair.Rank {
				err := upsertRankLeague(b, []byte(summoner.Id), rank)
				if err != nil {
					return fmt.Errorf("failed to update %s rank for %s [%s]: %s", queueType, summoner.Name, summoner.Id, err)
				}
			}
		}

		// last-modified timestamp for dealing with caching
		upsertRankTimestamp(b)

		return nil
	})
}
