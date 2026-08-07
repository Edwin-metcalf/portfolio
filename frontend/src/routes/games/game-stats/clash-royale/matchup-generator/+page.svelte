<script lang="ts">
	import { fetchAPI } from "$lib/api";
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
<div>
    <h1>
        Generate your matchup against you and your friend!
    </h1>
    <h3>The Clash Royale API only holds last 25 games per person. So you must have played each other in your last 25 games</h3>
    <div class="tag-submisssion">
        <input type="text" bind:value={matchupEntry.tag1} placeholder="Your Clash Royale Tag" class="tag-input"/>
        <input type="text" bind:value={matchupEntry.tag2} placeholder="Opponents Clash Royale Tag" class="tag-input"/>

        {#if !tagsSubmitted}
            <button class="submit-button" onclick={handleSubmit} disabled={isLoading}>
                {isLoading ? 'loading...' : 'Submit'} 
            </button>
        {:else}
            {#if !friendlyStats}
                <h2> No return of stats from api </h2>
            {:else if (friendlyStats.games.length) < 1}
                <h2>No games found between tags in last 25 games of each player</h2>
            {:else}
                <section class="friendly-panel">
                <h2 class="panel panel-label">Matchup stats for {friendlyStats.myTag} VS {friendlyStats.friendTag}</h2>
                <p class="friendly-tally">
                    {friendlyStats.myTag} {friendlyStats.wins}W's &middot; {friendlyStats.myTag} {friendlyStats.losses}W's &middot; 
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