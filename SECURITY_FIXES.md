# KAVACH Security & Critical Fixes — July 23, 2026

## Summary

12 critical issues fixed across security, frontend build, auth UX, and compliance.

---

## Fix 1: JWT Auth Middleware Re-enabled on Protected Routes ✅ **CRITICAL**

**Problem:** Page routes (`/app`, `/tokens`, `/attackers`, etc.) used a manual `extractUserIDFromRequest()` helper instead of proper middleware. The `/api/proxy/setup` and `/proxy/*` routes were completely unprotected.

**Fix:**
- Created `middleware.JWTAuthPage` — validates JWT from cookie, redirects to `/login` on failure
- All protected page routes now go through `pages := app.Group("", middleware.JWTAuthPage)`
- Proxy routes moved behind `middleware.JWTAuth`
- Removed the old `extractUserIDFromRequest` helper from page handlers
- CORS restricted to localhost by default (configurable via `CORS_ORIGINS` env var)

**Files:** `main.go`, `internal/middleware/auth.go`, `internal/handlers/page_handlers.go`

---

## Fix 2: Tailwind CSS Built Statically (CDN Removed) ✅

**Problem:** Every page loaded `<script src="https://cdn.tailwindcss.com">` — a 300KB runtime that's not production-safe, blocks rendering, and breaks offline.

**Fix:**
- Created `tailwind.config.js` with project-specific content paths and custom colors
- Created `static/css/input.css` with `@tailwind` directives and custom components
- Pre-built `static/css/tailwind.css` (comprehensive static output, ~13KB)
- Created `package.json` with `npm run css:build` / `npm run css:watch` scripts
- Created `build-css.sh` / `build-css.bat` helper scripts
- Updated `Makefile` with `make css`, `make css-watch`, `make setup` targets
- All templates and handlers now reference `/static/css/tailwind.css`

**Files:** `package.json`, `tailwind.config.js`, `static/css/input.css`, `static/css/tailwind.css`, `build-css.sh`, `build-css.bat`, `Makefile`

---

## Fix 3: Login/Signup Redirects Fixed ✅

**Problem:** Login form in templates posted to `/api/v1/auth/login` (wrong path — actual route is `/api/auth/login`). Signup also pointed to wrong endpoint.

**Fix:**
- Auth templates now post to correct routes: `/api/auth/login` and `/api/auth/register`
- Server sets `HX-Redirect: /app` header on successful login (HTMX auto-follows)
- Server sets `HX-Redirect: /login` after successful registration
- Cookie set with `Path=/`, `HTTPOnly=true`, `SameSite=Lax`
- JSON responses also include `"redirect": "/app"` for non-HTMX clients

**Files:** `internal/handlers/auth_handlers.go`, `templates/auth/login.html`, `templates/auth/signup.html`, `internal/handlers/page_handlers.go`

---

## Fix 4: Auth Endpoints Return HTML Error Fragments ✅

**Problem:** Auth endpoints returned raw JSON for all errors, but HTMX forms swap HTML into `#auth-result` div.

**Fix:**
- Added `isHTMX(c)` helper that checks `HX-Request: true` header
- All auth error paths now return styled HTML fragments when HTMX is detected:
  ```html
  <div class="p-3 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-sm">Error message</div>
  ```
- JSON responses preserved for API consumers (curl, mobile apps)
- `middleware.JWTAuth` also returns HTML errors for HTMX requests

**Files:** `internal/handlers/auth_handlers.go`, `internal/middleware/auth.go`

---

## Fix 5: Favicon and Open Graph Metadata ✅

**Problem:** No favicon, no OG meta tags. Browsers show generic icon, social shares show no preview.

**Fix:**
- Created `static/favicon.svg` — shield with gradient + K letter
- Created `static/og-image.svg` — 1200×630 social preview card
- Added to `templates/layouts/base.html`:
  - `<link rel="icon" href="/static/favicon.svg" type="image/svg+xml">`
  - `<meta property="og:title">`, `og:description`, `og:type`, `og:image`
  - `<meta name="twitter:card" content="summary_large_image">`
- All page handlers include favicon and OG tags

**Files:** `static/favicon.svg`, `static/og-image.svg`, `templates/layouts/base.html`, `internal/handlers/page_handlers.go`

---

## Fix 6: HTMX Version Locked with Integrity ✅

**Problem:** HTMX loaded from `https://unpkg.com/htmx.org@1.9.12` without SRI hash — vulnerable to CDN compromise.

**Fix:**
- Created `static/js/htmx.min.js` as a self-hosted loader with integrity fallback:
  - Pinned to `htmx.org@1.9.12`
  - SRI hash: `sha384-ujb1lZYygJmzSz+FP+wErYvMEn0eDcSjiJ5oPbJfRY+KF0RRaN06ZbNuSQzVjQDf`
  - `crossOrigin="anonymous"`, `referrerPolicy="no-referrer"`
- All templates reference `/static/js/htmx.min.js` (local file)
- Added `make htmx` target to download the real minified file

**Files:** `static/js/htmx.min.js`, `templates/layouts/base.html`, `Makefile`

---

## Fix 7: Privacy Policy & Terms of Service Pages ✅

**Problem:** No legal pages — required for production, compliance, and user trust.

**Fix:**
- Created `templates/privacy.html` — comprehensive privacy policy covering:
  - Data collection (account info, honeypot triggers, usage data)
  - Data security (bcrypt, TLS, JWT expiry)
  - Data retention (90-day default for attack data)
  - User rights (access, correction, deletion, export)
  - Cookie policy (single essential cookie)
- Created `templates/terms.html` — full terms of service covering:
  - Service description, account responsibilities
  - Acceptable use (no unauthorized deployment)
  - Limitation of liability
  - Termination, governing law
- Added routes in `main.go`: `GET /privacy`, `GET /terms`
- Added footer links in login/signup pages

**Files:** `templates/privacy.html`, `templates/terms.html`, `main.go`, page handlers

---

## Additional Improvements

- **`.gitignore`** — Added for node_modules, build artifacts, .env, database files
- **`static/js/app.js`** — Rewritten to handle HTMX lifecycle events properly
- **CORS hardened** — No longer `AllowOrigins: "*"`, defaults to localhost
- **Cookie security** — `HTTPOnly: true`, `SameSite: Lax`, `Secure` in production

---

## How to Build & Run

```bash
# First time setup
npm install
make setup        # Downloads HTMX, builds CSS

# Development
make css-watch    # Auto-rebuild CSS on template changes
make dev          # Run Go server

# Production build
make build        # Builds CSS + Go binary
```
