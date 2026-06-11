<script lang="ts">
    import {onMount, tick} from 'svelte';
    import Chart from 'chart.js/auto'
	import { fetchAPI } from '$lib/api';
    import {createWinLoseChart, type DotaStatsReturn, type MatchupStats} from '../stats'
	import { color } from 'chart.js/helpers';



    let chartCanvasOverall: HTMLCanvasElement;
    let chartCanvasRecent: HTMLCanvasElement;

    let chartInstanceOverall: Chart | null  = null;
    let chartInstanceRecent: Chart | null = null;
    let loading: boolean = true;
    let dotaData: DotaStatsReturn | null = null;
    let sortBy: 'winRate' | 'games' = 'winRate';

    async function getDotaWinLose(): Promise<DotaStatsReturn | null> {
        try {
            const result = await fetchAPI('/api/games/dota', {
                method: 'GET'
            });
            return result as DotaStatsReturn;
        } catch (err) {
            console.error('Error fetching Dota stats', err);
            return null;
        }
    }


    onMount(() => {
        (async () => {
            dotaData = await getDotaWinLose();
            if (!dotaData) return;

            
            console.log('Matchup Stats:', dotaData.matchupStats); // DEBUG
            console.log('All data:', dotaData); // DEBUG
            loading = false;
            await tick();
            

            if (chartCanvasOverall) {
                chartInstanceOverall = createWinLoseChart(chartCanvasOverall, dotaData.wins, dotaData.losses)
            }
            if (chartCanvasRecent) {
                chartInstanceRecent = createWinLoseChart(chartCanvasRecent, dotaData.recentWins, dotaData.recentLosses)

            }
        })();
        

        return() => {
            if (chartInstanceOverall) {
                chartInstanceOverall.destroy();
            }
            if(chartInstanceRecent) {
                chartInstanceRecent.destroy();
            }
        };
    });

    function formatPercentage(decimal: number) {
        return "%"+ (decimal*100).toFixed(2);
    }

    function getSortedMatchups(): [string, MatchupStats][] {
        if (!dotaData?.matchupStats) return [];
        
        const matchups = Object.entries(dotaData.matchupStats);

        if (sortBy === 'winRate') {
            return matchups.sort((a,b) => b[1].winRate - a[1].winRate);
        } else {
            const top_games_matchups = matchups.sort((a,b) => b[1].games - a[1].games).slice(0,10)
            return top_games_matchups.sort((a,b) => b[1].winRate - a[1].winRate)

        }
    }
    function getBestMatchups(): [string, MatchupStats][] {
        return getSortedMatchups().slice(0,5);
    }
    function getWorstMatchups(): [string, MatchupStats][] {
        return getSortedMatchups().slice(-5).reverse();
    }

    function getMatchupColor(winRate: number): string {
        if (winRate >= 0.6) return '#22CB00';
        if (winRate >= 0.5) return '#429E9D';
        if (winRate >= 0.4) return '#FFA500';
        return '#C30000';
    }

</script>
<div class="game-stat-helper-page">
    <header class="stats-header">
        <h1>My Dota Stats</h1>
            <button class="tab-btn">
                <a href="./clash-royale">
                                    Clash Royale
                </a>
            </button>
    </header>

        {#if loading}
            <p style="font-size: 1.5rem; color: #fff;">Loading stats...</p>
        {:else if dotaData}
            <div class="stats-grid">
                <section class="win-lose-winrate">
                    <h2 class="section-title">All Time</h2>
                    <div class="stats-container">
                        <div class="stat-card">
                            <div class="stat-label">Wins</div>
                            <div class="stat-value">{dotaData.wins}</div>
                        </div>
                        <div class="stat-card">
                            <div class="stat-label">Loses</div>
                            <div class="stat-value">{dotaData.losses}</div>
                        </div>
                        <div class="stat-card highlight">
                            <div class="stat-label">Win Rate</div>
                            <div class="stat-value">{formatPercentage(dotaData.winRate)}</div>
                        </div>
                    </div>
                    <div class="chart-container">
                        <canvas bind:this={chartCanvasOverall}></canvas>
                    </div>
                </section>

                <section class="win-lose-winrate">
                    <h2 class="section-title">Recent Games</h2>
                    <div class="stats-container">
                        <div class="stat-card">
                            <div class="stat-label">Wins</div>
                            <div class="stat-value">{dotaData.recentWins}</div>
                        </div>
                        <div class="stat-card">
                            <div class="stat-label">Loses</div>
                            <div class="stat-value">{dotaData.recentLosses}</div>
                        </div>
                        <div class="stat-card highlight">
                            <div class="stat-label">Win Rate</div>
                            <div class="stat-value">{formatPercentage(dotaData.recentWinRate)}</div>
                        </div>
                    </div>
                    <div class="chart-container">
                        <canvas bind:this={chartCanvasRecent}></canvas>
                    </div>
                </section>
            </div>
            <!--recent hero sections played-->
            <div class="hero-section">
                <h2 class="section-title">Recent Hero Stats</h2>
                <div class="hero-cards-container">
                    {#each Object.entries(dotaData.recentHeroStats) as [heroName, heroStats]}
                        <div class="hero-card">
                            <h2 class="hero-name">{heroName}</h2>
                            <p>Wins: {heroStats.wins}</p>
                            <p>Losses: {heroStats.losses}</p>
                            <p>
                                Average KDA: {heroStats.avgKDA[0]} / {heroStats.avgKDA[1]} / {heroStats.avgKDA[2]}
                            </p>
                        </div>
                    {/each}
                </div>
            </div>

            <!-- matchup stats area -->

            <div class="matchup-section">
                <h2 class="section-title">Enemy Matchups</h2>

                <div class="matchup-controls">
                    <button class="sort-btn" class:active={sortBy === 'winRate'} on:click={() => sortBy = 'winRate'}>
                        Sort by Win Rate
                    </button>

                    <button class="sort-btn" class:active={sortBy === 'games'} on:click={() => sortBy = 'games'}>
                        Sort by Games Played
                    </button>

                </div>

                <div class="matchup-grid">
                    <div class="matchup-subsection">
                        <h3 class="subsection-title">Best Matchups</h3>

                        <div class="matchup-list">
                            {#each getBestMatchups() as [heroName, stats]}
                                <div class="matchup-card" style="border-left: 4px solid {getMatchupColor(stats.winRate)};">
                                    <div class="matchup-hero-name">{heroName}</div>
                                    <div class="matchup-stats">
                                        <span class="matchup-stat">{stats.wins}W - {stats.losses}L</span>
                                        <span class="matchup-winrate" style="color: {getMatchupColor(stats.winRate)};">
                                            {formatPercentage(stats.winRate)}
                                        </span>
                                    </div>
                                    <div class="matchup-games">{stats.games} games</div>
                                </div>
                            {/each}
                        </div>
                    </div>

                    <div class="matchup-subsection">
                        <h3 class="subsection-title">Worst Matchups</h3>

                        <div class="matchup-list">
                            {#each getWorstMatchups() as [heroName, stats]}
                                <div class="matchup-card" style="border-left: 4px solid {getMatchupColor(stats.winRate)};">
                                    <div class="matchup-hero-name">{heroName}</div>
                                    <div class="matchup-stats">
                                        <span class="matchup-stat">{stats.wins}W - {stats.losses}L</span>
                                        <span class="matchup-winrate" style="color: {getMatchupColor(stats.winRate)};">
                                            {formatPercentage(stats.winRate)}
                                        </span>
                                    </div>
                                    <div class="matchup-games">{stats.games} games</div>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>
            </div>
            
        {:else}
            <h2 style="color: white;">Something broke... Maybe internet issues?? then its not my fault</h2>
        {/if}
</div>


<style>
    :global(body) {
        margin: 0;
        padding: 0;
    }
    .game-stat-helper-page {
        background-color: #181818;
        margin: 0 auto;
        padding: 2rem;
    }
    .tab-btn {
        background: rgba(66, 158, 157, 0.3);
        border: 2px solid #429E9D;
        color: #fff;
        padding: 0.75rem 1.5rem;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.3s ease;
        font-weight: 500;
        font-size: 1rem;
    }

    .tab-btn:hover {
        background: rgba(66, 158, 157, 0.5);
    }

    .stats-header {
        text-align: center;
        margin-bottom: 2rem;
        padding: 1.5rem 0;
        border-bottom: 2px solid rgba(62, 184, 182, 0.3);
    }
    .stats-header h1 {
        font-size: 3rem;
        font-weight: 700;
        color: #30d5c8;
        margin: 0;

    }
    .stats-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 2rem;
        margin: 4rem 0;
    }
    .chart-container {
        position: relative;
        width: 100%;
        height: 250px;
        margin-top: 1rem;
    }
    .section-title {
        font-size: 1.25rem;
        color: #429E9D;
        margin-bottom: 1rem;
        padding-left: 0.5rem;
        border-left: 4px solid #429E9D;
    }
    .stats-container {
        display: flex;
        gap: 1rem;
        justify-content: center;
        flex-wrap: wrap;
    }
    .win-lose-winrate {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        /*background-color: #30d5c8; */
    }
    .stat-card {
        background: #429E9D;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        padding: 1rem 1.5rem;
        text-align: center;
        min-width: 65px;
        flex: 1;
        max-width: 100px;
    }
    .stat-label {
        font-size: 0.875rem;
        color: #fff;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        margin-bottom: 0.5rem;
    }
    .stat-value {
        font-size: 1.2rem;
        font-weight: bold;
        color: #fff;
    }
    .hero-section {
        display: flex;
        justify-content: center;
        gap: .5rem;
        margin-top: 4rem;
        margin-bottom: 4rem;
    }
    .hero-card {
        flex: 1 1 calc(25% - 1rem);
        min-width: 180px;
        max-width: 250px;
        color: #fff;
        background: rgba(66, 158, 157, 0.15);
        border: 1px solid rgba(66, 158, 157, 0.3);
        border-radius: 6px;
        padding: 1rem;
    }
    .matchup-section {
        margin: 4rem 0;
        padding: 2rem;
        background: rgba(66, 158, 157, 0.05);
        border-radius: 12px;
        border: 1px solid rgba(66, 158, 157, 0.2);
    }
 
    .matchup-controls {
        display: flex;
        gap: 1rem;
        margin: 1.5rem 0;
        justify-content: center;
    }
 
    .sort-btn {
        background: rgba(66, 158, 157, 0.3);
        border: 1px solid #429E9D;
        color: #fff;
        padding: 0.75rem 1.5rem;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.3s ease;
        font-weight: 500;
    }
 
    .sort-btn:hover {
        background: rgba(66, 158, 157, 0.5);
    }
 
    .sort-btn.active {
        background: #429E9D;
        color: #181818;
    }
 
    .matchup-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 2rem;
        margin-top: 1.5rem;
    }
 
    .matchup-subsection {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
 
    .subsection-title {
        font-size: 1.1rem;
        color: #30d5c8;
        margin: 0;
        padding-left: 0.5rem;
    }
 
    .matchup-list {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }
 
    .matchup-card {
        background: rgba(0, 0, 0, 0.3);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 6px;
        padding: 1rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        transition: all 0.2s ease;
    }
 
    .matchup-card:hover {
        background: rgba(0, 0, 0, 0.5);
        transform: translateX(4px);
    }
 
    .matchup-hero-name {
        font-size: 1rem;
        font-weight: 600;
        color: #fff;
        min-width: 120px;
    }
 
    .matchup-stats {
        display: flex;
        gap: 1.5rem;
        align-items: center;
        flex: 1;
        justify-content: center;
    }
 
    .matchup-stat {
        color: #ccc;
        font-size: 0.95rem;
    }
 
    .matchup-winrate {
        font-weight: bold;
        font-size: 1.1rem;
        min-width: 60px;
        text-align: right;
    }
 
    .matchup-games {
        font-size: 0.85rem;
        color: #999;
        min-width: 90px;
        text-align: right;
    }
 
    /* HERO SECTION */
    .hero-section {
        margin-top: 3rem;
    }
 
    .hero-cards-container {
        display: flex;
        justify-content: center;
        gap: 0.5rem;
        flex-wrap: wrap;
    }
 
    .hero-card {
        flex: auto;
        color: #fff;
        background: rgba(66, 158, 157, 0.15);
        border: 1px solid rgba(66, 158, 157, 0.3);
        border-radius: 6px;
        padding: 1rem;
        min-width: 150px;
    }
 
    .hero-name {
        margin: 0 0 0.5rem 0;
        font-size: 1.1rem;
        color: #30d5c8;
    }
 
    .hero-card p {
        margin: 0.25rem 0;
        font-size: 0.9rem;
    }
 
    @media (max-width: 768px) {
        .stats-grid {
            grid-template-columns: 1fr;
        }
 
        .matchup-grid {
            grid-template-columns: 1fr;
        }
 
        .matchup-card {
            flex-direction: column;
            align-items: flex-start;
            gap: 0.5rem;
        }
 
        .matchup-stats {
            justify-content: flex-start;
            width: 100%;
        }
 
        .matchup-games {
            text-align: left;
        }
        .hero-card {
            flex: 1 1 calc(50% - 0.5rem);
        }
    }
    @media (max-width: 1024px) {
        .hero-card {
            flex: 1 1 calc(33% - 1rem);
        }
    }
    @media (max-width: 480px) {
        .hero-card {
            flex: 1 1 100%;
            max-width: none;
        }
    }
</style>

