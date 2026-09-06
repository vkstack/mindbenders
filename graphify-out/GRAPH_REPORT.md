# Graph Report - mindbenders  (2026-09-06)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 323 nodes · 586 edges · 15 communities (12 shown, 2 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 41 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0bfd4ef5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- CoRelationId
- profiler
- ProfileType
- Eval
- getlogger
- testing.T
- error.go
- S3Uploader
- iConfig
- Option
- id.go
- base62.go
- gitlab.com/dotpe/mindbenders

## God Nodes (most connected - your core abstractions)
1. `Fields` - 32 edges
2. `profiler` - 24 edges
3. `dlogger` - 20 edges
4. `metrics` - 15 edges
5. `ProfileType` - 15 edges
6. `CoRelationId` - 14 edges
7. `Level` - 13 edges
8. `NewCorelId()` - 13 edges
9. `GetCorelationId()` - 11 edges
10. `emptyLogger` - 10 edges

## Surprising Connections (you probably didn't know these)
- `HeaderLogoption()` --references--> `Fields`  [EXTRACTED]
  middleware/uiware/header-logger.go → logging/fields.go
- `accessLogOptionBasic()` --calls--> `GetCorelationId()`  [EXTRACTED]
  logging/logOption.go → corel/util.go
- `NewPacket()` --calls--> `NewContext()`  [EXTRACTED]
  kafka/packet/packet.go → corel/context.go
- `GetFileConfigManager()` --calls--> `WrapMessage()`  [EXTRACTED]
  bootconfig/fileconf.go → errors/error.go
- `Packet` --references--> `CoRelationId`  [EXTRACTED]
  kafka/packet/packet.go → corel/corel.go

## Import Cycles
- None detected.

## Communities (15 total, 2 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.11
Nodes (14): context.Context, github.com/prometheus/client_golang/prometheus.CounterVec, go.uber.org/zap.Field, emptyLogger, caller(), canonicalFile(), Fields, Level (+6 more)

### Community 1 - "CoRelationId"
Cohesion: 0.10
Nodes (27): CoRelationId, NewCorelId(), NewCorelIdFromHttp(), corelstr, AmqpLoader(), AmqpUnloader(), HttpCorelLoader(), HttpCorelUnLoader() (+19 more)

### Community 2 - "profiler"
Cohesion: 0.11
Nodes (25): bytes.Buffer, runtime.MemStats, sync.Once, sync.WaitGroup, time.Duration, time.Time, collectionTooFrequent, Batch (+17 more)

### Community 3 - "ProfileType"
Cohesion: 0.09
Nodes (17): github.com/google/pprof/profile.Profile, io.Writer, Protobuf, ValueType, profiler, Config, Profile, Delta (+9 more)

### Community 4 - "Eval"
Cohesion: 0.09
Nodes (24): github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.Engine, github.com/gin-gonic/gin.HandlerFunc, github.com/zsais/go-gin-prometheus.Metric, github.com/zsais/go-gin-prometheus.Prometheus, golang.org/x/net/context.Context, dlogger, AccessLogOptionRequestBody() (+16 more)

### Community 5 - "getlogger"
Cohesion: 0.10
Nodes (20): Init(), WriteDBError(), database/sql.DB, github.com/rs/zerolog.Logger, go.uber.org/zap.Logger, DefaultLogWriter(), getlogger(), dlogger (+12 more)

### Community 6 - "testing.T"
Cohesion: 0.10
Nodes (18): NewChildContext(), NewContext(), NewOrphanContext(), TestNewCorelCtx(), Test_base_Error(), Test_base_String(), TestCause(), TestUnWrap() (+10 more)

### Community 7 - "error.go"
Cohesion: 0.13
Nodes (7): secretManager, BasicError, causer, Cause(), Code(), UnWrap(), WrapMessage()

### Community 8 - "S3Uploader"
Cohesion: 0.14
Nodes (12): github.com/aws/aws-sdk-go/aws.Config, github.com/aws/aws-sdk-go/aws/session.Session, github.com/aws/aws-sdk-go/service/s3.S3, io.ReadSeeker, GetFileSaver(), NewNullUploader(), NewUploaderWithBucket(), NewUploaderWithConfig() (+4 more)

### Community 9 - "iConfig"
Cohesion: 0.17
Nodes (8): conf, ConfigManager, GetFileConfigManager(), fileConfig, iConfig, Init(), MustInit(), GetSecretManager()

### Community 10 - "Option"
Cohesion: 0.27
Nodes (9): accessLogOption, logOption, accessLogOptionBasic(), Option, DisabledStdLogging(), dlogger, WithAccessLogOptions(), WithZap() (+1 more)

### Community 11 - "id.go"
Cohesion: 0.47
Nodes (4): github.com/bwmarrin/snowflake.Node, defaultInit(), GetGenerator(), SetNode()

## Knowledge Gaps
- **6 isolated node(s):** `dlogger`, `dlogger`, `corelstr`, `gitlab.com/dotpe/mindbenders`, `dlogger` (+1 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 62 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `WriteDBError()` connect `getlogger` to `context.Context`?**
  _High betweenness centrality (0.419) - this node is a cross-community bridge._
- **Why does `Option` connect `getlogger` to `profiler`?**
  _High betweenness centrality (0.382) - this node is a cross-community bridge._
- **Why does `Config` connect `ProfileType` to `profiler`?**
  _High betweenness centrality (0.219) - this node is a cross-community bridge._
- **What connects `dlogger`, `dlogger`, `corelstr` to the rest of the system?**
  _6 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.10852713178294573 - nodes in this community are weakly interconnected._
- **Should `CoRelationId` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `profiler` be split into smaller, more focused modules?**
  _Cohesion score 0.11095305832147938 - nodes in this community are weakly interconnected._