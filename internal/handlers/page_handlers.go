package handlers

import (
	"strconv"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jindal-parth/kavach/internal/database"
	"github.com/jindal-parth/kavach/internal/middleware"
	"github.com/jindal-parth/kavach/internal/services"
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

// extractUserIDFromRequest tries to get userID from JWT in cookie or Authorization header
func extractUserIDFromRequest(c *fiber.Ctx) string {
	// First try c.Locals (set by JWT middleware if route is protected)
	userID := middleware.GetUserID(c)
	if userID != "" {
		return userID
	}

	// Try to extract JWT from cookie
	tokenString := c.Cookies("token")
	if tokenString == "" {
		// Try to extract from Authorization header
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}
	}

	if tokenString == "" {
		return ""
	}

	// Validate token and extract userID
	claims, err := services.ValidateJWT(tokenString)
	if err != nil {
		return ""
	}
	return claims.UserID
}

// LoginPage renders the login page
func LoginPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Login - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-black text-white">
    <div class="min-h-screen flex items-center justify-center">
        <div class="max-w-md w-full">
            <h1 class="text-3xl font-bold text-center mb-8">KAVACH</h1>
            <p class="text-center text-gray-400 mb-8">Armor that fights back</p>
            <form id="loginForm" class="space-y-4" onsubmit="handleLogin(event)">
                <input type="email" id="email" placeholder="Email" required class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
                <input type="password" id="password" placeholder="Password" required class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
                <button type="submit" class="w-full px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded font-semibold">Sign In</button>
            </form>
            <p class="text-center text-gray-400 mt-4">Don't have an account? <a href="/signup" class="text-purple-400">Sign up</a></p>
        </div>
    </div>
    <script>
    function handleLogin(e) {
        e.preventDefault();
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        })
        .then(r => r.json())
        .then(d => {
            if(d.data && d.data.token) {
                // Cookie is set by server, just redirect
                setTimeout(() => { window.location = '/app'; }, 100);
            } else {
                alert('Login failed: ' + (d.message || 'Unknown error'));
            }
        })
        .catch(err => alert('Error: ' + err));
    }
    </script>
</body>
</html>`)
}

// SignupPage renders the signup page
func SignupPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Sign Up - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-black text-white">
    <div class="min-h-screen flex items-center justify-center">
        <div class="max-w-md w-full">
            <h1 class="text-3xl font-bold text-center mb-8">KAVACH</h1>
            <p class="text-center text-gray-400 mb-8">Create your account</p>
            <form id="signupForm" class="space-y-4" onsubmit="handleSignup(event)">
                <input type="text" id="full_name" placeholder="Full Name" required class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
                <input type="email" id="email" placeholder="Email" required class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
                <input type="password" id="password" placeholder="Password" required class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
                <button type="submit" class="w-full px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded font-semibold">Sign Up</button>
            </form>
            <p class="text-center text-gray-400 mt-4">Already have an account? <a href="/login" class="text-purple-400">Sign in</a></p>
        </div>
    </div>
    <script>
    function handleSignup(e) {
        e.preventDefault();
        const full_name = document.getElementById('full_name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        fetch('/api/auth/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ full_name, email, password })
        })
        .then(r => r.json())
        .then(d => {
            if(d.success) {
                alert('Account created! Redirecting to login...');
                window.location = '/login';
            } else {
                alert('Signup failed: ' + (d.message || 'Unknown error'));
            }
        })
        .catch(err => alert('Error: ' + err));
    }
    </script>
</body>
</html>`)
}

// LandingPage renders a landing page
func LandingPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>KAVACH - Security Deception Platform</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-black text-white">
    <div class="min-h-screen flex items-center justify-center">
        <div class="text-center">
            <h1 class="text-5xl font-bold mb-4">KAVACH</h1>
            <p class="text-xl text-gray-400 mb-8">Armor that fights back</p>
            <p class="text-gray-500 mb-8">Security deception platform for detecting and profiling attackers</p>
            <div class="space-x-4">
                <a href="/login" class="px-6 py-2 bg-purple-600 hover:bg-purple-700 rounded font-semibold">Login</a>
                <a href="/signup" class="px-6 py-2 bg-gray-800 hover:bg-gray-700 rounded font-semibold">Sign Up</a>
            </div>
        </div>
    </div>
</body>
</html>`)
}

// DashboardPage renders the main dashboard
func DashboardPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	log.Printf("[DASHBOARD] userID extracted: '%s'", userID)
	if userID == "" {
		return c.Redirect("/login")
	}

	// Get dashboard stats
	log.Printf("[DASHBOARD] Fetching stats for userID: %s", userID)
	stats, err := database.GetDashboardStats(userID)
	log.Printf("[DASHBOARD] GetDashboardStats result: err=%v, stats=%v", err, stats)
	if err != nil {
		return c.SendString(`<h1>Dashboard</h1><p>Error loading dashboard</p>`)
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Dashboard - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold">Dashboard</h1>
            <a href="/login" onclick="logout(); return false;" class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded">Logout</a>
        </div>
        <div class="grid grid-cols-4 gap-4">
            <div class="bg-gray-900 p-4 rounded border border-gray-700">
                <p class="text-gray-400">Total Tokens</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["total_tokens"]) + `</p>
            </div>
            <div class="bg-gray-900 p-4 rounded border border-gray-700">
                <p class="text-gray-400">Active Tokens</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["active_tokens"]) + `</p>
            </div>
            <div class="bg-gray-900 p-4 rounded border border-gray-700">
                <p class="text-gray-400">Total Attackers</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["total_attackers"]) + `</p>
            </div>
            <div class="bg-gray-900 p-4 rounded border border-gray-700">
                <p class="text-gray-400">Events (24h)</p>
                <p class="text-3xl font-bold text-purple-400">` + toString(stats["events_last_24h"]) + `</p>
            </div>
        </div>
        <div class="mt-8 space-x-4">
            <a href="/tokens" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded">Tokens</a>
            <a href="/attackers" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded">Attackers</a>
            <a href="/alerts" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded">Alerts</a>
        </div>
    </div>
    <script>
    function logout() {
        localStorage.removeItem('token');
    }
    </script>
</body>
</html>`)
}

// TokensPage renders the tokens list page
func TokensPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Tokens - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Tokens</h1>
        <a href="/tokens/new" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded">+ New Token</a>
        <p class="mt-8"><a href="/app" class="text-purple-400">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// NewTokenPage renders the create token form
func NewTokenPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Create Token - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Create New Token</h1>
        <form id="tokenForm" class="max-w-md space-y-4" onsubmit="handleCreateToken(event)">
            <select id="token_type" required class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
                <option value="">Select Token Type</option>
                <option value="url">URL Token</option>
                <option value="api_key">API Key</option>
                <option value="document">Document</option>
                <option value="dns">DNS</option>
                <option value="email">Email</option>
            </select>
            <input type="text" id="description" placeholder="Description" class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded text-white">
            <button type="submit" class="w-full px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded font-semibold" id="submitBtn">Create Token</button>
        </form>
        <p class="mt-8"><a href="/tokens" class="text-purple-400">← Back to Tokens</a></p>
    </div>
    <script>
    async function handleCreateToken(e) {
        e.preventDefault();
        const token_type = document.getElementById('token_type').value;
        const description = document.getElementById('description').value;
        const submitBtn = document.getElementById('submitBtn');
        
        if (!token_type) {
            alert('Please select a token type');
            return;
        }
        
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
                alert('Failed to create token: ' + (data.message || 'Unknown error'));
            }
        } catch (err) {
            alert('Error: ' + err);
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerText = 'Create Token';
        }
    }
    </script>
</body>
</html>`)
}

// TokenDetailPage renders a single token detail
func TokenDetailPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Token Detail - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Token Detail</h1>
        <p class="mt-8"><a href="/tokens" class="text-purple-400">← Back to Tokens</a></p>
    </div>
</body>
</html>`)
}

// AttackersPage renders the attackers list page
func AttackersPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Attackers - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Attackers</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// AttackerDetailPage renders a single attacker detail
func AttackerDetailPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Attacker Detail - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Attacker Detail</h1>
        <p class="mt-8"><a href="/attackers" class="text-purple-400">← Back to Attackers</a></p>
    </div>
</body>
</html>`)
}

// AlertsPage renders the alerts list page
func AlertsPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Alerts - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Alerts</h1>
        <a href="/alerts/new" class="px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded">+ New Alert</a>
        <p class="mt-8"><a href="/app" class="text-purple-400">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// IntegrationsPage renders the integrations page
func IntegrationsPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Integrations - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Integrations</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// SettingsPage renders the settings page
func SettingsPage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Settings - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Settings</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400">← Dashboard</a></p>
    </div>
</body>
</html>`)
}

// ProfilePage renders the user profile page
func ProfilePage(c *fiber.Ctx) error {
	userID := extractUserIDFromRequest(c)
	if userID == "" {
		return c.Redirect("/login")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html>
<head>
    <title>Profile - KAVACH</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-950 text-gray-100">
    <div class="p-8">
        <h1 class="text-3xl font-bold mb-8">Profile</h1>
        <p class="mt-8"><a href="/app" class="text-purple-400">← Dashboard</a></p>
    </div>
</body>
</html>`)
}
