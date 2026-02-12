// Samuel - Custom JavaScript

// Initialize when document is ready
document.addEventListener('DOMContentLoaded', function() {
  // Add copy feedback enhancement
  enhanceCopyButtons();

  // Add external link handling
  handleExternalLinks();
});

// Enhance copy button feedback
function enhanceCopyButtons() {
  document.querySelectorAll('.md-clipboard').forEach(function(button) {
    button.addEventListener('click', function() {
      // Button already has built-in feedback from Material
      // This is a placeholder for any additional enhancements
    });
  });
}

// Open external links in new tab
function handleExternalLinks() {
  document.querySelectorAll('a[href^="http"]').forEach(function(link) {
    // Skip internal links to the same domain
    if (!link.href.includes(window.location.hostname)) {
      link.setAttribute('target', '_blank');
      link.setAttribute('rel', 'noopener noreferrer');
    }
  });
}

// Optional: Track page views (placeholder for analytics)
// function trackPageView() {
//   // Add your analytics tracking code here
// }
