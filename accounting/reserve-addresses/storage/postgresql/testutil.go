package postgresql

// DeleteAllTables cleans up the tables created in test. This method should only be used in tests.
func (s *Storage) DeleteAllTables() error {
	_, err := s.db.Exec(`DROP TABLE "addresses","addresses_version"`)
	return err
}
