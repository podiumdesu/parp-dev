package manager

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
)

func (m *Manager) VerifyRequestWithSig(cID string, req request.RequestMsg) (bool, common.Hash) {
	// Have to verify both signatures

	requestBody := request.ReqBody{
		// ChannelID:      req.ChannelID,
		Amount:         req.Amount,
		LocalBlockHash: req.LocalBlockHash,
		ReqByte:        req.ReqByte,
	}

	reqBodyHash := requestBody.Keccak256Hash()

	paymentBody := request.PaymentBody{
		// ChannelID: req.ChannelID,
		Amount: req.Amount,
	}

	// I spent more than 4 hours here to debug, because the signature received from the client always got altered
	// I FIANLLY KNOW THE REASON:
	// When you modify reqBSig[64] -= 27 inside VerifyRequestWithSig, it alters the original slice req.SignedReqBody because slices are reference types.

	// reqBSig := req.SignedReqBody
	// payBSig := req.SignedPaymentBody

	reqBSig := append([]byte{}, req.SignedReqBody...)
	payBSig := append([]byte{}, req.SignedPaymentBody...)

	if reqBSig[64] == 27 || reqBSig[64] == 28 {
		reqBSig[64] -= 27
	}

	if payBSig[64] == 27 || payBSig[64] == 28 {
		payBSig[64] -= 27
	}

	pubKB := m.GetClient(cID).PubKeyB
	rbFlag := cryptoutil.Verify(pubKB, requestBody.Keccak256Hash().Bytes(), reqBSig)
	pbFlag := cryptoutil.Verify(pubKB, paymentBody.PreHashByte(), payBSig)
	log.Println("rbFlag: ", rbFlag, " pbFlag: ", pbFlag)
	return rbFlag && pbFlag, reqBodyHash
}
