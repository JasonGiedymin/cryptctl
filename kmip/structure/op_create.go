package structure

import (
	"github.com/HouzuoGuo/cryptctl/kmip/ttlv"
)

// KMIP request message 420078
type SCreateRequest struct {
	SRequestHeader    SRequestHeader    // IBatchCount is assumed to be 1 in serialisation operations
	SRequestBatchItem SRequestBatchItem // payload is SRequestPayloadCreate
}

func (createReq *SCreateRequest) SerialiseToTTLV() ttlv.Item {
	createReq.SRequestHeader.IBatchCount.Value = 1
	ret := ttlv.NewStructure(TagRequestMessage, createReq.SRequestHeader.SerialiseToTTLV(), createReq.SRequestBatchItem.SerialiseToTTLV())
	return ret
}
func (createReq *SCreateRequest) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagRequestMessage, TagRequestHeader, &createReq.SRequestHeader); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagRequestMessage, TagBatchItem, &createReq.SRequestBatchItem); err != nil {
		return err
	}
	return nil
}

// 420079
type SRequestPayloadCreate struct {
	EObjectType        ttlv.Enumeration   // 420057
	STemplateAttribute STemplateAttribute // 420091
}

func (createPayload *SRequestPayloadCreate) SerialiseToTTLV() ttlv.Item {
	createPayload.EObjectType.Tag = TagObjectType
	return ttlv.NewStructure(TagRequestPayload, &createPayload.EObjectType, createPayload.STemplateAttribute.SerialiseToTTLV())
}
func (createPayload *SRequestPayloadCreate) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagRequestPayload, TagObjectType, &createPayload.EObjectType); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagRequestPayload, TagTemplateAttribute, &createPayload.STemplateAttribute); err != nil {
		return err
	}
	return nil
}

// 42000b of a create request's payload attribute called "Name"
type SCreateRequestNameAttributeValue struct {
	KeyName ttlv.Text        // 420055
	KeyType ttlv.Enumeration // 420054
}

func (nameAttr *SCreateRequestNameAttributeValue) SerialiseToTTLV() ttlv.Item {
	nameAttr.KeyName.Tag = TagNameValue
	nameAttr.KeyType.Tag = TagNameType
	return ttlv.NewStructure(TagAttributeValue, &nameAttr.KeyName, &nameAttr.KeyType)
}
func (nameAttr *SCreateRequestNameAttributeValue) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagAttribute, TagNameValue, &nameAttr.KeyName); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagAttribute, TagNameType, &nameAttr.KeyType); err != nil {
		return err
	}
	return nil
}

// KMIP response message 42007b
type SCreateResponse struct {
	SHeader            SResponseHeader // IBatchCount is assumed to be 1 in serialisation operations
	SResponseBatchItem SResponseBatchItem
}

func (createResp *SCreateResponse) SerialiseToTTLV() ttlv.Item {
	createResp.SHeader.IBatchCount.Value = 1
	ret := ttlv.NewStructure(TagResponseMessage, createResp.SHeader.SerialiseToTTLV(), createResp.SResponseBatchItem.SerialiseToTTLV())
	return ret
}
func (createResp *SCreateResponse) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagResponseMessage, TagResponseHeader, &createResp.SHeader); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagResponseMessage, TagBatchItem, &createResp.SResponseBatchItem); err != nil {
		return err
	}
	return nil
}

// 42007c - response payload from a create response
type SResponsePayloadCreate struct {
	EObjectType ttlv.Enumeration // 420057
	TUniqueID   ttlv.Text        // 420094
}

func (createPayload *SResponsePayloadCreate) SerialiseToTTLV() ttlv.Item {
	createPayload.EObjectType.Tag = TagObjectType
	createPayload.TUniqueID.Tag = TagUniqueID
	return ttlv.NewStructure(TagResponsePayload, &createPayload.EObjectType, &createPayload.TUniqueID)
}
func (createPayload *SResponsePayloadCreate) DeserialiseFromTTLV(in ttlv.Item) error {
	if err := DecodeStructItem(in, TagResponsePayload, TagObjectType, &createPayload.EObjectType); err != nil {
		return err
	} else if err := DecodeStructItem(in, TagResponsePayload, TagUniqueID, &createPayload.TUniqueID); err != nil {
		return err
	}
	return nil
}
