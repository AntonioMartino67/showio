# ⚙️ Showio - Backend API

Questo modulo costituisce il motore backend di **Showio**, un'applicazione full-stack per il tracciamento centralizzato di serie TV, film e anime. 
Il server è interamente sviluppato in **Go (Golang)** ed è progettato per garantire prestazioni elevate, tipizzazione forte e zero cold-start, compilando in un singolo binario leggero.

Il backend si occupa dell'autenticazione degli utenti tramite JWT custom, dell'unificazione delle richieste verso servizi esterni (TMDB e AniList), del caching relazionale dei dati su PostgreSQL e dell'esecuzione asincrona dei cicli di sincronizzazione tramite GitHub Actions.

---

## 🚀 Funzionalità Principali

* **Autenticazione Sicura**: Registrazione e login con hashing delle password (`bcrypt`) e gestione delle sessioni tramite JWT.
* **Ricerca Unificata**: Interrogazione simultanea di **TMDB** (REST) e **AniList** (GraphQL) con salvataggio automatico nella cache locale (`media_items`) per ridurre le chiamate di rete.
* **Tracking dei Progressi**: Gestione della watchlist personale, stati di visione (`watching`, `completed`, `dropped`, `plan_to_watch`) e aggiornamento degli episodi visti.
* **Motore di Auto-Update (Delta Check)**: Quando una serie TV viene aggiunta alla watchlist, una goroutine scarica in background stagioni ed episodi da TMDB.
* **Calendario Personalizzato**: Calcolo automatico del "prossimo episodio non visto" per le serie TV in corso.
* **Sincronizzazione Schedulata**: Endpoint protetto per il sync globale, richiamato nottetempo da GitHub Actions per mantenere aggiornate le date di uscita (`air_date`) delle serie in corso.

---

## 🛠️ Stack Tecnico e Librerie Core

| Categoria | Tecnologia | Descrizione |
| :--- | :--- | :--- |
| **Linguaggio** | Go (Golang) | Backend compilato, concurrency nativa (goroutine). |
| **Routing HTTP** | `go-chi/chi/v5` | Router HTTP leggero, idiomatico e compatibile con `net/http`. |
| **Database** | PostgreSQL (Neon.tech) | DB relazionale serverless con *scale-to-zero*. |
| **Driver DB** | `jackc/pgx/v5` | Driver nativo per Postgres con pool di connessioni. |
| **Autenticazione** | `golang-jwt/jwt/v5` | Generazione e validazione di JSON Web Tokens. |
| **Crittografia** | `golang.org/x/crypto` | Hashing sicuro delle password tramite `bcrypt`. |
| **Config** | `joho/godotenv` | Gestione delle variabili d'ambiente in locale. |
| **API Esterne** | TMDB API / AniList GraphQL | Fetching di metadati, poster, stagioni ed episodi. |
| **Hosting** | Render (Web Service) | Container always-on (piano Free). |

---

## 📂 Architettura e Struttura delle Directory

Il progetto adotta la convenzione standard di Go per backend e API:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go             # Punto di ingresso: avvia il server, carica env, registra route.
├── internal/                   # Codice privato, non importabile da altri moduli.
│   ├── auth/                   # Logica JWT, middleware di protezione, hashing password.
│   ├── database/               # Connessione a Postgres (Neon Pool) e query helper.
│   ├── external/               # Client HTTP per TMDB (REST) e AniList (GraphQL).
│   ├── handlers/               # Handler HTTP (controllori) per ogni endpoint.
│   └── models/                 # Struct Go che mappano le tabelle del database.
├── .env                        # Variabili d'ambiente locali (NON committare).
├── .gitignore                  
├── go.mod                      # Definizione del modulo e dipendenze dirette.
└── go.sum                      # Checksum crittografico delle dipendenze.