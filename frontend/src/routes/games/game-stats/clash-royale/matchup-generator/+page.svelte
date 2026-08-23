<script lang="ts">
	import { fetchAPI } from "$lib/api";
    import { X } from 'lucide-svelte';
    import { type ClashRoyaleFriendlyLoadReturn, type matchupGeneratorTags, createHeadToHeadChart } from "../../stats";
    import Chart from 'chart.js/auto';


    let matchupEntry = $state<matchupGeneratorTags>({tag1: "", tag2: ""});
    let tagsSubmitted: boolean = $state(false);
    let friendlyStats = $state<ClashRoyaleFriendlyLoadReturn | null>(null);
    let isLoading = $state(false);

    let headToHeadCanvas = $state<HTMLCanvasElement | null>(null);
	let headToHeadInstance: Chart<'bar'> | null = null;

    $effect(() => {
        if (headToHeadCanvas && friendlyStats) {
            headToHeadInstance = createHeadToHeadChart (
                headToHeadCanvas,
				friendlyStats.wins,
				friendlyStats.losses
            );

            return () => {
                if(headToHeadInstance) {
                    headToHeadInstance.destroy()
                    headToHeadInstance = null
                }
            };
        }
    });



    async function getMatchUpData(tags: matchupGeneratorTags): Promise<ClashRoyaleFriendlyLoadReturn | null> {
        try {
            const result = await fetchAPI('/api.games/clash-royale/matchup-generator', { 
                method: 'POST',
                headers : {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(tags)

        });
        return result as ClashRoyaleFriendlyLoadReturn;
        } catch (err) {
            console.error("Error: ", err);
            return null;
        }

    }

    async function handleSubmit() {
        isLoading = true;
        friendlyStats = await getMatchUpData(matchupEntry);
    
        tagsSubmitted = true;
        isLoading = false;
    }

    function getBattleResult(result: number): 'win' | 'loss' | 'tie' {
		if (result > 0) return 'win';
		if (result === 0) return 'loss';
		return 'tie';
	}

    function formatPercentage(decimal: number) {
		return (decimal * 100).toFixed(2) + '%';
	}

</script>
<div class="game-stat-helper-page">
    <a href="/games" class="exit-button" aria-label="Exit games">
        <X size={24} />
    </a>

    <h1 class="page-title">
        Generate your matchup against you and your friend!
    </h1>

    <h3 class="page-subtitle">
        The Clash Royale API only holds last 25 games per person. So you must have played each other in your last 25 games and start your tag with "#" or it will not work.
    </h3>
    <div class="tag-submisssion">
        <input type="text" bind:value={matchupEntry.tag1} placeholder="Your Clash Royale Tag" class="tag-input"/>
        <input type="text" bind:value={matchupEntry.tag2} placeholder="Opponents Clash Royale Tag" class="tag-input"/>

        
        <button class="submit-button" onclick={handleSubmit} disabled={isLoading}>
            {isLoading ? 'loading...' : 'Submit'} 
        </button>
        {#if tagsSubmitted}
            {#if !friendlyStats}
                <h2> No return of stats from api </h2>
            {:else if (friendlyStats.games.length) < 1}
                <h2>No games found between tags in last 25 games of each player</h2>
            {:else}
                <section class="friendly-panel">
                <h2 class="panel panel-label">Matchup stats for {friendlyStats.myName} VS {friendlyStats.friendName}</h2>
                <p class="friendly-tally">
                    {friendlyStats.myName} {friendlyStats.wins}W's &middot; {friendlyStats.friendName} {friendlyStats.losses}W's &middot; 
                        {friendlyStats.ties}T's &middot;
                    <span class="accent">{formatPercentage(friendlyStats.winRate)} WR</span>
                </p>
                <div class="friendly-row">
                    <div class="chart-container">
                        <canvas bind:this={headToHeadCanvas}></canvas>
                    </div>
                    <div class="battle-history">
                        <div class="battle-squares">
                            {#each friendlyStats.games as game}
                                {@const outcome = getBattleResult(game.result)}
                                <span class="battle-square {outcome}">
                                    {outcome === 'win' ? 'W' : outcome === 'loss' ? 'L' : 'T'}
                                </span>
                            {/each}
                        </div>
                    </div>
                </div>
            </section>
            {/if}
        {/if}

    </div>
</div>

<style>
    .page-title {
        font-family: var(--font-display);
        font-size: 1.75rem;
        color: var(--text);
        margin: 0 0 0.5rem;
    }
    .page-subtitle {
        font-family: var(--font-mono);
        font-size: 0.8rem;
        color: var(--text-muted);
        margin: 0 0 2rem;
        max-width: 500px;
    }
    .tag-input {
		font-family: var(--font-mono);
		font-size: 0.85rem;
		background: var(--panel);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		color: var(--text);
		padding: 0.65rem 0.85rem;
		min-width: 220px;
	}
	.tag-input::placeholder {
		color: var(--text-muted);
	}
	.tag-input:focus {
		outline: none;
		border-color: var(--mint);
	}
	.submit-button {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		background: transparent;
		border: 1px solid var(--mint);
		border-radius: var(--radius);
		color: var(--mint);
		padding: 0.65rem 1.25rem;
		cursor: pointer;
		transition: background 0.2s ease;
	}
	.submit-button:hover:not(:disabled) {
		background: rgba(62, 207, 192, 0.12);
	}
	.submit-button:disabled {
		opacity: 0.5;
		cursor: default;
	}
    .friendly-panel {
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
		margin-top: 1.5rem;
	}
	.friendly-panel .panel-label {
		font-size: 1rem;
	}
	.friendly-tally {
		color: var(--text-muted);
		font-family: var(--font-mono);
		font-size: 0.85rem;
	}
	.friendly-row {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 1.5rem;
		align-items: center;
		justify-items: center;
		width: 100%;
		margin-top: 1rem;
	}
	.chart-container {
		position: relative;
		aspect-ratio: auto;
		width: 100%;
		height: 100%;
		max-width: 400px;
		margin: 0;
	}
	.battle-history {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
	}
	.battle-squares {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		justify-content: center;
		max-width: 320px;
	}
	.battle-square {
		width: 24px;
		height: 24px;
		border-radius: 4px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-family: var(--font-mono);
		font-size: 0.65rem;
		font-weight: 500;
	}
	.battle-square.win {
		background: var(--mint);
		color: var(--bg);
	}
	.battle-square.loss {
		border: 1px dashed rgba(255, 255, 255, 0.28);
		color: var(--text-muted);
	}
	.battle-square.tie {
		border: 1px solid rgba(255, 255, 255, 0.28);
		color: var(--text-muted);
	}

</style>