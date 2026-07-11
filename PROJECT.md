# Resuming — ATS Resume Checker

**Resuming** is an AI-powered ATS (Applicant Tracking System) resume checker that helps job seekers optimize their resumes for better visibility and higher match rates. The platform uses intelligent analysis to evaluate resumes against job descriptions, providing actionable feedback rather than opaque filtering.

---

## The Problem It Solves

The modern job market forces job seekers to use fragmented, expensive tools to check their resumes against ATS systems, while many ATS tools unfairly reject qualified candidates through rigid keyword/filtering rules. Resuming addresses this by:

- providing a dedicated, low-cost ATS checker that explains *why* a resume scores the way it does, rather than silently filtering candidates out;
- offering contextual and semantic analysis instead of outdated keyword matching, with constructive feedback;
- being designed as a cross-platform solution with planned native apps for mobile and desktop.

---

## Project Objectives

1. **A fair and intelligent ATS checker** — contextual/semantic analysis instead of outdated keyword matching, with constructive feedback.
2. **A simple, intuitive experience** — a plug-and-play UI usable by anyone, regardless of technical or design background.
3. **Affordable and accessible** — lowering financial barriers so students, graduates, career changers, and others can build professional materials.
4. **Cross-platform reach** — web-first with planned native apps for Android, iOS, Linux, macOS, and Windows.
5. **Extensible architecture** — a modular monorepo designed to grow into API and AI services.

---

## Architecture Overview

Resuming is a **Turborepo monorepo** composed of independent workspaces:

```
resuming/
├── apps/
│   ├── web/              # Frontend application (Astro + SSR)
│   ├── ai/               # AI/ML service (planned)
│   ├── api/              # Backend API (planned)
│   ├── android/          # Native Android app (planned)
│   ├── ios/              # Native iOS app (planned)
│   ├── linux/            # Linux desktop app (planned)
│   ├── macos/            # macOS desktop app (planned)
│   └── windows/          # Windows desktop app (planned)
└── packages/
    ├── config/           # Shared TypeScript configuration
    └── env/              # Shared environment variable validation
```

- The **web app** serves as the primary frontend, handling resume upload, job description input, and displaying ATS analysis results.
- Future **AI** and **API** services will provide backend processing for resume parsing, scoring, and semantic analysis.
- **Shared packages** (`config`, `env`) provide consistent tooling and validated configuration across all apps.

---

## Technology Stack

| Layer | Technology |
| --- | --- |
| Monorepo tool | **Turborepo** |
| Package manager | **Bun** |
| Language | **TypeScript** (strict mode) |
| Web framework | **Astro** v7 (SSR with Node adapter) |
| CSS framework | **TailwindCSS** v4 |
| Linting | **Oxlint** via Ultracite |
| Formatting | **Oxfmt** via Ultracite |
| Git hooks | **Lefthook** |
| Error tracking | **Sentry** |
| Env validation | **@t3-oss/env-core** with Zod |

---

## Repository Layout

```
resuming/
├── apps/
│   └── web/                  # Astro frontend application
│       ├── src/
│       │   ├── components/   # Reusable UI components
│       │   ├── layouts/      # Page layouts
│       │   ├── pages/        # Route pages
│       │   └── styles/       # Global styles
│       ├── astro.config.mjs
│       ├── package.json
│       └── tsconfig.json
├── packages/
│   ├── config/               # Shared TypeScript base config
│   └── env/                  # Environment variable schemas (Zod)
├── AGENTS.md                 # AI agent coding standards
├── PROJECT.md                # Project documentation
├── README.md                 # Quick start guide
├── lefthook.yml              # Pre-commit hooks
├── oxfmt.config.ts           # Formatter configuration
├── oxlint.config.ts          # Linter configuration
├── package.json              # Root workspace config
├── tsconfig.json             # Root TypeScript config
└── turbo.json                # Turborepo pipeline config
```
