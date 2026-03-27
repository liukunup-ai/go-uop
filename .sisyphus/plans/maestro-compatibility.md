# Maestro Compatibility for Go-UOP

## TL;DR

> **Quick Summary**: Add Maestro compatibility to Go-UOP by implementing a translator that converts Maestro YAML files into Go-UOP internal actions, plus a CLI tool to execute Maestro flows.
>
> **Deliverables**:
> - `maestro/` package - Maestro YAML parser and command translator
> - `maestro/commands/` - Per-command translation modules
> - `cmd/maestro/` - CLI tool with test/validate/help subcommands
> - 15 core Maestro commands implemented
>
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 3 waves
> **Critical Path**: Types → Translator base → Commands → CLI → Integration

---

## Context

### Original Request
兼容Maestro全部命令 (Support all Maestro commands), workspace configuration, and element selectors.

### Interview Summary

**Key Discussions**:
- Command scope: User selected "Standard automation set" (~15 commands) instead of all 40+
- Output type: Both library AND CLI tool
- Compatibility mode: Support BOTH Go-UOP native YAML AND existing Maestro YAML files
- Test strategy: Tests after implementation (no TDD)
- Priority commands: tapOn, swipe, inputText, launchApp, killApp, stopApp, assertVisible, waitForAnimationToEnd, back, pressKey, clearState, runFlow, takeScreenshot, scroll

**Research Findings**:
- Go-UOP has YAML parsing foundation but lacks end-to-end execution path
- `yaml/commands/control.go` only handles If/Foreach/While + basic Launch/Wait
- Internal actions (TapAction, SwipeAction, SendKeysAction) exist but not wired to YAML
- Selectors (ByText, ByID, ByXPath) exist but lack relational selectors (above/below)
- No CLI tool exists

### Metis Review

**Identified Gaps** (addressed):
- Translation layer (Maestro YAML → Internal Actions) NOT implemented
- CLI tool completely missing
- Workspace config.yaml not implemented
- 35+ commands not implemented (only 15 needed for Phase 1)

---

## Work Objectives

### Core Objective
Implement Maestro YAML compatibility as a Phase 1 deliverable with 15 core commands, enabling both library usage and CLI execution.

### Concrete Deliverables
- `maestro/maestro.go` - Maestro YAML parser
- `maestro/types.go` - Maestro-specific types
- `maestro/translator.go` - Core translation engine
- `maestro/commands/*.go` - Per-command translators
- `cmd/maestro/main.go` - CLI entry point
- `config.yaml` workspace configuration support

### Definition of Done
- [ ] `go build ./cmd/maestro` compiles successfully
- [ ] `go test ./maestro/... -v` passes all tests
- [ ] `./maestro test flow.maestro.yaml` executes sample flow
- [ ] `./maestro --help` displays CLI help

### Must Have
- 15 core commands translated correctly
- CLI tool with test/validate/help subcommands
- Maestro YAML file parsing (file extension: `.maestro.yaml`)
- Variable substitution (`${var}`, `${ENV()}`)
- Error handling with descriptive messages

### Must NOT Have (Guardrails)
- ❌ Relational selectors (above/below/leftOf/rightOf) in Phase 1
- ❌ Maestro Studio integration
- ❌ JavaScript expression evaluation (only ${var} substitution)
- ❌ Parallel flow execution
- ❌ Cloud upload features
- ❌ HTML/PDF reports
- ❌ Extended appState commands
- ❌ copyTextFrom / pasteFromClipboard

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: YES (Go testing framework)
- **Automated tests**: Tests-after
- **Framework**: Go standard `testing` package
- **Test coverage**: Each command translator has unit tests

### QA Policy
Every task MUST include agent-executed QA scenarios. Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **CLI Tool**: Use `Bash` - Run command, assert exit code and output
- **YAML Parsing**: Use `Bash` (go test) - Parse YAML, assert structures
- **Translator**: Use `Bash` (go test) - Unit tests for translation logic

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Foundation - MUST complete first):
├── Task 1: Create maestro package structure + types
├── Task 2: Implement Maestro YAML parser (maestro.go)
├── Task 3: Implement core translator engine (translator.go)
└── Task 4: Create cmd/maestro CLI skeleton

Wave 2 (Command Implementations - PARALLEL after Wave 1):
├── Task 5: tapOn command translator
├── Task 6: swipe command translator
├── Task 7: inputText command translator
├── Task 8: launchApp/killApp/stopApp translators
├── Task 9: assertVisible/assertNotVisible translators
├── Task 10: waitForAnimationToEnd/pressKey/back translators
├── Task 11: scroll/takeScreenshot translators
└── Task 12: runFlow subflow executor

Wave 3 (Integration + Workspace Config):
├── Task 13: Workspace config.yaml support
├── Task 14: CLI integration (wire translator to executor)
└── Task 15: End-to-end integration test

Wave FINAL (4 parallel reviews):
├── Task F1: Plan compliance audit
├── Task F2: Code quality review (build + test)
├── Task F3: CLI smoke test
└── Task F4: Scope fidelity check
```

### Dependency Matrix

| Task | Depends On | Blocks |
|------|------------|--------|
| 1 | - | 2, 3, 4 |
| 2 | 1 | 3 |
| 3 | 2 | 5-12 |
| 4 | 1 | 14 |
| 5 | 3 | 15 |
| 6 | 3 | 15 |
| 7 | 3 | 15 |
| 8 | 3 | 15 |
| 9 | 3 | 15 |
| 10 | 3 | 15 |
| 11 | 3 | 15 |
| 12 | 3 | 15 |
| 13 | 4 | 14 |
| 14 | 4, 13 | 15 |
| 15 | 5-12, 14 | F1-F4 |

### Agent Dispatch Summary

- **Wave 1**: **4 tasks** → `deep` (architectural foundation)
- **Wave 2**: **8 tasks** → `unspecified-high` (command implementations)
- **Wave 3**: **3 tasks** → `unspecified-high` (integration)
- **FINAL**: **4 tasks** → `oracle`, `unspecified-high`, `unspecified-high`, `deep`

---

## TODOs

---

## TODOs

- [x] 1. Create maestro package structure + types

  **What to do**:
  - Create `maestro/` directory with package structure
  - Define `maestro/types.go` with Maestro-specific types:
    - `MaestroFlow` - Flow metadata (appId, name, tags)
    - `MaestroCommand` - Base command interface
    - Selector types: `TextSelector`, `IDSelector`, `IndexSelector`, `PointSelector`
  - Define `maestro/errors.go` for translation errors
  - Create `maestro/maestro.go` with Maestro YAML detection (file extension `.maestro.yaml`)
  - Follow existing Go-UOP patterns from `yaml/command.go`

  **Must NOT do**:
  - Don't modify existing `yaml/command.go`
  - Don't create parallel YAML structures - Maestro types are Maestro-specific

  **Recommended Agent Profile**:
  > **Category**: `deep`
  > - Reason: Architectural foundation - establishes package structure and type contracts
  > **Skills**: []
  > - No specific skills needed - follows existing Go patterns

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 1 (with Tasks 2, 3, 4)
  - **Blocks**: Tasks 2, 3, 4
  - **Blocked By**: None (can start immediately)

  **References** (CRITICAL):
  - `yaml/command.go:1-100` - Existing YAML command structure patterns
  - `yaml/parser.go:1-50` - Parser function signatures to follow
  - `internal/action/action.go:1-80` - Action type definitions for reference

  **Acceptance Criteria**:
  - [ ] `maestro/types.go` defines all Maestro-specific types
  - [ ] `go build ./maestro/...` compiles without errors
  - [ ] `mkdir -p maestro && ls maestro/` shows correct structure

  **QA Scenarios**:

  \`\`\`
  Scenario: Maestro package builds successfully
    Tool: Bash
    Preconditions: Clean workspace, no maestro/ directory
    Steps:
      1. mkdir -p maestro
      2. Create maestro/types.go with basic types
      3. Run: go build ./maestro/...
    Expected Result: Build succeeds with no errors
    Evidence: .sisyphus/evidence/task-1-build.{ext}

  Scenario: Maestro types are properly defined
    Tool: Bash
    Preconditions: maestro/types.go exists
    Steps:
      1. Run: go vet ./maestro/...
      2. Run: go fmt ./maestro/...
    Expected Result: No vet errors, properly formatted
    Evidence: .sisyphus/evidence/task-1-vet.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] Build output: task-1-build.log
  - [ ] Vet output: task-1-vet.log

  **Commit**: YES
  - Message: `feat(maestro): add Maestro package structure and types`
  - Files: `maestro/types.go`, `maestro/errors.go`, `maestro/maestro.go`

---

- [x] 2. Implement Maestro YAML parser

  **What to do**:
  - Implement `maestro.ParseFlow(reader io.Reader) (*MaestroFlow, error)` in `maestro/parser.go`
  - Support Maestro YAML format:
    \`\`\`yaml
    appId: com.example.app
    name: Test Flow
    tags: [smoke, auth]
    ---
    - tapOn: "Login"
    - inputText: "user@test.com"
    \`\`\`
  - Parse flow metadata (appId, name, tags)
  - Parse commands array (shorthand and extended selector forms)
  - Return descriptive errors for malformed YAML

  **Must NOT do**:
  - Don't execute commands - only parse
  - Don't require appId if not present (warn but continue)

  **Recommended Agent Profile**:
  > **Category**: `unspecified-high`
  > - Reason: YAML parsing with complex nested structures
  > **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3, 4)
  - **Blocks**: Task 3
  - **Blocked By**: Task 1

  **References**:
  - `yaml/parser.go:ParseFlow` - Existing parsing pattern to follow
  - `yaml/evaluator.go:Context` - Variable substitution context
  - Maestro docs: YAML format with `---` separator

  **Acceptance Criteria**:
  - [ ] `maestro.ParseFlow` parses valid Maestro YAML
  - [ ] `maestro.ParseFlow` returns error for invalid YAML
  - [ ] `maestro.ParseFlow` extracts appId, name, tags correctly

  **QA Scenarios**:

  \`\`\`
  Scenario: Parse valid Maestro YAML with appId and commands
    Tool: Bash
    Preconditions: maestro/parser.go exists
    Steps:
      1. Create test file with Maestro YAML
      2. Run: go test ./maestro -run TestParseFlow -v
    Expected Result: Parses successfully, appId="com.example.app", name="Test Flow"
    Evidence: .sisyphus/evidence/task-2-parse-valid.{ext}

  Scenario: Parse invalid YAML returns descriptive error
    Tool: Bash
    Preconditions: maestro/parser.go exists
    Steps:
      1. Create invalid YAML file
      2. Run: go test ./maestro -run TestParseFlowInvalid -v
    Expected Result: Error message includes line number and reason
    Evidence: .sisyphus/evidence/task-2-parse-invalid.{ext}

  Scenario: Parse YAML with shorthand selectors
    Tool: Bash
    Preconditions: maestro/parser.go exists
    Steps:
      1. Create YAML with shorthand: - tapOn: "Login"
      2. Run: go test ./maestro -run TestParseShorthand -v
    Expected Result: TapOn command with TextSelector "Login"
    Evidence: .sisyphus/evidence/task-2-shorthand.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] Parse valid output: task-2-parse-valid.log
  - [ ] Parse invalid output: task-2-parse-invalid.log
  - [ ] Shorthand parse: task-2-shorthand.log

  **Commit**: YES
  - Message: `feat(maestro): implement YAML parser for Maestro flows`
  - Files: `maestro/parser.go`, `maestro/parser_test.go`

---

- [x] 3. Implement core translator engine

  **What to do**:
  - Create `maestro/translator.go` with translation engine
  - Implement `Translator` struct with context and device
  - Implement `TranslateCommand(cmd *MaestroCommand) (action.Action, error)`
  - Implement `TranslateFlow(flow *MaestroFlow) ([]action.Action, error)`
  - Wire up selector building using `internal/selector` package
  - Implement command routing (tapOn → TapAction, inputText → SendKeysAction, etc.)

  **Must NOT do**:
  - Don't implement all commands - only the routing framework
  - Don't execute actions - just build them

  **Recommended Agent Profile**:
  > **Category**: `deep`
  > - Reason: Core architectural component - translation routing logic
  > **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 4)
  - **Blocks**: Tasks 5-12
  - **Blocked By**: Task 2

  **References**:
  - `internal/selector/selector.go:ByText, ByID, ByXPath` - Selector building
  - `internal/action/action.go:TapAction, SwipeAction, SendKeysAction` - Action structures
  - `yaml/commands/control.go:Executor` - Existing command execution pattern

  **Acceptance Criteria**:
  - [ ] `translator.TranslateFlow` returns slice of Actions
  - [ ] Unknown commands return ErrUnsupportedCommand
  - [ ] Selector building uses correct ByText/ByID methods

  **QA Scenarios**:

  \`\`\`
  Scenario: Translate tapOn command to TapAction
    Tool: Bash
    Preconditions: translator.go exists with TranslateCommand
    Steps:
      1. Create tapOn command with text selector
      2. Run: go test ./maestro -run TestTranslateTap -v
    Expected Result: Returns TapAction with TextSelector matching "Login"
    Evidence: .sisyphus/evidence/task-3-tap-translate.{ext}

  Scenario: Unknown command returns error
    Tool: Bash
    Preconditions: translator.go exists
    Steps:
      1. Create unknown command type
      2. Run: go test ./maestro -run TestTranslateUnknown -v
    Expected Result: Returns ErrUnsupportedCommand
    Evidence: .sisyphus/evidence/task-3-unknown.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] Tap translate: task-3-tap-translate.log
  - [ ] Unknown command: task-3-unknown.log

  **Commit**: YES
  - Message: `feat(maestro): implement core translation engine`
  - Files: `maestro/translator.go`, `maestro/translator_test.go`

---

- [x] 4. Create cmd/maestro CLI skeleton

  **What to do**:
  - Create `cmd/maestro/main.go` CLI entry point
  - Implement command structure using subcommands:
    - `maestro test <file>` - Execute Maestro flow
    - `maestro validate <file>` - Validate YAML syntax
    - `maestro --help` - Show help
  - Use `cobra` or built-in flag parsing for CLI
  - Print version info with `--version`
  - Output format: `[STEP 1/5] tapOn: "Login"` to stdout
  - Exit codes: 0 success, 1 failure

  **Must NOT do**:
  - Don't import platform drivers (ios/, android/) directly
  - Don't implement execution logic in CLI - delegate to library

  **Recommended Agent Profile**:
  > **Category**: `quick`
  > - Reason: CLI scaffolding - straightforward structure
  > **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 3)
  - **Blocks**: Task 14
  - **Blocked By**: Task 1

  **References**:
  - Go standard `flag` package or `cobra` for CLI structure
  - Existing CLI patterns in Go (e.g., `kubectl`, `docker`)

  **Acceptance Criteria**:
  - [ ] `go build ./cmd/maestro` compiles
  - [ ] `./maestro --help` shows usage
  - [ ] `./maestro validate <file>` returns syntax status
  - [ ] `./maestro test <file>` attempts execution

  **QA Scenarios**:

  \`\`\`
  Scenario: CLI builds successfully
    Tool: Bash
    Preconditions: cmd/maestro/main.go exists
    Steps:
      1. Run: go build -o maestro ./cmd/maestro
      2. Run: ./maestro --version
    Expected Result: Binary exists, version prints
    Evidence: .sisyphus/evidence/task-4-build.{ext}

  Scenario: Help command works
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run: ./maestro --help
    Expected Result: Shows usage with test/validate subcommands
    Evidence: .sisyphus/evidence/task-4-help.{ext}

  Scenario: Validate command with valid file
    Tool: Bash
    Preconditions: Binary built, valid .maestro.yaml exists
    Steps:
      1. Run: ./maestro validate test.maestro.yaml
    Expected Result: Exit code 0, "Valid" output
    Evidence: .sisyphus/evidence/task-4-validate.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] Build output: task-4-build.log
  - [ ] Help output: task-4-help.log
  - [ ] Validate output: task-4-validate.log

  **Commit**: YES
  - Message: `feat(maestro): add CLI tool with test/validate subcommands`
  - Files: `cmd/maestro/main.go`

---

- [x] 5. tapOn command translator
- [x] 6. swipe command translator
- [x] 7. inputText command translator
- [x] 8. launchApp/killApp/stopApp translators
- [x] 9. assertVisible/assertNotVisible translators
- [x] 10. waitForAnimationToEnd/pressKey/back translators
- [x] 11. scroll/takeScreenshot translators
- [x] 12. runFlow subflow executor

---

- [ ] 13. Workspace config.yaml support

  **What to do**:
  - Create `maestro/config.go`
  - Implement config.yaml parsing:
    - `flows` - glob patterns for flow discovery (e.g., `**/*.maestro.yaml`)
    - `testOutputDir` - output directory for screenshots/logs
    - `includeTags` / `excludeTags` - flow filtering
    - `executionOrder` - sequential flow ordering
    - `platform.ios` / `platform.android` - platform-specific settings
  - Support config in project root or `.maestro/` directory

  **Must NOT do**:
  - Don't implement cloud config (notifications, baselineBranch)
  - Don't implement Maestro Studio integration

  **Recommended Agent Profile**:
  > **Category**: `unspecified-high`
  > - Reason: Config file parsing - straightforward
  > **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 14, 15)
  - **Blocks**: Task 14
  - **Blocked By**: Task 4

  **References**:
  - Maestro docs: workspace-configuration
  - Go `gopkg.in/yaml.v3` for YAML parsing (already in go.mod)

  **Acceptance Criteria**:
  - [ ] Parse valid config.yaml with flows pattern
  - [ ] Parse includeTags/excludeTags for filtering
  - [ ] Parse platform-specific settings
  - [ ] go test ./maestro -run TestConfig -v passes

  **QA Scenarios**:

  \`\`\`
  Scenario: Parse valid config.yaml
    Tool: Bash
    Preconditions: maestro/config.go exists
    Steps:
      1. Create config.yaml with flows, tags
      2. go test ./maestro -run TestConfig -v
    Expected Result: Config struct with Flows = ["**/*.maestro.yaml"], Tags = ["smoke"]
    Evidence: .sisyphus/evidence/task-13-config.{ext}

  Scenario: Parse platform-specific settings
    Tool: Bash
    Preconditions: maestro/config.go exists
    Steps:
      1. Create config.yaml with platform settings
      2. go test ./maestro -run TestConfigPlatform -v
    Expected Result: Config.Platform.iOS.DisableAnimations = true
    Evidence: .sisyphus/evidence/task-13-platform.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] Config parse: task-13-config.log
  - [ ] Platform: task-13-platform.log

  **Commit**: YES
  - Message: `feat(maestro): implement workspace config.yaml support`
  - Files: `maestro/config.go`, `maestro/config_test.go`

---

- [ ] 14. CLI integration (wire translator to executor)

  **What to do**:
  - Wire CLI `test` command to translator/executor
  - Implement flow execution loop:
    1. Parse Maestro YAML
    2. Translate to Actions
    3. Execute each action via device interface
    4. Print progress to stdout
    5. Capture screenshot on failure
  - Implement `validate` command:
    1. Parse YAML
    2. Return validation result without execution
  - Implement tag filtering from config
  - Support `--device` flag for device selection
  - Support `--output` flag for screenshot directory
  - Exit code: 0 = all passed, 1 = any failure

  **Must NOT do**:
  - Don't import ios/ or android/ drivers directly - use core.Device interface
  - Don't execute commands in parallel (Phase 1 is sequential)

  **Recommended Agent Profile**:
  > **Category**: `deep`
  > - Reason: Integration work - wiring all components together
  > **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 13, 15)
  - **Blocks**: Task 15
  - **Blocked By**: Task 4, Task 13

  **References**:
  - `core/device.go:Device` - Interface for execution
  - `cmd/maestro/main.go` - CLI structure from Task 4
  - `maestro/translator.go` - Translation from Task 3

  **Acceptance Criteria**:
  - [ ] `./maestro test flow.maestro.yaml` executes flow
  - [ ] `./maestro validate flow.maestro.yaml` returns syntax status
  - [ ] Exit code 0 on success, 1 on failure
  - [ ] Progress printed to stdout

  **QA Scenarios**:

  \`\`\`
  Scenario: Execute simple flow with tapOn
    Tool: Bash
    Preconditions: CLI wired, test flow exists
    Steps:
      1. Create simple flow with tapOn command
      2. Run: ./maestro test simple.maestro.yaml
    Expected Result: Exit code 0, "[1/1] tapOn: Login" printed
    Evidence: .sisyphus/evidence/task-14-execute.{ext}

  Scenario: Validate returns success for valid YAML
    Tool: Bash
    Preconditions: CLI wired
    Steps:
      1. Run: ./maestro validate valid.maestro.yaml
    Expected Result: Exit code 0, "Valid YAML" output
    Evidence: .sisyphus/evidence/task-14-validate.{ext}

  Scenario: Validate returns error for invalid YAML
    Tool: Bash
    Preconditions: CLI wired
    Steps:
      1. Run: ./maestro validate invalid.maestro.yaml
    Expected Result: Exit code 1, error message printed
    Evidence: .sisyphus/evidence/task-14-validate-error.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] Execute: task-14-execute.log
  - [ ] Validate success: task-14-validate.log
  - [ ] Validate error: task-14-validate-error.log

  **Commit**: YES
  - Message: `feat(maestro): integrate translator into CLI executor`
  - Files: `cmd/maestro/main.go` (updated), `maestro/executor.go`

---

- [ ] 15. End-to-end integration test

  **What to do**:
  - Create `maestro/integration_test.go` (build tag: `//go:build integration`)
  - Test complete flow:
    1. Create sample Maestro YAML flow with multiple commands
    2. Parse → Translate → Execute (mock device)
    3. Verify all commands executed in order
  - Test error recovery:
    1. Flow with invalid selector
    2. Verify error returned with command context
  - Test config.yaml integration:
    1. Load config with flows pattern
    2. Verify flow discovery works
  - Create sample flows in `examples/` directory

  **Must NOT do**:
  - Don't run against real devices in unit tests (use mocks)
  - Don't skip any verification step

  **Recommended Agent Profile**:
  > **Category**: `unspecified-high`
  > - Reason: Integration testing - straightforward end-to-end
  > **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO (final integration)
  - **Parallel Group**: Wave 3 (with Tasks 13, 14)
  - **Blocks**: Tasks F1-F4
  - **Blocked By**: Tasks 5-12, 14

  **References**:
  - `maestro/parser.go` - Parsing from Task 2
  - `maestro/translator.go` - Translation from Task 3
  - `internal/action/action.go` - Action types for mocking

  **Acceptance Criteria**:
  - [ ] `go test ./maestro -tags=integration -v` passes
  - [ ] Sample flow executes all 5 commands
  - [ ] Error recovery provides context
  - [ ] Example flows exist in `examples/`

  **QA Scenarios**:

  \`\`\`
  Scenario: End-to-end flow execution
    Tool: Bash
    Preconditions: Integration test exists
    Steps:
      1. go test ./maestro -tags=integration -run TestEndToEnd -v
    Expected Result: All commands executed, exit code 0
    Evidence: .sisyphus/evidence/task-15-e2e.{ext}

  Scenario: Error recovery with context
    Tool: Bash
    Preconditions: Integration test exists
    Steps:
      1. go test ./maestro -tags=integration -run TestErrorRecovery -v
    Expected Result: Error includes command name and selector
    Evidence: .sisyphus/evidence/task-15-error.{ext}
  \`\`\`

  **Evidence to Capture**:
  - [ ] E2E: task-15-e2e.log
  - [ ] Error: task-15-error.log

  **Commit**: YES
  - Message: `feat(maestro): add end-to-end integration tests`
  - Files: `maestro/integration_test.go`, `examples/sample.maestro.yaml`

---

## Final Verification Wave

> 4 review agents run in PARALLEL. ALL must APPROVE. Present consolidated results to user and get explicit "okay" before completing.

- [x] F1. **Plan Compliance Audit** — `oracle`
  Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, check function). For each "Must NOT Have": search codebase for forbidden patterns — reject with file:line if found. Check evidence files exist in .sisyphus/evidence/.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [x] F2. **Code Quality Review** — `unspecified-high`
  Run `go vet ./maestro/... && go fmt ./maestro/... && go build ./cmd/maestro && go test ./maestro/...`. Review all changed files for: `as any`/`@ts-ignore` equivalent (not applicable in Go), empty catches (no error handling), commented-out code, unused imports. Check AI slop: excessive comments, over-abstraction, generic names.
  Output: `Build [PASS/FAIL] | Vet [PASS/FAIL] | Tests [N pass/N fail] | VERDICT`

- [x] F3. **CLI Smoke Test** — `unspecified-high`
  Start from clean state. Execute:
  - `./maestro --help` → expected output
  - `./maestro validate examples/sample.maestro.yaml` → exit 0
  - `./maestro test examples/sample.maestro.yaml` → exit 0
  - `./maestro validate examples/invalid.yaml` → exit 1
  Save to `.sisyphus/evidence/final-qa/`.
  Output: `CLI [PASS/FAIL] | Scenarios [N/N pass] | VERDICT`

- [x] F4. **Scope Fidelity Check** — `deep`
  For each task: read "What to do", read actual diff (git log/diff). Verify 1:1 — everything in spec was built (no missing), nothing beyond spec was built (no creep). Check "Must NOT do" compliance. Detect cross-task contamination: Task N touching Task M's files. Flag unaccounted changes.
  Output: `Tasks [N/N compliant] | Contamination [CLEAN/N issues] | Unaccounted [CLEAN/N files] | VERDICT`

---

## Commit Strategy

### Wave 1 (Foundation)
- `feat(maestro): add Maestro package structure and types` — maestro/types.go, maestro/errors.go, maestro/maestro.go
- `feat(maestro): implement YAML parser for Maestro flows` — maestro/parser.go, maestro/parser_test.go
- `feat(maestro): implement core translation engine` — maestro/translator.go, maestro/translator_test.go
- `feat(maestro): add CLI tool with test/validate subcommands` — cmd/maestro/main.go

### Wave 2 (Commands)
- `feat(maestro): implement tapOn command translator` — maestro/commands/tap.go, maestro/commands/tap_test.go
- `feat(maestro): implement swipe command translator` — maestro/commands/swipe.go, maestro/commands/swipe_test.go
- `feat(maestro): implement inputText command translator` — maestro/commands/input.go, maestro/commands/input_test.go
- `feat(maestro): implement app lifecycle command translators` — maestro/commands/app.go, maestro/commands/app_test.go
- `feat(maestro): implement assertion command translators` — maestro/commands/assert.go, maestro/commands/assert_test.go
- `feat(maestro): implement navigation command translators` — maestro/commands/navigation.go, maestro/commands/navigation_test.go
- `feat(maestro): implement scroll and takeScreenshot translators` — maestro/commands/media.go, maestro/commands/media_test.go
- `feat(maestro): implement runFlow subflow executor` — maestro/commands/flow.go, maestro/commands/flow_test.go

### Wave 3 (Integration)
- `feat(maestro): implement workspace config.yaml support` — maestro/config.go, maestro/config_test.go
- `feat(maestro): integrate translator into CLI executor` — cmd/maestro/main.go (updated), maestro/executor.go
- `feat(maestro): add end-to-end integration tests` — maestro/integration_test.go, examples/sample.maestro.yaml

---

## Success Criteria

### Verification Commands
```bash
go build ./cmd/maestro           # Build CLI
./maestro --help                 # Shows CLI help
go test ./maestro/... -v         # All unit tests pass
go build ./maestro/...            # Package builds
go vet ./maestro/...             # No vet errors
```

### Final Checklist
- [ ] All 15 commands implemented with translators
- [ ] CLI with test/validate/help subcommands works
- [ ] Both Maestro YAML and Go-UOP native YAML supported
- [ ] Workspace config.yaml parsing works
- [ ] All tests pass
- [ ] No forbidden features implemented (relational selectors, cloud, etc.)
- [ ] Example flows created in examples/
- [ ] Documentation updated (README)
