# PC Asset Tracker

A comprehensive IT asset management solution designed for automated hardware inventory and centralized monitoring. This project features a high-performance collection agent built with Go and a modern web-based dashboard developed with Svelte.

## 📊 Live Project Links

| Component | URL |
| :--- | :--- |
| **Interactive Dashboard** | [pc-tracker-ui.vercel.app](https://pc-tracker-ui.vercel.app) |
| **Google Looker Studio Report** | [View Analytics Report](https://lookerstudio.google.com/s/r480m8LFc4w) |
| **Raw Data Source** | [Google Sheets Sample Data](https://docs.google.com/spreadsheets/d/1WK_dJKXyquSS0xivuej8Ic8WxPDBLpYsSvUVVFP3eM8/edit?usp=sharing) |

---

## 🏗️ Architecture & Data Flow

```mermaid
graph TD
    subgraph Client [Client Node]
        Agent[Go Agent\n collector-gui.exe]
        LocalDB[(SQLite Local Hash)]
        Config[config.json\n Asset Tags]
        
        Agent <-->|1. Check for Hardware Changes| LocalDB
        Config -->|2. Append Location Data| Agent
    end

    subgraph Cloud [Cloud Database Layer]
        API[Google Sheets API v4]
        Sheets[(Google Sheets\n Master Database)]
        
        Agent -->|3. Secure HTTPS POST| API
        API --> Sheets
    end

    subgraph Presentation [Presentation & Analytics Layer]
        Vercel[Svelte Dashboard\n Hosted on Vercel]
        Looker[Google Looker Studio\n Enterprise Analytics]
        
        Sheets -->|4a. Live Data Fetch| Vercel
        Sheets -->|4b. Native Data Connector| Looker
    end

    classDef default fill:#1e293b,stroke:#3b82f6,stroke-width:2px,color:#f8fafc;
    classDef database fill:#0f172a,stroke:#10b981,stroke-width:2px,color:#f8fafc;
    class LocalDB,Sheets database;
```

---

## 🛠️ Technical Stack

* **Backend & Agent:** Go (Golang) using the [Wails](https://wails.io/) framework.
* **Frontend:** Svelte with Tailwind CSS for the dashboard and agent UI.
* **Database:** SQLite for local state management and delta-tracking.
* **Cloud Integration:** Google Sheets API v4 via Service Account authentication.
* **Analytics:** Google Looker Studio for enterprise-grade visualization.
* **Hosting:** Vercel for the web dashboard.

---

## ✨ Key Features

* **Deep Hardware Scanning:** Automatically retrieves OS details, CPU model, RAM capacity/modules, Storage drive info, and BIOS serial numbers.
* **Silent Background Mode:** Can be executed with a `--silent` flag to perform scans invisibly, ideal for deployment via Windows Task Scheduler.
* **Intelligent Syncing:** Uses SHA-256 hashing to compare current system state against a local database. Data is only uploaded to the cloud if a hardware change is detected, minimizing API traffic.
* **Asset Tagging:** Allows administrators to assign custom metadata (Department, Location, and Type) that persists across scans.
* **Two-Column Dashboard:** A high-visibility interface designed for IT managers to audit systems at a glance.

---

## ⚙️ How It Works

1.  **Local Execution:** The agent runs on the target PC (either manually or as a background service).
2.  **Hardware Audit:** The application queries the system kernel and WMI (Windows Management Instrumentation) for detailed specs.
3.  **Delta Check:** The results are hashed and compared to the last entry in a local SQLite database.
4.  **Cloud Update:** If changes are found, the agent pushes a new row to the master Google Sheet.
5.  **Visualization:** The Svelte web app and Looker Studio report provide real-time insights into the entire fleet of hardware.

---

## 📦 Installation & Build

### Prerequisites
* Go 1.21+
* Node.js & npm
* Wails CLI

### Build Instructions

```bash
# Clone the repository
git clone [https://github.com/zhengheetong/pc-asset-tracker.git](https://github.com/zhengheetong/pc-asset-tracker.git)

# Build the Collector Agent
cd collector-gui
wails build -clean

# Build/Deploy the Dashboard
cd dashboard
npm install
npm run build
```

---

## 🔐 Security & Configuration

To enable the cloud sync feature, a `service-account.json` file from the Google Cloud Console must be placed in the application root directory. This file contains the credentials required to authenticate with the Google Sheets API.

**Note:** Ensure `service-account.json` and `config.json` are added to your `.gitignore` to prevent leaking sensitive credentials or local environment settings.

---
### 👨‍💻 Credits
* **Lead Developer / Architect:** [zhengheetong](https://github.com/zhengheetong)
* **AI Pair Programmer:** Google Gemini