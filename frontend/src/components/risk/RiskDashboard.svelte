<script>
  import { onMount } from 'svelte';
  import RiskForm from './RiskForm.svelte';
  import RiskAnalytics from './RiskAnalytics.svelte';
  
  let assessments = [];
  let index = 0;
  let current = {};
  let activeTab = 'form'; // 'form' or 'analytics'
  
  onMount(async () => {
    try {
      const result = await window.go.main.App.GetRiskAssessments();
      assessments = result || [];
      if (assessments.length) {
        loadAt(assessments.length - 1); // Load most recent by default
      } else {
        current = {
          id: '',
          date: new Date().toISOString().split('T')[0],
          emotionalScore: 5,
          fomoScore: 5,
          biasScore: 5,
          overallScore: 5,
          notes: ''
        };
      }
    } catch (error) {
      console.error("Failed to load risk assessments:", error);
    }
  });
  
  function loadAt(i) {
    if (i >= 0 && i < assessments.length) {
      index = i;
      current = {...assessments[i]};
      // Format date for input type="date" if needed
      if (current.date && typeof current.date === 'string' && !current.date.includes('-')) {
        const date = new Date(current.date);
        current.date = date.toISOString().split('T')[0];
      }
    }
  }
  
  function next() {
    if (index < assessments.length - 1) loadAt(index + 1);
  }
  
  function prev() {
    if (index > 0) loadAt(index - 1);
  }
  
  async function saveAssessment(event) {
    try {
      const assessment = event.detail;
      await window.go.main.App.SaveRiskAssessment(assessment);
      const result = await window.go.main.App.GetRiskAssessments();
      assessments = result || [];
      
      // Find the index of the saved assessment to keep position
      const savedIndex = assessments.findIndex(a => a.id === assessment.id);
      if (savedIndex >= 0) {
        loadAt(savedIndex);
      } else {
        loadAt(assessments.length - 1); // Load the latest if it's a new one
      }
      
      showToast("Assessment saved successfully!");
    } catch (error) {
      console.error("Failed to save assessment:", error);
      showToast("Error saving assessment", true);
    }
  }
  
  function showToast(message, isError = false) {
    const toast = document.createElement('div');
    toast.className = `toast ${isError ? 'error' : 'success'}`;
    toast.textContent = message;
    document.body.appendChild(toast);
    
    setTimeout(() => {
      toast.classList.add('show');
      setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
          document.body.removeChild(toast);
        }, 300);
      }, 3000);
    }, 100);
  }
</script>

<div class="dashboard">
  <div class="tabs">
    <button 
      class:active={activeTab === 'form'} 
      on:click={() => activeTab = 'form'}
    >
      Risk Assessment
    </button>
    <button 
      class:active={activeTab === 'analytics'} 
      on:click={() => activeTab = 'analytics'}
    >
      Analytics
    </button>
  </div>
  
  {#if activeTab === 'form'}
    <div class="nav-arrows">
      <button 
        on:click={prev} 
        disabled={index === 0 || assessments.length === 0}
        class="nav-button"
      >
        ‹ Previous
      </button>
      <span class="date-display">
        {assessments.length ? `${index + 1} of ${assessments.length}` : 'New Entry'}
      </span>
      <button 
        on:click={next} 
        disabled={index === assessments.length - 1 || assessments.length === 0}
        class="nav-button"
      >
        Next ›
      </button>
    </div>
    
    <RiskForm 
      bind:data={current} 
      on:save={saveAssessment} 
    />
  {:else if activeTab === 'analytics'}
    <RiskAnalytics {assessments} />
  {/if}
</div>

<style>
  .dashboard {
    margin: 0 auto;
    padding: 1rem;
  }
  
  .tabs {
    display: flex;
    margin-bottom: 1rem;
    border-bottom: 1px solid #ddd;
  }
  
  .tabs button {
    padding: 0.75rem 1.5rem;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
  }
  
  .tabs button.active {
    font-weight: bold;
    border-bottom: 3px solid #4caf50;
  }
  
  .nav-arrows {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .nav-button {
    background-color: #f0f0f0;
    border: none;
    padding: 0.5rem 1rem;
    cursor: pointer;
    border-radius: 4px;
  }
  
  .nav-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .date-display {
    font-weight: bold;
  }
  
  :global(.toast) {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    padding: 1rem 2rem;
    background-color: #4caf50;
    color: white;
    border-radius: 4px;
    opacity: 0;
    transition: opacity 0.3s;
    z-index: 1000;
  }
  
  :global(.toast.error) {
    background-color: #f44336;
  }
  
  :global(.toast.show) {
    opacity: 1;
  }
</style> 