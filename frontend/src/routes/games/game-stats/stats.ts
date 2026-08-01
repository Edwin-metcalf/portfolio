import Chart, { type ChartConfiguration } from 'chart.js/auto'
import 'chartjs-adapter-date-fns'
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
interface Clan {
    tag: string;
    name: string;
}
interface Arena {
    id: number;
    name: string;
}
interface PlayerProfile {
    tag: string;
    name: string;
    trophies: number;
    bestTrophies: number;
    wins: number;
    losses: number;
    winRate: number;
    battleCount: number;
    threeCrownWins: number;
    clan: Clan;
    arena: Arena;
}
interface CRLadderChartPoint {
    battleTime: string;
    trophies: number;
}
export interface ClashRoyaleLoadReturn {
    profile: PlayerProfile;
    battleLog: CRLadderChartPoint[];
    friendly: ClashRoyaleFriendlyLoadReturn
}
export interface ClashRoyaleStatsReturn {
    wins: number;
    losses: number;
    winRate: number;
    currentTrophies: number;
    trophyProgress: number[];
}
//friendly game stuff
interface CRFriendlyGame {
    battleTime: string;
    result: number;
    myDeck: unknown; //this is not implemented yet but ready to be
    enemyDeck: unknown;
}
export interface ClashRoyaleFriendlyLoadReturn {
    myTag: string;
    friendTag: string;
    wins: number;
    losses: number;
    ties: number;
    winRate: number;
    games: CRFriendlyGame[];

}
//typing for the date object in the line chart
export type TrophyPoint = {
    x: Date;
    y: number;
};
export interface matchupGeneratorTags {
    tag1: string;
    tag2: string;
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

function parseClashRoyaleDate(text: string): Date {
    //need this function because clash royale sends back a wierd date format
    const formatted = text.replace(
        /^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})(\d{2})/,
        '$1-$2-$3T$4:$5:$6'
    );
    return new Date(formatted);
}
export function createTrophyLineChart(canvas: HTMLCanvasElement, trophyData: CRLadderChartPoint[]): Chart<'line', TrophyPoint[]> {
    console.log('raw trohpydata 1: ', trophyData[0])
    const pointData: TrophyPoint[] = trophyData.map(entry => ({
        x: parseClashRoyaleDate(entry.battleTime),
        y: entry.trophies
    }));

    console.log('parsed pointdata 1: ', pointData[0])

    const config: ChartConfiguration<'line', TrophyPoint[]> = {
        type: 'line',
        data: {
            //labels: "trophies",
            datasets: [{
                label: 'Trophy Progression',
                data: pointData,
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
            },
            scales: {
                x: {
                    type: 'time'
                }
            }
        }

    };

    return new Chart(canvas, config);
}

//would be cool to pass names into this for future ability to have a matchup calculator for anyone
export function createHeadToHeadChart(canvas:  HTMLCanvasElement, wins: number, losses: number): Chart<'bar'> {
    const winPct = (wins / (wins + losses)) * 100
    const lossPct = (losses / (wins + losses)) * 100
    const counts = [wins, losses]
    const config: ChartConfiguration<'bar'> = {
        type: 'bar',
        data: {
            labels: ['My Wins'],
            datasets: [
                {
                    label: 'Wins',
                    data: [winPct],
                    backgroundColor: 'rgba(75, 192, 192, 0.85)',
                },
                {
                    label: 'Losses',
                    data: [lossPct],
                    backgroundColor: 'rgba(255, 99, 132, 0.85)'
                }
            ]
        }, 
        options: {
            indexAxis: 'y',
            scales: {
                x: {
                    stacked: true,
                    max: 100, 
                    ticks: {
                        callback: (value) => value + '%'
                    }
                },
                y: {
                    stacked: true
                }
            },
            plugins: {
                tooltip: {
                    callbacks: {
                        label: (context) => {
                            const dataset = context.dataset;
                            const dataIndex = context.dataIndex;

                            const label = dataset.label;
                            const percentage = Number(context.raw).toFixed(1);
                            const count = counts[context.dataIndex];

                            return `${label}: ${percentage}% (${count} ${label})`;
                        }
                    }
                }
            }
        }
    }
    return new Chart(canvas, config)
}