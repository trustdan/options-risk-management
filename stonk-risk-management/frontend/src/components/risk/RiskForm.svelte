<script>
  import { createEventDispatcher } from 'svelte';
  
  export let data = {
    id: '',
    date: new Date().toISOString().split('T')[0],
    emotionalScore: 5,
    fomoScore: 5,
    biasScore: 5,
    overallScore: 5,
    notes: ''
  };
  
  const dispatch = createEventDispatcher();
  
  function calculateOverall() {
    data.overallScore = Math.round((data.emotionalScore + data.fomoScore + data.biasScore) / 3);
  }
  
  function handleSubmit() {
    calculateOverall();
    dispatch('save', data);
  }
</script>

<div class="risk-form">
  <h2>Daily Risk Assessment</h2>
  
  <div class="form-group">
    <label for="date">Date</label>
    <input 
      type="date" 
      id="date" 
      bind:value={data.date}
    />
  </div>
  
  <div class="form-group">
    <label for="emotional">Emotional Control (1-10)</label>
    <input 
      type="range" 
      id="emotional" 
      min="1" 
      max="10" 
      bind:value={data.emotionalScore}
      on:change={calculateOverall}
    />
    <span class="score">{data.emotionalScore}</span>
  </div>
  
  <div class="form-group">
    <label for="fomo">FOMO Resistance (1-10)</label>
    <input 
      type="range" 
      id="fomo" 
      min="1" 
      max="10" 
      bind:value={data.fomoScore}
      on:change={calculateOverall}
    />
    <span class="score">{data.fomoScore}</span>
  </div>
  
  <div class="form-group">
    <label for="bias">Cognitive Bias Awareness (1-10)</label>
    <input 
      type="range" 
      id="bias" 
      min="1" 
      max="10" 
      bind:value={data.biasScore}
      on:change={calculateOverall}
    />
    <span class="score">{data.biasScore}</span>
  </div>
  
  <div class="form-group overall">
    <label>Overall Risk Score</label>
    <div class="overall-score">{data.overallScore}</div>
  </div>
  
  <div class="form-group">
    <label for="notes">Notes</label>
    <textarea 
      id="notes" 
      bind:value={data.notes}
      placeholder="Any additional thoughts about your trading mindset today?"
    ></textarea>
  </div>
  
  <button on:click={handleSubmit}>Save Assessment</button>
</div>

<style>
  .risk-form {
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 1.5rem;
  }
  
  .form-group {
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
  }
  
  label {
    margin-bottom: 0.5rem;
    font-weight: bold;
  }
  
  input[type="range"] {
    margin-right: 1rem;
    flex-grow: 1;
  }
  
  .score {
    font-weight: bold;
    font-size: 1.2rem;
  }
  
  .overall {
    margin-top: 2rem;
    text-align: center;
  }
  
  .overall-score {
    font-size: 2rem;
    font-weight: bold;
    padding: 1rem;
    background-color: #f0f0f0;
    border-radius: 50%;
    width: 4rem;
    height: 4rem;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto;
  }
  
  button {
    background-color: #4caf50;
    color: white;
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    margin-top: 1rem;
    display: block;
    width: 100%;
  }
  
  button:hover {
    background-color: #45a049;
  }
  
  textarea {
    min-height: 100px;
    padding: 0.5rem;
  }
</style> 