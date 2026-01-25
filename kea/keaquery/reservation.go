package keaquery

import (
	"encoding/json"
	"fmt"
)

type queryKind int

const (
	rqUnknown queryKind = iota
	rqByIP
	rqByIdentifier
)

// ReservationQuery is a single type representing either:
//   - { "subnet-id", "ip-address" }
//   - { "subnet-id", "identifier-type", "identifier" }
type ReservationQuery struct {
	SubnetID uint32 `json:"subnet-id"`

	kind queryKind

	ip             string
	identifier     string
	identifierType string
}

func ReservationByIP(subnetID uint32, ip string) ReservationQuery {
	return ReservationQuery{
		SubnetID: subnetID,
		kind:     rqByIP,
		ip:       ip,
	}
}

func ReservationByIdentifier(subnetID uint32, idType, identifier string) ReservationQuery {
	return ReservationQuery{
		SubnetID:       subnetID,
		kind:           rqByIdentifier,
		identifierType: idType,
		identifier:     identifier,
	}
}

func (r ReservationQuery) IP() (ip string, ok bool) {
	if r.kind != rqByIP {
		return "", false
	}
	return r.ip, true
}

func (r ReservationQuery) Identifier() (idType, id string, ok bool) {
	if r.kind != rqByIdentifier {
		return "", "", false
	}
	return r.identifierType, r.identifier, true
}

func (r ReservationQuery) MarshalJSON() ([]byte, error) {
	switch r.kind {
	case rqByIP:
		aux := struct {
			SubnetID  uint32 `json:"subnet-id"`
			IPAddress string `json:"ip-address"`
		}{
			SubnetID:  r.SubnetID,
			IPAddress: r.ip,
		}
		return json.Marshal(aux)
	case rqByIdentifier:
		aux := struct {
			SubnetID       uint32 `json:"subnet-id"`
			IdentifierType string `json:"identifier-type"`
			Identifier     string `json:"identifier"`
		}{
			SubnetID:       r.SubnetID,
			IdentifierType: r.identifierType,
			Identifier:     r.identifier,
		}
		return json.Marshal(aux)
	default:
		return nil, fmt.Errorf("ReservationQuery: invalid variant (neither ip nor identifier set)")
	}
}
