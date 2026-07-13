# Resuming — Portfolio Builder with Integrated ATS Checker and Resume Generator

**Resuming** is an integrated web-based platform that combines three essential career tools into a single ecosystem: a professional portfolio builder, an AI-powered ATS (Applicant Tracking System) checker, and an automated resume generator. The platform helps job seekers create, optimize, and manage their professional application materials efficiently, with seamless data synchronization across all components.

---

## The Problem It Solves

The modern job market forces job seekers to use multiple fragmented, expensive tools for portfolios, resumes, and ATS checking, while many ATS systems unfairly reject qualified candidates through rigid keyword filtering. Resuming addresses this by:

- consolidating portfolio building, resume generation, and ATS checking into a single unified platform, eliminating redundant data entry across disconnected services;
- using a portfolio as the single source of truth — resumes are auto-populated from portfolio data, ensuring consistency and saving time;
- providing an AI-powered ATS checker that performs contextual and semantic analysis instead of outdated keyword matching, with actionable feedback;
- being completely free and accessible to all job seekers, removing financial barriers to professional career tools.

---

## Project Objectives

1. **A fully integrated career toolkit** — combine portfolio builder, resume generator, and ATS checker into one seamless ecosystem where data flows naturally between components.
2. **A simple and intuitive user experience** — a plug-and-play interface usable by anyone regardless of technical or design background.
3. **Affordable and accessible** — offered at no cost to democratize access to high-quality career tools for students, graduates, career changers, and others.
4. **A fair and intelligent ATS checker** — use AI for contextual understanding and semantic analysis instead of rigid keyword matching, providing constructive feedback.
5. **Seamless data synchronization** — the portfolio serves as the single source of truth, with resumes and other materials automatically pulling from this centralized data.

---

## System Features

### Admin Features

- **User Account Management** — manage, ban, and delete user accounts
- **Content and Template Management** — upload, update, and remove portfolio/resume templates; moderate user-generated content
- **Announcements** — post announcements with configurable end dates

### User Features

- **Portfolio Management** — build and customize professional online portfolios with projects, skills, contact info, and a personalized URL; choose from multiple templates; control privacy with optional password protection
- **Resume Generation** — automatically populate resumes from portfolio data; export to PDF; create multiple tailored versions for different job applications
- **ATS Optimization** — analyze resume content and structure; check for ATS formatting compatibility; receive a compliance score and detailed report; get real-time suggestions while editing
- **General and Account Settings** — secure authentication (register/login); dark/light mode theme; user dashboard for managing all content

---

## Architecture Overview

Resuming is a **Turborepo monorepo** composed of independent workspaces:

```
resuming/
├── apps/
│   ├── web/              # Frontend application (Astro + Svelte)
│   ├── ai-ml/            # AI/ML service (Python)
│   └── api/              # Backend API (planned)
└── packages/
    ├── config/           # Shared TypeScript configuration
    └── env/              # Shared environment variable validation
```

- The **web app** serves as the primary frontend, handling portfolio creation, resume generation, and displaying ATS analysis results.
- The **ai-ml** workspace provides Python-based AI/ML processing for resume parsing, scoring, and semantic analysis.
- **Shared packages** (`config`, `env`) provide consistent TypeScript tooling and validated configuration for the web and future API workspaces.

The portfolio acts as the single source of truth — all user data is entered once in the portfolio and automatically propagated to resume generation and ATS analysis, eliminating redundant data entry.

---

## Technology Stack

| Layer | Technology |
| --- | --- |
| Monorepo tool | **Turborepo** |
| Package manager | **Bun** |
| Language | **TypeScript** (strict mode) |
| Web framework | **Astro** v7 (SSR with Cloudflare adapter) |
| UI framework | **Svelte** v5 |
| CSS framework | **TailwindCSS** v4 |
| Linting | **Oxlint** via Ultracite |
| Formatting | **Oxfmt** via Ultracite |
| Git hooks | **Lefthook** |
| Error tracking | **Sentry** |
| Deployment | **Cloudflare** (via Wrangler) |
| Env validation | **Zod** (via `@t3-oss/env-core`) |
| AI/ML | **Python** (apps/ai-ml) |

---

## Repository Layout

```
resuming/
├── apps/
│   ├── web/                  # Astro + Svelte frontend application
│   │   ├── src/
│   │   │   ├── components/   # Reusable UI components
│   │   │   ├── layouts/      # Page layouts
│   │   │   ├── pages/        # Route pages
│   │   │   └── styles/       # Global styles
│   │   ├── astro.config.mjs
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── ai-ml/                # Python AI/ML service
│   └── api/                  # Backend API (planned)
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
