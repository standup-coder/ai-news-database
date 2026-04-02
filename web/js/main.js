/**
 * News4Coder Web UI - Interactive Functions
 * CNCF Design System
 */

document.addEventListener('DOMContentLoaded', () => {
  // Initialize all modules
  initNavigation();
  installTabs();
  copyToClipboard();
  animateOnScroll();
  smoothScroll();
});

/**
 * Navigation
 * - Fixed header on scroll
 * - Mobile menu toggle
 * - Active link highlighting
 */
function initNavigation() {
  const navbar = document.querySelector('.navbar');
  const mobileMenuBtn = document.querySelector('.mobile-menu-btn');
  const navLinks = document.querySelector('.nav-links');
  let lastScroll = 0;

  // Header shadow on scroll
  window.addEventListener('scroll', () => {
    const currentScroll = window.pageYOffset;
    
    if (currentScroll > 50) {
      navbar.style.boxShadow = '0 4px 20px rgba(0, 0, 0, 0.08)';
    } else {
      navbar.style.boxShadow = 'none';
    }
    
    lastScroll = currentScroll;
  });

  // Mobile menu toggle
  if (mobileMenuBtn) {
    mobileMenuBtn.addEventListener('click', () => {
      mobileMenuBtn.classList.toggle('active');
      navLinks.classList.toggle('active');
      
      // Animate hamburger to X
      const spans = mobileMenuBtn.querySelectorAll('span');
      if (mobileMenuBtn.classList.contains('active')) {
        spans[0].style.transform = 'rotate(45deg) translate(5px, 5px)';
        spans[1].style.opacity = '0';
        spans[2].style.transform = 'rotate(-45deg) translate(5px, -5px)';
      } else {
        spans[0].style.transform = 'none';
        spans[1].style.opacity = '1';
        spans[2].style.transform = 'none';
      }
    });
  }

  // Close mobile menu on link click
  document.querySelectorAll('.nav-link').forEach(link => {
    link.addEventListener('click', () => {
      navLinks.classList.remove('active');
      mobileMenuBtn?.classList.remove('active');
      const spans = mobileMenuBtn?.querySelectorAll('span');
      if (spans) {
        spans[0].style.transform = 'none';
        spans[1].style.opacity = '1';
        spans[2].style.transform = 'none';
      }
    });
  });
}

/**
 * Install Tabs
 * - Switch between OS install commands
 */
function installTabs() {
  const tabs = document.querySelectorAll('.install-tab');
  const commandEl = document.getElementById('install-command');
  
  const commands = {
    mac: 'curl -sSL https://get.news4coder.dev | bash',
    linux: 'curl -sSL https://get.news4coder.dev | bash',
    windows: 'irm https://get.news4coder.dev | iex'
  };

  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      // Remove active from all tabs
      tabs.forEach(t => t.classList.remove('active'));
      
      // Add active to clicked tab
      tab.classList.add('active');
      
      // Update command
      const os = tab.dataset.os;
      if (commandEl && commands[os]) {
        // Fade out
        commandEl.style.opacity = '0';
        
        setTimeout(() => {
          commandEl.textContent = commands[os];
          commandEl.style.opacity = '1';
        }, 150);
      }
    });
  });
}

/**
 * Copy to Clipboard
 * - Copy install command
 */
function copyToClipboard() {
  const copyBtn = document.querySelector('.copy-btn');
  const commandEl = document.getElementById('install-command');
  
  if (copyBtn && commandEl) {
    copyBtn.addEventListener('click', async () => {
      try {
        await navigator.clipboard.writeText(commandEl.textContent);
        
        // Visual feedback
        const originalHTML = copyBtn.innerHTML;
        copyBtn.innerHTML = `
          <svg viewBox="0 0 24 24" fill="none" stroke="#2ECC71" stroke-width="2">
            <path d="M20 6L9 17l-5-5"/>
          </svg>
        `;
        copyBtn.style.color = '#2ECC71';
        
        setTimeout(() => {
          copyBtn.innerHTML = originalHTML;
          copyBtn.style.color = '';
        }, 2000);
      } catch (err) {
        console.error('Failed to copy:', err);
      }
    });
  }
}

/**
 * Animate on Scroll
 * - Intersection Observer for scroll animations
 */
function animateOnScroll() {
  const animatedElements = document.querySelectorAll(
    '.feature-card, .source-card, .workflow-step, .section-header'
  );
  
  const observerOptions = {
    root: null,
    rootMargin: '0px',
    threshold: 0.1
  };
  
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.style.opacity = '1';
        entry.target.style.transform = 'translateY(0)';
        observer.unobserve(entry.target);
      }
    });
  }, observerOptions);
  
  animatedElements.forEach((el, index) => {
    // Set initial state
    el.style.opacity = '0';
    el.style.transform = 'translateY(20px)';
    el.style.transition = `opacity 0.5s ease ${index * 0.05}s, transform 0.5s ease ${index * 0.05}s`;
    
    observer.observe(el);
  });
}

/**
 * Smooth Scroll
 * - Smooth scroll for anchor links
 */
function smoothScroll() {
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      const href = this.getAttribute('href');
      if (href === '#') return;
      
      const target = document.querySelector(href);
      if (target) {
        e.preventDefault();
        
        const headerOffset = 80;
        const elementPosition = target.getBoundingClientRect().top;
        const offsetPosition = elementPosition + window.pageYOffset - headerOffset;
        
        window.scrollTo({
          top: offsetPosition,
          behavior: 'smooth'
        });
      }
    });
  });
}

/**
 * Terminal Typing Effect (Optional Enhancement)
 * - Simulates typing in the terminal
 */
function terminalTypingEffect() {
  const cursor = document.querySelector('.terminal-cursor');
  
  if (cursor) {
    // Ensure cursor keeps blinking
    setInterval(() => {
      cursor.style.opacity = cursor.style.opacity === '0' ? '1' : '0';
    }, 530);
  }
}

// Run terminal effect
terminalTypingEffect();

/**
 * Stats Counter Animation
 * - Animate numbers when in view
 */
function animateStats() {
  const stats = document.querySelectorAll('.stat-value');
  
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const el = entry.target;
        const finalValue = el.textContent;
        
        // If it's a number, animate it
        if (!isNaN(parseInt(finalValue))) {
          const num = parseInt(finalValue);
          let current = 0;
          const increment = num / 30;
          const timer = setInterval(() => {
            current += increment;
            if (current >= num) {
              el.textContent = finalValue;
              clearInterval(timer);
            } else {
              el.textContent = Math.floor(current) + '+';
            }
          }, 30);
        }
        
        observer.unobserve(el);
      }
    });
  }, { threshold: 0.5 });
  
  stats.forEach(stat => observer.observe(stat));
}

// Initialize stats animation
document.addEventListener('DOMContentLoaded', animateStats);
