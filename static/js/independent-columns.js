/**
 * Independent scrollable columns functionality
 * - When user scrolls a column, it scrolls independently
 * - When one column reaches the bottom and user continues scrolling, focus moves to the other column
 * - If both columns are at the bottom, the page scrolls normally
 */
document.addEventListener('DOMContentLoaded', function() {
  // Only run on desktop
  if (window.innerWidth >= 640) {
    const galleryColumn = document.querySelector('.gallery-column');
    const infoColumn = document.querySelector('.info-column');
    
    if (!galleryColumn || !infoColumn) return;
    
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
      
      const galleryRect = galleryColumn.getBoundingClientRect();
      const infoRect = infoColumn.getBoundingClientRect();
      
      // Determine which column the cursor is over
      const isOverGallery = e.clientX >= galleryRect.left && e.clientX <= galleryRect.right &&
                          e.clientY >= galleryRect.top && e.clientY <= galleryRect.bottom;
      
      const isOverInfo = e.clientX >= infoRect.left && e.clientX <= infoRect.right &&
                       e.clientY >= infoRect.top && e.clientY <= infoRect.bottom;
      
      if (isOverGallery) {
        // Scrolling down
        if (e.deltaY > 0) {
          if (isAtBottom(galleryColumn)) {
            // Gallery column is at bottom, scroll info column
            e.preventDefault();
            infoColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
            galleryScrolled = true;
          }
        } 
        // Scrolling up
        else if (e.deltaY < 0) {
          if (isAtTop(galleryColumn) && !isAtTop(infoColumn)) {
            // Gallery column is at top, scroll info column from bottom
            e.preventDefault();
            infoColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
          }
        }
      } 
      else if (isOverInfo) {
        // Scrolling down
        if (e.deltaY > 0) {
          if (isAtBottom(infoColumn)) {
            // Info column is at bottom, both columns are scrolled
            infoScrolled = true;
            
            // Both columns at bottom, allow page to scroll
            if (galleryScrolled && infoScrolled) {
              // Do nothing, let the page scroll naturally
            }
          }
        } 
        // Scrolling up
        else if (e.deltaY < 0) {
          if (isAtTop(infoColumn)) {
            // Info column is at top, scroll gallery column
            e.preventDefault();
            galleryColumn.scrollBy({ top: e.deltaY, behavior: 'auto' });
            infoScrolled = false;
          }
        }
      }
      
      // Reset scroll processing flag after a small delay
      setTimeout(function() {
        isProcessingScroll = false;
      }, 10);
    }, { passive: false });
    
    // Update indicators when gallery column scrolls
    galleryColumn.addEventListener('scroll', function() {
      galleryScrolled = isAtBottom(this);
    });
    
    // Update indicators when info column scrolls
    infoColumn.addEventListener('scroll', function() {
      infoScrolled = isAtBottom(this);
    });
  }
});
