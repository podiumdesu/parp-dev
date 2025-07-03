package handlers

const contractABI = `[
    {
        "inputs": [
            {
                "internalType": "address",
                "name": "_depositContractAddress",
                "type": "address"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "constructor"
    },
    {
        "inputs": [],
        "name": "ECDSAInvalidSignature",
        "type": "error"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "length",
                "type": "uint256"
            }
        ],
        "name": "ECDSAInvalidSignatureLength",
        "type": "error"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "s",
                "type": "bytes32"
            }
        ],
        "name": "ECDSAInvalidSignatureS",
        "type": "error"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "internalType": "bytes32",
                "name": "channelId",
                "type": "bytes32"
            }
        ],
        "name": "ChannelOpened",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "internalType": "address",
                "name": "recoveredSigner",
                "type": "address"
            },
            {
                "indexed": true,
                "internalType": "bytes32",
                "name": "channelId",
                "type": "bytes32"
            }
        ],
        "name": "DebugInfo",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "",
                "type": "uint256"
            }
        ],
        "name": "LogUint",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "internalType": "bytes32",
                "name": "ChannelId",
                "type": "bytes32"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "Amount",
                "type": "uint256"
            },
            {
                "indexed": false,
                "internalType": "bytes",
                "name": "Signature",
                "type": "bytes"
            }
        ],
        "name": "decodedPaymentBody",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "internalType": "bytes32",
                "name": "msgHash",
                "type": "bytes32"
            }
        ],
        "name": "msgHash",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "internalType": "address",
                "name": "signer",
                "type": "address"
            }
        ],
        "name": "signerAdd",
        "type": "event"
    },
    {
        "inputs": [
            {
                "internalType": "address",
                "name": "addr",
                "type": "address"
            }
        ],
        "name": "balance",
        "outputs": [
            {
                "internalType": "uint256",
                "name": "bal",
                "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes",
                "name": "serializedPaymentBody",
                "type": "bytes"
            }
        ],
        "name": "closeChan",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "channelId",
                "type": "bytes32"
            },
            {
                "internalType": "uint256",
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "confirmClosure",
        "outputs": [],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes",
                "name": "serializedPaymentBody",
                "type": "bytes"
            }
        ],
        "name": "decodePaymentBody",
        "outputs": [
            {
                "components": [
                    {
                        "internalType": "bytes32",
                        "name": "ChannelId",
                        "type": "bytes32"
                    },
                    {
                        "internalType": "uint256",
                        "name": "Amount",
                        "type": "uint256"
                    },
                    {
                        "internalType": "bytes",
                        "name": "Signature",
                        "type": "bytes"
                    }
                ],
                "internalType": "struct paychan.RequestPaymentBody",
                "name": "",
                "type": "tuple"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "depositContract",
        "outputs": [
            {
                "internalType": "contract IDepositContract",
                "name": "",
                "type": "address"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "channelID",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
            }
        ],
        "name": "encodeData",
        "outputs": [
            {
                "internalType": "bytes32",
                "name": "",
                "type": "bytes32"
            }
        ],
        "stateMutability": "pure",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "ChannelId",
                "type": "bytes32"
            },
            {
                "internalType": "uint256",
                "name": "Amount",
                "type": "uint256"
            }
        ],
        "name": "getPaymentHash",
        "outputs": [
            {
                "internalType": "bytes32",
                "name": "",
                "type": "bytes32"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "address",
                "name": "to",
                "type": "address"
            },
            {
                "internalType": "uint256",
                "name": "senderDeposit",
                "type": "uint256"
            }
        ],
        "name": "openChan",
        "outputs": [
            {
                "internalType": "bytes32",
                "name": "channelId",
                "type": "bytes32"
            }
        ],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "channelId",
                "type": "bytes32"
            }
        ],
        "name": "paychanCheck",
        "outputs": [
            {
                "components": [
                    {
                        "internalType": "bytes32",
                        "name": "id",
                        "type": "bytes32"
                    },
                    {
                        "internalType": "address payable",
                        "name": "sender",
                        "type": "address"
                    },
                    {
                        "internalType": "address payable",
                        "name": "recipient",
                        "type": "address"
                    },
                    {
                        "internalType": "uint256",
                        "name": "senderDeposit",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "startTime",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "status",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "fee",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "disputeStartTime",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "disputeDuration",
                        "type": "uint256"
                    },
                    {
                        "internalType": "bool",
                        "name": "senderConfirm",
                        "type": "bool"
                    },
                    {
                        "internalType": "bool",
                        "name": "recipientConfirm",
                        "type": "bool"
                    }
                ],
                "internalType": "struct paychan.PayChan",
                "name": "",
                "type": "tuple"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "channelId",
                "type": "bytes32"
            }
        ],
        "name": "paychanSelectedArguments",
        "outputs": [
            {
                "internalType": "address",
                "name": "sender",
                "type": "address"
            },
            {
                "internalType": "address",
                "name": "rec",
                "type": "address"
            },
            {
                "internalType": "uint256",
                "name": "status",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "senderB",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "fee",
                "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "id",
                "type": "bytes32"
            },
            {
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
            },
            {
                "internalType": "bytes",
                "name": "sig",
                "type": "bytes"
            },
            {
                "internalType": "address",
                "name": "signer",
                "type": "address"
            }
        ],
        "name": "verifyPayment",
        "outputs": [
            {
                "internalType": "bool",
                "name": "",
                "type": "bool"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`
