<script lang="ts">
    import {type ClashRoyaleStatsReturn, createWinLoseChart, createTrophyLineChart} from "../stats"
    import Chart from 'chart.js/auto'
    import { fetchAPI } from '$lib/api';
    import {onMount, tick} from 'svelte';
    import { gameStatsCache } from '$lib/store/GameStatsCache.svelte';



    let clashRoyaleData = $state<ClashRoyaleStatsReturn | null>(null);
    let chartCanvasOverall = $state<HTMLCanvasElement | null>(null);

    let chartInstanceOverall: Chart | null  = null;
    let loading = $state<boolean>(true);

     async function getClashRoyaleData(): Promise<ClashRoyaleStatsReturn | null> {
        try {
            const result = await fetchAPI('/api/games/clash-royale', {
                method: 'GET'
            });
            return result as ClashRoyaleStatsReturn;
        } catch (err) {
            console.error('Error fetching Clash Royale stats', err);
            return null;
        }
    }

    onMount(() => {
        (async () => {
            if (!gameStatsCache.clashRoyaleData.fetched) {
                const data = await getClashRoyaleData();
                if(!data) {
                    loading = false;
                    return;
                }
                gameStatsCache.setClashRoyaleData(data);
                
            } else console.log('used the cache')
            clashRoyaleData = gameStatsCache.clashRoyaleData.data;

            if (!clashRoyaleData) return;

            loading = false;
            await tick();
            
            //this is the win loss chart
            if (chartCanvasOverall) {
                chartInstanceOverall = createWinLoseChart(chartCanvasOverall, clashRoyaleData.wins, clashRoyaleData.losses)
            }
        
        })();

        return () => {
            if (chartInstanceOverall) {
                chartInstanceOverall.destroy();
            }
        };
    });


    function formatPercentage(decimal: number) {
        return "%"+ (decimal*100).toFixed(2);
    }
</script>
<header class="stats-header">
        <h1>My Clash Royale Stats</h1>
            <a href="./dota">
                Dota
            </a>

    </header>
{#if loading}
           <p style="font-size: 1.5rem; color: #fff;">Loading stats...</p>
{:else}
    {#if clashRoyaleData}
            <h1>Clash Royale</h1>
            <div class= "CR-section">
                <section class="win-lose-winrate">
                    <h2 class="section-title">All Time</h2>
                    <div class="stats-container">
                        <div class="stat-card">
                            <div class="stat-label">Wins</div>
                            <div class="stat-value">{clashRoyaleData.wins}</div>
                        </div>
                        <div class="stat-card">
                            <div class="stat-label">Loses</div>
                            <div class="stat-value">{clashRoyaleData.losses}</div>
                        </div>
                        <div class="stat-card highlight">
                            <div class="stat-label">Win Rate</div>
                            <div class="stat-value">{formatPercentage(clashRoyaleData.winRate)}</div>
                        </div>
                    </div>
                    <div class="chart-container">
                        <canvas bind:this={chartCanvasOverall}></canvas>
                    </div>
                </section>
            </div>
        {:else}
            <h2 style="color: white;">Something is very very broken</h2>
        {/if}
{/if}

<style>
    .CR-section {
        color: #fff;
        padding: 1rem;
    }

    .stats-container {
        display: flex;
        gap: 1rem;
        margin: 1rem 0;
    }

    .stat-card {
        background: #429E9D;
        border-radius: 8px;
        padding: 1rem 1.5rem;
        text-align: center;
        min-width: 80px;
    }

    .stat-label {
        font-size: 0.8rem;
        text-transform: uppercase;
        margin-bottom: 0.25rem;
    }

    .stat-value {
        font-size: 1.4rem;
        font-weight: bold;
    }

    .chart-container {
        width: 300px;
        height: 300px;
        margin-top: 1rem;
    }
</style>