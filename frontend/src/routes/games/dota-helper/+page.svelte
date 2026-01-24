<script lang="ts">
    import {onMount, tick} from 'svelte';
    import Chart, { type ChartConfiguration } from 'chart.js/auto'
	import { fetchAPI } from '$lib/api';


    interface winLosePercentage {
        wins: number,
        lose: number,
        Percentage: number
    }

    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart | null  = null;
    let loading: boolean = true;
    let data;

    async function getDotaWinLose() {
        try {
            const result = await fetchAPI('/api/dota-helper/Get', {
                method: 'GET'
            });
            data = await result;
            console.log(data)
            return data;
        } catch (err) {
            console.error('Error fetching stats', err);
            return null;
        }
    }

    onMount(() => {
        (async () => {
        const winLoseData = await getDotaWinLose()
            loading = false;
            await tick();
            

            if (winLoseData && chartCanvas) {
                const config: ChartConfiguration<'doughnut'> = {   
                    type: 'doughnut',
                    data: {
                        labels: [
                            'Wins',
                            'Loses'
                        ],
                        datasets: [{
                            data: [winLoseData.wins, winLoseData.lose],
                            backgroundColor: [
                                'rgb(24,255,44)',
                                'rgb(255, 15, 60)'
                            ],
                            hoverOffset: 5
                        }]
                    }
                };
                chartInstance = new Chart(chartCanvas, config);
                console.log(chartInstance)
            }
        })();
        

        return () => {
            if (chartInstance) {
                chartInstance.destroy();
            }
        };
    });

</script>

{#if loading}
    <p>Loading...</p>
{:else}
    <div class="chart-container">
        <canvas bind:this={chartCanvas}></canvas>
    </div>
{/if}


<style>
    .chart-container {
        position: relative;
        width: 100%;
        height: 400px;
    }
</style>

