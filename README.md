# 🎬 Showio

Showio è un tracker personale full-stack per serie TV, anime e film. Nasce dall'esigenza di avere un'unica piattaforma pulita e centralizzata che combini i database occidentali (TMDB) e quelli dedicati all'animazione giapponese (AniList). 

Questo repository è un **monorepo** che contiene sia il codice frontend che quello backend.

## 🌟 Funzionalità Principali
* **Ricerca Unificata:** Cerca simultaneamente show su TMDB (tramite REST API) e AniList (tramite GraphQL).
* **Gestione Libreria Personale:** Aggiungi titoli alla tua lista, aggiorna lo stato (es. "Watching", "Completed") e tieni traccia dell'ultimo episodio visto.
* **Calendario Uscite:** Scopri quando andrà in onda il prossimo episodio delle serie che stai seguendo.
* **Sincronizzazione in Background:** Aggiornamenti automatici notturni per mantenere il database allineato con le nuove uscite.

## 🏗️ Architettura e Tech Stack
Il progetto è diviso in due componenti principali, progettati per essere leggeri, veloci e ospitabili su servizi cloud gratuiti o a basso costo.

* **Backend (`/backend`):** Sviluppato in **Go (Golang)** per massimizzare le performance e ridurre i tempi di avvio (cold start). Utilizza PostgreSQL come database.
* **Frontend (`/frontend`):** Sviluppato in **Angular**, offre una Single Page Application (SPA) reattiva e fortemente tipizzata.

## 📂 Struttura del Repository
```text
showio/
├── backend/       # API Server in Go, logica di business e connessione DB
├── frontend/      # Web App in Angular
└── .github/       # Workflow di GitHub Actions (CI/CD e Cron Jobs)