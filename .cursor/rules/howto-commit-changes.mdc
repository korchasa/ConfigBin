---
description: How to commit changes
globs:
alwaysApply: false
---
## HOW TO COMMIT CHANGES

- Code must be checked with `./run check` before any commit
- The AI assistant must not decide which files or changes to include or exclude from a commit: all verified non-PII changes in the working tree must be staged together for the current commit, and only when potential PII is present (secrets, tokens, personal data, etc.) should it pause and ask the user for explicit confirmation before proceeding.
- Commit messages must fully comply with the 'Conventional Commits' v1.0.0 specification with strict profile, including the definition of breaking changes.
- Package updates and the addition of new ones should be done in a separate commit before the main one.
- Commit messages must be only in english.
- Use git commands only with `GIT_PAGER=cat` env variable. For example, `GIT_PAGER=cat git diff`.

Conventional Commits 1.0.0 with strict profiles:
```
Commit messages must follow **Conventional Commits 1.0.0 (Strict Profile)** rules:

### Structure
```
<type>(<scope>)!: <description>
[blank line if body/footers]
<body>
[blank line if footers]
<footers>
```
- UTF-8, LF line breaks only.
- No trailing spaces.
- Header ≤100 chars, description ≤72 chars.

### Allowed Types
`feat`, `fix`, `perf`, `refactor`, `docs`, `style`, `test`,
`build`, `ci`, `chore`, `revert`, `prompts` (no others allowed).

### Scope
Optional, in `()` after type. Must match regex like `core`, `ui/header`, `parser.json`.

### Breaking Changes
- Marked with `!` in header **or** `BREAKING CHANGE:` footer.
- Must explain what broke and how to migrate.

### Footers
Format: `Token: Value`. Allowed tokens:
`Closes`, `Fixes`, `Refs`, `Co-authored-by`, `Signed-off-by`,
`Reviewed-by`, `BREAKING CHANGE`, `X-*` (custom).
- Continuations start with a space.
- `revert` commits must include `Reverts: <sha>` in body.

### Versioning
- Breaking change → MAJOR bump.
- `feat` → MINOR bump.
- `fix` → PATCH bump.
- Highest rule applies.

### Valid Examples
```
feat(api)!: remove deprecated v1 endpoints

BREAKING CHANGE: v1 endpoints removed. Use /v2/*.
Refs: #120, org/payments#45
```

```
fix(auth): correctly refresh tokens on 401

Closes: #512
```

```
revert: feat(api): add tokens cache

Reverts: 1a2b3c4d...
```

### Invalid Examples
- `Feat: add button` → type not lowercase.
- `fix : typo` → space before colon.
- `feat(ui/header)!` → missing description.
- `chore: update deps` + `BREAKING CHANGE: updated react` → vague.
- `feat(api,db): ...` → multiple scopes not allowed.
- `docs: fix typo.` → description ends with period.
- `fix: ...` + `References: #1` → invalid footer token.

### Style
- Imperative mood (“add”, “fix”).
- Lowercase unless proper noun.
- Body lines wrap at 72 chars.
- Scope in kebab-case or path format.

### Compliance
1. Header matches regex.
2. Length limits met.
3. Footers valid.
4. Breaking change descriptions meaningful.
5. `revert` includes `Reverts: <sha>`.
```
