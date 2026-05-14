package kafapi

import (
	"errors"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Meta
// ---------------------------------------------------------------------------------------------------------------------

type Meta = struct {
	Id          uint
	Type        string
	Version     uint
	GroupId     string
	OperationId string
	OccurredAt  time.Time
}

func mapToMeta(rawData *RawData) (*Meta, error) {
	if err := errors.Join(
		assert.Num[uint]().Positive().Check(rawData.Id, "Id"),
		assert.Str().NotEmpty().Check(rawData.Type, "Type"),
		assert.Num[uint]().Positive().Check(rawData.Version, "Version"),
		assert.Str().NotEmpty().Check(rawData.OutboxGroupId, "GroupId"),
		assert.Str().NotEmpty().Check(rawData.OutboxOperationId, "OperationId"),
		assert.Time().NotZero().Check(rawData.OccurredAt, "OccurredAt"),
	); err != nil {
		return nil, err
	}

	return &Meta{
		Id:          rawData.Id,
		Type:        rawData.Type,
		Version:     rawData.Version,
		GroupId:     rawData.OutboxGroupId,
		OperationId: rawData.OutboxOperationId,
		OccurredAt:  rawData.OccurredAt,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Data
//
// Data-типы нельзя делать alias'ами, поскольку некоторые из них могут содержать одинаковый набор полей --
// при таком совпадении switch v.(type) будет расценивать их как одинаковые (alias'ы же) и сигнализировать о дубликатах.
// ---------------------------------------------------------------------------------------------------------------------

const dataTypeAccountCreated = "account_created"

type DataAccountCreatedV1 struct {
	AccountId string
	Email     std.Email
}

func mapToAccountCreatedV1(rawData *RawData) (*DataAccountCreatedV1, error) {
	if err := errors.Join(
		assert.Str().NotEmpty().Check(rawData.ExtraAccountId, "ExtraAccountId"),
		assert.Cmp[std.Email]().NotEq(std.EmailNil).Check(rawData.ExtraAccountEmail, "ExtraAccountEmail"),
	); err != nil {
		return nil, err
	}

	return &DataAccountCreatedV1{
		AccountId: rawData.ExtraAccountId,
		Email:     rawData.ExtraAccountEmail,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

const dataTypeAccountDeactivated = "account_deactivated"

type DataAccountDeactivatedV1 struct {
	AccountId string
}

func mapToAccountDeactivatedV1(rawData *RawData) (*DataAccountDeactivatedV1, error) {
	if err := assert.Str().NotEmpty().Check(rawData.ExtraAccountId, "ExtraAccountId"); err != nil {
		return nil, err
	}

	return &DataAccountDeactivatedV1{
		AccountId: rawData.ExtraAccountId,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

const dataTypeAccountIpWhitelistChanged = "account_ip_whitelist_changed"

type DataIpWhitelistChangedV1 struct {
	AccountId      string
	NewIpWhitelist []netip.Addr
}

func mapToIpWhitelistChangedV1(rawData *RawData) (*DataIpWhitelistChangedV1, error) {
	if err := errors.Join(
		assert.Str().NotEmpty().Check(rawData.ExtraAccountId, "ExtraAccountId"),
		assert.TrueCheck(rawData.FlagExtraAccountIpWhitelistIsSet, "FlagExtraAccountIpWhitelistIsSet"),
	); err != nil {
		return nil, err
	}

	return &DataIpWhitelistChangedV1{
		AccountId:      rawData.ExtraAccountId,
		NewIpWhitelist: rawData.ExtraAccountIpWhitelist,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

const dataTypeSessionCreated = "session_created"

type DataSessionCreatedV1 struct {
	SessionId string
	AccountId string
}

func mapToSessionCreatedV1(rawData *RawData) (*DataSessionCreatedV1, error) {
	if err := errors.Join(
		assert.Str().NotEmpty().Check(rawData.ExtraSessionId, "ExtraSessionId"),
		assert.Str().NotEmpty().Check(rawData.ExtraAccountId, "ExtraAccountId"),
	); err != nil {
		return nil, err
	}

	return &DataSessionCreatedV1{
		SessionId: rawData.ExtraSessionId,
		AccountId: rawData.ExtraAccountId,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

const dataTypeSessionClosed = "session_closed"

type DataSessionClosedV1 struct {
	SessionId string
	AccountId string
}

func mapToSessionClosedV1(rawData *RawData) (*DataSessionClosedV1, error) {
	if err := errors.Join(
		assert.Str().NotEmpty().Check(rawData.ExtraSessionId, "ExtraSessionId"),
		assert.Str().NotEmpty().Check(rawData.ExtraAccountId, "ExtraAccountId"),
	); err != nil {
		return nil, err
	}

	return &DataSessionClosedV1{
		SessionId: rawData.ExtraSessionId,
		AccountId: rawData.ExtraAccountId,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
