// Intersection Observer for fade-in animations on scroll
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.classList.add('fade-in');
            observer.unobserve(entry.target);
        }
    });
}, observerOptions);

// Apply observer to all sections
document.querySelectorAll('section, .feature-card, .metric').forEach(el => {
    observer.observe(el);
});

// Smooth header background on scroll
window.addEventListener('scroll', () => {
    const header = document.querySelector('header');
    if (window.scrollY > 50) {
        header.style.background = 'rgba(18, 14, 36, 0.9)';
        header.style.backdropFilter = 'blur(20px)';
    } else {
        header.style.background = 'rgba(18, 14, 36, 0.7)';
        header.style.backdropFilter = 'blur(12px)';
    }
});

// Glow effect on button hover
document.querySelectorAll('button, .btn, a.btn').forEach(btn => {
    btn.addEventListener('mouseenter', function() {
        this.style.boxShadow = '0 0 20px rgba(124, 58, 237, 0.5)';
    });
    btn.addEventListener('mouseleave', function() {
        this.style.boxShadow = '0 0 10px rgba(124, 58, 237, 0.2)';
    });
});

console.log('✨ Animations loaded!');
