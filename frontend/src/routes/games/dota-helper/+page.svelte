<script lang="ts">
    import {onMount, tick} from 'svelte';
    import Chart, { type ChartConfiguration } from 'chart.js/auto'
	import { fetchAPI } from '$lib/api';
    import {createWinLoseChart, type WinLoseData} from './stats'



    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart | null  = null;
    let loading: boolean = true;
    let data;
    let winLoseData: WinLoseData;

    async function getDotaWinLose(): Promise<WinLoseData | null> {
        try {
            const result = await fetchAPI('/api/dota-helper/Get', {
                method: 'GET'
            });
            console.log(data)
            return result as WinLoseData;
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
            

            if (chartCanvas) {
                chartInstance = createWinLoseChart(chartCanvas, winLoseData)
            }
        })();
        

        return () => {
            if (chartInstance) {
                chartInstance.destroy();
            }
        };
    });

    function formatPercentage(decimal: number) {
        return (decimal*100).toFixed(2);
    }

</script>
<div class="Dota-helper-page">
    <header class="stats-header">
        <h1>My Dota Stats</h1>
    </header>

    {#if loading}
        <p>Loading...</p>
    {:else}
        <section class="win-lose-winrate">
            <h2 class="section-title">All Time</h2>
            <div class="stats-container">
                <div class="stat-card">
                    <div class="stat-label">Wins</div>
                    <div class="stat-value">{winLoseData.wins}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Loses</div>
                    <div class="stat-value">{winLoseData.lose}</div>
                </div>
                <div class="stat-card highlight">
                    <div class="stat-label">Win Rate</div>
                    <div class="stat-value">{formatPercentage(winLoseData.winRate)}</div>
                </div>
            </div>
        </section>

        <div class="chart-container">
                <canvas bind:this={chartCanvas}></canvas>
        </div>
    {/if}
</div>


<style>
      
    .chart-container {
        position: relative;
        width: 40%;
        height: 400px;
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
        gap: 1.5rem;
        justify-content: left;
    }
    .win-lose-winrate {
        display: flex;
        gap: 1.5rem;
        justify-content: left;
        margin: 2rem 0;
        /*background-color: #30d5c8; */
    }
    .stat-card {
        background: #429E9D;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        padding: 1rem 1.5rem;
        text-align: center;
        min-width: 75px;
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

