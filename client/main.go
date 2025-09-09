package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Bid struct {
	Bid string
}

func main() {

	fmt.Println("Iniciando")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Requisição iniciada com sucesso")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Tempo da execução insuficiente")
		os.Exit(1)
		//panic(err)
	}
	defer res.Body.Close()
	//io.Copy(os.Stdout, res.Body)
	fmt.Println("Requisição realizada com sucesso")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", body)

	var dolar Bid
	err = json.Unmarshal(body, &dolar)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Criando arquivo com valor %s\n", dolar.Bid)
	criarArquivo(dolar.Bid)
}

func criarArquivo(valor string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("Valor: " + valor)
	_, err = f.WriteString("Dolar: {" + valor + "}")
	if err != nil {
		panic(err)
	}
	fmt.Println("Arquivo gerado com sucesso")

	f.Close()
}
