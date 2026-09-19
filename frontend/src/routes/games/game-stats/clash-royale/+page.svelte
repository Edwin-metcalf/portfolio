<script lang="ts">
	import {
		type ClashRoyaleLoadReturn,
		createWinLoseChart,
		createTrophyLineChart,
		type TrophyPoint,
		createHeadToHeadChart
	} from '../stats';
	import Chart from 'chart.js/auto';
	import { fetchAPI } from '$lib/api';
	import { onMount, tick } from 'svelte';
	import { gameStatsCache } from '$lib/store/GameStatsCache.svelte';
    import { X } from 'lucide-svelte';

	let clashRoyaleData = $state<ClashRoyaleLoadReturn | null>(null);
	let chartCanvasOverall = $state<HTMLCanvasElement | null>(null);

	let chartInstanceOverall: Chart | null = null;

	let chartCanvasTrophies = $state<HTMLCanvasElement | null>(null);
	let chartInstanceTrophies: Chart<'line', TrophyPoint[]> | null = null;

	let headToHeadCanvas = $state<HTMLCanvasElement | null>(null);
	let headToHeadInstance: Chart<'bar'> | null = null;

	let loading = $state<boolean>(true);
	const leagueNumToNameMap = new Map<number, string>([
		[1, 'Master I'],
		[2, 'Master II'],
		[3, 'Master III'],
		[4, 'Champion'],
		[5, 'Grand Champion'],
		[6, 'Royal Champion'],
		[7, 'Ultimate Champion'],
	]);



	async function getClashRoyaleData(): Promise<ClashRoyaleLoadReturn | null> {
		try {
			const result = await fetchAPI('/api/games/clash-royale', {
				method: 'GET'
			});
			return result as ClashRoyaleLoadReturn;
		} catch (err) {
			console.error('Error fetching Clash Royale stats', err);
			return null;
		}
	}

	onMount(() => {
		(async () => {
			if (!gameStatsCache.clashRoyaleData.fetched) {
				const data = await getClashRoyaleData();
				if (!data) {
					loading = false;
					return;
				}
				gameStatsCache.setClashRoyaleData(data);
			} else console.log('used the cache');
			clashRoyaleData = gameStatsCache.clashRoyaleData.data;

			if (!clashRoyaleData) return;

			loading = false;
			await tick();

			//this is the win loss chart
			if (chartCanvasOverall) {
				chartInstanceOverall = createWinLoseChart(
					chartCanvasOverall,
					clashRoyaleData.profile.wins,
					clashRoyaleData.profile.losses
				);
			}
			//gonna have to make this either straight up by date or with an option to switch back and forth
			if (chartCanvasTrophies) {
				chartInstanceTrophies = createTrophyLineChart(
					chartCanvasTrophies,
					clashRoyaleData.battleLog
				);
			}
			//head to head chart
			if (headToHeadCanvas) {
				headToHeadInstance = createHeadToHeadChart(
					headToHeadCanvas,
					clashRoyaleData.friendly.wins,
					clashRoyaleData.friendly.losses
				);
			}
		})();

		return () => {
			if (chartInstanceOverall) {
				chartInstanceOverall.destroy();
			}
			if (chartInstanceTrophies) {
				chartInstanceTrophies.destroy();
			}
			if (headToHeadInstance) {
				headToHeadInstance.destroy();
			}
		};
	});

	function formatPercentage(decimal: number) {
		return (decimal * 100).toFixed(2) + '%';
	}

	function getBattleResult(result: number): 'win' | 'loss' | 'tie' {
		if (result > 0) return 'win';
		if (result === 0) return 'loss';
		return 'tie';
	}

    function getArenaNumByTrophies(trophies: number): number {
        //this is kinda a hack solution especially at lower arenas
        if (trophies < 5000) {
            if (trophies < 1000) return Math.floor(trophies / 300) + 1;
            return Math.floor(trophies / 400) + 3;
        }

        return 15 + Math.floor((trophies - 5000) / 500);
    }
</script>

<div class="game-stat-helper-page">
    <a href="/games" class="exit-button" aria-label="Exit games">
        <X size={24} />
    </a>
	<nav class="game-tabs">
		<a href="./dota" class="tab">Dota 2</a>
		<a href="./clash-royale" class="tab active">Clash Royale</a>
	</nav>

	{#if loading}
		<p class="loading-texts">Loading stats...</p>
	{:else if clashRoyaleData}
		<section class="profile-panel">
			<h2 class="panel-label">My Clash Royale Stats</h2>

			<p class="player-name">{clashRoyaleData.profile.name}</p>
			<p class="player-meta">
				{clashRoyaleData.profile.tag} &middot; 
				Arena {getArenaNumByTrophies(clashRoyaleData.profile.trophies)}: 
				{clashRoyaleData.profile.arena.name} &middot; 
				{leagueNumToNameMap.get(clashRoyaleData.profile.currentPathOfLegendSeasonResult.leagueNumber)}
			</p>

			<div class="profile-stats">
				<div class="stat">
					<span class="stat-label">Wins</span>
					<span class="stat-value">{clashRoyaleData.profile.wins}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Loses</span>
					<span class="stat-value">{clashRoyaleData.profile.losses}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Win Rate</span>
					<span class="stat-value accent">{formatPercentage(clashRoyaleData.profile.winRate)}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Best Trophies</span>
					<span class="stat-value">{clashRoyaleData.profile.bestTrophies}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Current Trophies</span>
					<span class="stat-value">{clashRoyaleData.profile.trophies}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Current Rank</span>
					<span class="stat-value">{leagueNumToNameMap.get(clashRoyaleData.profile.currentPathOfLegendSeasonResult.leagueNumber)}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Best Rank</span>
					<span class="stat-value">{leagueNumToNameMap.get(clashRoyaleData.profile.bestPathOfLegendSeasonResult.leagueNumber)}</span>
				</div>
			</div>
		</section>
		<div class="chart-row">
			<section class="winloss-panel">
				<h2 class="panel panel-label">Win / Loss breakdown</h2>
				<div class="chart-container">
					<canvas bind:this={chartCanvasOverall}></canvas>
				</div>
			</section>

			<section class="trophy-panel">
				<h2 class="panel panel-label">Trophy progression</h2>
				<div class="stat">
					<span class="stat-label">{clashRoyaleData.battleLog.length} games included</span>
				</div>
				<div class="chart-container">
					<canvas bind:this={chartCanvasTrophies}></canvas>
				</div>
			</section>
		</div>

		<section class="ranked-panel">
			<h2 class="panel panel-label">Ranked Stats</h2>
			<div class="profile-stats">
				<div class="stat">
					<span class="stat-label">Wins</span>
					<span class="stat-value">{clashRoyaleData.ranked.wins}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Loses</span>
					<span class="stat-value">{clashRoyaleData.ranked.losses}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Win Rate</span>
					<span class="stat-value accent">{formatPercentage(clashRoyaleData.ranked.winRate)}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Current League</span>
					<span class="stat-value">{leagueNumToNameMap.get(clashRoyaleData.profile.currentPathOfLegendSeasonResult.leagueNumber)}</span>
				</div>
			</div>
			<!--THINGS TO ADD
			section titles like this deck, maybe add the evo information lowkey 
			question if we need the images 
			color code the win rates show that its this deck vs these cards ie this deck and matchups
			then do it to the matchups looks a bit dumb right now
			-->
			{#if clashRoyaleData.ranked.deckCardWinRates?.length}
				{#each clashRoyaleData.ranked.deckCardWinRates as deckStats}
					<div class="ranked-deck">
						<div class="deck-cards">
							{#each deckStats.deck as card}
								<img src={card.iconUrls.medium} alt={card.name} title={card.name} class="card-icon"/>
							{/each}
						</div>

						<div class="card-matchups">
							{#each Object.entries(deckStats.cardStats) as [cardName, record]}
							<div class="matchup-chip">
								<span class="matchup-name">{cardName}</span>
								<span class="matchup-rate">{formatPercentage(record.winRate)}</span>
							</div>
							{/each}
						</div>
					</div>
				{/each}
			{/if}

		</section>

		<section class="friendly-panel">
			<h2 class="panel panel-label">Recent Vs Ryan</h2>
			<p class="friendly-tally">
				{clashRoyaleData.friendly.wins}W &middot; {clashRoyaleData.friendly.losses}L &middot; {clashRoyaleData
					.friendly.ties}T &middot;
				<span class="accent">{formatPercentage(clashRoyaleData.friendly.winRate)} WR</span>
			</p>
            <div class="friendly-row">
                <div class="chart-container">
                    <canvas bind:this={headToHeadCanvas}></canvas>
                </div>
                <div class="battle-history">
                    <div class="battle-squares">
                        {#each clashRoyaleData.friendly.games as game}
                            {@const outcome = getBattleResult(game.result)}
                            <span class="battle-square {outcome}">
                                {outcome === 'win' ? 'W' : outcome === 'loss' ? 'L' : 'T'}
                            </span>
                        {/each}
                    </div>
                </div>
            </div>
		</section>

		<section class="matchup-panel">
			<p class="matchup-text">Want to see your own matchup against a friend?</p>
			<a href="/games/game-stats/clash-royale/matchup-generator" class="matchup-button">
				check your matchup
			</a>

		</section>

	{:else}
		<h2 style="color: white;">Something is very very broken</h2>
	{/if}
</div>

<style>
	.ranked-deck {
		margin-top: 1.5rem;
	}
	.deck-cards {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}
	.card-icon {
		width: 48px;
		height: 56px;
		object-fit: contain;
	}
	.card-matchups {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
		margin-top: 1rem;
	}
	.matchup-chip {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2px;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 0.4rem 0.6rem;
		font-family: var(--font-mono);
		font-size: 0.7rem;
	}
	.matchup-name {
		color: var(--text-muted);
	}
	.matchup-rate {
		color: var(--mint);
	}
.matchup-panel {
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
		text-align: center;
	}
	.matchup-text {
		font-family: var(--font-mono);
		color: var(--text-muted);
		margin: 0 0 1rem;
	}
	.matchup-button {
		display: inline-block;
		font-family: var(--font-mono);
		background: var(--bg);
		border: 1px solid var(--mint);
		border-radius: var(--radius);
		padding: 0.75rem 1.5rem;
		color: var(--mint);
		text-decoration: none;
		font-weight: 500;
		transition: transform 0.15s ease, background 0.15s ease;
	}
	.matchup-button:hover {
		transform: translateY(-2px);
		background: var(--mint);
		color: var(--bg);
	}
	.chart-container {
        position: relative;
		width: 100%;
		height: 100%;
        max-width: 300px;
        aspect-ratio: 1 / 1;
		margin: 0 auto;
		margin-top: 1rem;
    }
	.profile-panel,
	.winloss-panel,
	.trophy-panel,
	.friendly-panel,
	.ranked-panel {
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
		margin-bottom: 1.5rem;
	}
	.friendly-panel .panel-label {
		font-size: 1rem;
	}
    .friendly-row .chart-container {
        aspect-ratio: auto;
        width: 100%;
        max-width: 400px;
        margin: 0;
    }
	.trophy-panel .chart-container {
		max-width: 100%;
		aspect-ratio: auto;
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
	.chart-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.5rem;
	}
	@media (max-width: 700px) {
		.chart-row {
			grid-template-columns: 1fr;
		}
	}
	.player-name {
		font-family: var(--font-display);
		font-size: 1.75rem;
		color: var(--text);
		margin: 0;
	}
	.player-meta {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		color: var(--text-muted);
		margin: 0.25rem 0 0;
	}
	.profile-stats {
		display: flex;
		gap: 2.5rem;
		margin-top: 1.25rem;
		flex-wrap: wrap;
	}
	.friendly-tally {
		color: var(--text-muted);
		font-family: var(--font-mono);
		font-size: 0.85rem;
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
