Absolutely—let’s bake those navigation and analytics features in from day one. Below is an expanded plan that:

1. **Integrates date-toggle and simple analytics pages** into the Risk and Stock dashboards  
2. **Updates the 7-week schedule** to include these tasks  
3. Provides **new Cucumber scenarios** for date cycling and analytics views  
4. Supplies **targeted pseudocode** showing how you’d wire up toggles, data series, and analytics calculations

---

## 1. Updated 7-Week Development Schedule

| Week | Focus                                        | Deliverables                                                                                                              |
|:----:|:---------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------|
| **1** | Foundation & CI/CD                            | • Wails/Svelte/Badger scaffold<br>• Git workflow & CI (lint, tests)<br>• Logger + error-wrapper library                   |
| **2** | Risk Dashboard Core + Navigation              | • CRUD `RiskAssessment` endpoints<br>• Date-toggle UI (Prev/Next)<br>• “Export/Import JSON” backup<br>• Unit tests        |
| **3** | Risk Analytics & Stock CRUD + Tabs            | • Risk Analytics view (trend chart) + endpoint or frontend calc<br>• CRUD `StockRating`<br>• Day-tab UI for ratings     |
| **4** | Stock Analytics & Trade Calendar Core         | • Stock Analytics view (sector + symbol trends)<br>• Calendar CRUD + filters<br>• Notifications stub                      |
| **5** | Advanced Features & Performance               | • Pagination endpoints<br>• Badger bloom filters & compression toggle<br>• Archive old entries<br>• UI loading/toasts    |
| **6** | Testing & Documentation                       | • Go unit tests (>80% coverage)<br>• Svelte component tests (Vitest)<br>• Inline GoDoc + user guide in `/docs`           |
| **7** | Packaging & Release                           | • Cross-platform installers with code signing<br>• Nightly CI builds & artifact uploads<br>• Final QA on clean VMs        |

---

## 2. Cucumber / Gherkin Scenarios

### 2.1 Risk Management + Navigation + Analytics

```gherkin
Feature: Risk Management Dashboard

  Scenario: Save a new risk assessment and cycle through dates
    Given I open the Risk Management view
    When I set all sliders and click "Save"
    Then a new record "risk_<id>" is stored
    And I see it displayed immediately
    When I click the "Next" arrow
    Then the next saved date’s assessment appears
    When I click the "Previous" arrow
    Then the previous date’s assessment appears

  Scenario: View risk trends analytics
    Given I navigate to the "Risk Analytics" tab
    When the view loads
    Then I see a chart of overallScore vs. date
    And I see average Emotional, FOMO, Bias displayed for the last 7 days
```

### 2.2 Stock Rating + Tabs + Analytics

```gherkin
Feature: Stock Rating Dashboard

  Scenario: Tab through daily ratings
    Given I open the Stock Rating view
    When I click "Next Day"
    Then the ratings for the next date appear
    When I click "Previous Day"
    Then the ratings for the prior date appear

  Scenario: View stock rating trends analytics
    Given I navigate to the "Rating Analytics" tab
    When the view loads
    Then I see a multi-line chart: each sector’s average sentiment over time
    And I can toggle on/off individual symbols to overlay their sentiment trends
```

---

## 3. Pseudocode Examples

### 3.1 Risk Dashboard: Date Toggle & Analytics

```javascript
// src/components/RiskDashboard.svelte
<script>
  import { onMount } from 'svelte';
  let assessments = [], index = 0;
  let current = {};

  onMount(async () => {
    assessments = await window.backend.GetRiskAssessments();  // returns sorted by date asc
    if (assessments.length) loadAt(assessments.length - 1);
  });

  function loadAt(i) {
    index = i;
    current = assessments[i];
  }

  function next() {
    if (index < assessments.length - 1) loadAt(index + 1);
  }
  function prev() {
    if (index > 0) loadAt(index - 1);
  }

  // Analytics data series: overallScore over time
  const analyticsSeries = assessments.map(a => ({
    date: a.date,
    score: a.overallScore
  }));
</script>

<div class="nav-arrows">
  <button on:click={prev} disabled={index===0}>‹</button>
  <button on:click={next} disabled={index===assessments.length-1}>›</button>
</div>

<RiskForm bind:data={current} on:save={async () => {
  await window.backend.SaveRiskAssessment(current);
  assessments = await window.backend.GetRiskAssessments();
}}/>

<!-- Analytics Tab -->
{#if activeTab === 'analytics'}
  <LineChart data={analyticsSeries} x="date" y="score" />
  <SummaryStats items={assessments}/>
{/if}
```

### 3.2 Backend: Risk Analytics Endpoint (Optional)

```go
// pkg/database/database.go
func GetRiskAnalytics() ([]struct{ Date string; AvgScore int }, error) {
  recs, err := GetAllRiskAssessments()
  if err != nil { return nil, err }
  var out []struct{ Date string; AvgScore int }
  for _, r := range recs {
    out = append(out, struct{ Date string; AvgScore int }{
      Date: r.Date, AvgScore: r.OverallScore,
    })
  }
  return out, nil
}
```

### 3.3 Stock Rating: Day Tabs & Trend Data

```javascript
// src/components/StockDashboard.svelte
<script>
  import { onMount } from 'svelte';
  let ratings = [], dateIndex = 0, currentDay = [], analytics=[];

  onMount(async () => {
    ratings = await window.backend.GetStockRatings(); 
    // group by date:
    const byDate = groupBy(ratings, r => r.date);
    dates = Object.keys(byDate).sort();
    loadDate(dates.length - 1);
    analytics = computeSectorTrends(ratings);
  });

  function loadDate(idx) {
    dateIndex = idx;
    currentDay = byDate[dates[idx]];
  }
  function nextDay() { if (dateIndex < dates.length-1) loadDate(dateIndex+1) }
  function prevDay() { if (dateIndex > 0) loadDate(dateIndex-1) }

  function computeSectorTrends(allRatings) {
    // returns [{ sector, series:[{date,avg}] }]
    const sectors = [...new Set(allRatings.map(r=>r.sector))];
    return sectors.map(sec => ({
      sector: sec,
      series: dates.map(d => {
        const slice = allRatings.filter(r=>r.date===d && r.sector===sec);
        const avg = slice.reduce((sum,r)=>sum+r.stockSentiment,0)/slice.length||0;
        return { date: d, avg };
      })
    }));
  }
</script>

<div class="nav-arrows">
  <button on:click={prevDay} disabled={dateIndex===0}>‹</button>
  <button on:click={nextDay} disabled={dateIndex===dates.length-1}>›</button>
</div>
<RatingList items={currentDay}/>

{#if activeTab==='analytics'}
  {#each analytics as line}
    <LineChart data={line.series} x="date" y="avg" label={line.sector}/>
  {/each}
{/if}
```

---

## 4. Next Steps

1. **Week 2 Tasks**:  
   - Wire up Prev/Next in Risk UI and backend stub  
   - Build the Risk Analytics chart component; decide frontend vs. backend data prep  
2. **Week 3 Tasks**:  
   - Add Day Tabs in Stock UI  
   - Implement `computeSectorTrends` (frontend) or `GetStockAnalytics` (backend)  
3. **Continuous**:  
   - Add toast notifications on save/update/delete  
   - Expand unit tests for navigation flows & analytics data  

With these toggles and lightweight analytics views baked in from the start, your dashboard will not only record data but empower traders to **review their journey** and **spot trends** at a glance. Let me know if you’d like more detail on any piece!