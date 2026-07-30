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

	let clashRoyaleData = $state<ClashRoyaleLoadReturn | null>(null);
	let chartCanvasOverall = $state<HTMLCanvasElement | null>(null);

	let chartInstanceOverall: Chart | null = null;

	let chartCanvasTrophies = $state<HTMLCanvasElement | null>(null);
	let chartInstanceTrophies: Chart<'line', TrophyPoint[]> | null = null;

	let headToHeadCanvas = $state<HTMLCanvasElement | null>(null);
	let headToHeadInstance: Chart<'bar'> | null = null;

	let loading = $state<boolean>(true);

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
</script>

<div class="game-stat-helper-page">
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
				{clashRoyaleData.profile.tag} &middot; {clashRoyaleData.profile.arena.name}
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
				<div class="chart-container">
					<canvas bind:this={chartCanvasTrophies}></canvas>
				</div>
			</section>
		</div>

		<section class="friendly-panel">
			<h2 class="panel-label">Recent Vs Ryan</h2>
			<p class="friendly-tally">
				{clashRoyaleData.friendly.wins}W &middot; {clashRoyaleData.friendly.losses}L &middot; {clashRoyaleData
					.friendly.ties}T &middot;
				<span class="accent">{formatPercentage(clashRoyaleData.friendly.winRate)} WR</span>
			</p>
			<div class="chart-container">
				<canvas bind:this={headToHeadCanvas}></canvas>
			</div>

			<div class="battle-squares">
				{#each clashRoyaleData.friendly.games as game}
					{@const outcome = getBattleResult(game.result)}
					<span class="battle-square {outcome}">
						{outcome === 'win' ? 'W' : outcome === 'loss' ? 'L' : 'T'}
					</span>
				{/each}
			</div>
		</section>
	{:else}
		<h2 style="color: white;">Something is very very broken</h2>
	{/if}
</div>

<style>
	.game-stat-helper-page {
		max-width: 960px;
		margin: 0 auto;
		padding: 2.5rem 2rem 4rem;
	}
	.chart-container {
		width: 300px;
		height: 300px;
		margin-top: 1rem;
	}
	.game-tabs {
		display: flex;
		gap: 24px;
		padding-bottom: 14px;
		border-bottom: 1px solid var(--border);
		margin-bottom: 2rem;
	}
	.tab {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		letter-spacing: 0.5px;
		text-transform: uppercase;
		color: var(--text-muted);
		text-decoration: none;
		padding-bottom: 10px;
	}
	.tab.active {
		color: var(--mint);
		border-bottom: 2px solid var(--mint);
	}
	.profile-panel,
	.winloss-panel,
	.trophy-panel,
	.friendly-panel {
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
		margin-bottom: 1.5rem;
	}
	.panel-label {
		font-family: var(--font-mono);
		font-size: 0.7rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		color: var(--text-muted);
		margin: 0 0 1rem;
	}
	.friendly-panel .panel-label {
		font-family: var(--font-display);
		font-size: 1.1rem;
		text-transform: none;
		letter-spacing: normal;
		color: var(--text);
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
	.panel {
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1.5rem;
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
	.stat-label {
		display: block;
		font-size: 0.7rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		color: var(--text-muted);
	}
	.stat-value {
		font-family: var(--font-mono);
		font-size: 1.1rem;
		color: var(--text);
	}
	.accent {
		color: var(--mint);
	}
	.loading-texts {
		color: var(--text-muted);
		font-family: var(--font-mono);
		font-size: 0.9rem;
	}
	.friendly-tally {
		color: var(--text-muted);
		font-family: var(--font-mono);
		font-size: 0.85rem;
	}
	.battle-squares {
		display: flex;
		gap: 6px;
		margin: 0.75rem 0;
	}
	.battle-square {
		width: 22px;
		height: 22px;
		border-radius: 2px;
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
