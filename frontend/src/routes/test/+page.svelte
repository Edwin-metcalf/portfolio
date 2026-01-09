<script lang="ts">
    import { fetchAPI } from "$lib/api";
    import { onMount } from "svelte";

    let healthStatus = 'Loading...';
    let leaderboard: Array<{name: string, score: number}> = [];

    onMount(async () => {
        try {
            const health = await fetchAPI('/api/health');
            healthStatus = health.message;

            const data = await fetchAPI('/api/space-invaders/leaderboard');
            leaderboard = data;
        } catch (err) {
            healthStatus = 'Error: ' + (err as Error).message;
        }
    });

    async function updateScores() {
        try {
            const data = await fetchAPI('/api/space-invaders/leaderboard');
            leaderboard = data;
        }
        catch (err) {
            console.error('Error', err);
        }
    }

    async function submitScore() {
        try {
            const result = await fetchAPI('/api/space-invaders/entry', {
                method: 'POST',
                body: JSON.stringify({ name: 'test player', score: 1000})
            });

            console.log('submitted: ', result)
        } catch (err) {
            console.error('Error: ', err);
        }
    }
</script>

<div>
    <h1>backend status: {healthStatus}</h1>
    <h2>leaderboard</h2>

    <ul>
        {#each leaderboard as entry}
            <li>{entry.name}: {entry.score}</li>

        {/each}
    </ul>

    <button on:click={submitScore}> Test submit</button>
    <button on:click={updateScores}>Update scores</button>
</div>