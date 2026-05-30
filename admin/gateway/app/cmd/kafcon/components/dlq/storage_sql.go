package dlq

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/admin/gateway/internal/infra/sql/db"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ StorageInterface = StorageSQL{}

type StorageSQL struct {
	txr txr.TxrInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewStorageSQL
//
// Паникует при нулевых аргументах.
func NewStorageSQL(txr txr.TxrInterface) StorageSQL {
	assert.NotNilDeepMust(txr)

	return StorageSQL{
		txr: txr,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Add
//
// Паникует при нулевых аргументах.
//
// Для маркировки потенциально любой группы топика используется AnyGroup.
//
// Ошибки:
//   - std.ErrorRuntime
func (s StorageSQL) Add(ctx context.Context, topic string, groupId string, kMsg *kafka.Message) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)
	assert.Str().NotEmpty().Must(groupId)
	assert.NotNilDeepMust(kMsg)

	_m_ := "Add"

	if err := s.txr.Tx(ctx, func(ctx context.Context) error {
		tx := txr.TxFromCtx(ctx).(*sql.Tx)
		return db.SysDlqTable.Insert(ctx, tx, &db.SysDlqRow{
			Topic:     topic,
			GroupId:   groupId,
			Key:       kMsg.Key,
			Value:     kMsg.Value,
			Partition: kMsg.TopicPartition.Partition,
			Offset:    int64(kMsg.TopicPartition.Offset),
			Metadata: func() sql.NullString {
				if kMsg.TopicPartition.Metadata == nil {
					return sql.NullString{}
				}
				return sql.NullString{
					String: *kMsg.TopicPartition.Metadata,
					Valid:  true,
				}
			}(),
			HeadersKeysJson: func() []byte {
				keys := make([]string, len(kMsg.Headers))
				for i, h := range kMsg.Headers {
					keys[i] = h.Key
				}
				res, err := json.Marshal(keys)
				if err != nil {
					panic(err)
				}
				return res
			}(),
			HeadersValuesJson: func() []byte {
				values := make([][]byte, len(kMsg.Headers))
				for i, h := range kMsg.Headers {
					values[i] = h.Value
				}
				res, err := json.Marshal(values)
				if err != nil {
					panic(err)
				}
				return res
			}(),
			Timestamp:     kMsg.Timestamp,
			TimestampType: int(kMsg.TimestampType),
			CreatedAt:     time.Now(),
		})
	}); err != nil {
		return std.WrapErrorToRuntime(err, s, _m_)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// IsGroupPoisoned
//
// Паникует при нулевых аргументах.
//
// Для маркировки потенциально любой группы топика используется AnyGroup.
//
// Ошибки:
//   - std.ErrorRuntime
func (s StorageSQL) IsGroupPoisoned(ctx context.Context, topic string, groupId string) (bool, error) {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)
	assert.Str().NotEmpty().Must(groupId)

	_m_ := "IsGroupPoisoned"

	var is bool
	if err := s.txr.Tx(ctx, func(ctx context.Context) error {
		tx := txr.TxFromCtx(ctx).(*sql.Tx)
		_is, err := db.SysDlqTable.HasTopicGroup(ctx, tx, topic, groupId)
		is = _is
		return err
	}); err != nil {
		return false, std.WrapErrorToRuntime(err, s, _m_)
	}

	return is, nil
}

// ---------------------------------------------------------------------------------------------------------------------
