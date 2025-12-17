package sqlite

import (
	"configBin/pkg"
	"configBin/pkg/storage"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	log "github.com/sirupsen/logrus"
)

var ErrBinNotFound = errors.New("no bin found")

type Storage struct {
	db        *sql.DB
	encryptor storage.Encryptor
}

func NewSqliteStorage(dbPath string, encryptor storage.Encryptor) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	return &Storage{
		db:        db,
		encryptor: encryptor,
	}, nil
}

func (s *Storage) ApplySchema() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS bins
        (
            uuid       TEXT,
            data       TEXT,
            version    INTEGER,
            created_at TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS bins_uuid_version_index ON bins (uuid, version);
        `)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

func (s *Storage) CreateBin(binID uuid.UUID, pass string, unencryptedData string) error {
	encryptedData, err := s.encryptor.Encrypt(unencryptedData, pass)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO bins (uuid, data, version, created_at) 
		VALUES (?, ?, 0, ?)
		`,
		binID, encryptedData, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert bin: %w", err)
	}
	return nil
}

func (s *Storage) GetBin(binID uuid.UUID, pass string) (*pkg.Bin, error) {
	bin := &pkg.Bin{}
	rows, err := s.db.Query(`
		SELECT 
		    uuid, data, version, created_at 
		FROM bins 
		WHERE uuid = ? 
		ORDER BY version DESC
		`,
		binID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bin: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to query bin rows: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("failed to close rows: %v", err)
		}
	}()

	if rows.Next() {
		var data string
		err = rows.Scan(&bin.ID, &data, &bin.Version, &bin.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan first row: %w", err)
		}
		bin.Data, err = s.encryptor.Decrypt(data, pass)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt first row data: %w", err)
		}
	} else {
		return nil, nil
	}

	bin.History = make([]pkg.Bin, 0)
	for rows.Next() {
		prevBin := pkg.Bin{}
		prevData := ""
		err = rows.Scan(&prevBin.ID, &prevData, &prevBin.Version, &prevBin.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		prevBin.Data, err = s.encryptor.Decrypt(prevData, pass)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data of version `%d`: %w", bin.Version, err)
		}
		bin.History = append(bin.History, prevBin)
	}

	return bin, nil
}

func (s *Storage) UpdateBin(binID uuid.UUID, pass string, unencryptedData string) error {
	bin, err := s.GetBin(binID, pass)
	if err != nil {
		return fmt.Errorf("failed to get bin: %w", err)
	}
	if bin == nil {
		return ErrBinNotFound
	}

	data, err := s.encryptor.Encrypt(unencryptedData, pass)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO bins (uuid, data, version, created_at) 
		VALUES (?, ?, (SELECT COALESCE(MAX(version), 0) + 1 
        FROM bins 
        WHERE uuid = ?), ?)`,
		bin.ID, data, bin.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert bin: %w", err)
	}
	return nil
}

func (s *Storage) Close() {
	err := s.db.Close()
	if err != nil {
		log.Warnf("error closing database: %v", err)
	}
}

func (s *Storage) IsReady() bool {
	err := s.db.Ping()
	if err != nil {
		log.Errorf("failed to ping db: %v", err)
		return false
	}
	return true
}
