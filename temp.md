# Building a High-Performance Demand-Side Platform: Architecture and Implementation  

## Introduction  

In the digital advertising ecosystem, Demand-Side Platforms (DSPs) play a crucial role by enabling advertisers to bid on ad inventory in real time. This article outlines the architecture and technical considerations behind building a scalable DSP capable of processing **3 billion bid requests daily** (peaking at **40,000 requests per second**) while maintaining sub-100ms response times.  

## System Requirements  

The primary technical objectives included:  

- **High Throughput**: Handle **40,000 OpenRTB v2.5 bid requests per second** from multiple ad exchanges.  
- **Low Latency**: Respond with optimized ads in **under 100 milliseconds**.  
- **Frequency Capping**: Enforce server-side impression limits per user.  
- **Dynamic Supply-Side Handling**: Unlike Supply-Side Platforms (SSPs), DSPs lack prior knowledge of publishers, requiring real-time traffic source differentiation for bid optimization.  

## Infrastructure Design  

### Core Components  

To meet these demands, the system was structured as follows:  

1. **Ad Servers (x3)**:  
   - Load-balanced to distribute incoming bid requests.  
   - Horizontally scalable—additional servers can be added for higher traffic.  

2. **Redis Caching Service**:  
   - Centralized store for **real-time frequency capping checks**.  
   - Capable of **100,000+ reads/writes per second**, ensuring the 40k RPS target is met.  

3. **MySQL Database**:  
   - Persistent storage for campaign metadata, audience targeting rules, and ledger data.  

4. **Log Aggregator & Analytics Server**:  
   - Collects and processes bid, win, and click logs for reporting and optimization.  

### Scaling Strategy  

- **Bottleneck Analysis**: Frequency capping checks in Redis were identified as the primary constraint.  
- **User-Based Sharding**: For higher traffic, requests can be partitioned across sub-clusters using **user IDs** to distribute Redis load.  

---

## Key Functional Modules  

### 1. Audience Targeting  

Advertisers define campaigns using OpenRTB v2.5-compliant targeting segments:  

- **Geographic**: Country, city, ISP, connection speed.  
- **Demographic**: Age, gender, language.  
- **Device**: Manufacturer, OS, browser.  
- **Temporal**: Day of week, hour of day.  
- **Content Categories**: Whitelists/blacklists for brand safety.  

### 2. Frequency Capping  

To prevent ad fatigue, the system enforces:  

- **Impression Caps**: Maximum views per user within a time window.  
- **Click Caps**: Limits on clicks (though these are prioritized due to their high value).  

#### Data Structures  

```go
type Cap struct {
    CapNumber    uint8  // Max impressions allowed
    CapPeriod    uint16 // Time window (minutes)
    CapThrottle  uint16 // Minimum delay between impressions (minutes)
    ClickNumber  uint8  // Max clicks allowed
    ClickPeriod  uint16 // Click tracking window (minutes)
}

type Fcap struct {
    Total     uint8  // Total impressions/clicks served
    StartYM   uint8  // Start timestamp (year-month)
    StartDHM  uint16 // Start timestamp (day-hour-minute)
    Last      uint16 // Minutes since last impression
}

type BothCap struct {
    Imp Fcap // Impression tracking
    Cli Fcap // Click tracking
}
```  

- **Redis Integration**: Stores `BothCap` records per user for real-time checks.  

### 3. Campaign Management & Cache Updates  

- **Challenge**: Database queries are too slow for real-time bidding.  
- **Solution**: A **NATS-based pub/sub system** synchronizes campaign data across servers every **10 minutes**.  
  - Improves upon flat-file caching (used in prior projects with 30-minute refresh cycles).  
  - Also used to forward logs (bids, wins, clicks) to the aggregator.  

---

## System Workflow  

1. **Bid Request Handling**:  
   - Received by `unify` (HTTPS service) on ad servers.  
   - Audience targeting and frequency checks are applied.  

2. **Cache Synchronization**:  
   - `spread` services on each node update local caches via NATS.  
   - `redis-cache` publishes updates every 10 minutes.  

3. **Log Processing**:  
   - `nats-client` on the aggregator collects logs.  
   - `ledger` aggregates data into MySQL for reporting.  

---

## Performance Outcomes  

- **Throughput**: Sustained **40,000 RPS** with room to scale via sharding.  
- **Latency**: Median response time **<100ms**.  
- **Cache Refresh**: Campaign updates propagate in **10 minutes** (vs. 30 minutes in legacy systems).  

## Lessons Learned  

1. **Bottlenecks Shift with Scale**: At high RPS, frequency capping dominates performance.  
2. **Decoupling is Critical**: Separating caching (Redis), messaging (NATS), and persistence (MySQL) ensures scalability.  
3. **Standards Matter**: OpenRTB compliance simplified integration with exchanges.  

## Future Enhancements  

- **Machine Learning**: Real-time bid price optimization.  
- **Edge Caching**: Reduce latency for global traffic.  

---

This architecture demonstrates how modern distributed systems principles—combined with careful component selection (Redis, NATS, Go)—can meet the stringent demands of programmatic advertising.
