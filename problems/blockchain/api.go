package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type API struct {
	blockchain     *Blockchain
	actionc        chan func()
	quitc          chan struct{}
	nodeIdentifier string
	nodeRegistry   map[string]int64
	requestMethod  func(string) io.ReadCloser
}

func NewAPI(
	blockchain *Blockchain,
	nodeIdentifier string,
	requestMethod func(string) io.ReadCloser,
) *API {
	api := API{
		blockchain:     blockchain,
		actionc:        make(chan func()),
		nodeIdentifier: nodeIdentifier,
		requestMethod:  requestMethod,
		nodeRegistry:   make(map[string]int64),
	}
	return &api
}

func (api *API) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting API down.")
			return ctx.Err()
		case f := <-api.actionc:
			f()
		}
	}
}

func (api *API) Cron(ctx context.Context, cronBeat chan bool) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cronBeat:
			mineResponse, _ := api.Mine()
			bytes, _ := json.Marshal(mineResponse)
			log.Println(string(bytes))
			resolveResponse, _ := api.NodeResolve()
			bytes, _ = json.Marshal(resolveResponse)
			log.Println(string(bytes))
		}
	}
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/transactions/new" {
		api.handleNewTransaction(w, r)
	} else if r.Method == "GET" && r.URL.Path == "/mine" {
		api.handleMine(w, r)
	} else if r.Method == "GET" && r.URL.Path == "/chain" {
		api.handleChain(w, r)
	} else if r.Method == "POST" && r.URL.Path == "/nodes/register" {
		api.handleNodeRegister(w, r)
	} else if r.Method == "GET" && r.URL.Path == "/nodes/resolve" {
		api.handleNodeResolve(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("%s %s\n", r.Method, r.URL.Path)
}

func UnixTime() int64 {
	return time.Now().Unix()
}

// new transaction
type transactionResponse struct {
	Status      string      `json:"status"`
	Transaction Transaction `json:"transaction"`
}

func (api *API) handleNewTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		createOk = make(chan transactionResponse)
		errc     = make(chan error)
	)
	api.actionc <- func() {
		var request Transaction
		err := json.NewDecoder(r.Body).Decode(&request)
		defer r.Body.Close()
		if err != nil {
			errc <- err
			return
		}
		transaction := api.blockchain.NewTransaction(
			request.Sender,
			request.Recipient,
			request.Amount,
		)
		createOk <- transactionResponse{
			Status:      "transaction created",
			Transaction: *transaction,
		}
	}
	select {
	case response := <-createOk:
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusCreated)
	case err := <-errc:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// mine
type MineResponse struct {
	Status       string        `json:"status"`
	Index        uint          `json:"index"`
	Transactions []Transaction `json:"transaction"`
	Proof        uint          `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

func (api *API) handleMine(w http.ResponseWriter, r *http.Request) {
	response, err := api.Mine()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *API) Mine() (MineResponse, error) {
	var (
		mineOk = make(chan MineResponse)
		errc   = make(chan error)
	)
	api.actionc <- func() {
		lastProof := api.blockchain.LastBlock().Proof
		proof := api.blockchain.hashCash.ProofOfWork(lastProof)
		api.blockchain.NewTransaction(
			api.nodeIdentifier,
			api.nodeIdentifier,
			1,
		)
		block := api.blockchain.NewBlock(proof, nil)
		mineOk <- MineResponse{
			Status:       "new block created",
			Index:        block.Index,
			Transactions: block.Transactions,
			Proof:        block.Proof,
			PreviousHash: block.PreviousHash,
		}
	}
	select {
	case response := <-mineOk:
		return response, nil
	case err := <-errc:
		return MineResponse{}, err
	}
}

// chain
type chainResponse struct {
	Chain  []Block `json:"chain"`
	Length int     `json:"length"`
}

func (api *API) handleChain(w http.ResponseWriter, r *http.Request) {
	var (
		chainOk = make(chan chainResponse)
	)
	api.actionc <- func() {
		chainOk <- chainResponse{
			Chain:  api.blockchain.Chain(),
			Length: len(api.blockchain.Chain()),
		}
	}
	select {
	case response := <-chainOk:
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// node register
type nodeRequest struct {
	Address string `json:"address"`
}

type nodeResponse struct {
	Address   string `json:"address"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

func (api *API) handleNodeRegister(w http.ResponseWriter, r *http.Request) {
	var (
		registryOk = make(chan nodeResponse)
		noOpOk     = make(chan nodeResponse)
		errc       = make(chan error)
	)
	api.actionc <- func() {
		var request nodeRequest
		json.NewDecoder(r.Body).Decode(&request)
		if request.Address == api.nodeIdentifier {
			noOpOk <- nodeResponse{
				Address:   request.Address,
				Status:    "node cannot register itself, no operation",
				Timestamp: UnixTime(),
			}
			return
		}
		_, ok := api.nodeRegistry[request.Address]
		if !ok {
			api.nodeRegistry[request.Address] = UnixTime()
		}
		registryOk <- nodeResponse{
			Address:   request.Address,
			Status:    "node registered",
			Timestamp: api.nodeRegistry[request.Address],
		}
	}
	select {
	case response := <-registryOk:
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusCreated)
	case response := <-noOpOk:
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
	case err := <-errc:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// node register
type ResolveResponse struct {
	Status string  `json:"status"`
	Chain  []Block `json:"chain"`
}

func (api *API) handleNodeResolve(w http.ResponseWriter, r *http.Request) {
	response, err := api.NodeResolve()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func (api *API) NodeResolve() (ResolveResponse, error) {
	var (
		resolveOk = make(chan ResolveResponse)
		errc      = make(chan error)
	)
	api.actionc <- func() {
		replaced := api.ResolveConflicts()

		var status string
		if replaced {
			status = "chain replaced"
		} else {
			status = "chain authoritative"
		}
		resolveOk <- ResolveResponse{
			Status: status,
			Chain:  api.blockchain.Chain(),
		}
	}
	select {
	case response := <-resolveOk:
		return response, nil
	case err := <-errc:
		return ResolveResponse{}, err
	}
}

func (api *API) ResolveConflicts() bool {
	var neighbours []string

	for k, _ := range api.nodeRegistry {
		neighbours = append(neighbours, k)
	}

	var newChain []Block

	// we're only looking for chains longer than ours
	maxLen := len(api.blockchain.chain)

	// get and verify the chains from all the nodes in our network
	for _, node := range neighbours {
		url := fmt.Sprintf("http://%s/chain", node)
		response := api.requestMethod(url)
		var cr chainResponse
		json.NewDecoder(response).Decode(&cr)
		response.Close()

		length := len(cr.Chain)

		if length > maxLen {
			newChain = cr.Chain
		}
	}

	if len(newChain) > 0 {
		api.blockchain.chain = newChain
		return true
	}

	return false
}
