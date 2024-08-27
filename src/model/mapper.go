package model

import (
	wallet "github.com/PPrydorozhnyi/wallet/proto"
	"github.com/googleapis/go-type-adapters/adapters"
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

func toWalletDto(currency string, walletEntry *wallet.Wallet_WalletEntry) *WalletDto {
	balances := make([]*BalanceDto, 0, len(walletEntry.Balances))
	for id, b := range walletEntry.Balances {
		balances = append(balances, toBalanceDto(id, b))
	}
	return &WalletDto{
		Currency: currency,
		Balances: balances,
	}
}

func toBalanceDto(balanceId string, balance *wallet.Wallet_Balance) *BalanceDto {
	//todo handle possible error
	amount, _ := adapters.ProtoDecimalToFloat(balance.Amount)
	return &BalanceDto{
		Id:        balanceId,
		Type:      balance.Type,
		Vertical:  balance.Vertical,
		Amount:    amount,
		CreatedAt: balance.CreatedAt,
		UpdatedAt: balance.UpdatedAt,
	}
}
