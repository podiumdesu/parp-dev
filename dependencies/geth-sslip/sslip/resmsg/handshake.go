package resmsg

type HandshakeMsg struct {
	Type             string `json:"type"`
	ServerPublicKeyB []byte `json:"serverPublicKeyB"`
	Signature        []byte `json:"signature"`
}

type HandshakeMsgBody struct {
	Type             string `json:"type"`
	ServerPublicKeyB []byte `json:"serverPublicKeyB"`
}

func (h *HandshakeMsg) HashBytes() []byte {
	data := HandshakeMsgBody{
		Type:             h.Type,
		ServerPublicKeyB: h.ServerPublicKeyB,
	}
	return hashData(data)
}
func (h *HandshakeMsg) Bytes() []byte {
	return marshalToJson(h)
}
