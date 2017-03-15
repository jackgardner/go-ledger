package main

import (
	"github.com/namsral/flag"
	"net/http"
	"github.com/playlyfe/go-graphql"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "github.com/jackgardner/go-ledger/proto"
	"encoding/json"
	"fmt"
	"log"
)

var (
	port int
	ledgerServiceAddr string
)

type Ledger struct {
	LedgerId string `json:"ledgerId"`
}
func main() {
	schema := `
	type Transaction {
		amountInPence: Int
		sourceLedgerId: String
		destinationLedgerId: String
	}
	type Ledger {
		ledgerId: String
		transactions: [ Transaction ]
	}
	type QueryRoot {
		ledgers(ledgerId: String): [ Ledger ]
	}
	`

	flag.IntVar(&port, "port", 8080, "Service port")
	flag.StringVar(&ledgerServiceAddr, "LEDGER_SERVICE_ADDR", "", "Address of ledger service")
	flag.Parse()

	opts := []grpc.DialOption{ grpc.WithInsecure() }
	conn, err := grpc.Dial(ledgerServiceAddr, opts...)

	if err != nil {
		grpclog.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewLedgerClient(conn)

	resolvers := map[string]interface{}{}
	resolvers["QueryRoot/ledgers"] = func(params *graphql.ResolveParams) (interface{}, error) {
		var ledgerId string
		ledgerId = params.Args["ledgerId"].(string)

		return []map[string]interface{}{
			{
				"ledgerId": ledgerId,
			},
		}, nil
	}
	resolvers["Ledger/transactions"] = func(params *graphql.ResolveParams) (interface{}, error) {


		if obj, ok := params.Source.(map[string]interface{}); ok {

			transactionRequest := pb.ListTransactionsRequest{SourceLedgerId: obj["ledgerId"].(string)}
			stream, err := client.GetTransactions(context.Background(), &transactionRequest)

			if err != nil {
				grpclog.Fatalf("Failed to uery GetTransactions, err: %v", err)
			}
			// TODO Can use streaming
			return stream.Transactions, nil
		}

		return []*pb.Transaction{}, nil
	}

	ctx := map[string]interface{}{}
	variables := map[string]interface{}{}
	executor, err := graphql.NewExecutor(schema, "QueryRoot", "", resolvers)
	executor.ResolveType = func(value interface{}) string {
		switch value.(type) {
		case *Ledger:
			return "Ledger"
		}

		return ""
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()["query"][0]
		result, _ := executor.Execute(ctx, query, variables, "")
		json.NewEncoder(w).Encode(result)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
