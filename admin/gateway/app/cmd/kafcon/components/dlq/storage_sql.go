package dlq

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/admin/gateway/internal/infra/sql/db"
	"fmt"
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

// Remove
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (s StorageSQL) Remove(ctx context.Context, kMsg *kafka.Message) error {
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(kMsg)

	_m_ := "Remove"

	if err := s.txr.Tx(ctx, func(ctx context.Context) error {
		tx := txr.TxFromCtx(ctx).(*sql.Tx)
		return db.SysDlqTable.DeleteByTopicPartitionOffsetTimestamp(
			ctx,
			tx,
			*kMsg.TopicPartition.Topic,
			kMsg.TopicPartition.Partition,
			int64(kMsg.TopicPartition.Offset),
			kMsg.Timestamp,
		)
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

func (s StorageSQL) rowToKMsg(row *db.SysDlqRow) (*kafka.Message, error) {
	_m_ := "rowToKMsg"

	// --

	var headersKeys []string
	err := json.Unmarshal(row.HeadersKeysJson, &headersKeys)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, s, _m_, "headersKeys", fmt.Sprintf("Id=%d", row.Id))
	}

	var headersValues [][]byte
	err = json.Unmarshal(row.HeadersValuesJson, &headersValues)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, s, _m_, "headersValues", fmt.Sprintf("Id=%d", row.Id))
	}

	if len(headersKeys) != len(headersValues) {
		return nil, std.NewErrorRuntimeFf("len(headersKeys) != len(headersValues), Id=%d", row.Id)
	}

	headers := make([]kafka.Header, len(headersKeys))
	for i, key := range headersKeys {
		value := headersValues[i]
		headers[i] = kafka.Header{
			Key:   key,
			Value: value,
		}
	}

	// --

	var metadata *string = nil
	if row.Metadata.Valid {
		metadata = &row.Metadata.String
	}

	// --

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:       &row.Topic,
			Partition:   row.Partition,
			Offset:      kafka.Offset(row.Offset),
			Metadata:    metadata,
			Error:       nil,
			LeaderEpoch: nil,
		},
		Value:         row.Value,
		Key:           row.Key,
		Timestamp:     row.Timestamp,
		TimestampType: kafka.TimestampType(row.TimestampType),
		Opaque:        nil,
		Headers:       headers,
	}, nil
}

// GetMessages
//
// В возвращенный канал будет последовательно записано `limit` сообщений (или все при `limit` == 0),
// если таковые найдутся по топику `topic` и группе `groupId`,
// в порядке их поступления в DLQ, после чего канал закроется.
//
// Ошибка извлечения будет передана в канал отдельным элементом, после чего извлечение прервется, и канал закроется.
//
// Завершение контекста прерывает извлечение и закрывает канал.
//
// Паникует при нулевых аргументах:
//   - ctx
//   - topic
//   - groupId
//
// Ошибки:
//   - std.ErrorRuntime
func (s StorageSQL) GetMessages(
	ctx context.Context,
	topic string,
	groupId string,
	limit uint,
) <-chan GetMessagesItem {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)
	assert.Str().NotEmpty().Must(groupId)

	_m_ := "GetMessages"

	rChan := make(chan GetMessagesItem)

	go func() {
		defer close(rChan)

		var totalFetched uint = 0
		var batch []*db.SysDlqRow
		var batchSize uint = 10
		var afterId uint = 0
		for {
			if err := s.txr.Tx(ctx, func(ctx context.Context) error {
				tx := txr.TxFromCtx(ctx).(*sql.Tx)
				rows, err := db.SysDlqTable.QueryByTopicGroupAfterIdAscId(ctx, tx, topic, groupId, afterId, batchSize)
				if err != nil {
					return err
				}
				batch = rows
				return nil
			}); err != nil {
				select {
				case <-ctx.Done():
				case rChan <- getMessagesMakeItemError(std.WrapErrorToRuntime(err, s, _m_, "query")):
				}
				return
			}

			if len(batch) == 0 {
				return
			}

			for _, row := range batch {
				kMsg, err := s.rowToKMsg(row)
				if err != nil {
					select {
					case <-ctx.Done():
					case rChan <- getMessagesMakeItemError(std.WrapErrorToRuntime(err, s, _m_, "rowToKMsg")):
					}
					return
				}

				select {
				case <-ctx.Done():
					return
				case rChan <- getMessagesMakeItemMsg(kMsg):
					// pass
				}

				afterId = row.Id

				totalFetched += 1
				if totalFetched == limit {
					return
				}
			}
		}
	}()

	return rChan
}

func getMessagesMakeItemError(err error) GetMessagesItem {
	return GetMessagesItem{
		KMsg: nil,
		Err:  err,
	}
}

func getMessagesMakeItemMsg(kMsg *kafka.Message) GetMessagesItem {
	return GetMessagesItem{
		KMsg: kMsg,
		Err:  nil,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
