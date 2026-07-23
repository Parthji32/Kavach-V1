package handlers

import (
	"strconv"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
)

// Helper to convert interface to string
func toString(val interface{}) string {
	if val == nil {
		return "0"
	}
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.Itoa(int(v))
	case string:
		return v
	default:
		return "0"
	}
}

// LoginPage renders the login page
func LoginPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign In — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <meta property="og:title" content="Sign In — Kavach">
    <meta property="og:description" content="Sign in to your Kavach security deception dashboard.">
    <meta property="og:type" content="website">
    <link rel="stylesheet" href="/static/css/tailwind.css">
    <script src="/static/js/htmx.min.js"></script>
</head>
<body class="bg-[#0A0A14] text-gray-300 min-h-screen">
    <div class="min-h-screen flex items-center justify-center px-4">
        <div class="w-full max-w-md">
            <div class="text-center mb-8">
                <a href="/" class="inline-block">
                    <h1 class="text-3xl font-bold bg-gradient-to-r from-purple-500 to-cyan-500 bg-clip-text text-transparent">Kavach</h1>
                </a>
                <p class="text-gray-500 mt-1">Armor that fights back</p>
            </div>

            <div class="bg-[#0D0B1A] border border-[#1E1A30] rounded-xl p-8">
                <h2 class="text-xl font-semibold text-white mb-6">Sign in to your account</h2>

                <form hx-post="/api/auth/login" hx-target="#auth-result" hx-swap="innerHTML" hx-disabled-elt="find button[type=submit]">
                    <div class="space-y-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-300 mb-1.5">Email</label>
                            <input type="email" name="email" required
                                placeholder="you@company.com"
                                class="w-full bg-[#0A0A14] border border-[#1E1A30] rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-600/50 outline-none transition">
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-300 mb-1.5">Password</label>
                            <input type="password" name="password" required
                                placeholder="Enter your password"
                                class="w-full bg-[#0A0A14] border border-[#1E1A30] rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-600/50 outline-none transition">
                        </div>
                    </div>

                    <button type="submit" class="w-full mt-6 bg-purple-600 hover:bg-purple-700 text-white py-2.5 rounded-lg text-sm font-semibold transition">
                        Sign In
                    </button>

                    <div id="auth-result" class="mt-4"></div>
                </form>

                <p class="text-center text-sm text-gray-500 mt-6">
                    No account yet?
                    <a href="/signup" class="text-purple-400 hover:text-purple-300 transition">Create one</a>
                </p>
            </div>

            <p class="text-center text-xs text-gray-600 mt-6">
                <a href="/privacy" class="hover:text-gray-400">Privacy</a> · <a href="/terms" class="hover:text-gray-400">Terms</a>
            </p>
        </div>
    </div>
</body>
</html>`)
}

// SignupPage renders the signup page
func SignupPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign Up — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <meta property="og:title" content="Sign Up — Kavach">
    <meta property="og:description" content="Create your Kavach account — security deception platform.">
    <meta property="og:type" content="website">
    <link rel="stylesheet" href="/static/css/tailwind.css">
    <script src="/static/js/htmx.min.js"></script>
</head>
<body class="bg-[#0A0A14] text-gray-300 min-h-screen">
    <div class="min-h-screen flex items-center justify-center px-4">
        <div class="w-full max-w-md">
            <div class="text-center mb-8">
                <a href="/" class="inline-block">
                    <h1 class="text-3xl font-bold bg-gradient-to-r from-purple-500 to-cyan-500 bg-clip-text text-transparent">Kavach</h1>
                </a>
                <p class="text-gray-500 mt-1">Create your account</p>
            </div>

            <div class="bg-[#0D0B1A] border border-[#1E1A30] rounded-xl p-8">
                <h2 class="text-xl font-semibold text-white mb-6">Get started for free</h2>

                <form hx-post="/api/auth/register" hx-target="#auth-result" hx-swap="innerHTML" hx-disabled-elt="find button[type=submit]">
                    <div class="space-y-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-300 mb-1.5">Full Name</label>
                            <input type="text" name="full_name" required
                                placeholder="Jane Smith"
                                class="w-full bg-[#0A0A14] border border-[#1E1A30] rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-600/50 outline-none transition">
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-300 mb-1.5">Email</label>
                            <input type="email" name="email" required
                                placeholder="you@company.com"
                                class="w-full bg-[#0A0A14] border border-[#1E1A30] rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-600/50 outline-none transition">
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-300 mb-1.5">Password</label>
                            <input type="password" name="password" required minlength="8"
                                placeholder="Min 8 characters"
                                class="w-full bg-[#0A0A14] border border-[#1E1A30] rounded-lg px-4 py-2.5 text-sm text-white placeholder-gray-500 focus:border-purple-600 focus:ring-1 focus:ring-purple-600/50 outline-none transition">
                        </div>
                    </div>

                    <button type="submit" class="w-full mt-6 bg-purple-600 hover:bg-purple-700 text-white py-2.5 rounded-lg text-sm font-semibold transition">
                        Create Account
                    </button>

                    <div id="auth-result" class="mt-4"></div>
                </form>

                <p class="text-center text-sm text-gray-500 mt-6">
                    Already have an account?
                    <a href="/login" class="text-purple-400 hover:text-purple-300 transition">Sign in</a>
                </p>
            </div>

            <p class="text-center text-xs text-gray-600 mt-6">
                <a href="/privacy" class="hover:text-gray-400">Privacy</a> · <a href="/terms" class="hover:text-gray-400">Terms</a>
            </p>
        </div>
    </div>
</body>
</html>`)
}

// LandingPage renders a landing page
func LandingPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>KAVACH — Security Deception Platform</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <meta property="og:title" content="KAVACH — Security Deception Platform">
    <meta property="og:description" content="Proactive security through deception. Deploy honeypot tokens to detect, profile, and neutralize attackers.">
    <meta property="og:type" content="website">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-[#0A0A14] text-white">
    <div class="min-h-screen flex items-center justify-center">
        <div class="text-center px-4">
            <h1 class="text-5xl font-bold mb-4 bg-gradient-to-r from-purple-500 to-cyan-500 bg-clip-text text-transparent">KAVACH</h1>
            <p class="text-xl text-gray-400 mb-8">Armor that fights back</p>
            <p class="text-gray-500 mb-8 max-w-md mx-auto">Security deception platform for detecting and profiling attackers using honeypot tokens.</p>
            <div class="space-x-4">
                <a href="/login" class="px-6 py-2.5 bg-purple-600 hover:bg-purple-700 rounded-lg font-semibold transition">Login</a>
                <a href="/signup" class="px-6 py-2.5 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg font-semibold transition">Sign Up</a>
            </div>
        </div>
    </div>
</body>
</html>`)
}

// DashboardPage renders the main dashboard (protected by JWTAuthPage middleware)
func DashboardPage(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	log.Printf("[DASHBOARD] userID: '%s'", userID)

	// Get dashboard stats
	stats, err := database.GetDashboardStats(userID)
	if err != nil {
		log.Printf("[DASHBOARD] Error fetching stats: %v", err)
		stats = map[string]interface{}{
			"total_tokens":    0,
			"active_tokens":   0,
			"total_attackers": 0,
			"events_last_24h": 0,
		}
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
    <script src="/static/js/htmx.min.js"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold">Dashboard</h1>
            <button onclick="logout()" class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-sm font-medium transition">Logout</button>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="bg-gray-900 p-4 rounded-lg border border-gray-700">
                <p class="text-gray-400 text-sm">Total Tokens</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["total_tokens"]) + `</p>
            </div>
            <div class="bg-gray-900 p-4 rounded-lg border border-gray-700">
                <p class="text-gray-400 text-sm">Active Tokens</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["active_tokens"]) + `</p>
            </div>
            <div class="bg-gray-900 p-4 rounded-lg border border-gray-700">
                <p class="text-gray-400 text-sm">Total Attackers</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["total_attackers"]) + `</p>
            </div>
            <div class="bg-gray-900 p-4 rounded-lg border border-gray-700">
                <p class="text-gray-400 text-sm">Events (24h)</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["events_last_24h"]) + `</p>
            </div>
        </div>
        <div class="mt-8 flex flex-wrap gap-3">
            <a href="/tokens" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg text-sm font-medium transition">Tokens</a>
            <a href="/attackers" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg text-sm font-medium transition">Attackers</a>
            <a href="/alerts" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg text-sm font-medium transition">Alerts</a>
            <a href="/settings" class="px-4 py-2 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg text-sm font-medium transition">Settings</a>
        </div>
    </div>
    <script>
    function logout() {
        fetch('/api/auth/logout', { method: 'POST', credentials: 'include' }).finally(() => {
            document.cookie = 'token=; Max-Age=0; path=/;';
            window.location.href = '/login';
        });
    }
    </script>
</body>
</html>`)
}

// TokensPage renders the tokens list page (protected by JWTAuthPage middleware)
func TokensPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tokens — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
    <script src="/static/js/htmx.min.js"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold">Tokens</h1>
            <a href="/tokens/new" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg text-sm font-medium transition">+ New Token</a>
        </div>
        <p class="mt-8"><a href="/app" class="text-purple-400 hover:text-purple-300">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// NewTokenPage renders the create token form (protected by JWTAuthPage middleware)
func NewTokenPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Create Token — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
    <script src="/static/js/htmx.min.js"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Create New Token</h1>
        <form id="tokenForm" class="max-w-md space-y-4" onsubmit="handleCreateToken(event)">
            <select id="token_type" required class="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-white">
                <option value="">Select Token Type</option>
                <option value="url">URL Token</option>
                <option value="api_key">API Key</option>
                <option value="document">Document</option>
                <option value="dns">DNS</option>
                <option value="email">Email</option>
            </select>
            <input type="text" id="description" placeholder="Description" class="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-white">
            <button type="submit" class="w-full px-4 py-2.5 bg-purple-600 hover:bg-purple-700 rounded-lg font-semibold transition" id="submitBtn">Create Token</button>
        </form>
        <p class="mt-8"><a href="/tokens" class="text-purple-400 hover:text-purple-300">← Back to Tokens</a></p>
    </div>
    <script>
    async function handleCreateToken(e) {
        e.preventDefault();
        const token_type = document.getElementById('token_type').value;
        const description = document.getElementById('description').value;
        const submitBtn = document.getElementById('submitBtn');
        if (!token_type) { alert('Please select a token type'); return; }
        submitBtn.disabled = true;
        submitBtn.innerText = 'Creating...';
        try {
            const resp = await fetch('/api/tokens', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ token_type, description })
            });
            const data = await resp.json();
            if (data.success && data.data) {
                alert('Token created! Value: ' + data.data.token_value);
                window.location = '/tokens';
            } else {
                alert('Failed: ' + (data.message || 'Unknown error'));
            }
        } catch (err) { alert('Error: ' + err); }
        finally { submitBtn.disabled = false; submitBtn.innerText = 'Create Token'; }
    }
    </script>
</body>
</html>`)
}

// TokenDetailPage renders a single token detail (protected by JWTAuthPage middleware)
func TokenDetailPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Token Detail — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Token Detail</h1>
        <p class="mt-8"><a href="/tokens" class="text-purple-400 hover:text-purple-300">← Back to Tokens</a></p>
    </div>
</body>
</html>`)
}

// AttackersPage renders the attackers list page (protected by JWTAuthPage middleware)
func AttackersPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Attackers — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Attackers</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400 hover:text-purple-300">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// AttackerDetailPage renders a single attacker detail (protected by JWTAuthPage middleware)
func AttackerDetailPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Attacker Detail — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Attacker Detail</h1>
        <p class="mt-8"><a href="/attackers" class="text-purple-400 hover:text-purple-300">← Back to Attackers</a></p>
    </div>
</body>
</html>`)
}

// AlertsPage renders the alerts list page (protected by JWTAuthPage middleware)
func AlertsPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Alerts — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold">Alerts</h1>
            <a href="/alerts/new" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-lg text-sm font-medium transition">+ New Alert</a>
        </div>
        <p class="mt-8"><a href="/app" class="text-purple-400 hover:text-purple-300">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// IntegrationsPage renders the integrations page (protected by JWTAuthPage middleware)
func IntegrationsPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Integrations — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Integrations</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400 hover:text-purple-300">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// SettingsPage renders the settings page (protected by JWTAuthPage middleware)
func SettingsPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Settings</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400 hover:text-purple-300">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// ProfilePage renders the user profile page (protected by JWTAuthPage middleware)
func ProfilePage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Profile — Kavach</title>
    <link rel="icon" href="/static/favicon.svg" type="image/svg+xml">
    <link rel="stylesheet" href="/static/css/tailwind.css">
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Profile</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400 hover:text-purple-300">← Dashboard</a></p>
    </div>
</body>
</html>`)
}
