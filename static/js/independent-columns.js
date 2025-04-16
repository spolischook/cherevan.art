/**
 * Independent scrollable columns functionality
 * - When user scrolls a column, it scrolls independently
 * - When one column reaches the bottom and user continues scrolling, focus moves to the other column
 * - If both columns are at the bottom, the page scrolls normally
 */
document.addEventListener('DOMContentLoaded', function() {
  // Only run on desktop
  if (window.innerWidth >= 640) {
    // Get the gallery wrapper (left column) and the scrollable content inside it
    // Use simpler selectors based on position rather than class names with colons
    const columnsContainer = document.querySelector('section.flex-col.sm\\:flex-row');
    const galleryWrapper = columnsContainer.querySelector('div:first-child');
    const galleryContent = galleryWrapper.querySelector('.desktop-gallery');
    const infoColumn = columnsContainer.querySelector('div:last-child');
    
    if (!galleryWrapper || !infoColumn || !galleryContent) return;
    
    let galleryScrolled = false;
    let infoScrolled = false;
    let isProcessingScroll = false;
    
    // Check if an element has reached its bottom
    function isAtBottom(element) {
      // Allow for small rounding errors with +1
      return element.scrollHeight - element.scrollTop - element.clientHeight <= 1;
    }
    
    // Check if an element is at its top
    function isAtTop(element) {
      return element.scrollTop <= 0;
    }
    
    // Handle wheel events for the entire container
    document.querySelector('.art-columns-container').addEventListener('wheel', function(e) {
      if (isProcessingScroll) return;
      isProcessingScroll = true;
      
      const galleryRect = galleryWrapper.getBoundingClientRect();
      const infoRect = infoColumn.getBoundingClientRect();
      
      // Determine which column the cursor is over
      const isOverGallery = e.clientX >= galleryRect.left && e.clientX <= galleryRect.right &&
                          e.clientY >= galleryRect.top && e.clientY <= galleryRect.bottom;
      
      const isOverInfo = e.clientX >= infoRect.left && e.clientX <= infoRect.right &&
                       e.clientY >= infoRect.top && e.clientY <= infoRect.bottom;
      
      if (isOverGallery) {
        // Scrolling down
        if (e.deltaY > 0) {
          if (isAtBottom(galleryContent)) {
            // Gallery reached bottom, scroll info column
            e.preventDefault();
            infoColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
            galleryScrolled = true;
          } else {
            // Scroll gallery normally
            e.preventDefault();
            galleryContent.scrollBy({ top: e.deltaY, behavior: 'auto' });
          }
        } 
        // Scrolling up
        else if (e.deltaY < 0) {
          if (isAtTop(galleryContent) && !isAtTop(infoColumn)) {
            // Gallery at top, scroll info column instead
            e.preventDefault();
            infoColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
          } else {
            // Scroll gallery normally
            e.preventDefault();
            galleryContent.scrollBy({ top: e.deltaY, behavior: 'auto' });
          }
        }
      } 
      else if (isOverInfo) {
        // Scrolling down
        if (e.deltaY > 0) {
          if (isAtBottom(infoColumn)) {
            // Info column at bottom
            infoScrolled = true;
            
            // Both columns at bottom, allow page to scroll
            if (galleryScrolled && infoScrolled) {
              // Do nothing, let the page scroll naturally
            }
          } else {
            // Scroll info column normally
            e.preventDefault();
            infoColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
          }
        } 
        // Scrolling up
        else if (e.deltaY < 0) {
          if (isAtTop(infoColumn)) {
            // Info column at top, scroll gallery instead
            e.preventDefault();
            galleryContent.scrollBy({ top: e.deltaY, behavior: 'auto' });
            infoScrolled = false;
          } else {
            // Scroll info column normally
            e.preventDefault();
            infoColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
          }
        }
      }
      
      // Reset scroll processing flag after a small delay
      setTimeout(function() {
        isProcessingScroll = false;
      }, 10);
    }, { passive: false });
    
    // Update indicators when gallery reaches bottom
    galleryContent.addEventListener('scroll', function() {
      galleryScrolled = isAtBottom(this);
    });
    
    // Update indicators when info column reaches bottom
    infoColumn.addEventListener('scroll', function() {
      infoScrolled = isAtBottom(this);
    });
  }
});
