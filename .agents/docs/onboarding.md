# Agent layout

Shared, agent-neutral setup. `AGENTS.md` is the index every agent should follow.

## Tree

- `AGENTS.md` — canonical index: one-liner, commands, layout, doc links, plans convention.
- `CLAUDE.md` — `@AGENTS.md` pointer. Claude Code reads `CLAUDE.md`, not `AGENTS.md`.
- `.agents/docs/` — agent-facing docs. Package docs stay in `GO-Manager/docs/`, `Translator/docs/`, `web-ui/docs/` and are linked from `AGENTS.md`, not moved.
- `.agents/skills/` — project skills (none yet). Cursor reads this dir natively. Claude: one relative symlink per skill, `.claude/skills/<name>` → `../../.agents/skills/<name>`.
- `.agents/plans/` — working plans, gitignored. Do not commit.

## Cursor

Reads root `AGENTS.md` (and a nested `AGENTS.md` when working in that subtree). No extra entry file. Skills: `.agents/skills/`. New plans: `.agents/plans/`. A leftover local plan may still sit under `.cursor/plans/`; do not add more there.

## Claude Code

Keep `@AGENTS.md` as the first line of `CLAUDE.md`. Claude-only notes may follow that import. Skills via the symlinks above.

## Any other agent

Point it at `AGENTS.md` via that agent's entry-file mechanism, or re-run `/agentify-project` to research and wire it.

## Windows symlinks

Committed skill links are relative git symlinks. Enable Developer Mode and `git config core.symlinks true` (or clone with symlink support). Without that, links become text files and skills do not load. The `@AGENTS.md` import in `CLAUDE.md` is not a symlink.
