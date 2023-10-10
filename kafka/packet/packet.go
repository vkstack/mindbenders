package packet

import (
	"context"
	"encoding/json"
	"reflect"

	"gitlab.com/dotpe/mindbenders/corel"
)

type Packet interface {
	Bytes() []byte
	Ctx() context.Context
	CorelId() *corel.CoRelationId
}

type packet struct {
	corelid *corel.CoRelationId
	ctx     context.Context

	Correlation string      `json:"correlation"`
	Payload     interface{} `json:"data"`
}

func (p *packet) Ctx() context.Context {
	return p.ctx
}

func (p *packet) CorelId() *corel.CoRelationId {
	return p.corelid
}

func (p *packet) Bytes() []byte {
	raw, _ := json.Marshal(p)
	return raw
}

// consumer side of the app logic will use this constructor
func NewPacket(raw []byte, dst interface{}) (Packet, error) {
	p := packet{
		Payload: dst,
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	p.corelid = corel.DecodeCorelationId(p.Correlation).Sibling()
	p.ctx = corel.NewContext(p.corelid)
	return &p, nil
}

func NewPacketFromStr(raw string, dst interface{}) (Packet, error) {
	return NewPacket([]byte(raw), dst)
}

// producer side of the app logic will use this constructor
func NewPacketFromEntity(ctx context.Context, entity interface{}) Packet {
	corelid, _ := corel.ReadCorelId(ctx)
	if corelid == nil {
		corelid = corel.NewCorelId("anonym:production-" + reflect.TypeOf(entity).Name())
	} else {
		corelid = corelid.Child()
	}
	p := packet{
		Payload:     entity,
		ctx:         ctx,
		corelid:     corelid,
		Correlation: corel.EncodeCorel(corelid),
	}
	return &p
}
