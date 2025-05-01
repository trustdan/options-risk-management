<script>
  export let text = ''; // The text to display in the tooltip
  let visible = false;
</script>

<div 
  class="tooltip-container" 
  on:mouseenter={() => visible = true} 
  on:mouseleave={() => visible = false}
  on:focusin={() => visible = true}  
  on:focusout={() => visible = false} 
>
  <slot></slot> <!-- The element that triggers the tooltip -->
  {#if visible && text}
    <div class="tooltip-box" role="tooltip">
      {@html text.replace(/\n/g, '<br/>')} <!-- Display text, converting newlines -->
    </div>
  {/if}
</div>

<style>
  .tooltip-container {
    display: inline-block; /* or block, depending on context */
  }

  .tooltip-box {
    position: absolute;
    bottom: 110%; /* Position above the trigger element */
    left: 50%;
    transform: translateX(-50%);
    background-color: var(--tooltip-bg, #333); /* Use CSS variable or default */
    color: var(--tooltip-text, #fff);
    padding: 8px 12px;
    border-radius: 4px;
    font-size: 0.85rem;
    white-space: pre-wrap; /* Respect newlines in the text prop */
    z-index: 10; /* Ensure it's above other elements */
    box-shadow: 0 2px 5px rgba(0,0,0,0.2);
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.2s ease-in-out, visibility 0.2s ease-in-out;
    min-width: 150px; /* Adjust as needed */
    text-align: left;
  }

  .tooltip-container:hover .tooltip-box,
  .tooltip-container:focus-within .tooltip-box {
    opacity: 1;
    visibility: visible;
  }
  
  /* Optional: Add an arrow */
  .tooltip-box::after {
    content: '';
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: var(--tooltip-bg, #333) transparent transparent transparent;
  }

  /* Dark mode adjustments if needed */
  :global(body.dark-mode) .tooltip-box {
     --tooltip-bg: #555;
     --tooltip-text: #eee;
  }
</style> 