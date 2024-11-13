package manager

import (
	"log"
	"poc-client/msg/request"
	"poc-client/utils/cryptoutil"
)

func (m *Manager) VerifyRequestWithSig(cID string, req request.RequestMsg) bool {
	// Have to verify both signatures

	requestBody := request.ReqBody{
		// ChannelID:      req.ChannelID,
		Amount:         req.Amount,
		LocalBlockHash: req.LocalBlockHash,
		ReqByte:        req.ReqByte,
	}

	paymentBody := request.PaymentBody{
		// ChannelID: req.ChannelID,
		Amount: req.Amount,
	}

	reqBSig := req.SignedReqBody
	payBSig := req.SignedPaymentBody

	if reqBSig[64] == 27 || reqBSig[64] == 28 {
		reqBSig[64] -= 27
	}

	if payBSig[64] == 27 || payBSig[64] == 28 {
		payBSig[64] -= 27
	}

	pubKB := m.GetClient(cID).PubKeyB
	rbFlag := cryptoutil.Verify(pubKB, requestBody.Keccak256Hash().Bytes(), req.SignedReqBody)
	pbFlag := cryptoutil.Verify(pubKB, paymentBody.PreHashByte(), req.SignedPaymentBody)
	log.Println("rbFlag: ", rbFlag, " pbFlag: ", pbFlag)
	return rbFlag && pbFlag
}
