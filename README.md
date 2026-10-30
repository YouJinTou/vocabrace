# Wordlines (vocabrace)

**A competitive multiplayer word game built around a shared board.**

Wordlines is a crossword-inspired game in which players take turns placing letter tiles to form words and earn points. The project combines a Go game engine with a real-time, serverless multiplayer backend and an Angular interface.

## What it does

- **Multiplayer wordplay:** Players place words on a shared board and compete for points.
- **Game rules and scoring:** The Go engine handles tiles, move validation, scoring, turn order, and end-of-game logic.
- **Turn actions:** Players can place a word, exchange tiles, or pass; the interface also handles turn timeouts.
- **Alternative layouts:** The game engine includes classic and spiral board layouts.
- **Real-time sessions:** WebSocket messaging connects players to games, with matchmaking, connection handling, and reconnection flows.

## Technology

- **Game engine and services:** Go 1.15
- **Frontend:** Angular 10, TypeScript, Angular Material
- **Communication:** Amazon API Gateway WebSocket APIs
- **Backend infrastructure:** AWS Lambda, DynamoDB, SNS
- **Infrastructure as code:** Terraform

## Repository structure

```text
src/
  games/wordlines/       Game rules, board, tiles, validation, and scoring
  services/iam/          Authentication-related service
  services/pooling/      Matchmaking and multiplayer connection services
  services/com/          Shared communication and game-state logic
  tools/                 Shared Go utilities
ui/                      Angular frontend
tf/                      Terraform modules and deployment scripts
```

## Development

The Go game logic includes unit tests:

```bash
cd src
go test ./...
```

To start the frontend development server:

```bash
cd ui
npm install
npm start
```

The frontend relies on configured authentication and WebSocket services for multiplayer functionality. A complete working deployment requires AWS resources and environment configuration; launching the Angular app alone does **not** start the backend.

The repository also contains historical build and Terraform deployment scripts in `tf/`. Review their configuration and resource costs before using them.

## Status

**Historical project (2020).** Wordlines is preserved as an example of multiplayer game development, Go application logic, and event-driven AWS infrastructure. It is not presented as a currently maintained public game.
