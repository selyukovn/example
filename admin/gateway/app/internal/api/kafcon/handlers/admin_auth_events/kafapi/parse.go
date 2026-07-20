package kafapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"net/netip"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

type RawData = struct {
	Id                               uint
	OccurredAt                       time.Time
	Type                             string
	Version                          uint
	ExtraAccountEmail                std.Email
	ExtraAccountId                   string
	ExtraAccountIpWhitelist          []netip.Addr
	FlagExtraAccountIpWhitelistIsSet bool
	ExtraSessionId                   string
	CreatedAt                        time.Time
	OutboxGroupId                    string
	OutboxOperationId                string
}

// ---------------------------------------------------------------------------------------------------------------------

// parse
//
// Ошибки:
//   - ErrorDecoding
//   - ErrorMapping
//   - ErrorUnsupported
func parse(kMsg *kafka.Message) (*Meta, any, error) {
	// ---- decode ----
	rawData, err := parseDecode(kMsg)
	if err != nil {
		return nil, nil, newErrorDecoding(err)
	}

	meta, err := mapToMeta(rawData)
	if err != nil {
		// Внимание!
		// Отсутствие метаданных или ошибка их маппинга приравниваются к ошибке декодирования.
		return nil, nil, newErrorDecoding(err)
	}
	// ---- /decode ----

	// ---- account created ----
	if rawData.Type == dataTypeAccountCreated && rawData.Version >= 1 {
		data, err := mapToAccountCreatedV1(rawData)
		if err != nil {
			return nil, nil, newErrorMapping(err, meta, rawData)
		}
		return meta, data, nil
	}
	// ---- /account created ----

	// ---- account deactivated ----
	if rawData.Type == dataTypeAccountDeactivated && rawData.Version >= 1 {
		data, err := mapToAccountDeactivatedV1(rawData)
		if err != nil {
			return nil, nil, newErrorMapping(err, meta, rawData)
		}
		return meta, data, nil
	}
	// ---- /account deactivated ----

	// ---- ip whitelist changed ----
	if rawData.Type == dataTypeAccountIpWhitelistChanged && rawData.Version >= 1 {
		data, err := mapToIpWhitelistChangedV1(rawData)
		if err != nil {
			return nil, nil, newErrorMapping(err, meta, rawData)
		}
		return meta, data, nil
	}
	// ---- /ip whitelist changed ----

	// ---- session created ----
	if rawData.Type == dataTypeSessionCreated && rawData.Version >= 1 {
		data, err := mapToSessionCreatedV1(rawData)
		if err != nil {
			return nil, nil, newErrorMapping(err, meta, rawData)
		}
		return meta, data, nil
	}
	// ---- /session created ----

	// ---- session closed ----
	if rawData.Type == dataTypeSessionClosed && rawData.Version >= 1 {
		data, err := mapToSessionClosedV1(rawData)
		if err != nil {
			return nil, nil, newErrorMapping(err, meta, rawData)
		}
		return meta, data, nil
	}
	// ---- /session closed ----

	// undefined
	return nil, nil, newErrorUnsupported(meta, rawData)
}

func parseDecode(kMsg *kafka.Message) (*RawData, error) {
	_om_ := "kafapi.parseDecode"

	encodedDataMess := &struct {
		Payload struct {
			After struct {
				Id                          uint   `json:"id"`
				OccurredAt                  int    `json:"occurred_at"`
				Type                        string `json:"type"`
				Version                     uint   `json:"version"`
				ExtraAccountEmail           string `json:"extra_account_email"`
				ExtraAccountId              string `json:"extra_account_id"`
				ExtraAccountIpWhitelistJson string `json:"extra_account_ip_whitelist_json"`
				ExtraSessionId              string `json:"extra_session_id"`
				CreatedAt                   int    `json:"created_at"`
				OutboxGroupId               string `json:"outbox_group_id"`
				OutboxOperationId           string `json:"outbox_operation_id"`
			} `json:"after"`
		} `json:"payload"`
	}{}

	err := json.Unmarshal(kMsg.Value, encodedDataMess)
	if err != nil {
		return nil, fmt.Errorf("%s/%s: %w", _om_, "Unmarshal", err)
	}

	// ...
	// "id": 7,
	// "occurred_at": 1777892107000,
	// "type": "c2Vzc2lvbl9jcmVhdGVk",
	// "version": 1,
	// "extra_account_email": null,
	// "extra_account_id": "MmQ0NjRiYzctNDRkYS0xMWYxLWI4NjYtMDI0MmFjMTgwMDAy",
	// "extra_account_ip_whitelist_json": null,
	// "extra_session_id": "NGExYzY2MmEtM2I3NC00YmY2LTg2ODctZTA2YmJiMmQ3ZWNk",
	// "created_at": 1777892107000,
	// "outbox_group_id": "MmQ0NjRiYzctNDRkYS0xMWYxLWI4NjYtMDI0MmFjMTgwMDAy",
	// "outbox_operation_id": "M2RlZjgwYWFmODRlNGJjY2JkMmI4ZmY5OGEzMmZhNWE="
	// ...
	encodedData := encodedDataMess.Payload.After

	// ---- Id ----
	Id := encodedData.Id
	// ---- /Id ----

	// ---- OccurredAt ----
	OccurredAt := time.Unix(int64(encodedData.OccurredAt), 0)
	// ---- /OccurredAt ----

	// ---- CreatedAt ----
	CreatedAt := time.Unix(int64(encodedData.CreatedAt), 0)
	// ---- /CreatedAt ----

	// ---- Type ----
	typeBytes, err := base64.StdEncoding.DecodeString(encodedData.Type)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "Type", "DecodeString", err)
	}
	Type := string(typeBytes)
	// ---- /Type ----

	// ---- ExtraAccountEmail ----
	ExtraAccountEmail := std.EmailNil
	extraAccountEmailBytes, err := base64.StdEncoding.DecodeString(encodedData.ExtraAccountEmail)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraAccountEmail", "DecodeString", err)
	}
	extraAccountEmailString := string(extraAccountEmailBytes)
	if extraAccountEmailString != "" {
		ExtraAccountEmail, err = std.EmailFromString(extraAccountEmailString)
		if err != nil {
			return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraAccountEmail", "EmailFromString", err)
		}
	}
	// ---- /ExtraAccountEmail ----

	// ---- ExtraAccountId ----
	extraAccountIdBytes, err := base64.StdEncoding.DecodeString(encodedData.ExtraAccountId)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraAccountId", "DecodeString", err)
	}
	ExtraAccountId := string(extraAccountIdBytes)
	// ---- /ExtraAccountId ----

	// ---- ExtraAccountIpWhitelist ----
	var ExtraAccountIpWhitelist []netip.Addr = nil
	FlagExtraAccountIpWhitelistIsSet := encodedData.ExtraAccountIpWhitelistJson != ""
	extraAccountIpWhitelistJsonBytes, err := base64.StdEncoding.DecodeString(encodedData.ExtraAccountIpWhitelistJson)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraAccountIpWhitelist", "DecodeString", err)
	}
	if len(extraAccountIpWhitelistJsonBytes) > 0 && string(extraAccountIpWhitelistJsonBytes) != "[]" {
		extraAccountIpWhitelistSliceStr := make([]string, 0)
		err = json.Unmarshal(extraAccountIpWhitelistJsonBytes, &extraAccountIpWhitelistSliceStr)
		if err != nil {
			return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraAccountIpWhitelist", "Unmarshal", err)
		}
		ExtraAccountIpWhitelist = make([]netip.Addr, len(extraAccountIpWhitelistSliceStr))
		for i, ipStr := range extraAccountIpWhitelistSliceStr {
			ExtraAccountIpWhitelist[i], err = netip.ParseAddr(ipStr)
			if err != nil {
				return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraAccountIpWhitelist", "ParseAddr", err)
			}
		}
	}
	// ---- /ExtraAccountIpWhitelist ----

	// ---- ExtraSessionId ----
	extraSessionIdBytes, err := base64.StdEncoding.DecodeString(encodedData.ExtraSessionId)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "ExtraSessionId", "DecodeString", err)
	}
	ExtraSessionId := string(extraSessionIdBytes)
	// ---- /ExtraSessionId ----

	// ---- OutboxGroupId ----
	outboxGroupIdBytes, err := base64.StdEncoding.DecodeString(encodedData.OutboxGroupId)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "OutboxGroupId", "DecodeString", err)
	}
	OutboxGroupId := string(outboxGroupIdBytes)
	// ---- /OutboxGroupId ----

	// ---- OutboxOperationId ----
	outboxOperationIdBytes, err := base64.StdEncoding.DecodeString(encodedData.OutboxOperationId)
	if err != nil {
		return nil, fmt.Errorf("%s/%s/%s: %w", _om_, "OutboxOperationId", "DecodeString", err)
	}
	OutboxOperationId := string(outboxOperationIdBytes)
	// ---- /OutboxOperationId ----

	return &RawData{
		Id:                               Id,
		OccurredAt:                       OccurredAt,
		Type:                             Type,
		Version:                          encodedData.Version,
		ExtraAccountEmail:                ExtraAccountEmail,
		ExtraAccountId:                   ExtraAccountId,
		ExtraAccountIpWhitelist:          ExtraAccountIpWhitelist,
		FlagExtraAccountIpWhitelistIsSet: FlagExtraAccountIpWhitelistIsSet,
		ExtraSessionId:                   ExtraSessionId,
		CreatedAt:                        CreatedAt,
		OutboxGroupId:                    OutboxGroupId,
		OutboxOperationId:                OutboxOperationId,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
