package data

import (
  "fmt"
  "encoding/json"

  bolt "go.etcd.io/bbolt"
)

func upsertContestant(c *Contestant) func(tx *bolt.Tx) error {
  return func(tx *bolt.Tx) error {
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
  }
}

func (s *Store) UpdateContestant(c *Contestant) error {
	return s.storage.Update(upsertContestant(c))
}

func (s *Store) UpdateContestants(contestants []*Contestant) error {
	return s.storage.Update(func(tx *bolt.Tx) error {
		for _, c := range contestants {
			if c == nil {
				continue
			}
      
      err := upsertContestant(c)

      // todo: handle error in a bag?
      fmt.Println(err)
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

