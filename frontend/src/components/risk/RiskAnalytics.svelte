<script>
  import { onMount } from 'svelte';
  
  export let assessments = [];
  let chartCanvas;
  let chart;
  let summary = {
    avgEmotional: 0,
    avgFOMO: 0,
    avgBias: 0,
    avgOverall: 0,
    lastWeek: []
  };
  
  $: {
    if (assessments && assessments.length) {
      calculateSummary();
      if (chartCanvas) {
        renderChart();
      }
    }
  }
  
  onMount(() => {
    loadChartLibrary().then(() => {
      renderChart();
    });
  });
  
  function loadChartLibrary() {
    return new Promise((resolve) => {
      if (window.Chart) {
        resolve();
        return;
      }
      
      const script = document.createElement('script');
      script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
      script.onload = () => resolve();
      document.head.appendChild(script);
    });
  }
  
  function calculateSummary() {
    // Calculate overall averages
    let totalEmotional = 0;
    let totalFOMO = 0;
    let totalBias = 0;
    let totalOverall = 0;
    
    assessments.forEach(a => {
      totalEmotional += a.emotionalScore;
      totalFOMO += a.fomoScore;
      totalBias += a.biasScore;
      totalOverall += a.overallScore;
    });
    
    const count = assessments.length;
    summary.avgEmotional = Math.round((totalEmotional / count) * 10) / 10;
    summary.avgFOMO = Math.round((totalFOMO / count) * 10) / 10;
    summary.avgBias = Math.round((totalBias / count) * 10) / 10;
    summary.avgOverall = Math.round((totalOverall / count) * 10) / 10;
    
    // Calculate last 7 days summary
    const oneWeekAgo = new Date();
    oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
    
    summary.lastWeek = assessments.filter(a => {
      const assessmentDate = new Date(a.date);
      return assessmentDate >= oneWeekAgo;
    });
    
    // Set default last week stats if no entries
    if (summary.lastWeek.length === 0) {
      summary.lastWeekEmotional = 'N/A';
      summary.lastWeekFOMO = 'N/A';
      summary.lastWeekBias = 'N/A';
    } else {
      let lwEmotional = 0, lwFOMO = 0, lwBias = 0;
      
      summary.lastWeek.forEach(a => {
        lwEmotional += a.emotionalScore;
        lwFOMO += a.fomoScore;
        lwBias += a.biasScore;
      });
      
      const lwCount = summary.lastWeek.length;
      summary.lastWeekEmotional = Math.round((lwEmotional / lwCount) * 10) / 10;
      summary.lastWeekFOMO = Math.round((lwFOMO / lwCount) * 10) / 10;
      summary.lastWeekBias = Math.round((lwBias / lwCount) * 10) / 10;
    }
  }
  
  function renderChart() {
    if (!window.Chart || !chartCanvas || assessments.length === 0) return;
    
    // Prepare data for chart
    const dates = assessments.map(a => {
      const date = new Date(a.date);
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    });
    
    const overallScores = assessments.map(a => a.overallScore);
    const emotionalScores = assessments.map(a => a.emotionalScore);
    const fomoScores = assessments.map(a => a.fomoScore);
    const biasScores = assessments.map(a => a.biasScore);
    
    // Clean up previous chart if it exists
    if (chart) {
      chart.destroy();
    }
    
    // Create the chart
    chart = new window.Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: dates,
        datasets: [
          {
            label: 'Overall',
            data: overallScores,
            borderColor: 'rgb(54, 162, 235)',
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            tension: 0.1,
            fill: false
          },
          {
            label: 'Emotional',
            data: emotionalScores,
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.2)',
            tension: 0.1,
            fill: false,
            hidden: true
          },
          {
            label: 'FOMO',
            data: fomoScores,
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.2)',
            tension: 0.1,
            fill: false,
            hidden: true
          },
          {
            label: 'Bias',
            data: biasScores,
            borderColor: 'rgb(255, 159, 64)',
            backgroundColor: 'rgba(255, 159, 64, 0.2)',
            tension: 0.1,
            fill: false,
            hidden: true
          }
        ]
      },
      options: {
        responsive: true,
        scales: {
          y: {
            beginAtZero: true,
            max: 10
          }
        },
        plugins: {
          title: {
            display: true,
            text: 'Risk Assessment Trends'
          }
        }
      }
    });
  }
</script>

<div class="analytics">
  <h2>Risk Analytics</h2>
  
  {#if assessments.length === 0}
    <p class="no-data">No risk assessments yet. Add some assessments to see analytics.</p>
  {:else}
    <div class="chart-container">
      <canvas bind:this={chartCanvas}></canvas>
    </div>
    
    <div class="stats">
      <div class="stats-section">
        <h3>Overall Averages</h3>
        <div class="stat-grid">
          <div class="stat-item">
            <span class="stat-label">Emotional:</span>
            <span class="stat-value">{summary.avgEmotional}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">FOMO:</span>
            <span class="stat-value">{summary.avgFOMO}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Bias:</span>
            <span class="stat-value">{summary.avgBias}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Overall:</span>
            <span class="stat-value">{summary.avgOverall}</span>
          </div>
        </div>
      </div>
      
      <div class="stats-section">
        <h3>Last 7 Days</h3>
        <div class="stat-grid">
          <div class="stat-item">
            <span class="stat-label">Emotional:</span>
            <span class="stat-value">{summary.lastWeekEmotional}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">FOMO:</span>
            <span class="stat-value">{summary.lastWeekFOMO}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Bias:</span>
            <span class="stat-value">{summary.lastWeekBias}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Entries:</span>
            <span class="stat-value">{summary.lastWeek.length}</span>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .analytics {
    max-width: 800px;
    margin: 0 auto;
    padding: 1rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 1.5rem;
  }
  
  .no-data {
    text-align: center;
    font-style: italic;
    color: #666;
    padding: 2rem;
  }
  
  .chart-container {
    width: 100%;
    height: 300px;
    margin-bottom: 2rem;
  }
  
  .stats {
    display: flex;
    flex-wrap: wrap;
    gap: 2rem;
    justify-content: space-around;
  }
  
  .stats-section {
    flex: 1;
    min-width: 250px;
    padding: 1rem;
    background-color: #f9f9f9;
    border-radius: 8px;
  }
  
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    text-align: center;
    font-size: 1.2rem;
  }
  
  .stat-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .stat-label {
    font-weight: bold;
    color: #555;
  }
  
  .stat-value {
    font-size: 1.5rem;
    font-weight: bold;
  }
</style> 