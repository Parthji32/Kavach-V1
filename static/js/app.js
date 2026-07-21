// KAVACH Frontend App
// Handles client-side interactions, HTMX integration, and UI state

document.addEventListener('DOMContentLoaded', function() {
    // Initialize app
    initializeApp();
    initializeFormHandlers();
    initializeClipboard();
});

function initializeApp() {
    console.log('🛡️ KAVACH app initialized');
    
    // HTMX event listeners
    document.body.addEventListener('htmx:afterRequest', function(event) {
        if (event.detail.xhr.status === 401) {
            window.location.href = '/login';
        }
    });

    // HTMX error handling
    document.body.addEventListener('htmx:responseError', function(event) {
        console.error('HTMX Error:', event.detail);
    });
}

function initializeFormHandlers() {
    // Login form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            // Let HTMX handle it
        });
    }

    // Token creation
    const tokenForm = document.getElementById('token-form');
    if (tokenForm) {
        tokenForm.addEventListener('submit', function(e) {
            // HTMX will post to /api/tokens
        });
    }
}

function initializeClipboard() {
    // Copy token to clipboard on click
    const copyButtons = document.querySelectorAll('[data-copy-text]');
    copyButtons.forEach(button => {
        button.addEventListener('click', function() {
            const text = this.getAttribute('data-copy-text');
            navigator.clipboard.writeText(text).then(() => {
                // Show feedback
                const originalText = this.innerText;
                this.innerText = '✓ Copied!';
                setTimeout(() => {
                    this.innerText = originalText;
                }, 2000);
            });
        });
    });
}

// Utility: Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}

// Utility: Format risk score color
function getRiskColor(score) {
    if (score > 75) return '#EF4444'; // Red
    if (score > 50) return '#F59E0B'; // Amber
    if (score > 25) return '#3B82F6'; // Blue
    return '#10B981'; // Green
}

// Show toast notification
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `fixed bottom-4 right-4 px-4 py-3 rounded-lg text-white font-medium z-50 ${
        type === 'error' ? 'bg-red-500' : 
        type === 'success' ? 'bg-green-500' : 
        'bg-blue-500'
    }`;
    toast.innerText = message;
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transition = 'opacity 0.3s ease';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// Logout
function logout() {
    localStorage.removeItem('jwt_token');
    window.location.href = '/login';
}
