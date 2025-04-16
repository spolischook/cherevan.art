/**
 * Independent scrollable columns functionality (modern approach)
 * - Each column scrolls independently when hovered
 * - If the column can’t scroll further, the page scrolls naturally
 */
document.addEventListener('DOMContentLoaded', function() {
  if (window.innerWidth >= 640) {
    const columnsContainer = document.querySelector('section.art-columns-container');
    if (!columnsContainer) return;
    const galleryWrapper = columnsContainer.querySelector('div:first-child');
    const galleryContent = galleryWrapper.querySelector('.desktop-gallery');
    const infoColumn = columnsContainer.querySelector('div:last-child');
    if (!galleryContent || !infoColumn) return;

    function canScroll(el, deltaY) {
      if (deltaY > 0) {
        // Scrolling down
        return el.scrollTop + el.clientHeight < el.scrollHeight;
      } else if (deltaY < 0) {
        // Scrolling up
        return el.scrollTop > 0;
      }
      return false;
    }

    function wheelHandler(e, el) {
      if (canScroll(el, e.deltaY)) {
        e.preventDefault();
        el.scrollBy({ top: e.deltaY, behavior: 'auto' });
      }
      // If the container can’t scroll, event is not prevented so the page scrolls
    }

    galleryContent.addEventListener('wheel', e => wheelHandler(e, galleryContent), { passive: false });
    infoColumn.addEventListener('wheel', e => wheelHandler(e, infoColumn), { passive: false });
  }
});
