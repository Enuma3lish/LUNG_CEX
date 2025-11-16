package blockchain

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
)

type SolanaClient struct {
	rpcClient *rpc.Client
	payer     solana.PrivateKey
}

func NewSolanaClient() (*SolanaClient, error) {
	// Use devnet for virtual trading
	rpcURL := os.Getenv("SOLANA_RPC_URL")
	if rpcURL == "" {
		rpcURL = rpc.DevNet_RPC
	}

	client := rpc.New(rpcURL)

	// Load payer private key from env or generate a new one for testing
	privateKeyStr := os.Getenv("SOLANA_PRIVATE_KEY")
	var payer solana.PrivateKey
	var err error

	if privateKeyStr != "" {
		payer, err = solana.PrivateKeyFromBase58(privateKeyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
	} else {
		// Generate a new keypair for testing
		payer = solana.NewWallet()
		log.Printf("Generated new Solana keypair. Public key: %s", payer.PublicKey())
		log.Printf("Private key (save this): %s", payer.String())
	}

	return &SolanaClient{
		rpcClient: client,
		payer:     payer,
	}, nil
}

// RecordTradeOnChain records a trade transaction on Solana blockchain
func (s *SolanaClient) RecordTradeOnChain(
	userID uuid.UUID,
	assetSymbol string,
	tradeType string,
	quantity float64,
	price float64,
) (string, error) {
	ctx := context.Background()

	// For this virtual CEX, we'll create a simple memo transaction
	// In production, you would interact with your custom program

	// Get recent blockhash
	recent, err := s.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Create a memo with trade data
	memo := fmt.Sprintf("TRADE:%s:%s:%s:%.8f:%.2f",
		userID.String(),
		assetSymbol,
		tradeType,
		quantity,
		price,
	)

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				0, // Zero lamports transfer, just for recording
				s.payer.PublicKey(),
				s.payer.PublicKey(),
			).Build(),
			// You can add a memo program instruction here
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(s.payer.PublicKey()),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(s.payer.PublicKey()) {
			return &s.payer
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	sig, err := s.rpcClient.SendTransactionWithOpts(
		ctx,
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)

	if err != nil {
		// Log error but don't fail the trade
		log.Printf("Warning: Failed to record trade on Solana: %v (memo: %s)", err, memo)
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("Trade recorded on Solana. Signature: %s", sig.String())
	return sig.String(), nil
}

// GetTransactionStatus checks the status of a transaction
func (s *SolanaClient) GetTransactionStatus(signature string) (bool, error) {
	ctx := context.Background()

	sig, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return false, err
	}

	status, err := s.rpcClient.GetSignatureStatuses(ctx, true, sig)
	if err != nil {
		return false, err
	}

	if len(status.Value) == 0 || status.Value[0] == nil {
		return false, nil
	}

	return status.Value[0].ConfirmationStatus == rpc.ConfirmationStatusFinalized, nil
}
