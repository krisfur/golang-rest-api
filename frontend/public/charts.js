import {
    Chart,
    CategoryScale,
    LinearScale,
    PointElement,
    ScatterController,
    Title,
    Tooltip,
    Legend
} from 'https://esm.sh/chart.js';

const clusterColors = [
    '#f38ba8',
    '#a6e3a1',
    '#89b4fa',
    '#f9e2af',
    '#cba6f7',
    '#94e2d5',
];

Chart.register(
    CategoryScale,
    LinearScale,
    PointElement,
    ScatterController,
    Title,
    Tooltip,
    Legend
);

export function renderCharts(data) {
    const chartsContainer = document.getElementById('charts');
    chartsContainer.innerHTML = '';

    data.forEach((item, idx) => {
        // Create a wrapper div to make canvas responsive in grid
        const chartWrapper = document.createElement('div');
        chartWrapper.style.position = 'relative';
        chartWrapper.style.width = '100%';
        chartWrapper.style.height = '250px';
        chartsContainer.appendChild(chartWrapper);

        const canvas = document.createElement('canvas');
        chartWrapper.appendChild(canvas);

        // Group points by cluster
        const clusters = {};
        item.result.points.forEach(p => {
            if (!clusters[p.cluster]) {
                clusters[p.cluster] = [];
            }
            clusters[p.cluster].push({ x: p.x, y: p.y });
        });

        const datasets = Object.keys(clusters).map(clusterId => ({
            label: `Cluster ${parseInt(clusterId) + 1}`,
            data: clusters[clusterId],
            backgroundColor: clusterColors[clusterId % clusterColors.length],
            pointRadius: 3,
        }));

        new Chart(canvas.getContext('2d'), {
            type: 'scatter',
            data: { datasets },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        labels: {
                            color: '#cdd6f4'
                        }
                    },
                    title: {
                        display: true,
                        text: `Dataset ${idx + 1}`,
                        color: '#cdd6f4'
                    }
                },
                scales: {
                    x: {
                        title: {
                            display: true,
                            text: 'X',
                            color: '#cdd6f4'
                        },
                        ticks: {
                            color: '#cdd6f4'
                        }
                    },
                    y: {
                        title: {
                            display: true,
                            text: 'Y',
                            color: '#cdd6f4'
                        },
                        ticks: {
                            color: '#cdd6f4'
                        }
                    }
                }
            }
        });
    });
}