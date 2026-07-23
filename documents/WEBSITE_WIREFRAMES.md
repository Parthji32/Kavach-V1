# WEBSITE WIREFRAMES: KAVACH v1

**Date:** July 22, 2026  
**Design System:** Purple dark theme (#7C3AED), glassmorphism, Tailwind CSS  
**Breakpoints:** Desktop (1024px+), Tablet (768px-1023px), Mobile (<768px)

---

## PAGE 1: HOMEPAGE (UPDATED)

### Desktop Layout (1024px+)

```
┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  [KAVACH Logo]  Products  How It Works  Use Cases  Pricing  │
│                            Docs  [Start Free]  [Login]      │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│         HERO SECTION                                          │
│         ═════════════════════════════════════════             │
│                                                               │
│         Catch Attackers the Moment They Move                 │
│         ═══════════════════════════════════════               │
│                                                               │
│         Deploy invisible honeypots across your infrastructure.│
│         The instant an attacker touches one — you know who    │
│         they are, where they came from, exactly what they try │
│                                                               │
│         [Start Free Trial]  [Watch 2-Min Demo]               │
│                                                               │
│                                                               │
│         (Background: animated orbs + gradient)                │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  3 CORE MESSAGES (3-column grid)                             │
│  ═══════════════════════════════════════════════════════════ │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │ 🚫 FIREWALL  │  │ 100% SIGNAL  │  │ ⚡ 5-MINUTE  │       │
│  │ CAN'T CATCH  │  │ ZERO NOISE   │  │ DEPLOY       │       │
│  │ WHAT'S      │  │              │  │              │       │
│  │ INSIDE      │  │ Every alert  │  │ Docker +     │       │
│  │              │  │ = confirmed  │  │ dashboard    │       │
│  │ Attackers    │  │ threat       │  │ ready in 5   │       │
│  │ bypass edge  │  │              │  │ minutes      │       │
│  │ defenses     │  │ No false     │  │              │       │
│  │ & move       │  │ positives    │  │ No agents,   │       │
│  │ undetected   │  │              │  │ no headache  │       │
│  │ (avg 197 days)                    │              │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  QUICK "HOW IT WORKS" PREVIEW (4 steps, animated icons)     │
│  ═══════════════════════════════════════════════════════════ │
│                                                               │
│  1. DEPLOY          2. SCATTER        3. DETECT      4. ALERT │
│  ───────────────    ─────────────     ─────────────  ──────── │
│  Docker +           Generate tokens   Attacker      Slack /   │
│  config             in dashboard      touches token Webhook   │
│                                                               │
│  [Learn More →]                                              │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  SOCIAL PROOF                                                │
│  ═══════════════════════════════════════════════════════════ │
│                                                               │
│  [Logo 1]  [Logo 2]  [Logo 3]  [Logo 4]  [Logo 5]           │
│                                                               │
│  "KAVACH reduced our MTTD from 6 months to 2 days."          │
│  — Sarah Chen, CISO at TechCorp                              │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  FINAL CTA                                                   │
│  ═══════════════════════════════════════════════════════════ │
│                                                               │
│  Ready to catch attackers?                                   │
│                                                               │
│  [Start Free Trial - 7 Days, 5 Tokens, No CC]                │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Footer: Links | Docs | Blog | Contact | Social Icons        │
└─────────────────────────────────────────────────────────────┘
```

### Mobile Layout (<768px)

```
┌────────────────────────────────┐
│ ☰ [KAVACH Logo]      [Login]   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Catch Attackers the Moment     │
│ They Move                       │
│                                │
│ Deploy invisible honeypots...  │
│                                │
│ [Start Free Trial]             │
│ [Watch Demo]                   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ 3 Core Messages (stacked)      │
│ ┌──────────────────────────┐   │
│ │ Firewall Can't Catch     │   │
│ │ What's Inside            │   │
│ └──────────────────────────┘   │
│ ┌──────────────────────────┐   │
│ │ 100% Signal, Zero Noise  │   │
│ └──────────────────────────┘   │
│ ┌──────────────────────────┐   │
│ │ 5-Minute Deploy          │   │
│ └──────────────────────────┘   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ How It Works (scrollable)       │
│ 1. DEPLOY → 2. SCATTER →       │
│ 3. DETECT → 4. ALERT           │
│ [Learn More]                   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Social Proof                   │
│ (Logos stacked)                │
│ Logo 1, Logo 2, Logo 3...      │
│                                │
│ Testimonial text...            │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Footer (compact)               │
└────────────────────────────────┘
```

---

## PAGE 2: HOW IT WORKS (/how-it-works)

### Desktop Layout

```
┌─────────────────────────────────────────────────────────────┐
│ [KAVACH Logo]  Products  How It Works  Use Cases  Pricing   │
│                            Docs  [Start Free]  [Login]      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  How It Works                                                │
│  ═════════════════════════════════════════════════════════   │
│  From deployment to detection in 4 simple steps              │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  STEP 1: DEPLOY                                              │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  [Docker Container Screenshot]        Run KAVACH in Docker   │
│                                        Takes 5 minutes        │
│                                        Works on any          │
│                                        infrastructure         │
│                                                               │
│  $ docker run -p 3000:3000 kavach    Command line setup      │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  STEP 2: SCATTER                                             │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│                         Generate tokens in dashboard         │
│  [Dashboard Screenshot]  Pick token type:                    │
│                          • API Key (sk_...)                  │
│                          • URL (https://...)                 │
│                          • Document (word.docx)              │
│                          • DNS record                        │
│                          • Email credential                  │
│                                                               │
│                          Place them anywhere:                │
│                          • Git repos                         │
│                          • Email                             │
│                          • Slack channels                    │
│                          • Config files                      │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  STEP 3: DETECT                                              │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  Attacker finds token            KAVACH detects access       │
│  (reconnaissance, breach)        Classifies threat           │
│           ↓                                ↓                 │
│  [Network Diagram]                [Classifier Output]        │
│                                                               │
│  Generates fingerprint:                                      │
│  • IP address & geolocation                                  │
│  • Browser & OS info                                         │
│  • Device fingerprint                                        │
│  • Request headers                                           │
│  • Behavior anomalies                                        │
│  • Risk score (0-100)                                        │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  STEP 4: ALERT                                               │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  Real-time notification sent to your team:                  │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │   WEBHOOK    │  │    SLACK     │  │    EMAIL     │       │
│  │              │  │              │  │              │       │
│  │ POST to your │  │ @security    │  │ [alert@...]  │       │
│  │ infrastructure                 │                │       │
│  │ Full attacker│  │ Attacker:    │  │ Subject:     │       │
│  │ payload in   │  │ 192.168.1.5  │  │ Honeypot     │       │
│  │ JSON         │  │ Risk: 95%    │  │ Alert: Critical       │
│  │              │  │              │  │              │       │
│  │ Custom       │  │ Action:      │  │ Attacker IP: │       │
│  │ integration  │  │ Block now?   │  │ [Details]    │       │
│  │ ready        │  │              │  │              │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│                                                               │
│  All alerts → Dashboard for analysis                         │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  ARCHITECTURE OVERVIEW                                       │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  Internet → [Your Network]                                  │
│                    ↓                                         │
│              KAVACH Container                                │
│              ┌──────────────────┐                           │
│              │ • Token Manager   │ → Generate honeypots      │
│              │ • Classifier      │ → Score threats           │
│              │ • Alert Engine    │ → Send notifications      │
│              │ • Dashboard       │ → Visualize attacks       │
│              └──────────────────┘                           │
│                    ↓                                         │
│              SQLite Database                                 │
│              (runs locally, zero external calls)             │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ [Start Free Trial]                                           │
└─────────────────────────────────────────────────────────────┘
```

### Mobile Layout

```
┌────────────────────────────────┐
│ ☰ [Logo]              [Login]   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ How It Works                   │
│ 4 simple steps to catch attacks │
└────────────────────────────────┘

┌────────────────────────────────┐
│ STEP 1: DEPLOY                 │
│ [Screenshot]                   │
│ Docker container, 5 minutes    │
│ [Learn more ▼]                 │
└────────────────────────────────┘

┌────────────────────────────────┐
│ STEP 2: SCATTER                │
│ [Screenshot]                   │
│ Generate tokens, place them    │
│ [Learn more ▼]                 │
└────────────────────────────────┘

┌────────────────────────────────┐
│ STEP 3: DETECT                 │
│ [Screenshot]                   │
│ Automatic threat classification│
│ [Learn more ▼]                 │
└────────────────────────────────┘

┌────────────────────────────────┐
│ STEP 4: ALERT                  │
│ [Screenshot]                   │
│ Instant notifications          │
│ [Learn more ▼]                 │
└────────────────────────────────┘

┌────────────────────────────────┐
│ [Start Free Trial]             │
└────────────────────────────────┘
```

---

## PAGE 3: PRICING (/pricing)

### Desktop Layout

```
┌─────────────────────────────────────────────────────────────┐
│ [KAVACH Logo]  Products  How It Works  Use Cases  Pricing   │
│                            Docs  [Start Free]  [Login]      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  Pricing                                                     │
│  ═════════════════════════════════════════════════════════   │
│  Choose the plan that scales with your security needs        │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  [Annual Billing] [Monthly Billing]  (10% discount annually)│
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  ┌──────────────────┐  ┌──────────────────┐  ┌───────────── │
│  │ STARTER          │  │ PROFESSIONAL     │  │ ENTERPRISE   │
│  │                  │  │                  │  │              │
│  │ $2,000/month     │  │ $5,000/month     │  │ Custom       │
│  │ ($24K/year)      │  │ ($60K/year)      │  │ (Call sales) │
│  │                  │  │ ⭐ Most Popular  │  │              │
│  │                  │  │                  │  │              │
│  ├──────────────────┤  ├──────────────────┤  ├───────────── │
│  │ Features:        │  │ Everything in    │  │ Everything   │
│  │ • 5 tokens       │  │ Starter, plus:   │  │              │
│  │ • 1 user         │  │ • Unlimited      │  │ Plus:        │
│  │ • Basic alerts   │  │   tokens         │  │ • Dedicated  │
│  │   (webhook only) │  │ • 3 users        │  │   support    │
│  │ • Dashboard      │  │ • Webhook        │  │ • White-label│
│  │ • 7-day history  │  │ • Slack          │  │   dashboard  │
│  │                  │  │ • Email alerts   │  │ • SLA 99.9%  │
│  │                  │  │ • 90-day history │  │ • Custom     │
│  │                  │  │ • Priority       │  │   integrations
│  │                  │  │   support        │  │              │
│  │                  │  │                  │  │              │
│  │ [Start Free]     │  │ [Start Free]     │  │ [Talk to     │
│  │ [Buy]            │  │ [Buy]            │  │  Sales]      │
│  └──────────────────┘  └──────────────────┘  └───────────── │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  COMPARISON TABLE                                            │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  Feature                    │ Starter │ Pro │ Enterprise    │
│  ─────────────────────────────────────────────────────────   │
│  Token Types                │ All 5   │ All │ All           │
│  Users                       │ 1       │ 3   │ Unlimited     │
│  API Rate Limit              │ 1K/day  │ 10K │ Unlimited     │
│  Integrations                │ Webhook │ All │ Custom        │
│  Response Time (SLA)         │ Best    │ 99% │ 99.9%         │
│  Support                     │ Chat    │ Pri │ Dedicated     │
│  Deployment                  │ Cloud   │ Cloud│Self/Cloud    │
│  Training                    │ -       │ -   │ Included      │
│  Custom Development          │ -       │ -   │ Available     │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  ROI CALCULATOR                                              │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  Average cost of a data breach: $4.29M                       │
│  Average time to detect: 197 days                            │
│  KAVACH detection time: Same day (24 hours)                  │
│                                                               │
│  With KAVACH:                                                │
│  ┌─────────────────────────────────────────────┐             │
│  │ Early detection = $2M+ in prevented costs    │             │
│  │ Year 1 ROI with Starter plan: 8,500%        │             │
│  │ Payback period: < 1 week                    │             │
│  └─────────────────────────────────────────────┘             │
│                                                               │
│  (Based on Ponemon Institute 2024 Breach Report)             │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  FAQ                                                         │
│  ═════════════════════════════════════════════════════════   │
│  Q: Can I upgrade/downgrade anytime?                         │
│  A: Yes, change your plan any time. Prorated billing.        │
│                                                               │
│  Q: What if I need more tokens?                              │
│  A: Add-on packs: 10 tokens = $100/month                     │
│                                                               │
│  Q: Do you offer discounts for annual billing?               │
│  A: Yes, 10% discount when you prepay annually.              │
│                                                               │
│  Q: Is there a free trial?                                   │
│  A: Yes, 7-day free trial. 5 tokens, full access.            │
│     No credit card required.                                 │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  Still have questions?                                       │
│  [Talk to Sales] [Read Docs] [Contact Support]               │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Mobile Layout

```
┌────────────────────────────────┐
│ ☰ [Logo]              [Login]   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Pricing                        │
│ Simple, transparent plans      │
└────────────────────────────────┘

┌────────────────────────────────┐
│ STARTER   PROFESSIONAL  ENTERPRISE
│ $2K/mo    $5K/mo       Custom
│           ⭐            │
│ [Start]   [Start] [Talk]        │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Starter Includes:              │
│ • 5 tokens                     │
│ • 1 user                       │
│ • Webhook alerts               │
│ • Dashboard                    │
│ • 7-day history                │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Professional Includes:         │
│ • Unlimited tokens             │
│ • 3 users                      │
│ • Webhook + Slack + Email      │
│ • 90-day history               │
│ • Priority support             │
│                                │
│ (Swipe for more)              │
└────────────────────────────────┘

┌────────────────────────────────┐
│ ROI Calculator                 │
│ Avg breach cost: $4.29M        │
│ KAVACH: Detects in 24 hours    │
│ Savings: $2M+                  │
│ ROI: 8,500% Year 1             │
└────────────────────────────────┘

┌────────────────────────────────┐
│ [Start Free Trial]             │
└────────────────────────────────┘
```

---

## PAGE 4: USE CASES (/use-cases)

### Desktop Layout

```
┌─────────────────────────────────────────────────────────────┐
│ [KAVACH Logo]  Products  How It Works  Use Cases  Pricing   │
│                            Docs  [Start Free]  [Login]      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  Use Cases                                                   │
│  ═════════════════════════════════════════════════════════   │
│  See how KAVACH catches attackers in real scenarios          │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  USE CASE 1: LATERAL MOVEMENT DETECTION                      │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  The Problem:                     How KAVACH Solves It:      │
│  Attacker breaches edge           Honeypot tokens            │
│  security, moves inside →         scatter across network      │
│  undetected for months            Attacker explores →        │
│  (avg 197 days)                   finds honeypot → CAUGHT    │
│                                                               │
│  [Network diagram]                Detection Time: <1 hour    │
│                                   Status: ✅ Threat          │
│                                            Eliminated        │
│                                                               │
│  Industries Most at Risk: Financial, Healthcare, Tech        │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  USE CASE 2: INSIDER THREAT DETECTION                        │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  The Problem:                     How KAVACH Solves It:      │
│  Disgruntled employee             Place honeypot creds       │
│  accesses systems they            in "sensitive" areas       │
│  shouldn't → IT can't              (git repos, shared drives) │
│  distinguish from normal use       Insider accesses →        │
│                                    IMMEDIATELY flagged        │
│                                                               │
│  [Employee scenario diagram]      Detection: Instant         │
│                                   Action: Notify SOC,        │
│                                           Lock account       │
│                                                               │
│  Industries Most at Risk: Financial, Defense, Tech           │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  USE CASE 3: CREDENTIAL STUFFING ATTACKS                     │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  The Problem:                     How KAVACH Solves It:      │
│  Attacker uses stolen              Deploy honeypot API keys  │
│  credentials in mass attack        to mock endpoints         │
│  across services → generates       Attacker mass-tests →     │
│  alert fatigue (false positives)   honeypots fire (100%      │
│                                    confirmed threat)         │
│                                                               │
│  [Attack flow diagram]             Confidence: 100%          │
│                                    Risk: Critical            │
│                                    Action: Rotate creds,     │
│                                            Force MFA          │
│                                                               │
│  Industries Most at Risk: SaaS, Retail, Tech                 │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  USE CASE 4: SUPPLY CHAIN ATTACK DETECTION                   │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  The Problem:                     How KAVACH Solves It:      │
│  3rd-party vendor compromised,     Deploy honeypot creds     │
│  attacker accesses your            in vendor integration     │
│  systems through integration       points & data pipelines   │
│  → months to detect                Vendor compromise         │
│                                    detected same day         │
│                                                               │
│  [Supply chain diagram]            Early Warning: 6 months   │
│                                    ahead of detection        │
│                                    Impact: Prevented          │
│                                            lateral move      │
│                                                               │
│  Industries Most at Risk: Finance, Healthcare, Tech          │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  USE CASE 5: COMPLIANCE VALIDATION                           │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  The Problem:                     How KAVACH Solves It:      │
│  Auditors ask: "Can you prove     KAVACH provides:          │
│  you'd detect a breach early?"     • Audit trails of all     │
│  No single tool proves this        honeypot accesses        │
│                                    • Incident response       │
│                                      playbook (automated)    │
│                                    • Risk scoring &          │
│                                      prioritization          │
│                                                               │
│  [Compliance checklist]            Status: ✅ Audit Pass      │
│                                    Frameworks:               │
│                                    SOC 2, ISO 27001,        │
│                                    HIPAA, PCI-DSS           │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  By the Numbers                                              │
│  ═════════════════════════════════════════════════════════   │
│                                                               │
│  197 days    → Average time to detect breach                 │
│  <1 hour     → KAVACH detection time                         │
│  $4.29M      → Average breach cost                           │
│  $2M+        → Potential savings with KAVACH                 │
│  99.8%       → KAVACH detection accuracy                     │
│                                                               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  Ready to apply KAVACH to your security strategy?            │
│  [Start Free Trial] [Book Demo] [Contact Sales]              │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Mobile Layout (Accordion/Expandable)

```
┌────────────────────────────────┐
│ ☰ [Logo]              [Login]   │
└────────────────────────────────┘

┌────────────────────────────────┐
│ Use Cases                      │
│ 5 real scenarios where KAVACH  │
│ catches attackers             │
└────────────────────────────────┘

┌────────────────────────────────┐
│ ▶ Lateral Movement Detection   │
│   [Tap to expand]              │
└────────────────────────────────┘

┌────────────────────────────────┐
│ ▶ Insider Threat Detection     │
│   [Tap to expand]              │
└────────────────────────────────┘

┌────────────────────────────────┐
│ ▶ Credential Stuffing          │
│   [Tap to expand]              │
└────────────────────────────────┘

┌────────────────────────────────┐
│ ▶ Supply Chain Attacks         │
│   [Tap to expand]              │
└────────────────────────────────┘

┌────────────────────────────────┐
│ ▶ Compliance Validation        │
│   [Tap to expand]              │
└────────────────────────────────┘

┌────────────────────────────────┐
│ By the Numbers:                │
│ 197 days → Avg detection time  │
│ <1 hour → KAVACH time          │
│ $4.29M → Avg breach cost       │
│ $2M+ → Potential savings       │
└────────────────────────────────┘

┌────────────────────────────────┐
│ [Start Free Trial]             │
└────────────────────────────────┘
```

---

## PAGE 5: LOGIN (/login)

### Desktop Layout

```
┌─────────────────────────────────────────────────────────────┐
│ [KAVACH Logo]                                               │
│                                                               │
│                         Don't have an account?               │
│                         [Sign Up Free]                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                                                          │ │
│  │          Welcome Back                                   │ │
│  │          ═════════════════════════════════════════      │ │
│  │          Sign in to your KAVACH dashboard              │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────┐   │ │
│  │  │ Email Address                                  │   │ │
│  │  │ [_____________________________________]       │   │ │
│  │  └────────────────────────────────────────────────┘   │ │
│  │                                                          │ │
│  │  ┌────────────────────────────────────────────────┐   │ │
│  │  │ Password                                       │   │ │
│  │  │ [_____________________________________]       │   │ │
│  │  │                          [Show]                │   │ │
│  │  └────────────────────────────────────────────────┘   │ │
│  │                                                          │ │
│  │  ☐ Remember me for 30 days                           │ │
│  │                                                          │ │
│  │  [Sign In]                                             │ │
│  │                                                          │ │
│  │  ─────────────────────────────────────────────────    │ │
│  │                                                          │ │
│  │  [Forgot Password?] | [Contact Support]                 │ │
│  │                                                          │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                               │
│  Background: Subtle animated orbs                            │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Mobile Layout

```
┌────────────────────────────────┐
│ [KAVACH Logo]                  │
│                                │
│ Don't have an account?         │
│ [Sign Up Free]                 │
└────────────────────────────────┘

┌────────────────────────────────┐
│                                │
│ Welcome Back                   │
│ Sign in to your dashboard      │
│                                │
│ Email Address                  │
│ [________________]             │
│                                │
│ Password                       │
│ [________________]  [Show]     │
│                                │
│ ☐ Remember me                  │
│                                │
│ [Sign In]                      │
│                                │
│ ──────────────────────────     │
│                                │
│ [Forgot Password?]             │
│ [Contact Support]              │
│                                │
└────────────────────────────────┘
```

---

## DESIGN SYSTEM NOTES

### Colors
- **Primary (Purple):** #7C3AED
- **Primary Hover:** #8B5CF6
- **Background:** #0A0A14
- **Surface:** #0D0B1A
- **Panel:** #120E24
- **Border:** rgba(124, 58, 237, 0.1)
- **Text:** #E0E7FF
- **Success:** #10B981
- **Warning:** #F59E0B
- **Error:** #EF4444

### Typography
- **Hero Title:** 48px, bold, gradient (purple→cyan)
- **Section Title:** 32px, bold, white
- **Subheading:** 20px, medium, light
- **Body:** 16px, regular, light-gray
- **Small Text:** 14px, regular, dimmed

### Components
- **Buttons:** 12px padding vertical, 24px horizontal, rounded-lg, glow on hover
- **Cards:** rounded-xl, border: 1px rgba(124, 58, 237, 0.2), shadow-lg
- **Inputs:** rounded-lg, border: 1px rgba(124, 58, 237, 0.3), padding: 12px
- **Transitions:** all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)

### Spacing
- Hero section: 100px top/bottom padding
- Section padding: 80px vertical, 40px horizontal
- Card gap: 24px
- Mobile padding: 20px horizontal

### Animations
- Fade-in on scroll
- Button hover: translate-y(-2px), box-shadow glow
- Links hover: text color shift to cyan
- Mobile menu: slide-in from left

---

## IMPLEMENTATION PRIORITY

1. **Homepage (Updated)** — Fix hero messaging, remove 3 products section
2. **How It Works** — 4-step visual guide (highest conversion impact)
3. **Pricing** — Drive purchasing decisions
4. **Use Cases** — Show relevance to different personas
5. **Login** — Allow existing users to access

---

**Next Step:** Review these wireframes, approve the layout/flow, then begin building! 🚀
