# 📦 Sistema de Emissão de Notas Fiscais - Desafio Técnico KORP

Sistema completo de emissão de notas fiscais e controle de estoque construído em arquitetura de microsserviços, focado em simplicidade, resiliência, alta disponibilidade e código limpo.

---

## 🛠️ Tecnologias Utilizadas

* **Frontend:** Angular 17+ (Componentes Standalone, RxJS, HttpClientModule, FormsModule, Angular Material / Toastr)
* **Backend:** Golang (Gorilla Mux, Go Modules)
* **Banco de Dados:** PostgreSQL (`database/sql` & `lib/pq`) com sistema de **Fallback em Memória (`sync.Mutex`)**
* **Comunicação:** REST HTTP com suporte nativo a CORS e tratamento de requisições *Preflight* (`OPTIONS`)

---

## 🏗️ Arquitetura e Microsserviços

O sistema é dividido em 3 camadas independentes:

┌───────────────────────────────┐
              │       Frontend Angular        │
              │    http://localhost:4200      │
              └──────────────┬────────────────┘
                             │
             ┌───────────────┴───────────────┐
             │                               │
             ▼                               ▼
┌─────────────────────────────────┐   ┌─────────────────────────────────┐
│       Serviço de Estoque        │   │     Serviço de Faturamento      │
│     http://localhost:8081       │   │     http://localhost:8082       │
└────────────────┬────────────────┘   └────────────────┬────────────────┘
│                                     │
└──────────────────┬──────────────────┘
▼
┌────────────────────┐
│ Banco de Dados     │
│ PostgreSQL / Mem   │
└────────────────────┘

*O serviço iniciará na porta `8081`.*

---

### 2. Iniciar o Serviço de Faturamento (Backend)

Em outro terminal:
```powershell
cd backend/faturamento-service
go run main.go
```
*O serviço iniciará na porta `8082`.*

---

### 3. Iniciar o Frontend (Angular)

Em um terceiro terminal:
```powershell
cd frontend
npm install
ng serve
```
Acesse a aplicação no navegador através do endereço: **`http://localhost:4200`**

---

## 👤 Autor

Desenvolvido por **Gabriel Naoki Uto Turigoe** como parte do desafio técnico da KORP.

* 🐙 **GitHub:** [@GabrielNaokiUT](https://github.com/GabrielNaokiUT)
