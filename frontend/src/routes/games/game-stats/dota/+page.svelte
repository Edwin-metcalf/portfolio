<script lang="ts">
	import { onMount, tick } from 'svelte';
	import Chart from 'chart.js/auto';
	import { fetchAPI } from '$lib/api';
	import { createWinLoseChart, type DotaStatsReturn, type MatchupStats } from '../stats';
	import { gameStatsCache } from '$lib/store/GameStatsCache.svelte';
    import { X } from 'lucide-svelte';


	let chartCanvasOverall = $state<HTMLCanvasElement | null>(null);
	let chartCanvasRecent = $state<HTMLCanvasElement | null>(null);

	let chartInstanceOverall: Chart | null = null;
	let chartInstanceRecent: Chart | null = null;

	let loading = $state<boolean>(true);
	let dotaData = $state<DotaStatsReturn | null>(null);
	let sortBy = $state<'winRate' | 'games'>('winRate');

	const sortedMatchups = $derived(getSortedMatchups());
	const bestMatchups = $derived(sortedMatchups.slice(0, 5));
	const worstMatchups = $derived(sortedMatchups.slice(-5).reverse());

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
			if (!gameStatsCache.dotaData.fetched) {
				const data = await getDotaWinLose();
				if (!data) {
					loading = false;
					return;
				}
				gameStatsCache.setDotaData(data);
			} else console.log('used the cache');
			dotaData = gameStatsCache.dotaData.data;
			if (!dotaData) return;

			console.log('Matchup Stats:', dotaData.matchupStats); // DEBUG
			console.log('All data:', dotaData); // DEBUG
			loading = false;
			await tick();

			if (chartCanvasOverall) {
				chartInstanceOverall = createWinLoseChart(
					chartCanvasOverall,
					dotaData.wins,
					dotaData.losses
				);
			}
			if (chartCanvasRecent) {
				chartInstanceRecent = createWinLoseChart(
					chartCanvasRecent,
					dotaData.recentWins,
					dotaData.recentLosses
				);
			}
		})();

		return () => {
			if (chartInstanceOverall) {
				chartInstanceOverall.destroy();
			}
			if (chartInstanceRecent) {
				chartInstanceRecent.destroy();
			}
		};
	});

	function formatPercentage(decimal: number) {
		return '%' + (decimal * 100).toFixed(2);
	}

	function getSortedMatchups(): [string, MatchupStats][] {
		if (!dotaData?.matchupStats) return [];

		const matchups = Object.entries(dotaData.matchupStats);

		if (sortBy === 'winRate') {
			return matchups.sort((a, b) => b[1].winRate - a[1].winRate);
		} else {
			const top_games_matchups = matchups.sort((a, b) => b[1].games - a[1].games).slice(0, 10);
			return top_games_matchups.sort((a, b) => b[1].winRate - a[1].winRate);
		}
	}

	function getMatchupColor(winRate: number): string {
		if (winRate >= 0.6) return '#22CB00';
		if (winRate >= 0.5) return '#429E9D';
		if (winRate >= 0.4) return '#FFA500';
		return '#C30000';
	}
</script>

<div class="game-stat-helper-page">
    <a href="/games" class="exit-button" aria-label="Exit games">
		<X size={24} />
	</a>
	<nav class="game-tabs">
		<a href="./dota" class="tab active">Dota 2</a>
		<a href="./clash-royale" class="tab">Clash Royale</a>
	</nav>

	{#if loading}
		<p style="font-size: 1.5rem; color: #fff;">Loading stats...</p>
	{:else if dotaData}
		<div class="panel-row">
			<section class="panel">
				<h2 class="panel-label">My All Time Dota Stats</h2>
				<div class="stat-row">
					<div class="stat">
						<span class="stat-label">Wins</span>
						<span class="stat-value">{dotaData.wins}</span>
					</div>
					<div class="stat">
						<span class="stat-label">Losses</span>
						<span class="stat-value">{dotaData.losses}</span>
					</div>
					<div class="stat">
						<span class="stat-label">Win Rate</span>
						<span class="stat-value accent">{formatPercentage(dotaData.winRate)}</span>
					</div>
				</div>
				<div class="chart-container">
					<canvas bind:this={chartCanvasOverall}></canvas>
				</div>
			</section>

			<section class="panel">
				<h2 class="panel-label">My Recent Dota Stats</h2>
				<div class="stat-row">
					<div class="stat">
						<span class="stat-label">Wins</span>
						<span class="stat-value">{dotaData.recentWins}</span>
					</div>
					<div class="stat">
						<span class="stat-label">Losses</span>
						<span class="stat-value">{dotaData.recentLosses}</span>
					</div>
					<div class="stat">
						<span class="stat-label">Win Rate</span>
						<span class="stat-value accent">{formatPercentage(dotaData.recentWinRate)}</span>
					</div>
				</div>
				<div class="chart-container">
					<canvas bind:this={chartCanvasRecent}></canvas>
				</div>
			</section>
		</div>

		<!--recent hero sections played-->
		<section class="panel">
			<h2 class="panel-label">Recent Hero Stats</h2>
			<div class="hero-cards-container">
				{#each Object.entries(dotaData.recentHeroStats) as [heroName, heroStats]}
					<div class="hero-card">
						<h2 class="hero-name">{heroName}</h2>
						<div class="stat-row">
							<div class="stat">
								<span class="stat-label">Wins</span>
								<span class="stat-value">{heroStats.wins}</span>
							</div>
							<div class="stat">
								<span class="stat-label">Losses</span>
								<span class="stat-value">{heroStats.losses}</span>
							</div>
						</div>
						<p class="hero-kda">
							Avg KDA: {heroStats.avgKDA[0]} / {heroStats.avgKDA[1]} / {heroStats.avgKDA[2]}
						</p>
					</div>
				{/each}
			</div>
		</section>
		<!-- matchup stats area -->

		<section class="panel">
			<h2 class="panel-label">Enemy Matchups</h2>

			<div class="matchup-controls">
				<button
					class="sort-btn"
					class:active={sortBy === 'winRate'}
					onclick={() => (sortBy = 'winRate')}
				>
					Sort by Win Rate
				</button>
				<button
					class="sort-btn"
					class:active={sortBy === 'games'}
					onclick={() => (sortBy = 'games')}
				>
					Sort by Games Played
				</button>
			</div>

			<div class="matchup-grid">
				<div class="matchup-subsection">
					<h3 class="subsection-title">Best Matchups</h3>
					<div class="matchup-list">
						{#each bestMatchups as [heroName, stats]}
							<div
								class="matchup-card"
								style="border-left: 4px solid {getMatchupColor(stats.winRate)};"
							>
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
						{#each worstMatchups as [heroName, stats]}
							<div
								class="matchup-card"
								style="border-left: 4px solid {getMatchupColor(stats.winRate)};"
							>
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
		</section>
	{:else}
		<h2 style="color: white;">Something broke... Maybe internet issues?? then its not my fault</h2>
	{/if}
</div>

<style>
	.game-stat-helper-page {
		margin: 0 auto;
		padding: 2rem;
	}
	.chart-container {
		position: relative;
		width: 100%;
		height: 250px;
		margin-top: 1rem;
		margin-bottom: 0.5rem;
	}
	.hero-card {
		background: var(--panel);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1rem;
		min-width: 180px;
		flex: 1 1 calc(25% - 1rem);
		max-width: 250px;
	}
	.hero-kda {
		margin: 0.75rem 0 0;
		font-size: 0.85rem;
		color: var(--text-muted);
		font-family: var(--font-mono);
	}

	.matchup-controls {
		display: flex;
		gap: 1rem;
		margin: 1.5rem 0;
		justify-content: center;
	}
	.sort-btn {
		background: transparent;
		border: 1px solid var(--mint);
		color: var(--text);
		padding: 0.75rem 1.5rem;
		border-radius: var(--radius);
		cursor: pointer;
		transition: all 0.3s ease;
		font-weight: 500;
		font-family: var(--font-mono);
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.sort-btn:hover {
		background: rgba(62, 207, 192, 0.15);
	}

	.sort-btn.active {
		background: var(--mint);
		color: var(--bg);
		border-color: var(--mint);
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
		font-family: var(--font-mono);
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		color: var(--text-muted);
		margin: 0;
	}

	.matchup-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.matchup-card {
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		transition: all 0.2s ease;
	}

	.matchup-hero-name {
		font-size: 0.8rem;
		font-weight: 100;
		color: var(--text);
		min-width: 120px;
		font-family: var(--font-mono);
	}

	.matchup-stats {
		display: flex;
		gap: 1.5rem;
		align-items: center;
		flex: 1;
		justify-content: center;
	}

	.matchup-stat {
		color: var(--text-muted);
		font-size: 0.95rem;
		font-family: var(--font-mono);
	}

	.matchup-winrate {
		font-weight: bold;
		font-size: 1.1rem;
		min-width: 60px;
		text-align: right;
		font-family: var(--font-mono);
	}

	.matchup-games {
		font-size: 0.85rem;
		color: var(--text-muted);
		min-width: 90px;
		text-align: right;
		font-family: var(--font-mono);
	}

	.hero-cards-container {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.hero-name {
		color: var(--mint);
		margin: 0 0 0.5rem 0;
		font-family: var(--font-mono);
		font-size: 1rem;
	}

	.hero-card p {
		margin: 0.25rem 0;
		font-size: 0.9rem;
	}

	@media (max-width: 768px) {

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
