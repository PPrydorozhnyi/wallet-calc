package processor

import (
	"github.com/PPrydorozhnyi/wallet/db"
	"github.com/PPrydorozhnyi/wallet/model"
	"github.com/PPrydorozhnyi/wallet/proto"
	"github.com/PPrydorozhnyi/wallet/util"
	"github.com/google/uuid"
	bigDecimal "github.com/shopspring/decimal"
	"time"
)

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

	ledger, err := buildLedger(outcomes, request, accountId)

	if err != nil {
		return nil, err
	}

	return ledger, db.PersistCommandResult(acc, ledger)
}

func applyActions(request *model.TransactionRequest, ws *wallet.Wallet, err error) ([]*wallet.LedgerRecord_Outcome,
	error) {
	outcomes := make([]*wallet.LedgerRecord_Outcome, len(*request.Actions))

	for i, action := range *request.Actions {
		balance := ws.Wallets[action.Currency].Balances[action.BalanceId]

		originalBalance, e := bigDecimal.NewFromString(balance.Amount.GetValue())
		if e != nil {
			return nil, err
		}

		// todo in request it should be big.Float also
		// todo result has values like "200.4920000000000000000000000000000000005"
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

func buildLedger(outcomes []*wallet.LedgerRecord_Outcome, request *model.TransactionRequest,
	accountId string) (*model.Ledger,
	error) {
	ledgerRecord := &wallet.LedgerRecord{
		Outcomes: outcomes,
		Reason: &wallet.LedgerRecord_Reason{
			Id:        request.Reason.Id,
			Name:      request.Reason.Name,
			Reference: request.Reason.Reference,
			Meta:      request.Reason.Meta,
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
		CommandId:    request.CommandId,
		ClientId:     0,
		CommandType:  "TRANSACTION",
	}, nil
}

func buildOutcome(amountAfter bigDecimal.Decimal, action *model.Action) *wallet.LedgerRecord_Outcome {
	txTypeValue := wallet.LedgerRecord_TransactionType_value[action.TransactionType]
	txType := wallet.LedgerRecord_TransactionType(txTypeValue)

	id, _ := uuid.NewV7() //todo add exception handling

	return &wallet.LedgerRecord_Outcome{
		Id:              id.String(),
		BalanceId:       action.BalanceId,
		Currency:        action.Currency,
		TransactionType: &txType,
		Amount:          util.BigDecimalToDecimal(action.Amount),
		BalanceAfter:    util.BigDecimalToDecimal(amountAfter),
	}
}
