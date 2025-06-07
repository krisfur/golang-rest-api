# Concurrent KMeans Lab

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)    [![Go](https://img.shields.io/badge/Go-1.24.4-blue)](https://go.dev/)   ![Node.js](https://img.shields.io/badge/Node.js-339933?style=flat&logo=node.js&logoColor=white) [![Node.js](https://img.shields.io/badge/Node.js-22.15.1-brightgreen)](https://nodejs.org/)

An interactive web app that demonstrates **concurrent KMeans clustering** in Go.  

![Screenshot](screenshot.png)

## 🚀 Features
- **Go Backend:** Concurrent clustering of datasets using goroutines
- **Frontend:** Node.js with vanilla JavaScript and Chart.js
- **Responsive Layout:** Charts scale nicely to different screen sizes
- **Catppuccin Mocha Theme:** Visually appealing, dark-mode friendly

## ⚙️ How it works
- 9 unique datasets generated randomly
- Dynamic selection of cluster count (K)
- Run clustering on existing data or generate new data with a single click
- Visual comparison of clustering results in a 3x3 responsive grid

## 🛠️ Getting Started

1. Clone the repo:
    ```bash
    git clone https://github.com/YOUR_USERNAME/concurrent-kmeans-lab.git
    cd concurrent-kmeans-lab
    ```

2. Build and run the Go backend:
    ```bash
    cd api
    go run main.go
    ```

3. Start the frontend:
    ```bash
    cd frontend
    npm install
    npm start
    ```

4. Visit [http://localhost:3000](http://localhost:3000) in your browser.

## 📝 License

MIT License.