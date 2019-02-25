package postgresql

import "fmt"

func (s *Storage) DeleteAllTables() error {
	var logger = s.sugar.With("func", "accounting/reserve-addresses/storage/postgresql.Storage.Create")
	logger.Debugw("deleting table", "name", addressesTableName)
	_, err := s.db.Exec(fmt.Sprintf(`DROP TABLE "%s"`, addressesTableName))
	return err
}
