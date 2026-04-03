<script>
  export let currentPage = 1;
  export let totalPages = 1;
  export let onPrev = () => {};
  export let onNext = () => {};
</script>

<div class="pagination-controls">
  <button class="page-btn" on:click={onPrev} disabled={currentPage === 1}>
    &larr; Previous
  </button>

  <span class="page-info">
    Page <strong>{currentPage}</strong> of <strong>{totalPages}</strong>
  </span>

  <button class="page-btn" on:click={onNext} disabled={currentPage === totalPages}>
    Next &rarr;
  </button>
</div>
