package postgresql

import "fmt"

// DeleteAllTables cleans up the tables created in test. This method should only be used in tests.
func (s *Storage) DeleteAllTables() error {
	var logger = s.sugar.With("func", "accounting/reserve-addresses/storage/postgresql.Storage.Create")
	logger.Debugw("deleting table", "name", addressesTableName)
	_, err := s.db.Exec(fmt.Sprintf(`DROP TABLE "%s","%s"`, addressesTableName, addressVersionTableName))
	return err
}
