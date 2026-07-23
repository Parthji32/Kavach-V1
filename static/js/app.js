// KAVACH Frontend App v1.1.0
// Handles HTMX integration, auth state, and UI interactions

document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    initializeClipboard();
});

function initializeApp() {
    console.log('🛡️ KAVACH app initialized');
    
    // HTMX: Handle 401 responses by redirecting to login
    document.body.addEventListener('htmx:responseError', function(event) {
        if (event.detail.xhr && event.detail.xhr.status === 401) {
            // Session expired - redirect to login
            window.location.href = '/login';
        }
    });

    // HTMX: After successful swap, check for redirect signals
    document.body.addEventListener('htmx:afterRequest', function(event) {
        const xhr = event.detail.xhr;
        if (!xhr) return;
        
        // If the server sends HX-Redirect, htmx handles it automatically.
        // For non-htmx fetch calls with 401, redirect:
        if (xhr.status === 401 && !event.detail.elt) {
            window.location.href = '/login';
        }
    });

    // HTMX: Show loading state
    document.body.addEventListener('htmx:beforeRequest', function(event) {
        const btn = event.detail.elt.querySelector('button[type=submit]');
        if (btn) {
            btn.dataset.originalText = btn.innerText;
            btn.innerText = 'Loading...';
        }
    });

    // HTMX: Restore button after response
    document.body.addEventListener('htmx:afterRequest', function(event) {
        const btn = event.detail.elt.querySelector('button[type=submit]');
        if (btn && btn.dataset.originalText) {
            btn.innerText = btn.dataset.originalText;
            delete btn.dataset.originalText;
        }
    });
}

function initializeClipboard() {
    // Copy token to clipboard on click
    document.addEventListener('click', function(e) {
        const btn = e.target.closest('[data-copy-text]');
        if (!btn) return;
        
        const text = btn.getAttribute('data-copy-text');
        navigator.clipboard.writeText(text).then(() => {
            const originalText = btn.innerText;
            btn.innerText = '✓ Copied!';
            setTimeout(() => { btn.innerText = originalText; }, 2000);
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
function showToast(message, type) {
    type = type || 'info';
    const toast = document.createElement('div');
    const bgColor = type === 'error' ? 'bg-red-500' : type === 'success' ? 'bg-green-500' : 'bg-blue-500';
    toast.className = 'fixed bottom-4 right-4 px-4 py-3 rounded-lg text-white font-medium z-50 ' + bgColor;
    toast.innerText = message;
    document.body.appendChild(toast);
    
    setTimeout(function() {
        toast.style.opacity = '0';
        toast.style.transition = 'opacity 0.3s ease';
        setTimeout(function() { toast.remove(); }, 300);
    }, 3000);
}

// Logout: clear cookie and redirect
function logout() {
    document.cookie = 'token=; Max-Age=0; path=/;';
    window.location.href = '/login';
}
