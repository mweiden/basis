package blockchain

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
	"io"
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
		quitc:          make(chan struct{}),
		nodeIdentifier: nodeIdentifier,
		requestMethod: requestMethod,
		nodeRegistry:   make(map[string]int64),
	}
	go api.loop()
	return &api
}

func (api *API) loop() {
	for {
		select {
		case <-api.quitc:
			return
		case f := <-api.actionc:
			f()
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
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusCreated)
	case err := <-errc:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// mine
type mineResponse struct {
	Status       string        `json:"status"`
	Index        uint          `json:"index"`
	Transactions []Transaction `json:"transaction"`
	Proof        uint          `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

func (api *API) handleMine(w http.ResponseWriter, r *http.Request) {
	var (
		mineOk = make(chan mineResponse)
		errc   = make(chan error)
	)
	api.actionc <- func() {
		lastProof := api.blockchain.LastBlock().Proof
		proof := api.blockchain.hashCash.ProofOfWork(lastProof)
		api.blockchain.NewTransaction(
			"0",
			api.nodeIdentifier,
			1,
		)
		block := api.blockchain.NewBlock(proof, nil)
		mineOk <- mineResponse{
			Status:       "new block created",
			Index:        block.Index,
			Transactions: block.Transactions,
			Proof:        block.Proof,
			PreviousHash: block.PreviousHash,
		}
	}
	select {
	case response := <-mineOk:
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	case err := <-errc:
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		errc       = make(chan error)
	)
	api.actionc <- func() {
		var request nodeRequest
		json.NewDecoder(r.Body).Decode(&request)
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
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusCreated)
	case err := <-errc:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// node register
type resolveResponse struct {
	Status    string `json:"status"`
	Chain 	[]Block `json:"chain"`
}

func (api *API) handleNodeResolve(w http.ResponseWriter, r *http.Request) {
	var (
		resolveOk = make(chan resolveResponse)
		errc       = make(chan error)
	)
	api.actionc <- func() {
		replaced := api.ResoveConflicts()

		var status string
		if replaced {
			status = "chain replaced"
		} else {
			status = "chain authoritative"
		}
		resolveOk <- resolveResponse{
			Status: status,
			Chain: api.blockchain.Chain(),
		}
	}
	select {
	case response := <-resolveOk:
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	case err := <-errc:
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (api *API) ResoveConflicts() bool {
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
