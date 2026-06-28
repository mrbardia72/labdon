# Architecture

## Layers

API
↓
Service
↓
Repository
↓
Database

Rules:
- Services must not access DB directly
- Repository handles persistence
- API layer only handles request/response