# Edwin Metcalf — Portfolio
Personal site and playground for side projects, built with a Go backend and SvelteKit frontend.

🔗 [edwinmetcalf.com](https://edwinmetcalf.com)

## Features
- **Portfolio** - information about pass projects, tools and languages I am familiar with, past classes and more
- **Game Stats Dashboard** — personal Dota 2 and Clash Royale stats, synced from each game's public API into a local database with historical trend charts
- **Space Invaders** — a from-scratch playable mix of classic space invaders and vampires survivors with custom art and enemies
- More in progress...

## Tech Stack
- **Frontend:** SvelteKit, Svelte 5 (runes), Chart.js
- **Backend:** Go
- **Database:** SQLite (dev) / PostgreSQL (prod)

## In Progress
- **Games:** Working on a minimetro style casual game
- **Head to Head CR:** an interface to put you and your friends CR tags in and see your last 25 game matchups

## Structure
- `/frontend` - Svelte frontend
- `/backend` - Go backend

## Development

### Frontend
```bash
cd frontend
npm install
npm run dev
```

### Backend
```bash
cd backend
go run main.go
```
