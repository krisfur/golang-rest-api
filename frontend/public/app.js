import { renderCharts } from '/static/charts.js';

document.addEventListener('DOMContentLoaded', () => {
    const fetchDataBtn = document.getElementById('fetchDataBtn');
    const generateDataBtn = document.getElementById('generateDataBtn');
    const kSelector = document.getElementById('k-selector');
    const chartsContainer = document.getElementById('charts');

    async function fetchData() {
        chartsContainer.innerHTML = '<p>Loading data...</p>';
        const k = kSelector.value;

        try {
            const response = await fetch(`http://localhost:8080/aggregate?k=${k}`);
            if (!response.ok) throw new Error('API error: ' + response.status);
            const data = await response.json();
            renderCharts(data);
        } catch (error) {
            console.error(error);
            chartsContainer.innerHTML = '<p style="color: var(--red)">Error fetching data. Is the Go API running?</p>';
        }
    }

    async function generateData() {
        try {
            const response = await fetch(`http://localhost:8080/generate`, { method: 'POST' });
            if (!response.ok) throw new Error('API error: ' + response.status);
            // Immediately fetch data using the current K value
            fetchData();
        } catch (error) {
            console.error(error);
            alert("Error generating new data.");
        }
    }

    fetchDataBtn.addEventListener('click', fetchData);
    generateDataBtn.addEventListener('click', generateData);

    // Allow pressing Enter on the input field
    kSelector.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            fetchData();
        }
    });
});