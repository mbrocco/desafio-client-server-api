package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CotacaoDolar struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/", BuscaDolarHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaDolarHandler(w http.ResponseWriter, r *http.Request) {
	// cria o contexto para realizar a requisição da API para buscar a cotação do dolar
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	// if r.URL.Path != "/" {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	cotacao, error := buscaCotacaoDolar(ctx)
	if error != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// result, err := json.Marshal(cep)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(result)
	json.NewEncoder(w).Encode(cotacao)

	/*
		file, err := os.Create("cidade.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("CEP: %s, Localidade: %s, UF: %s", data.Cep, data.Localidade, data.Uf))

		ctx := r.Context()

		log.Println("Request iniciada")

		defer log.Println("Request finalizada")

		select {
		case <-time.After(200 * time.Millisecond):
			// Imprime no comand line stdout
			log.Println("Request processada com sucesso")
			// Imprime no browser
			w.Write([]byte("Request processada com sucesso"))
		case <-ctx.Done():
			// Imprime no comand line stdout
			log.Println("Request cancelada pelo cliente")
		}*/
}

func buscaCotacaoDolar(ctx context.Context) (*CotacaoDolar, error) {

	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Println("Erro ao fazer requisição")
		return nil, err
	}

	time.Sleep(time.Second)
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta")
		return nil, err
	}
	var data CotacaoDolar
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Println("Erro ao fazer parse da resposta")
		return nil, err
	}

	select {
	case <-ctx.Done():
		fmt.Println("Requisição interrompida. Tempo limite atingido.")
		return nil, fmt.Errorf("ss")
	default:
		fmt.Println("Requisição OK")
		return &data, nil
	}

}
