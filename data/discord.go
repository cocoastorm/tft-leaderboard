package data

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func (s *Store) LookupDiscordSequence(authId string) (uint64, error) {
	var lookup []byte

	err := s.storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("discord"))
		lookup = b.Get([]byte(authId))
		return nil
	})

	if lookup == nil {
		return 0, fmt.Errorf("failed to lookup: %s", authId)
	}

	sequence, err := btoi(lookup)
	if err != nil {
		return 0, fmt.Errorf("sequence id failed: %s", err)
	}

	return sequence, nil
}

func (s *Store) UpdateDiscord(sequenceId uint64, authId string) error {
	return s.storage.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("discord"))

		id, err := itob(sequenceId)
		if err != nil {
			return err
		}

		return b.Put(id, []byte(authId))
	})
}
