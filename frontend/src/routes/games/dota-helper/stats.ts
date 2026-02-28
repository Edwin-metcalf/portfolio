import Chart, { type ChartConfiguration } from 'chart.js/auto'
//need to change this to take in the correct now overall and recent data
export interface HeroStats {
    wins: number;
    losses: number;
    games: number;
    kills: number;
    deaths: number;
    assists: number;
}
export interface DotaStatsReturn {
    wins: number;
    losses: number;
    winRate: number;
    recentWins: number;
    recentLosses: number;
    recentWinRate: number;
    avgKDA: number[];
    recentHeroStats: Record<string, HeroStats>
}

export function createWinLoseChart(canvas: HTMLCanvasElement, wins: number, losses: number): Chart {
    const config: ChartConfiguration<'doughnut'> = {   
        type: 'doughnut',
        data: {
        labels: [
            'Wins',
            'Loses'
            ],
        datasets: [{
            data: [wins, losses],
            backgroundColor: [
                'rgb(34, 203, 0)',
                'rgb(195, 0, 0)'
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