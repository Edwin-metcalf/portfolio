import Chart, { type ChartConfiguration } from 'chart.js/auto'

export interface WinLoseData {
    wins: number;
    lose: number;
    winRate: number
}

export function createWinLoseChart(canvas: HTMLCanvasElement, data: WinLoseData): Chart {
    const config: ChartConfiguration<'doughnut'> = {   
        type: 'doughnut',
        data: {
        labels: [
            'Wins',
            'Loses'
            ],
        datasets: [{
            data: [data.wins, data.lose],
            backgroundColor: [
                'rgb(24,255,44)',
                'rgb(255, 15, 60)'
            ],
            hoverOffset: 5
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'bottom'
                }
            }
        }
    };
    return new Chart(canvas, config);
}