package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// #####################################################################################################################
// ROW
// #####################################################################################################################

type SysDlqRow = struct {
	Id                uint
	Topic             string
	GroupId           string
	Key               []byte
	Value             []byte
	Partition         int32
	Offset            int64
	Metadata          sql.NullString
	HeadersKeysJson   []byte
	HeadersValuesJson []byte
	Timestamp         time.Time
	TimestampType     int
	CreatedAt         time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var SysDlqTable = sysDlqTable{}

// ---------------------------------------------------------------------------------------------------------------------

type sysDlqTable struct{}

func (t sysDlqTable) Name() string                    { return "sys_dlq" }
func (t sysDlqTable) ColumnId() string                { return "id" }
func (t sysDlqTable) ColumnTopic() string             { return "topic" }
func (t sysDlqTable) ColumnGroupId() string           { return "group_id" }
func (t sysDlqTable) ColumnKey() string               { return "m_key" }
func (t sysDlqTable) ColumnValue() string             { return "m_value" }
func (t sysDlqTable) ColumnPartition() string         { return "m_partition" }
func (t sysDlqTable) ColumnOffset() string            { return "m_offset" }
func (t sysDlqTable) ColumnMetadata() string          { return "m_metadata" }
func (t sysDlqTable) ColumnHeadersKeysJson() string   { return "m_headers_keys_json" }
func (t sysDlqTable) ColumnHeadersValuesJson() string { return "m_headers_values_json" }
func (t sysDlqTable) ColumnTimestamp() string         { return "m_timestamp" }
func (t sysDlqTable) ColumnTimestampType() string     { return "m_timestamp_type" }
func (t sysDlqTable) ColumnCreatedAt() string         { return "created_at" }

// ---------------------------------------------------------------------------------------------------------------------

// Insert
//
// Паникует при нулевых аргументах.
func (t sysDlqTable) Insert(ctx context.Context, tx *sql.Tx, row *SysDlqRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*SysDlqRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnId(),
			t.ColumnTopic(),
			t.ColumnGroupId(),
			t.ColumnKey(),
			t.ColumnValue(),
			t.ColumnPartition(),
			t.ColumnOffset(),
			t.ColumnMetadata(),
			t.ColumnHeadersKeysJson(),
			t.ColumnHeadersValuesJson(),
			t.ColumnTimestamp(),
			t.ColumnTimestampType(),
			t.ColumnCreatedAt(),
		).
		Values(
			row.Id,
			row.Topic,
			row.GroupId,
			row.Key,
			row.Value,
			row.Partition,
			row.Offset,
			row.Metadata,
			row.HeadersKeysJson,
			row.HeadersValuesJson,
			row.Timestamp,
			row.TimestampType,
			row.CreatedAt,
		).
		String(),
		// --
		row.Id,
		row.Topic,
		row.GroupId,
		row.Key,
		row.Value,
		row.Partition,
		row.Offset,
		row.Metadata,
		row.HeadersKeysJson,
		row.HeadersValuesJson,
		row.Timestamp,
		row.TimestampType,
		row.CreatedAt,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------

// HasTopicGroup
//
// Паникует при нулевых аргументах.
func (t sysDlqTable) HasTopicGroup(ctx context.Context, tx *sql.Tx, topic string, groupId string) (bool, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Str().NotEmpty().Must(topic)
	assert.Str().NotEmpty().Must(groupId)

	row := &SysDlqRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(t.ColumnId()).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ? AND %s = ?", t.ColumnTopic(), t.ColumnGroupId())).
			String(),
			// --
			topic,
			groupId,
		).
		Scan(&row.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return row.Id > 0, err
}

// ---------------------------------------------------------------------------------------------------------------------
