# Options Trading Dashboard

![image](https://github.com/user-attachments/assets/9a56ba9f-fb22-42a8-b858-21f61950ca3f)


## Trade With Confidence Through Better Psychology

**Options Trading Dashboard** is a powerful desktop application designed to help stock and options traders take control of their trading psychology, manage risk effectively, and make more informed trading decisions.

> "The difference between being a consistently profitable trader and a struggling one often comes down to psychology and risk management." 

## Why We Created This App

Trading is as much about managing your mental state as it is about market analysis. This app was born from the real experiences of traders who know that:

- Emotional states dramatically affect trading decisions
- Position sizing needs to adapt to your mental condition
- FOMO can lead to devastating losses
- Having a structured approach to trading psychology is essential

## Core Features

### 📊 Daily Risk Assessment Dashboard
- **Psychological State Tracking**: Gauge your emotional state, FOMO level, and market bias
- **Smart Position Sizing**: Get data-driven recommendations on position size based on your current mental state
- **Warning Flags**: Automatic alerts for potential euphoria, revenge trading, or FOMO
- **Trading Psychology Guidelines**: Built-in reminders of essential trading psychology principles

![Risk Assessment](https://placeholder-for-risk-screenshot.png)

### 📝 Trading Journal
- Document trades with emotional context
- Track win/loss ratios and P&L
- Reflect on lessons learned
- Build self-awareness by connecting emotions to trading outcomes

### 📅 Options Trading Calendar
- Visualize your options strategies by expiration week
- Color-coded strategy organization
- Multi-sector trade planning
- Comprehensive options strategy selection

### 📈 Stock Rating Dashboard
- Market and sector sentiment tracking
- Individual stock ratings
- Chart pattern identification
- Historical ratings view

## Getting Started

1. **Download and Install**: Grab the latest version from the releases page
2. **Daily Assessment**: Start each trading day with a quick psychological self-assessment
3. **Follow the Guidance**: Use the position sizing recommendations to structure your trades
4. **Journal and Reflect**: Record your trades and emotions to build self-awareness
5. **Improve Over Time**: Use the analytics to track your psychological progress

## Key Benefits

- **Trade with appropriate size**: Avoid overtrading on your worst days
- **Recognize dangerous patterns**: Catch yourself before making emotional mistakes
- **Build consistency**: Develop a structured approach to trading psychology
- **Trade more confidently**: Know that your decisions are backed by self-awareness
- **Protect your capital**: The first rule of trading is to preserve your ability to trade tomorrow

---

## Technical Details and Development Information

### Installation

#### Windows
1. Download the latest installer (`options-trading-dashboard-setup-x.x.x.exe`) from the releases page
2. Run the installer and follow the on-screen instructions
3. Launch the application from the Start menu or desktop shortcut

#### Building from Source
To build the application from source, you'll need:
- Go 1.18 or higher
- Node.js 16 or higher
- Wails (https://wails.io)

Clone the repository and run:
```
wails build
```

### Data Storage
All your data is stored locally in `~/.options-trading-dashboard` directory.

### Supported Options Strategies
The application supports a comprehensive range of options trading strategies, including:

#### Basic Spreads
- Long Call/Put - Directional bets on market movement
- Covered Call - Owning stock and selling calls against it for premium income

#### Vertical Spreads
- Bull Call Spread - Bullish strategy with limited risk/reward
- Bear Call Spread - Bearish strategy with limited risk/reward
- Bull Put Spread - Bullish strategy with limited risk/reward
- Bear Put Spread - Bearish strategy with limited risk/reward

#### Advanced Strategies
- Calendar/Horizontal Spreads - Time-based spreads with different expiration dates
- Diagonal Spreads - Combining different strikes and expiration dates
- Butterfly Spreads - Neutral strategies for limited price movement
- Iron Condors/Butterflies - Range-bound strategies for selling premium
- Ratio Spreads (Backspreads) - Unbalanced spreads with more long than short options

#### Danger - naked options ahead
- Short Call - Bearish directional bet with unlimited risk potential
- Short Put - Bullish directional bet with significant downside risk
- Cash-Secured Put - Selling puts with cash to secure the potential stock purchase (risk of capital misallocation)
- Call Ratio Spread - Neutral to moderately bearish with naked call exposure (more short than long options)
- Put Ratio Spread - Neutral to moderately bullish with naked put exposure (more short than long options)

### Development

#### Prerequisites
- Go 1.18+
- Node.js 16+
- Wails CLI
- Svelte knowledge

#### Dev Environment
Run the development server with:
```
wails dev
```

### Common Issues & Solutions

#### Date Handling in JavaScript/Go Applications
When working with dates between JavaScript frontend and Go backend:

1. **ISO String Format**: Always convert JavaScript Date objects to ISO strings before sending to Go:
   ```javascript
   // Correct way to send dates to Go backend
   const dateStr = new Date(date).toISOString().split('T')[0];
   const formattedDate = `${dateStr}T00:00:00Z`; // Use complete ISO format
   ```

2. **Date Comparison**: When comparing dates, convert both to the same format first:
   ```javascript
   // Convert both dates to YYYY-MM-DD format before comparing
   const date1 = new Date(firstDate).toISOString().split('T')[0];
   const date2 = new Date(secondDate).toISOString().split('T')[0];
   const datesMatch = date1 === date2;
   ```

3. **Using Wails Models**: When using generated model classes, create objects with `createFrom()`:
   ```javascript
   import { models } from '../wailsjs/go/models';
   
   // Create data object
   const ratingData = {
     id: existingRating ? existingRating.id : '',
     date: formattedDate, // ISO string, not Date object
     // other properties...
   };
   
   // Create proper model instance
   const ratingToSave = models.StockRating.createFrom(ratingData);
   
   // Send to backend
   await SaveStockRating(ratingToSave);
   ```

These solutions address common issues with time.Time unmarshalling errors in Go and date comparison mismatches in JavaScript.

## Troubleshooting

If you encounter any issues with the application:

1. **Check the logs**: Logs are stored in `~/.options-trading-dashboard/logs`
2. **Restart the application**: Many issues can be resolved with a simple restart
3. **File an issue**: If problems persist, file an issue on the GitHub repository with:
   - The exact steps to reproduce the problem
   - Any error messages you received
   - Your operating system and app version

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License
MIT

## Credits
Created with ❤️ for traders who understand that psychology is the key to consistent profitability.
