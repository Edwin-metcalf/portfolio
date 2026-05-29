import Chart, { type ChartConfiguration } from 'chart.js/auto'
//need to change this to take in the correct now overall and recent data
export interface HeroStats {
    wins: number;
    losses: number;
    games: number;
    kills: number;
    deaths: number;
    assists: number;
    avgKDA: number[];
}
export interface MatchupStats {
    enemyHeroName: string;
    wins: number;
    losses: number;
    games: number;
    winRate: number;
}
export interface DotaStatsReturn {
    wins: number;
    losses: number;
    winRate: number;
    recentWins: number;
    recentLosses: number;
    recentWinRate: number;
    avgKDA: number[];
    recentHeroStats: Record<string, HeroStats>;
    matchupStats:   Record<string, MatchupStats>;
}
//stuff for the clash royale
export interface ClashRoyaleStatsReturn {
    wins: number;
    losses: number;
    winRate: number;
    currentTrophies: number;
    trophyProgress: number[];
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

//will have to see what the data comes back as could do it by time
// or could do it like this where its just all games 
export function createTrophyLineChart(canvas: HTMLCanvasElement, trophyData: number[]){
    const maxTrophies = Math.max(...trophyData);

    const labels = [];
    for (let i = 0; i <= maxTrophies; i += 1000){
        labels.push(i)
    }
    const config: ChartConfiguration<'line'> = {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Trophy Progression',
                data: trophyData,
                fill: true,
                tension: 0.1
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