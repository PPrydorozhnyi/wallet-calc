package processor

import (
	"github.com/PPrydorozhnyi/wallet/db"
	"github.com/PPrydorozhnyi/wallet/model"
	pb "github.com/PPrydorozhnyi/wallet/proto"
	"github.com/PPrydorozhnyi/wallet/util"
	"github.com/google/uuid"
	bigDecimal "github.com/shopspring/decimal"
	"time"
)

func CreateAccount(request *model.CreateWalletRequest) (*model.Ledger, error) {
	wrs := request.Wallets
	wallets := make(map[string]*pb.Wallet_WalletEntry, len(wrs))
	for _, wr := range wrs {
		wallets[wr.Currency] = createWallet(wr)
	}

	outcomes := buildWalletOutcomes(wallets, len(wrs))

	acc := &model.Account{
		Id: request.AccountId,
		WalletState: &pb.Wallet{
			Wallets: wallets,
		},
		Version: db.InitWalletVersion,
	}

	ledger, err := buildLedger(outcomes, request.Reason, request.AccountId, request.CommandId)

	if err != nil {
		return nil, err
	}

	return ledger, db.PersistCommandResult(acc, ledger)
}

func buildWalletOutcomes(wallets map[string]*pb.Wallet_WalletEntry, capacity int) []*pb.LedgerRecord_Outcome {
	outcomes := make([]*pb.LedgerRecord_Outcome, 0, capacity)

	for currency, wr := range wallets {
		for balanceId, balance := range wr.Balances {
			id, _ := uuid.NewV7() //todo add exception handling
			outcomes = append(outcomes, &pb.LedgerRecord_Outcome{
				Id:           id.String(),
				BalanceId:    balanceId,
				Currency:     currency,
				Vertical:     &balance.Vertical,
				BalanceType:  &balance.Type,
				Amount:       balance.Amount,
				BalanceAfter: balance.Amount,
			})
		}
	}

	return outcomes
}

func createWallet(wallet *model.WalletDefinition) *pb.Wallet_WalletEntry {
	return &pb.Wallet_WalletEntry{
		Balances:  createBalances(wallet),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
}

func createBalances(wallet *model.WalletDefinition) map[string]*pb.Wallet_Balance {
	balances := make(map[string]*pb.Wallet_Balance, len(wallet.Balances))

	for _, bd := range wallet.Balances {
		balanceId, _ := uuid.NewV7() // todo add error handling
		balances[balanceId.String()] = &pb.Wallet_Balance{
			Type:      bd.BalanceType,
			Vertical:  bd.Vertical,
			Amount:    util.BigDecimalToDecimal(bd.Amount),
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		}
	}

	return balances
}

func ApplyTransaction(accountId string, request *model.TransactionRequest) (*model.Ledger, error) {
	acc, err := db.GetWallet(accountId)

	if err != nil {
		return nil, err
	}

	ws := acc.WalletState

	outcomes, err := applyActions(request, ws, err)
	if err != nil {
		return nil, err
	}

	ledger, err := buildLedger(outcomes, request.Reason, accountId, request.CommandId)

	if err != nil {
		return nil, err
	}

	return ledger, db.PersistCommandResult(acc, ledger)
}

func applyActions(request *model.TransactionRequest, ws *pb.Wallet, err error) ([]*pb.LedgerRecord_Outcome,
	error) {
	outcomes := make([]*pb.LedgerRecord_Outcome, len(*request.Actions))

	for i, action := range *request.Actions {
		balance := ws.Wallets[action.Currency].Balances[action.BalanceId]

		originalBalance, e := bigDecimal.NewFromString(balance.Amount.GetValue())
		if e != nil {
			return nil, err
		}

		var resultBalance bigDecimal.Decimal

		if action.TransactionType == model.CREDIT {
			resultBalance = originalBalance.Add(action.Amount)
		} else {
			resultBalance = originalBalance.Sub(action.Amount)
		}

		outcomes[i] = buildOutcome(resultBalance, &action)
		balance.Amount = outcomes[i].BalanceAfter
	}
	return outcomes, nil
}

func buildLedger(outcomes []*pb.LedgerRecord_Outcome, reason *model.Reason,
	accountId string, commandId string) (*model.Ledger,
	error) {
	ledgerRecord := &pb.LedgerRecord{
		Outcomes: outcomes,
		Reason: &pb.LedgerRecord_Reason{
			Id:        reason.Id,
			Name:      reason.Name,
			Reference: reason.Reference,
			Meta:      reason.Meta,
		},
	}
	cid, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}
	return &model.Ledger{
		Id:           cid,
		AccountId:    accountId,
		LedgerRecord: ledgerRecord,
		CreatedAt:    time.Now(),
		CommandId:    commandId,
		ClientId:     0,
		CommandType:  "TRANSACTION",
	}, nil
}

func buildOutcome(amountAfter bigDecimal.Decimal, action *model.Action) *pb.LedgerRecord_Outcome {
	txTypeValue := pb.LedgerRecord_TransactionType_value[action.TransactionType]
	txType := pb.LedgerRecord_TransactionType(txTypeValue)

	id, _ := uuid.NewV7() //todo add exception handling

	return &pb.LedgerRecord_Outcome{
		Id:              id.String(),
		BalanceId:       action.BalanceId,
		Currency:        action.Currency,
		TransactionType: &txType,
		Amount:          util.BigDecimalToDecimal(action.Amount),
		BalanceAfter:    util.BigDecimalToDecimal(amountAfter),
	}
}
