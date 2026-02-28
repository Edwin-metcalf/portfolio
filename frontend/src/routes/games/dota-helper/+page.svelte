<script lang="ts">
    import {onMount, tick} from 'svelte';
    import Chart, { type ChartConfiguration } from 'chart.js/auto'
	import { fetchAPI } from '$lib/api';
    import {createWinLoseChart, type DotaStatsReturn} from './stats'



    let chartCanvasOverall: HTMLCanvasElement;
    let chartCanvasRecent: HTMLCanvasElement;

    let chartInstanceOverall: Chart | null  = null;
    let chartInstanceRecent: Chart | null = null;
    let loading: boolean = true;
    let data;
    let winLoseData: DotaStatsReturn;

    async function getDotaWinLose(): Promise<DotaStatsReturn | null> {
        try {
            const result = await fetchAPI('/api/dota-helper/Get', {
                method: 'GET'
            });
            console.log(data)
            return result as DotaStatsReturn;
        } catch (err) {
            console.error('Error fetching stats', err);
            return null;
        }
    }

    onMount(() => {
        (async () => {
            const data = await getDotaWinLose();
            if (!data) return;

            winLoseData = data 
            loading = false;
            await tick();
            

            if (chartCanvasOverall) {
                chartInstanceOverall = createWinLoseChart(chartCanvasOverall, winLoseData.wins, winLoseData.losses)
            }
            if (chartCanvasRecent) {
                chartInstanceRecent = createWinLoseChart(chartCanvasRecent, winLoseData.recentWins, winLoseData.recentLosses)

            }
        })();
        

        return () => {
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

</script>
<div class="dota-helper-page">
    <header class="stats-header">
        <h1>My Dota Stats</h1>
    </header>

    {#if loading}
        <p style="font-size: 1.5rem; color: #fff;">Loading stats...</p>
    {:else}
        <div class="stats-grid">
            <section class="win-lose-winrate">
                <h2 class="section-title">All Time</h2>
                <div class="stats-container">
                    <div class="stat-card">
                        <div class="stat-label">Wins</div>
                        <div class="stat-value">{winLoseData.wins}</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Loses</div>
                        <div class="stat-value">{winLoseData.losses}</div>
                    </div>
                    <div class="stat-card highlight">
                        <div class="stat-label">Win Rate</div>
                        <div class="stat-value">{formatPercentage(winLoseData.winRate)}</div>
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
                        <div class="stat-value">{winLoseData.recentWins}</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Loses</div>
                        <div class="stat-value">{winLoseData.recentLosses}</div>
                    </div>
                    <div class="stat-card highlight">
                        <div class="stat-label">Win Rate</div>
                        <div class="stat-value">{formatPercentage(winLoseData.recentWinRate)}</div>
                    </div>
                </div>
                <div class="chart-container">
                    <canvas bind:this={chartCanvasRecent}></canvas>
                </div>
            </section>
        </div>
        <div class="hero-section">
            {#each Object.entries(winLoseData.recentHeroStats) as [heroName, heroStats]}
                <div class="hero-card">
                    <h2>{heroName}</h2>
                    <p>Wins: {heroStats.wins}</p>
                    <p>Losses: {heroStats.losses}</p>
                </div>
            {/each}
        </div>
        
    {/if}
</div>


<style>
    :global(body) {
        margin: 0;
        padding: 0;
    }
    .dota-helper-page {
        background-color: #181818;
        margin: 0 auto;
        padding: 2rem;
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
        margin: 2rem 0;
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
</style>

