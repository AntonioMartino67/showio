# 🎬 Showio
### Personal Media Tracker | TV Shows, Anime & Movies

**Versione:** 1.0.0
**Stato:** In Sviluppo
**Autore:** Antonio
**Licenza:** Da definire (consigliata MIT per uso personale, rivedere se si passa a uso commerciale)
**Costo Operativo Attuale:** $0 (uso personale, esclusivamente Free Tier)
**Repository:** https://github.com/AntonioMartino67/showio

---

## 📑 1. Panoramica del Progetto

**Showio** è un'applicazione web personale per il tracciamento di contenuti multimediali (serie TV, anime, film). Nasce per superare i limiti delle app commerciali esistenti (es. TV Time), offrendo un'interfaccia pulita, senza pubblicità invadente e con un sistema di aggiornamento automatico dei palinsesti globali.

L'obiettivo è permettere all'utente di:
- Tracciare i propri progressi di visione episodio per episodio
- Ricevere aggiornamenti precisi su quando uscirà il prossimo episodio delle serie seguite
- Avere tutto in un'unica dashboard, attingendo da database mondiali (TMDB, AniList)

Il progetto è pensato per partire come **strumento a uso personale**, con un'architettura strutturata fin dall'inizio in modo da poter evolvere in futuro verso un prodotto aperto ad altri utenti (anche con eventuale monetizzazione), senza dover essere riscritto da zero.

---

## 🎯 2. Funzionalità (Feature List)

### ✅ MVP (Minimum Viable Product)
- **Autenticazione:** Registrazione e Login sicuro (JWT gestito manualmente lato backend).
- **Ricerca Unificata:** Barra di ricerca che interroga simultaneamente database di Serie TV, Film e Anime.
- **Tracking Granulare:** Segnare episodi come "Visti", calcolo automatico della % di completamento per stagione e serie.
- **Dashboard Personale:** Visualizzazione rapida delle serie "In Corso", "Completate" e "Piano di visione" (Watchlist).
- **Calendario Uscite:** Sezione dedicata che mostra i prossimi episodi in uscita per le serie seguite dall'utente.
- **Cache Locale:** Salvataggio dei metadati (trama, poster, cast) nel DB locale per caricamenti istantanei e risparmio di chiamate API.

### 🔵 Funzionalità Avanzate (Roadmap Futura)
- **Notifiche Push:** Alert sul browser/dispositivo quando esce un nuovo episodio.
- **Import/Export:** Migrazione dati da Trakt.tv, TV Time o MyAnimeList.
- **Statistiche Personali:** Grafici su tempo di visione, generi preferiti, ecc.
- **PWA (Progressive Web App):** Installazione su smartphone e funzionalità offline per consultare la propria lista.
- **Condivisione Social:** Generazione di card grafiche (spoiler-free) dei propri progressi.
- **Piano "Supporter":** Eventuali funzionalità extra per utenti che vogliono sostenere il progetto (vedi sezione 8).

---

## 🛠️ 3. Stack Tecnologico

### Perché questo stack
Lo stack è stato scelto bilanciando due esigenze: restare a **costo zero** in fase di sviluppo/uso personale, e limitare il numero di tecnologie completamente nuove da imparare in un periodo (Esame di Stato) in cui il tempo è prezioso. Per questo si è scelto di tenere il **frontend su una tecnologia già nota** (Angular) e concentrare l'apprendimento nuovo sul **backend** (Go).

### Frontend (Client)
- **Framework:** Angular (già usato in altri progetti personali, es. ASTRO-POINTER e Vallauri Store).
- **Linguaggio:** TypeScript.
- **Styling:** Tailwind CSS (o Angular Material, da valutare in fase di setup).
- **State Management:** Signals nativi di Angular (o RxJS/Services, in base alla versione).
- **Routing:** Angular Router.
- **HTTP Client:** `HttpClient` di Angular con `provideHttpClient(withFetch())` (necessario per gestione corretta di cookie/JWT cross-origin, lezione imparata da progetti precedenti come Vallauri Store).
- **Hosting:** Vercel (Piano Hobby, gratuito) o Netlify.

### Backend (Server)
- **Linguaggio:** Go — scelto come principale investimento di apprendimento del progetto: binari singoli, footprint minimo, concorrenza nativa, cold start quasi nullo.
- **Router HTTP:** `chi` (leggero e idiomatico, buona scelta per chi arriva da altri framework più "battery-included").
- **Driver DB:** `pgx` (driver Postgres nativo, più performante del generico `database/sql` con driver esterno).
- **Autenticazione:** JWT gestito manualmente (`golang-jwt/jwt`) + hashing password con `bcrypt` (`golang.org/x/crypto`).
- **HTTP Client esterno:** libreria standard `net/http` per le chiamate a TMDB/AniList.
- **Config:** `godotenv` per variabili d'ambiente in locale.

### Automazione (Cron / Background Jobs)
- **Scheduler:** GitHub Actions (workflow schedulato gratuito) al posto di sistemi come Celery/Redis — richiama periodicamente un endpoint del backend per il controllo aggiornamenti, senza bisogno di un worker sempre attivo.

### Database & Infrastruttura
- **Database:** PostgreSQL, ospitato su **Neon.tech** (Piano Free) — scelto per il vero "scale-to-zero" senza costi quando il progetto è inattivo, branching del DB per testare migrazioni in sicurezza, e compatibilità diretta con `pgx`/`database/sql` senza librerie proprietarie.
- **Hosting Frontend:** Vercel (Piano Hobby).
- **Hosting Backend:** container always-on su piano free (es. Fly.io o Render Free) — preferito a soluzioni serverless "vere" per Go, dato che il cold start di un binario Go è comunque minimo e un container always-on evita complessità di deploy serverless non necessarie per un progetto mono-utente.

---

## 🏗️ 4. Architettura e Flusso Dati

L'architettura è **Client-Server disaccoppiata**. Il Frontend non parla mai direttamente con le API esterne (TMDB/AniList), ma solo con il Backend di Showio: questo protegge le chiavi API e centralizza la logica.

### Struttura del monorepo
```
showio/
├── backend/          # API Go
│   ├── cmd/
│   │   └── api/      # entry point (main.go)
│   ├── internal/
│   │   ├── handlers/ # gestione richieste HTTP
│   │   ├── models/   # struct corrispondenti alle tabelle DB
│   │   ├── database/ # connessione e query Postgres
│   │   ├── auth/     # generazione/verifica JWT, hashing password
│   │   └── external/ # client TMDB e AniList
│   └── go.mod
├── frontend/         # App Angular
│   └── src/
├── .github/
│   └── workflows/    # GitHub Actions per il cron notturno
└── .gitignore
```

### Scelte architetturali chiave
- **Backend always-on leggero:** un binario Go compilato ha un footprint di memoria minimo, quindi anche un piano free "always-on" (non serverless) resta sostenibile senza costi.
- **Niente Redis/Celery:** i job schedulati (es. controllo notturno nuovi episodi) sono gestiti da GitHub Actions, gratuito e senza infrastruttura da mantenere.
- **Database con keep-alive nativo:** Neon gestisce autonomamente lo scale-to-zero, quindi non serve un ping periodico artificiale come su altri provider.

### Il Motore di Aggiornamento Automatico (Auto-Update Engine)
1. **On-Demand (Cache):** quando l'utente cerca un contenuto, il backend lo cerca su TMDB/AniList, lo salva nel DB locale e lo restituisce.
2. **Scheduled Polling (GitHub Actions):** un workflow gira periodicamente (es. ogni notte) e richiama un endpoint dedicato del backend.
3. **Delta Check:** il backend interroga le API esterne chiedendo gli aggiornamenti per i contenuti presenti nel DB locale con status `AIRING`.
4. **Sync & Notify:** se vengono rilevati nuovi episodi o cambi di data, il DB locale viene aggiornato e l'episodio viene marcato per un'eventuale notifica futura.

---

## 🗄️ 5. Schema del Database (Concettuale)

Lo schema include già alcuni campi "future-proof" (soft delete, timestamp) pensati per facilitare un'eventuale evoluzione verso un uso multi-utente conforme al GDPR.

```sql
-- Tabella Utenti
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL -- soft delete, utile in ottica GDPR
);

-- Tabella Metadati Media (Cache locale da API esterne)
CREATE TABLE media_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(50) NOT NULL, -- ID di TMDB o AniList
    source VARCHAR(20) NOT NULL, -- 'tmdb', 'anilist'
    title VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'tv', 'movie', 'anime'
    poster_url TEXT,
    overview TEXT,
    status VARCHAR(20), -- 'airing', 'ended', 'upcoming'
    last_synced_at TIMESTAMP
);

-- Tabella Stagioni
CREATE TABLE seasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_item_id UUID REFERENCES media_items(id) ON DELETE CASCADE,
    season_number INT NOT NULL,
    episode_count INT,
    air_date DATE
);

-- Tabella Episodi
CREATE TABLE episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    season_id UUID REFERENCES seasons(id) ON DELETE CASCADE,
    episode_number INT NOT NULL,
    title VARCHAR(255),
    air_date DATE,
    is_watched BOOLEAN DEFAULT FALSE
);

-- Tabella Progressi Utente
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    media_item_id UUID REFERENCES media_items(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'watching', -- 'watching', 'completed', 'dropped', 'plan_to_watch'
    current_season INT DEFAULT 1,
    current_episode INT DEFAULT 0,
    rating INT CHECK (rating >= 0 AND rating <= 10),
    last_watched_at TIMESTAMP,
    UNIQUE(user_id, media_item_id)
);
```

---

## 💰 6. Costi: uso personale vs. uso commerciale

### Uso personale (situazione attuale)
Con un solo utente, l'intero stack resta a **$0**. Non ci sono limiti realistici di storage, connessioni DB o chiamate API che si possano raggiungere con un uso individuale.

### Se in futuro si apre ad altri utenti
Alcuni costi diventano probabili, non tanto per la tecnologia in sé quanto per lo scaling:

| Servizio | Quando diventa a pagamento |
|---|---|
| Vercel | Superati i limiti di banda/build del piano Hobby (e comunque il piano Hobby vieta uso commerciale) |
| Database (Neon) | Superato lo storage o le ore di compute del piano free |
| Hosting backend (Fly.io/Render) | Superate le ore/risorse incluse nel piano free |
| Notifiche push | Servizi come OneSignal hanno free tier ma con limiti di invii |

---

## ⚖️ 7. Considerazioni legali per un'eventuale versione commerciale

Da verificare **prima** di aprire il progetto ad altri utenti o di introdurre pubblicità/abbonamenti:

1. **ToS di TMDB e AniList:** entrambe richiedono attribuzione visibile ("Powered by TMDB", ecc.) e hanno clausole specifiche sull'uso commerciale dei dati oltre una certa scala — da rileggere direttamente sui loro siti prima di monetizzare.
2. **GDPR:** informativa privacy, possibilità di export/cancellazione dati utente, hosting DB in region conforme (Neon permette di scegliere la region in EU).
3. **ToS di Vercel:** il piano Hobby gratuito vieta esplicitamente l'uso commerciale.
4. **Pagamenti:** se si introduce un piano "supporter", va integrato un provider come Stripe o LemonSqueezy.

---

## 📈 8. Possibili strategie di sostenibilità futura

Se il progetto dovesse aprirsi ad altri utenti, le opzioni più realistiche (vista la natura "no-ads invasive" del progetto) sono:

- **Donazioni** (Ko-fi, Buy Me a Coffee, GitHub Sponsors)
- **Piano "Supporter"** con funzionalità extra (statistiche avanzate, export, temi personalizzati) invece di pubblicità
- **Pubblicità non invasiva** (es. banner statici orientati a sviluppatori/appassionati), da valutare solo con traffico consistente

---

## 🚀 9. Setup del Progetto

### Prerequisiti
- [Go](https://go.dev/dl/) 1.22+
- [Node.js](https://nodejs.org/) 20+ e Angular CLI (`npm install -g @angular/cli`)
- Un account [Neon.tech](https://neon.tech) (Postgres free tier)
- Chiavi API: [TMDB](https://www.themoviedb.org/documentation/api), [AniList](https://anilist.gitbook.io/anilist-apiv2-docs/) (AniList non richiede API key per le query pubbliche via GraphQL)

### Backend
```bash
cd backend
go mod init github.com/AntonioMartino67/showio/backend
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get github.com/joho/godotenv
go run cmd/api/main.go
```

### Frontend
```bash
cd frontend
ng new showio-frontend --directory=. --routing --style=css
npm install
ng serve
```

### Variabili d'ambiente (backend/.env — da creare, non versionato)
```
DATABASE_URL=postgresql://user:password@host/dbname?sslmode=require
JWT_SECRET=una_stringa_segreta_lunga_e_casuale
TMDB_API_KEY=la_tua_chiave_tmdb
PORT=8080
```

Sezioni da completare più avanti: configurazione del workflow GitHub Actions per il cron notturno, istruzioni di deploy su Fly.io/Render (backend) e Vercel (frontend).

---

## 📚 10. Riferimenti API Esterne

- **TMDB (The Movie Database):** https://www.themoviedb.org/documentation/api
- **AniList API:** https://anilist.gitbook.io/anilist-apiv2-docs/

---

## 📝 11. Note di Sviluppo

Questo README viene mantenuto aggiornato man mano che le decisioni architetturali evolvono. Cronologia decisioni principali:
- Stack iniziale valutato: React + FastAPI (Python)
- Stack alternativo valutato: Go + SvelteKit
- **Stack finale scelto:** Go (backend, nuovo apprendimento) + Angular (frontend, tecnologia già nota) — bilancia investimento in una nuova competenza con la necessità di consegnare risultati concreti durante il periodo dell'Esame di Stato.
