package model

import (
	pb "github.com/PPrydorozhnyi/wallet/proto"
	"github.com/PPrydorozhnyi/wallet/util"
)

func ToAccountResponse(account *Account) *AccountResponse {
	if account == nil {
		return nil
	}

	wallets := make(map[string]*WalletDto)
	for c, w := range account.WalletState.Wallets {
		wallets[c] = toWalletDto(c, w)
	}

	return &AccountResponse{
		AccountId: account.Id,
		Wallets:   wallets,
	}
}

func toWalletDto(currency string, walletEntry *pb.Wallet_WalletEntry) *WalletDto {
	balances := make([]*BalanceDto, 0, len(walletEntry.Balances))
	for id, b := range walletEntry.Balances {
		balances = append(balances, toBalanceDto(id, b))
	}
	return &WalletDto{
		Currency: currency,
		Balances: balances,
	}
}

func toBalanceDto(balanceId string, balance *pb.Wallet_Balance) *BalanceDto {
	//todo handle possible error
	amount, _ := util.DecimalToBigDecimal(balance.Amount)
	return &BalanceDto{
		Id:        balanceId,
		Type:      balance.Type,
		Vertical:  balance.Vertical,
		Amount:    amount,
		CreatedAt: balance.CreatedAt,
		UpdatedAt: balance.UpdatedAt,
	}
}

func ToTransactionResponse(ledger *Ledger, extended bool) *CommandResponse {
	return &CommandResponse{
		Id:          ledger.Id,
		ProcessedAt: ledger.CreatedAt.UnixMilli(),
		Actions:     toOutcomeDtos(ledger.LedgerRecord.Outcomes, extended),
	}
}

func toOutcomeDtos(outcomes []*pb.LedgerRecord_Outcome, extended bool) []*OutcomeDto {
	outcomeDtos := make([]*OutcomeDto, len(outcomes))

	for i, outcome := range outcomes {

		bAfter, _ := util.DecimalToBigDecimal(outcome.BalanceAfter) // todo add exception handling

		outcomeDtos[i] = &OutcomeDto{
			BalanceId:    outcome.BalanceId,
			BalanceAfter: bAfter,
			Currency:     outcome.Currency,
			ActionId:     outcome.Id,
		}

		if extended {
			outcomeDtos[i].BalanceType = outcome.BalanceType
			outcomeDtos[i].Vertical = outcome.Vertical
		}
	}

	return outcomeDtos
}
