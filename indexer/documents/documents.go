package documents

import (
	"fmt"
	"strings"
	"time"

	"github.com/meeypioneer/mey-indexer/indexer/category"
)

// DocType is an interface for structs to be used as database documents
type DocType interface {
	GetID() string
	SetID(string)
}

// BaseEsType implements DocType and contains the document's id
type BaseEsType struct {
	Id string `json:"-" db:"id"`
}

// GetID returns the document's id
func (m BaseEsType) GetID() string {
	return m.Id
}

// SetID sets the document's id
func (m BaseEsType) SetID(id string) {
	m.Id = id
}

// EsBlock is a block stored in the database
type EsBlock struct {
	*BaseEsType
	Timestamp     time.Time `json:"ts" db:"ts"`
	BlockNo       uint64    `json:"no" db:"no"`
	TxCount       uint      `json:"txs" db:"txs"`
	Size          int64     `json:"size" db:"size"`
	RewardAccount string    `json:"reward_account" db:"reward_account"`
	RewardAmount  string    `json:"reward_amount" db:"reward_amount"`
}

// EsTx is a transaction stored in the database
type EsTx struct {
	*BaseEsType
	Timestamp   time.Time           `json:"ts" db:"ts"`
	BlockNo     uint64              `json:"blockno" db:"blockno"`
	Account     string              `json:"from" db:"from"`
	Recipient   string              `json:"to" db:"to"`
	Amount      string              `json:"amount" db:"amount"`             // string of BigInt
	AmountFloat float32             `json:"amount_float" db:"amount_float"` // float for sorting
	Type        string              `json:"type" db:"type"`
	Category    category.TxCategory `json:"category" db:"category"`
	Payload    	string 				`json:"payload" db:"payload"`
}

// EsName is a name-address mapping stored in the database
type EsName struct {
	*BaseEsType
	Name        string `json:"name" db:"name"`
	Address     string `json:"address" db:"address"`
	UpdateBlock uint64 `json:"blockno" db:"blockno"`
	UpdateTx    string `json:"tx" db:"tx"`
}

// EsMappings contains the elasticsearch mappings
var EsMappings = map[string]string{
	"tx": `{
		"mappings":{
			"tx":{
				"properties":{
					"ts": {
						"type": "date"
					},
					"blockno": {
						"type": "long"
					},
					"from": {
						"type": "keyword"
					},
					"to": {
						"type": "keyword"
					},
					"amount": {
						"enabled": false
					},
					"amount_float": {
						"type": "float"
					},
					"type": {
						"type": "keyword"
					},
					"category": {
						"type": "keyword"
					}
				}
			}
		}
	}`,
	"block": `{
		"mappings":{
			"block":{
				"properties": {
					"ts": {
						"type": "date"
					},
					"no": {
						"type": "long"
					},
					"txs": {
						"type": "long"
					},
					"size": {
						"type": "long"
					},
					"reward_account": {
						"type": "keyword"
					},
					"reward_amount": {
						"enabled": false
					}
				}
			}
		}
	}`,
	"name": `{
		"mappings":{
			"name":{
				"properties": {
					"name": {
						"type": "keyword"
					},
					"address": {
						"type": "keyword"
					},
					"blockno": {
						"type": "long"
					},
					"tx": {
						"type": "keyword"
					}
				}
			}
		}
	}`,
}

func mapCategoriesToStr(categories []category.TxCategory) []string {
	vsm := make([]string, len(categories))
	for i, v := range categories {
		vsm[i] = fmt.Sprintf("'%s'", v)
	}
	return vsm
}

var categories = strings.Join(mapCategoriesToStr(category.TxCategories), ",")

// SQLSchemas contains schema for SQL backends
var SQLSchemas = map[string]string{
	"tx": `
		CREATE TABLE ` + "`" + `%indexName%` + "`" + ` (
			id CHAR(44) NOT NULL UNIQUE,
			ts DATETIME NOT NULL,
			blockno INTEGER UNSIGNED NOT NULL,
			` + "`" + `from` + "`" + ` VARCHAR(52) NOT NULL,
			` + "`" + `to` + "`" + ` VARCHAR(52),
			amount VARCHAR(30) NOT NULL,
			amount_float FLOAT(23) NOT NULL,
			type CHAR(1) NOT NULL,
			category ENUM(` + categories + `),
			PRIMARY KEY (id),
			INDEX tx_from (` + "`" + `from` + "`" + `(10)),
			INDEX tx_to (` + "`" + `to` + "`" + `(10)),
			INDEX tx_category (category),
			INDEX tx_blockno (blockno)
		);`,
	"block": `
		CREATE TABLE ` + "`" + `%indexName%` + "`" + ` (
			id CHAR(44) NOT NULL UNIQUE,
			ts DATETIME NOT NULL,
			no INTEGER UNSIGNED NOT NULL,
			txs MEDIUMINT UNSIGNED NOT NULL,
			size MEDIUMINT UNSIGNED NOT NULL,
			reward_account VARCHAR(52),
			reward_amount VARCHAR(30),
			PRIMARY KEY (id),
			INDEX block_no (no),
			INDEX reward_account (no)
		);`,
	"name": `
		CREATE TABLE ` + "`" + `%indexName%` + "`" + ` (
			id VARCHAR(60) NOT NULL UNIQUE,
			name VARCHAR(12) NOT NULL,
			address VARCHAR(52) NOT NULL,
			blockno INTEGER UNSIGNED NOT NULL,
			tx CHAR(44) NOT NULL,
			PRIMARY KEY (id),
			INDEX name_name (name),
			INDEX name_address (address)
		);`,
}
