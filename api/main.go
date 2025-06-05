package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// DataPoint represents a 2D point
type DataPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ClusteredPoint represents a 2D point assigned to a cluster
type ClusteredPoint struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Cluster int     `json:"cluster"`
}

// ClusterResult represents the result of clustering
type ClusterResult struct {
	K          int              `json:"k"`
	Iterations int              `json:"iterations"`
	Points     []ClusteredPoint `json:"points"`
}

// APIResult represents the API response
type APIResult struct {
	Source string        `json:"source"`
	Delay  int           `json:"delay"`
	Result ClusterResult `json:"result"`
}

var (
	currentDatasets [][]DataPoint
	dataMutex       sync.Mutex
)

func main() {
	// Initialize 9 datasets at startup
	numCharts := 9
	currentDatasets = make([][]DataPoint, numCharts)
	for i := range numCharts {
		currentDatasets[i] = generateData(100, rand.Intn(5)+2)
	}

	http.HandleFunc("/aggregate", aggregateHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Go API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	numCharts := 9
	newDatasets := make([][]DataPoint, numCharts)
	for i := 0; i < numCharts; i++ {
		newDatasets[i] = generateData(100, rand.Intn(5)+2)
	}

	dataMutex.Lock()
	currentDatasets = newDatasets
	dataMutex.Unlock()

	w.Write([]byte(`{"status":"new data generated"}`))
}

func aggregateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	kParam := r.URL.Query().Get("k")
	k := 3 // Default
	if parsedK, err := strconv.Atoi(kParam); err == nil {
		k = parsedK
	}
	if k < 2 {
		k = 2
	}
	if k > 6 {
		k = 6
	}

	// Copy current datasets safely
	dataMutex.Lock()
	dataCopy := make([][]DataPoint, len(currentDatasets))
	for i, dataset := range currentDatasets {
		copyDataset := make([]DataPoint, len(dataset))
		copy(copyDataset, dataset)
		dataCopy[i] = copyDataset
	}
	dataMutex.Unlock()

	results := parallelCluster(dataCopy, k)
	json.NewEncoder(w).Encode(results)
}

func parallelCluster(datasets [][]DataPoint, k int) []APIResult {
	var wg sync.WaitGroup
	results := make([]APIResult, len(datasets))

	for i, dataset := range datasets {
		wg.Add(1)
		go func(i int, data []DataPoint) {
			defer wg.Done()
			start := time.Now()
			_, clusteredPoints, iterations := kMeans(data, k, 100)
			delay := int(time.Since(start).Milliseconds())

			results[i] = APIResult{
				Source: "KMeans Clustering",
				Delay:  delay,
				Result: ClusterResult{
					K:          k,
					Iterations: iterations,
					Points:     clusteredPoints,
				},
			}
		}(i, dataset)
	}

	wg.Wait()
	return results
}

func generateData(n, clusters int) []DataPoint {
	data := make([]DataPoint, 0, n)
	for range clusters {
		centerX := rand.Float64() * 100
		centerY := rand.Float64() * 100
		for j := 0; j < n/clusters; j++ {
			x := centerX + rand.NormFloat64()*5
			y := centerY + rand.NormFloat64()*5
			data = append(data, DataPoint{x, y})
		}
	}
	return data
}

func kMeans(data []DataPoint, k, maxIter int) ([]DataPoint, []ClusteredPoint, int) {
	centroids := initializeCentroids(data, k)
	assignments := make([]int, len(data))

	for iter := range maxIter {
		for i, p := range data {
			idx := closestCentroid(p, centroids)
			assignments[i] = idx
		}

		newCentroids := make([]DataPoint, k)
		counts := make([]int, k)
		for i, idx := range assignments {
			newCentroids[idx].X += data[i].X
			newCentroids[idx].Y += data[i].Y
			counts[idx]++
		}

		for i := range k {
			if counts[i] > 0 {
				newCentroids[i].X /= float64(counts[i])
				newCentroids[i].Y /= float64(counts[i])
			} else {
				newCentroids[i] = centroids[i]
			}
		}

		if converged(centroids, newCentroids) {
			clusteredPoints := make([]ClusteredPoint, len(data))
			for i, p := range data {
				clusteredPoints[i] = ClusteredPoint{
					X:       p.X,
					Y:       p.Y,
					Cluster: assignments[i],
				}
			}
			return newCentroids, clusteredPoints, iter + 1
		}

		centroids = newCentroids
	}

	clusteredPoints := make([]ClusteredPoint, len(data))
	for i, p := range data {
		clusteredPoints[i] = ClusteredPoint{
			X:       p.X,
			Y:       p.Y,
			Cluster: assignments[i],
		}
	}

	return centroids, clusteredPoints, maxIter
}

func initializeCentroids(data []DataPoint, k int) []DataPoint {
	centroids := make([]DataPoint, k)
	perm := rand.Perm(len(data))
	for i := range k {
		centroids[i] = data[perm[i]]
	}
	return centroids
}

func closestCentroid(p DataPoint, centroids []DataPoint) int {
	minDist := math.MaxFloat64
	minIdx := 0
	for i, c := range centroids {
		dist := math.Hypot(p.X-c.X, p.Y-c.Y)
		if dist < minDist {
			minDist = dist
			minIdx = i
		}
	}
	return minIdx
}

func converged(a, b []DataPoint) bool {
	const threshold = 1e-3
	for i := range a {
		if math.Abs(a[i].X-b[i].X) > threshold || math.Abs(a[i].Y-b[i].Y) > threshold {
			return false
		}
	}
	return true
}
