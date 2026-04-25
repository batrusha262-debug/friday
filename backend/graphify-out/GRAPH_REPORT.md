# Graph Report - backend  (2026-04-25)

## Corpus Check
- 62 files · ~14,827 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 382 nodes · 825 edges · 18 communities detected
- Extraction: 46% EXTRACTED · 54% INFERRED · 0% AMBIGUOUS · INFERRED: 443 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]

## God Nodes (most connected - your core abstractions)
1. `New()` - 64 edges
2. `Suite` - 37 edges
3. `Handler` - 27 edges
4. `PgRepository` - 27 edges
5. `repoStub` - 27 edges
6. `stubService` - 26 edges
7. `Service` - 26 edges
8. `parseID()` - 23 edges
9. `JSON()` - 21 edges
10. `NewService()` - 15 edges

## Surprising Connections (you probably didn't know these)
- `TestSuite()` --calls--> `New()`  [INFERRED]
  integration/suite_test.go → pkg/postgres/postgres.go
- `NewGameID()` --calls--> `New()`  [INFERRED]
  internal/pack/domain/values/id_game.go → pkg/postgres/postgres.go
- `NewGameTeamID()` --calls--> `New()`  [INFERRED]
  internal/pack/domain/values/id_game.go → pkg/postgres/postgres.go
- `TestCreatePack_validation()` --calls--> `New()`  [INFERRED]
  internal/pack/domain/service/service_test.go → pkg/postgres/postgres.go
- `TestCreatePack_callsRepo()` --calls--> `New()`  [INFERRED]
  internal/pack/domain/service/service_test.go → pkg/postgres/postgres.go

## Communities

### Community 0 - "Community 0"
Cohesion: 0.08
Nodes (8): Category, decode(), parseID(), JSON(), NoContent(), Handler, Service, Service

### Community 1 - "Community 1"
Cohesion: 0.06
Nodes (20): TestError_statusCodes(), Handler(), Error(), errorResponse, supportID(), RoundTypeValues(), NewService(), repoStub (+12 more)

### Community 2 - "Community 2"
Cohesion: 0.18
Nodes (4): Application, New(), Client, Suite

### Community 3 - "Community 3"
Cohesion: 0.07
Nodes (13): NewCategoryID(), NewPackID(), NewQuestionID(), NewRoundID(), main(), Config, New(), TestRegister_allRoutesReachable() (+5 more)

### Community 4 - "Community 4"
Cohesion: 0.11
Nodes (8): PgRepository, ForeignKeyViolation(), IsForeignKeyViolation(), IsNotFound(), IsUniqueViolation(), NotFound(), pgErrCode(), UniqueViolation()

### Community 5 - "Community 5"
Cohesion: 0.14
Nodes (11): Config, getenv(), Load(), PostgresConfig, TestLoad_defaults(), TestLoad_envOverride(), NewClient(), NewHandler() (+3 more)

### Community 6 - "Community 6"
Cohesion: 0.15
Nodes (9): contextKeyLogger, contextKeyTraceID, TraceID, enum, EnrichLogger(), LoggerFromContext(), LoggerFromContextOrDefault(), WithLogger() (+1 more)

### Community 7 - "Community 7"
Cohesion: 0.19
Nodes (8): roundType, RoundTypeEnum, QuestionTypeValues(), validateQuestion(), ParseRoundType(), ParseRoundTypeEmpty(), ParseRoundTypeFold(), ParseRoundTypeOptional()

### Community 8 - "Community 8"
Cohesion: 0.23
Nodes (6): gameStatus, GameStatusEnum, ParseGameStatus(), ParseGameStatusEmpty(), ParseGameStatusFold(), ParseGameStatusOptional()

### Community 9 - "Community 9"
Cohesion: 0.25
Nodes (6): questionType, QuestionTypeEnum, ParseQuestionType(), ParseQuestionTypeEmpty(), ParseQuestionTypeFold(), ParseQuestionTypeOptional()

### Community 10 - "Community 10"
Cohesion: 0.22
Nodes (8): Category, Game, GameBoard, GameQuestionState, GameTeam, Pack, Question, Round

### Community 11 - "Community 11"
Cohesion: 0.32
Nodes (5): GameQuestionState, NewGameID(), NewGameTeamID(), GameID, GameTeamID

### Community 12 - "Community 12"
Cohesion: 0.67
Nodes (1): GameTeam

### Community 13 - "Community 13"
Cohesion: 0.67
Nodes (1): Pack

### Community 14 - "Community 14"
Cohesion: 0.67
Nodes (1): Question

### Community 15 - "Community 15"
Cohesion: 0.67
Nodes (1): Game

### Community 16 - "Community 16"
Cohesion: 0.67
Nodes (1): Round

### Community 17 - "Community 17"
Cohesion: 1.0
Nodes (1): Repository

## Knowledge Gaps
- **14 isolated node(s):** `Repository`, `Service`, `Pack`, `Round`, `Category` (+9 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 12`** (3 nodes): `GameTeam`, `.ToDomain()`, `game_team.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 13`** (3 nodes): `Pack`, `.ToDomain()`, `pack.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 14`** (3 nodes): `Question`, `.ToDomain()`, `question.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 15`** (3 nodes): `Game`, `.ToDomain()`, `game.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 16`** (3 nodes): `Round`, `.ToDomain()`, `round.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 17`** (2 nodes): `pack.go`, `Repository`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Community 3` to `Community 1`, `Community 2`, `Community 11`, `Community 5`?**
  _High betweenness centrality (0.247) - this node is a cross-community bridge._
- **Why does `PgRepository` connect `Community 4` to `Community 5`?**
  _High betweenness centrality (0.085) - this node is a cross-community bridge._
- **Why does `NewPgRepository()` connect `Community 5` to `Community 2`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Are the 62 inferred relationships involving `New()` (e.g. with `.TestCreateGame()` and `.TestGetGame()`) actually correct?**
  _`New()` has 62 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Repository`, `Service`, `Pack` to the rest of the system?**
  _14 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.08 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.06 - nodes in this community are weakly interconnected._