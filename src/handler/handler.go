package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PPrydorozhnyi/wallet/db"
	"github.com/PPrydorozhnyi/wallet/model"
	"github.com/PPrydorozhnyi/wallet/processor"
	"io"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostsHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetPosts(w, r)
	case "POST":
		handleCreatePost(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func AccountsHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	request := &model.CreateWalletRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, "Bad request. Cannot parse body", http.StatusBadRequest)
		return
	}

	account, err := processor.CreateAccount(request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	serializeTransactionResponse(w, account, true)
}

func AccountHandle(w http.ResponseWriter, r *http.Request) {
	idString := strings.Split(r.URL.Path, "/")[2]

	switch r.Method {
	case "GET":
		account := handleGetAccount(w, r, idString)
		serializeAccount(w, account)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func WalletsHandle(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("accountId")

	switch r.Method {
	case "GET":
		account := handleGetAccount(w, r, idString)
		serializeAccountResponse(w, account)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func TransactionsHandle(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("accountId")

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	request := &model.TransactionRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, "Bad request. Cannot parse body", http.StatusBadRequest)
		return
	}

	ledger, err := processor.ApplyTransaction(idString, request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	serializeTransactionResponse(w, ledger, false)
}

func handleGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(model.Posts())
	if err != nil {
		http.Error(w, "Cannot serialize posts", http.StatusInternalServerError)
	}
}

func handleGetAccount(w http.ResponseWriter, r *http.Request, id string) *model.Account {
	account, err := db.GetWallet(id)

	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return nil
	}

	return account
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	var p model.Post
	err = json.Unmarshal(body, &p)

	if err != nil {
		http.Error(w, "Cannot parse body", http.StatusBadRequest)
		return
	}

	post, cErr := model.CreatePost(p)

	if cErr != nil {
		http.Error(w, cErr.Error(), http.StatusInternalServerError)
		return
	}

	serializePost(w, post)
}

func serializePost(w http.ResponseWriter, post *model.Post) {
	err := json.NewEncoder(w).Encode(post)
	if err != nil {
		http.Error(w, "Cannot serialize post", http.StatusInternalServerError)
	}
}

func serializeAccount(w http.ResponseWriter, account *model.Account) {
	if account == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, "Cannot serialize account", http.StatusInternalServerError)
	}
}

func serializeAccountResponse(w http.ResponseWriter, account *model.Account) {
	if account == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(model.ToAccountResponse(account)); err != nil {
		http.Error(w, "Cannot serialize account", http.StatusInternalServerError)
	}
}

func serializeTransactionResponse(w http.ResponseWriter, ledger *model.Ledger, extended bool) {
	if err := json.NewEncoder(w).Encode(model.ToTransactionResponse(ledger, extended)); err != nil {
		http.Error(w, "Cannot serialize ledger response", http.StatusInternalServerError)
	}
}
