# Graph Report - mindbenders  (2026-09-06)

## Corpus Check
- 58 files · ~11,061 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 337 nodes · 596 edges · 23 communities (18 shown, 3 thin omitted)
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
- testing.T
- error.go
- getlogger
- S3Uploader
- iConfig
- Option
- prometheus.go
- logOption
- Eval
- Init
- Delta
- For Pushing to Topic
- id.go
- Errors
- base62.go
- logging/README.md
- gitlab.com/dotpe/mindbenders

## God Nodes (most connected - your core abstractions)
1. `Fields` - 32 edges
2. `profiler` - 24 edges
3. `dlogger` - 20 edges
4. `metrics` - 15 edges
5. `ProfileType` - 15 edges
6. `CoRelationId` - 14 edges
7. `NewCorelId()` - 13 edges
8. `Level` - 13 edges
9. `GetCorelationId()` - 11 edges
10. `emptyLogger` - 10 edges

## Surprising Connections (you probably didn't know these)
- `GetFileConfigManager()` --calls--> `WrapMessage()`  [EXTRACTED]
  bootconfig/fileconf.go → errors/error.go
- `NewPacket()` --calls--> `NewContext()`  [EXTRACTED]
  kafka/packet/packet.go → corel/context.go
- `accessLogOptionBasic()` --calls--> `GetCorelationId()`  [EXTRACTED]
  logging/logOption.go → corel/util.go
- `HeaderLogoption()` --references--> `Fields`  [EXTRACTED]
  middleware/uiware/header-logger.go → logging/fields.go
- `WriteDBError()` --references--> `IDotpeLogger`  [EXTRACTED]
  clients/mysql/db.go → logging/interface.go

## Import Cycles
- None detected.

## Communities (23 total, 3 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.11
Nodes (14): context.Context, github.com/prometheus/client_golang/prometheus.CounterVec, go.uber.org/zap.Field, emptyLogger, caller(), canonicalFile(), Fields, Level (+6 more)

### Community 1 - "CoRelationId"
Cohesion: 0.10
Nodes (27): CoRelationId, NewCorelId(), NewCorelIdFromHttp(), corelstr, AmqpLoader(), AmqpUnloader(), HttpCorelLoader(), HttpCorelUnLoader() (+19 more)

### Community 2 - "profiler"
Cohesion: 0.11
Nodes (24): bytes.Buffer, runtime.MemStats, sync.Once, sync.WaitGroup, time.Duration, time.Time, collectionTooFrequent, Batch (+16 more)

### Community 3 - "ProfileType"
Cohesion: 0.15
Nodes (12): profiler, Config, Option, defaultConfig(), WithService(), WithTargetSetter(), WithUploader(), collectGenericProfile() (+4 more)

### Community 4 - "testing.T"
Cohesion: 0.14
Nodes (13): Test_base_Error(), Test_base_String(), TestCause(), TestUnWrap(), DefaultMultiError(), NewMultiError(), Test_multierror_AddErrors(), TestDefaultMultiError() (+5 more)

### Community 5 - "error.go"
Cohesion: 0.13
Nodes (7): secretManager, BasicError, causer, Cause(), Code(), UnWrap(), WrapMessage()

### Community 6 - "getlogger"
Cohesion: 0.14
Nodes (15): WriteDBError(), github.com/rs/zerolog.Logger, go.uber.org/zap.Logger, DefaultLogWriter(), getlogger(), dlogger, LogWriter(), MustGet() (+7 more)

### Community 7 - "S3Uploader"
Cohesion: 0.14
Nodes (12): github.com/aws/aws-sdk-go/aws.Config, github.com/aws/aws-sdk-go/aws/session.Session, github.com/aws/aws-sdk-go/service/s3.S3, io.ReadSeeker, GetFileSaver(), NewNullUploader(), NewUploaderWithBucket(), NewUploaderWithConfig() (+4 more)

### Community 8 - "iConfig"
Cohesion: 0.17
Nodes (8): conf, ConfigManager, GetFileConfigManager(), fileConfig, iConfig, Init(), MustInit(), GetSecretManager()

### Community 9 - "Option"
Cohesion: 0.17
Nodes (11): NewChildContext(), NewContext(), NewOrphanContext(), TestNewCorelCtx(), Test_dlogger_Write(), Option, DisabledStdLogging(), dlogger (+3 more)

### Community 10 - "prometheus.go"
Cohesion: 0.22
Nodes (10): github.com/gin-gonic/gin.Engine, github.com/zsais/go-gin-prometheus.Metric, github.com/zsais/go-gin-prometheus.Prometheus, profOpt, customCollector(), initializeCollector(), SetPrometheusMetricsOnGin(), WithRouter() (+2 more)

### Community 11 - "logOption"
Cohesion: 0.20
Nodes (9): github.com/gin-gonic/gin.Context, accessLogOption, logOption, accessLogOptionBasic(), AccessLogOptionRequestBody(), WithAccessLogOptions(), Health(), PostJSONValidator() (+1 more)

### Community 12 - "Eval"
Cohesion: 0.31
Nodes (8): github.com/gin-gonic/gin.HandlerFunc, golang.org/x/net/context.Context, dlogger, Eval(), PostwareLimitUser(), PrewareLimitUser(), redis_rate.Limit, redis_rate.Limiter

### Community 13 - "Init"
Cohesion: 0.31
Nodes (5): Init(), database/sql.DB, DBCollectorLables, WithDBPoolMetrics(), Option

### Community 14 - "Delta"
Cohesion: 0.25
Nodes (5): github.com/google/pprof/profile.Profile, io.Writer, Protobuf, ValueType, Delta

### Community 15 - "For Pushing to Topic"
Cohesion: 0.29
Nodes (6): consumer side, COREL ***(Correlation in Distributed Logs)***, For Api Call, For Pushing to Topic, publisher side, Usages

### Community 16 - "id.go"
Cohesion: 0.47
Nodes (4): github.com/bwmarrin/snowflake.Node, defaultInit(), GetGenerator(), SetNode()

### Community 17 - "Errors"
Cohesion: 0.50
Nodes (3): Errors, Problem, Solutions

## Knowledge Gaps
- **12 isolated node(s):** `corelstr`, `causer`, `gitlab.com/dotpe/mindbenders`, `dlogger`, `dlogger` (+7 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 72 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `WriteDBError()` connect `getlogger` to `context.Context`, `Init`?**
  _High betweenness centrality (0.385) - this node is a cross-community bridge._
- **Why does `Option` connect `Init` to `profiler`?**
  _High betweenness centrality (0.351) - this node is a cross-community bridge._
- **Why does `Config` connect `ProfileType` to `profiler`?**
  _High betweenness centrality (0.201) - this node is a cross-community bridge._
- **What connects `corelstr`, `causer`, `gitlab.com/dotpe/mindbenders` to the rest of the system?**
  _12 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.10852713178294573 - nodes in this community are weakly interconnected._
- **Should `CoRelationId` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `profiler` be split into smaller, more focused modules?**
  _Cohesion score 0.10526315789473684 - nodes in this community are weakly interconnected._