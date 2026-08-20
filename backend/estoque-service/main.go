package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Estrutura base do Produto[cite: 1]
type Produto struct {
	Codigo    int    `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

var db *sql.DB

func main() {
	var err error
	// Conexão com o banco (ajuste o password se o seu postgres local usar outro)
	connStr := "user=postgres password=root dbname=estoque sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/produtos", CriarProduto).Methods("POST")
	r.HandleFunc("/produtos/{id}/baixar", BaixarEstoque).Methods("POST")

	log.Println("Serviço de Estoque rodando na porta 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}

// Cria um novo produto[cite: 1]
func CriarProduto(w http.ResponseWriter, r *http.Request) {
	var p Produto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Erro nos dados de entrada", http.StatusBadRequest) // Erro 400 para bad request
		return
	}

	err := db.QueryRow("INSERT INTO produtos (descricao, saldo) VALUES ($1, $2) RETURNING codigo", p.Descricao, p.Saldo).Scan(&p.Codigo)
	if err != nil {
		http.Error(w, "Erro ao salvar no banco", http.StatusInternalServerError) // Erro 500[cite: 2]
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

type BaixaRequest struct {
	Quantidade int `json:"quantidade"`
}

// Atualiza o saldo (chamado na impressão da nota)[cite: 1]
func BaixarEstoque(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codigo := vars["id"]

	var req BaixaRequest
	json.NewDecoder(r.Body).Decode(&req)

	_, err := db.Exec("UPDATE produtos SET saldo = saldo - $1 WHERE codigo = $2 AND saldo >= $1", req.Quantidade, codigo)
	if err != nil {
		http.Error(w, "Erro ao atualizar estoque", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
