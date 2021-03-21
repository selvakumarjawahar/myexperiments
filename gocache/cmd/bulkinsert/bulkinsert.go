package bulkinsert

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/lib/pq"
)

type BulkInsert struct {
	maxBatchSize         int
	maxBatchAge          int
	db                   *sql.DB
	table                string
	firstInsertTimestamp int64
	columnNames          []string
	values               *[][]string
	columnIndexes        map[string]int
}

// NewBulkInsert returns new BulkInsert
func NewBulkInsert(db *sql.DB, table string, maxBatchSize int, maxBatchAge int, columnNames []string) *BulkInsert {
	// Log columns to be inserted so that it is possible to debug outcome
	indexMap := make(map[string]int)
	log.Debug().Msgf("NewBulkInsert table %s, columns:", table)
	for i, col := range columnNames {
		log.Debug().Msgf("%d = %s", i, col)
		// set one based index because map returns zero when key is not found
		indexMap[col] = i + 1
	}

	return &BulkInsert{
		db:            db,
		table:         table,
		columnNames:   columnNames,
		maxBatchAge:   maxBatchAge,
		maxBatchSize:  maxBatchSize,
		columnIndexes: indexMap,
	}
}

// GetDB returns pointer to DB
func (bi *BulkInsert) GetDB() *sql.DB {
	return bi.db
}
// NewRow return column array for row intialized to nil values "/nil"
func (bi *BulkInsert) NewRow() *[]string {
	columns := make([]string, len(bi.columnNames), len(bi.columnNames))
	for i := range columns {
		columns[i] = "\nil"
	}
	return &columns
}

// SetValue sets value for column in row
func (bi *BulkInsert) SetValue(column string, value string, columns *[]string) {
	if len(*columns) != len(bi.columnNames) {
		log.Error().Msg("Column count differs from map column count")
		return
	}
	index := bi.columnIndexes[column]
	if index == 0 {
		log.Error().Msgf("Column was not found %s", column)
		return
	}
	(*columns)[index-1] = value
}

// writes bulk inserted data to db
func write(db *sql.DB, table string, columnNames []string, values *[][]string) {
	log.Debug().Msgf("Inserting %d rows into %s", len(*values), table)
	start := time.Now()
	txn, err := db.Begin()
	if err != nil {
		log.Error().Err(err).Msg("Prepare transaction failed")
		return
	}
	copyIn(txn, table, columnNames, values)
	err = txn.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed commiting transaction")
	}
	log.Debug().Msgf("Inserted %d rows into %s in %v s", len(*values), table, time.Now().Sub(start).Seconds())

}

// copyIn writes bulk inserted data to db with defined transaction
func copyIn(txn *sql.Tx, table string, columnNames []string, values *[][]string) {
	stmt, err := txn.Prepare(pq.CopyIn(table, columnNames...))
	if err != nil {
		log.Error().Err(err).Msg("Prepare transaction failed")
		return
	}

	for _, row := range *values {
		val := make([]interface{}, len(row))
		//add rows replacing "\nil" tags with nil value
		for i, item := range row {
			if item == "\nil" {
				val[i] = nil
			} else {
				val[i] = item
			}
		}
		_, err = stmt.Exec(val...)
		if err != nil {
			log.Error().Err(err).Msg("Failed executing statement")
			return
		}

	}
	_, err = stmt.Exec()
	if err != nil {
		log.Error().Err(err).Msg("Failed executing closing statement")
	}
	err = stmt.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed close")
	}
}

// Write inserts rows to db
func (bi *BulkInsert) Write() error {
	go write(bi.db, bi.table, bi.columnNames, bi.values)
	bi.values = nil
	return nil
}

// WriteAndCommit writes rows to DB and commits transaction
func (bi *BulkInsert) WriteAndCommit(txn *sql.Tx)  {
	go func(txn *sql.Tx, table string, columnNames []string, values *[][]string) {
		start := time.Now()
		copyIn(txn, table, columnNames, values)
		err := txn.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Failed commiting transaction")
		}
		log.Debug().Msgf("Inserted %d rows into %s in %v s", len(*values), table, time.Now().Sub(start).Seconds())
	} (txn, bi.table, bi.columnNames, bi.values)
	bi.values = nil
}

// Insert bulk inserts to pq db
func (bi *BulkInsert) Insert(values []string) error {
	var err error
	if len(values) != len(bi.columnNames) {
		log.Error().Msgf("Mismatch in columns/values len %d != %d", len(bi.columnNames), len(values))
		for i, col := range bi.columnNames {
			if len(values) < i {
				break
			}
			log.Debug().Msgf("%s = %s", col, values[i])

		}
		return errors.New("mismatch in columns/values")
	}

	if bi.firstInsertTimestamp == 0 {
		bi.firstInsertTimestamp = time.Now().Unix()
	}
	now := time.Now().Unix()

	// commit to db if batch age or size condition is met
	if (bi.maxBatchAge>0 && ((now-bi.firstInsertTimestamp) >= int64(bi.maxBatchAge))) || 
	   (bi.maxBatchSize>0 && (bi.values != nil && len(*bi.values) >= bi.maxBatchSize)) {
		err = bi.Write()
		if err != nil {
			return err
		}
		bi.firstInsertTimestamp = now
	}
	if bi.values == nil {
		bi.values = new([][]string)
	}
	*bi.values = append(*bi.values, values)
	return nil
}
