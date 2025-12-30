package internal

import (
	"fmt"
)

// TransactionStatus represents the state of a transaction
type TransactionStatus string

const (
	TxActive    TransactionStatus = "active"
	TxCommitted TransactionStatus = "committed"
	TxAborted   TransactionStatus = "aborted"
)

// Transaction represents a database transaction
type Transaction struct {
	ID         string            `json:"id"`
	StartTime  int64             `json:"startTime"`
	CommitTime *int64            `json:"commitTime,omitempty"`
	Status     TransactionStatus `json:"status"`
	ReadSet    []string          `json:"readSet"`
	WriteSet   []string          `json:"writeSet"`
}

// Version represents a version of a row
type Version struct {
	ID        string                 `json:"id"`
	RowID     string                 `json:"rowId"`
	Data      map[string]interface{} `json:"data"`
	CreatedBy string                 `json:"createdBy"`
	CreatedAt int64                  `json:"createdAt"`
	DeletedBy *string                `json:"deletedBy,omitempty"`
	DeletedAt *int64                 `json:"deletedAt,omitempty"`
	Prev      *string                `json:"prev,omitempty"`
}

// Row represents a logical row with its version chain
type Row struct {
	ID             string   `json:"id"`
	CurrentVersion string   `json:"currentVersion"`
	VersionChain   []string `json:"versionChain"`
}

// MVCCStore manages MVCC state
type MVCCStore struct {
	Transactions      map[string]*Transaction `json:"transactions"`
	Versions          map[string]*Version     `json:"versions"`
	Rows              map[string]*Row         `json:"rows"`
	GlobalTimestamp   int64                   `json:"globalTimestamp"`
	ActiveTransaction string                  `json:"activeTransaction,omitempty"`
	txSeq             int
	verSeq            int
}

// NewMVCCStore creates a new MVCC store
func NewMVCCStore() *MVCCStore {
	return &MVCCStore{
		Transactions:    make(map[string]*Transaction),
		Versions:        make(map[string]*Version),
		Rows:            make(map[string]*Row),
		GlobalTimestamp: 1,
	}
}

// Clone creates a deep copy of the store
func (s *MVCCStore) Clone() *MVCCStore {
	clone := &MVCCStore{
		Transactions:      make(map[string]*Transaction),
		Versions:          make(map[string]*Version),
		Rows:              make(map[string]*Row),
		GlobalTimestamp:   s.GlobalTimestamp,
		ActiveTransaction: s.ActiveTransaction,
		txSeq:             s.txSeq,
		verSeq:            s.verSeq,
	}

	for id, tx := range s.Transactions {
		txCopy := *tx
		txCopy.ReadSet = append([]string{}, tx.ReadSet...)
		txCopy.WriteSet = append([]string{}, tx.WriteSet...)
		clone.Transactions[id] = &txCopy
	}

	for id, ver := range s.Versions {
		verCopy := *ver
		verCopy.Data = make(map[string]interface{})
		for k, v := range ver.Data {
			verCopy.Data[k] = v
		}
		clone.Versions[id] = &verCopy
	}

	for id, row := range s.Rows {
		rowCopy := *row
		rowCopy.VersionChain = append([]string{}, row.VersionChain...)
		clone.Rows[id] = &rowCopy
	}

	return clone
}

// BeginTransaction starts a new transaction
func (s *MVCCStore) BeginTransaction() *Transaction {
	s.txSeq++
	txID := fmt.Sprintf("tx-%d", s.txSeq)

	tx := &Transaction{
		ID:        txID,
		StartTime: s.GlobalTimestamp,
		Status:    TxActive,
		ReadSet:   []string{},
		WriteSet:  []string{},
	}

	s.Transactions[txID] = tx
	s.ActiveTransaction = txID
	return tx
}

// Read reads a row within a transaction
func (s *MVCCStore) Read(txID string, rowID string) (*Version, error) {
	tx, ok := s.Transactions[txID]
	if !ok {
		return nil, fmt.Errorf("transaction %s not found", txID)
	}

	if tx.Status != TxActive {
		return nil, fmt.Errorf("transaction %s is not active", txID)
	}

	row, ok := s.Rows[rowID]
	if !ok {
		return nil, fmt.Errorf("row %s not found", rowID)
	}

	// Find visible version using snapshot isolation
	for _, verID := range row.VersionChain {
		ver := s.Versions[verID]
		if s.isVisible(tx, ver) {
			tx.ReadSet = append(tx.ReadSet, rowID)
			return ver, nil
		}
	}

	return nil, fmt.Errorf("no visible version for row %s", rowID)
}

// Write writes a new version within a transaction
func (s *MVCCStore) Write(txID string, rowID string, data map[string]interface{}) (*Version, error) {
	tx, ok := s.Transactions[txID]
	if !ok {
		return nil, fmt.Errorf("transaction %s not found", txID)
	}

	if tx.Status != TxActive {
		return nil, fmt.Errorf("transaction %s is not active", txID)
	}

	s.verSeq++
	verID := fmt.Sprintf("ver-%d", s.verSeq)

	row, exists := s.Rows[rowID]
	if !exists {
		// Create new row
		row = &Row{
			ID:           rowID,
			VersionChain: []string{},
		}
		s.Rows[rowID] = row
	}

	var prevVer *string
	if len(row.VersionChain) > 0 {
		prevVer = &row.VersionChain[0]
	}

	version := &Version{
		ID:        verID,
		RowID:     rowID,
		Data:      data,
		CreatedBy: txID,
		CreatedAt: s.GlobalTimestamp,
		Prev:      prevVer,
	}

	s.Versions[verID] = version
	row.VersionChain = append([]string{verID}, row.VersionChain...)
	row.CurrentVersion = verID
	tx.WriteSet = append(tx.WriteSet, rowID)

	return version, nil
}

// Delete marks a version as deleted
func (s *MVCCStore) Delete(txID string, rowID string) error {
	tx, ok := s.Transactions[txID]
	if !ok {
		return fmt.Errorf("transaction %s not found", txID)
	}

	if tx.Status != TxActive {
		return fmt.Errorf("transaction %s is not active", txID)
	}

	row, ok := s.Rows[rowID]
	if !ok {
		return fmt.Errorf("row %s not found", rowID)
	}

	if len(row.VersionChain) > 0 {
		ver := s.Versions[row.VersionChain[0]]
		deletedAt := s.GlobalTimestamp
		ver.DeletedBy = &txID
		ver.DeletedAt = &deletedAt
	}

	tx.WriteSet = append(tx.WriteSet, rowID)
	return nil
}

// Commit commits a transaction
func (s *MVCCStore) Commit(txID string) error {
	tx, ok := s.Transactions[txID]
	if !ok {
		return fmt.Errorf("transaction %s not found", txID)
	}

	if tx.Status != TxActive {
		return fmt.Errorf("transaction %s is not active", txID)
	}

	s.GlobalTimestamp++
	commitTime := s.GlobalTimestamp
	tx.CommitTime = &commitTime
	tx.Status = TxCommitted

	if s.ActiveTransaction == txID {
		s.ActiveTransaction = ""
	}

	return nil
}

// Abort aborts a transaction
func (s *MVCCStore) Abort(txID string) error {
	tx, ok := s.Transactions[txID]
	if !ok {
		return fmt.Errorf("transaction %s not found", txID)
	}

	if tx.Status != TxActive {
		return fmt.Errorf("transaction %s is not active", txID)
	}

	tx.Status = TxAborted

	// Remove uncommitted versions
	for _, rowID := range tx.WriteSet {
		row := s.Rows[rowID]
		newChain := []string{}
		for _, verID := range row.VersionChain {
			ver := s.Versions[verID]
			if ver.CreatedBy != txID {
				newChain = append(newChain, verID)
			} else {
				delete(s.Versions, verID)
			}
		}
		row.VersionChain = newChain
		if len(newChain) > 0 {
			row.CurrentVersion = newChain[0]
		} else {
			delete(s.Rows, rowID)
		}
	}

	if s.ActiveTransaction == txID {
		s.ActiveTransaction = ""
	}

	return nil
}

// isVisible checks if a version is visible to a transaction
func (s *MVCCStore) isVisible(tx *Transaction, ver *Version) bool {
	creatorTx := s.Transactions[ver.CreatedBy]

	// Version created by this transaction is visible
	if ver.CreatedBy == tx.ID {
		return true
	}

	// Version must be committed before our snapshot
	if creatorTx.Status != TxCommitted {
		return false
	}

	if creatorTx.CommitTime == nil || *creatorTx.CommitTime > tx.StartTime {
		return false
	}

	// Check if deleted before our snapshot
	if ver.DeletedBy != nil && ver.DeletedAt != nil {
		deleterTx := s.Transactions[*ver.DeletedBy]
		if deleterTx.Status == TxCommitted && *deleterTx.CommitTime <= tx.StartTime {
			return false
		}
	}

	return true
}

// GetVisibleVersions returns all versions visible to a transaction
func (s *MVCCStore) GetVisibleVersions(txID string) []string {
	tx, ok := s.Transactions[txID]
	if !ok {
		return nil
	}

	visible := []string{}
	for _, row := range s.Rows {
		for _, verID := range row.VersionChain {
			ver := s.Versions[verID]
			if s.isVisible(tx, ver) {
				visible = append(visible, verID)
				break // Only first visible version per row
			}
		}
	}

	return visible
}

// GarbageCollect removes old versions that are no longer needed
func (s *MVCCStore) GarbageCollect() []string {
	oldestActiveStart := s.GlobalTimestamp
	for _, tx := range s.Transactions {
		if tx.Status == TxActive && tx.StartTime < oldestActiveStart {
			oldestActiveStart = tx.StartTime
		}
	}

	removed := []string{}
	for _, row := range s.Rows {
		foundVisible := false
		newChain := []string{}

		for _, verID := range row.VersionChain {
			ver := s.Versions[verID]

			// Keep if might be needed by active transactions
			creatorTx := s.Transactions[ver.CreatedBy]
			if creatorTx.Status == TxCommitted && creatorTx.CommitTime != nil {
				if *creatorTx.CommitTime < oldestActiveStart && foundVisible {
					// Safe to remove
					removed = append(removed, verID)
					delete(s.Versions, verID)
					continue
				}
				foundVisible = true
			}

			newChain = append(newChain, verID)
		}

		row.VersionChain = newChain
	}

	return removed
}

// InsertInitialData adds some initial rows for testing
func (s *MVCCStore) InsertInitialData() {
	// Create a committed transaction for initial data
	tx := s.BeginTransaction()

	s.Write(tx.ID, "users:1", map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	})

	s.Write(tx.ID, "users:2", map[string]interface{}{
		"id":    2,
		"name":  "Bob",
		"email": "bob@example.com",
	})

	s.Write(tx.ID, "products:1", map[string]interface{}{
		"id":    1,
		"name":  "Widget",
		"price": 9.99,
	})

	s.Commit(tx.ID)
}
