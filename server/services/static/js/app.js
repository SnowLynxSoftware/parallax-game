// Parallax Custom JavaScript
$(document).ready(function () {
  console.log("Parallax initialized");

  // Add fade-in animation to elements on page load
  $("body").addClass("fade-in");

  // Intersection Observer for scroll animations
  if ("IntersectionObserver" in window) {
    const observerOptions = {
      threshold: 0.1,
      rootMargin: "0px 0px -50px 0px",
    };

    const observer = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (entry.isIntersecting) {
          entry.target.classList.add("fade-in");
          observer.unobserve(entry.target);
        }
      });
    }, observerOptions);

    // Observe elements that should animate on scroll
    $(".feature-card").each(function () {
      observer.observe(this);
    });
  }

  // Add ripple effect to buttons
  $(".btn").on("click", function (e) {
    const button = $(this);
    const ripple = $('<span class="ripple"></span>');

    const size = Math.max(button.outerWidth(), button.outerHeight());
    const x = e.pageX - button.offset().left - size / 2;
    const y = e.pageY - button.offset().top - size / 2;

    ripple.css({
      position: "absolute",
      width: size,
      height: size,
      left: x,
      top: y,
      background: "rgba(255, 255, 255, 0.3)",
      borderRadius: "50%",
      transform: "scale(0)",
      animation: "ripple 0.6s linear",
      pointerEvents: "none",
    });

    button.css("position", "relative").css("overflow", "hidden").append(ripple);

    setTimeout(() => {
      ripple.remove();
    }, 600);
  });

  // Add CSS for ripple animation
  if (!$("#ripple-styles").length) {
    $(
      '<style id="ripple-styles">@keyframes ripple { to { transform: scale(4); opacity: 0; } }</style>'
    ).appendTo("head");
  }

  // Preload images for better performance
  function preloadImages() {
    const images = [
      // Add any images you want to preload here
    ];

    images.forEach(function (src) {
      const img = new Image();
      img.src = src;
    });
  }

  preloadImages();

  // Add loading state management
  window.showLoading = function () {
    $("body").addClass("loading");
  };

  window.hideLoading = function () {
    $("body").removeClass("loading");
  };

  // Performance monitoring
  if ("performance" in window) {
    window.addEventListener("load", function () {
      setTimeout(function () {
        const perfData = window.performance.timing;
        const loadTime = perfData.loadEventEnd - perfData.navigationStart;
        console.log("Page load time:", loadTime + "ms");
      }, 0);
    });
  }
});

// Utility functions
window.Parallax = {
  utils: {
    // Throttle function for performance
    throttle: function (func, limit) {
      let inThrottle;
      return function () {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
          func.apply(context, args);
          inThrottle = true;
          setTimeout(() => (inThrottle = false), limit);
        }
      };
    },

    // Debounce function for search/input
    debounce: function (func, delay) {
      let timeoutId;
      return function () {
        const args = arguments;
        const context = this;
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(context, args), delay);
      };
    },

    // Format numbers with commas
    formatNumber: function (num) {
      return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    },

    // Validate email
    isValidEmail: function (email) {
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return re.test(email);
    },
  },
};
