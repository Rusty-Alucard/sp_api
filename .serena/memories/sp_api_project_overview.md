# SP API Project Overview

## Project Summary
SP API is a GraphQL-based soccer/football data management system built with Go. It follows Clean Architecture principles and manages tournaments, teams, players, matches, and related data.

## Technical Stack
- **Language**: Go 1.16
- **GraphQL**: gqlgen (99designs/gqlgen)
- **Database**: PostgreSQL 13.2
- **Dependency Injection**: Google Wire
- **Configuration**: Viper
- **Containerization**: Docker & Docker Compose
- **Additional Services**: Redis, HashiCorp Vault

## Architecture Layers
1. **Domain Layer** (`domain/`): Core business logic and entities
   - `model/`: Generated GraphQL models and enums
   - `repository/`: Repository interfaces

2. **Use Case Layer** (`usecase/`): Business use cases
   - Event and Team use cases with CRUD operations

3. **Infrastructure Layer** (`infrastructure/`): External concerns
   - `persistence/`: PostgreSQL data access
   - `provider/`: Database connection management

4. **GraphQL Layer** (`graph/`): API interface
   - `generated/`: Auto-generated GraphQL code
   - `resolver/`: Query resolvers

## Key Entities
- **Event**: Tournaments (FIFA World Cup, etc.)
- **Team**: National teams with FIFA codes and confederations
- **Player**: Individual players with positions and stats
- **Coach**: Team coaches
- **Match**: Individual matches with stadiums and timelines
- **Stadium**: Venues with timezone information
- **Timeline**: Match events (goals, cards, substitutions)

## Database Schema
- `events`: Tournament information
- `teams`: Team data with FIFA trigrammes
- `players`: Player information
- `coaches`: Coach data
- `matches`: Match details
- `stadiums`: Venue information
- `participants`: Event-team relationships
- `registered_players`: Player registrations
- `match_progresses`: Match timeline events

## Configuration
- Environment-based config using Viper
- YAML configuration files in `config/environments/`
- Database connection pooling
- Docker Compose for local development

## Development Workflow
- GraphQL schema-first development
- Code generation with gqlgen
- Dependency injection with Wire
- Docker-based infrastructure

## Current Status
- Basic CRUD operations for Events and Teams
- GraphQL playground available
- PostgreSQL with proper indexing
- Clean architecture implementation
- Docker development environment ready

