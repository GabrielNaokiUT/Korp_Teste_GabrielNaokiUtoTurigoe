package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Produto struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

type BaixaEstoqueRequest struct {
	ProdutoCodigo string `json:"produtoCodigo"`
	Quantidade    int    `json:"quantidade"`
}

var (
	db       *sql.DB
	produtos = make(map[string]*Produto)
	mu       sync.Mutex
)

func initDB() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=korp sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("PostgreSQL indisponível. Usando modo em memória.")
		return
	}
	if err = db.Ping(); err != nil {
		log.Println("PostgreSQL indisponível. Usando modo em memória.")
		db = nil
		return
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS produtos (
		codigo VARCHAR(50) PRIMARY KEY,
		descricao VARCHAR(255) NOT NULL,
		saldo INT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Println("Erro ao verificar tabela produtos:", err)
	} else {
		log.Println("Conectado ao PostgreSQL com sucesso no Estoque!")
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cadastrarProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p Produto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}

	if p.Codigo == "" {
		http.Error(w, `{"error":"Código é obrigatório"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		_, err := db.Exec(`
			INSERT INTO produtos (codigo, descricao, saldo) 
			VALUES ($1, $2, $3)
			ON CONFLICT (codigo) DO UPDATE 
			SET descricao = EXCLUDED.descricao, saldo = EXCLUDED.saldo`,
			p.Codigo, p.Descricao, p.Saldo)
		if err != nil {
			log.Println("Erro DB:", err)
			http.Error(w, `{"error":"Erro ao salvar no banco de dados"}`, http.StatusInternalServerError)
			return
		}
	}

	produtos[p.Codigo] = &p

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func baixarEstoque(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	codigo := vars["id"]

	var req BaixaEstoqueRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Quantidade <= 0 {
		req.Quantidade = 1
	}

	if codigo == "" {
		codigo = req.ProdutoCodigo
	}

	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		var saldoAtual int
		err := db.QueryRow("SELECT saldo FROM produtos WHERE codigo = $1", codigo).Scan(&saldoAtual)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error":"Produto não encontrado"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error":"Erro no banco de dados"}`, http.StatusInternalServerError)
			return
		}

		if saldoAtual < req.Quantidade {
			http.Error(w, `{"error":"Saldo insuficiente em estoque"}`, http.StatusBadRequest)
			return
		}

		novoSaldo := saldoAtual - req.Quantidade
		_, err = db.Exec("UPDATE produtos SET saldo = $1 WHERE codigo = $2", novoSaldo, codigo)
		if err != nil {
			http.Error(w, `{"error":"Erro ao atualizar saldo"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"codigo":     codigo,
			"novo_saldo": novoSaldo,
			"mensagem":   "Baixa realizada com sucesso",
		})
		return
	}

	// Armazenamento em memória caso o banco não esteja rodando
	p, exists := produtos[codigo]
	if !exists {
		p = &Produto{Codigo: codigo, Descricao: "Produto " + codigo, Saldo: 10}
		produtos[codigo] = p
	}

	if p.Saldo < req.Quantidade {
		http.Error(w, `{"error":"Saldo insuficiente em estoque"}`, http.StatusBadRequest)
		return
	}

	p.Saldo -= req.Quantidade
	json.NewEncoder(w).Encode(map[string]interface{}{
		"codigo":     codigo,
		"novo_saldo": p.Saldo,
		"mensagem":   "Baixa realizada com sucesso",
	})
}

func listarProdutos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		rows, err := db.Query("SELECT codigo, descricao, saldo FROM produtos")
		if err == nil {
			defer rows.Close()
			var lista []Produto
			for rows.Next() {
				var p Produto
				rows.Scan(&p.Codigo, &p.Descricao, &p.Saldo)
				lista = append(lista, p)
			}
			json.NewEncoder(w).Encode(lista)
			return
		}
	}

	var lista []Produto
	for _, p := range produtos {
		lista = append(lista, *p)
	}
	json.NewEncoder(w).Encode(lista)
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.Use(enableCORS)

	r.HandleFunc("/produtos", cadastrarProduto).Methods("POST", "OPTIONS")
	r.HandleFunc("/produtos", listarProdutos).Methods("GET", "OPTIONS")
	r.HandleFunc("/produtos/{id}/baixar", baixarEstoque).Methods("POST", "OPTIONS")

	fmt.Println("Serviço de Estoque rodando na porta :8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
