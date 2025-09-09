package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Cotacao struct {
	Cotacoes struct {
		Code        string `json:"code"`
		Codein      string `json:"codein"`
		Name        string `json:"name"`
		High        string `json:"high"`
		Low         string `json:"low"`
		VarBid      string `json:"varBid"`
		PctChange   string `json:"pctChange"`
		Bid         string `json:"bid"`
		Ask         string `json:"ask"`
		Timestamp   string `json:"timestamp"`
		Create_date string `json:"create_date"`
	} `json:"USDBRL"`
}

type Bid struct {
	Bid string
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/cotacao", getCotacao)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("health check")
}

func getCotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request getCotacao iniciada")
	defer log.Println("Request getCotacao finalizada")
	select {
	case <-time.After(200 * time.Millisecond):
		cotacao, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			panic(err)
		}
		defer cotacao.Body.Close()
		res, err := io.ReadAll(cotacao.Body)
		if err != nil {
			panic(err)
		}

		var data Cotacao
		err = json.Unmarshal(res, &data)
		if err != nil {
			panic(err)
		}
		// log.Printf("Parsed Data: %+v\n", data)
		log.Printf("Bid: %s\n", data.Cotacoes.Bid)

		inserirCotacao(&data.Cotacoes.Bid, r)

		// Imprime no comand line stdout
		log.Println("Request getCotacao processada com sucesso")

		// Imprime no browser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var bidJson = Bid{Bid: data.Cotacoes.Bid}

		json.NewEncoder(w).Encode(bidJson)

	case <-ctx.Done():
		// Imprime no comand line stdout
		log.Println("Request getCotacao cancelada pelo cliente")
	}
}

func inserirCotacao(valor *string, r *http.Request) {

	ctxDb := r.Context()
	select {
	case <-time.After(10 * time.Millisecond):

		db, err := sql.Open("mysql", "root:senha123@tcp(172.26.93.133:32002)/cotacao")
		if err != nil {
			panic(err)
		}

		stmt, err := db.Prepare("insert into cotacao (valor) values (?)")
		if err != nil {
			panic(err)
		}

		defer stmt.Close()
		defer db.Close()

		_, err = stmt.Exec(valor)
		if err != nil {
			panic(err)
		}

	case <-ctxDb.Done():
		log.Println("Requisição cancelada")
	}

}
