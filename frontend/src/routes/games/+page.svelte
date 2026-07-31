<script lang="ts">
	import { X } from 'lucide-svelte';

	const games = [
		{
			title: 'Space Invaders',
			description: 'My take on the classic space invaders',
			url: '/games/space-invaders',
			status: 'playable'
		},
		{
			title: 'Video Game Stats',
			description: 'Dota 2 and Clash Royale Stats + head to head matchup against Ryan',
			url: '/games/game-stats',
			status: 'playable'
		},
		{
			title: 'Under Construction',
			description: 'Always working on something',
			url: '#',
			status: 'coming-soon'
		}
	];
</script>

<div class="games-section">
	<a href="/" class="exit-button" aria-label="Exit games">
		<X size={24} />
	</a>
	<h1>Games & Projects</h1>

	<div class="games-grid">
		{#each games as game}
			<a href={game.url} class="game-card" class:coming-soon={game.status === 'coming-soon'}>
				<div class="game-content">
					<h2>{game.title}</h2>
					<p>{game.description}</p>
					{#if game.status === 'playable'}
						<span class="play-button">Play Now</span>
					{:else}
						<span class="coming-soon">Coming soon</span>
					{/if}
				</div>
			</a>
		{/each}
	</div>
</div>

<style>
	.games-section {
		position: relative;
		margin: 0 auto;
		padding: 4rem 2rem;
		min-height: 100vh;
		background-color: #161b22;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
	h1 {
		font-family: var(--font-display);
		font-size: 3rem;
		margin-bottom: 2rem;
		text-align: center;
		color: var(--mint);
		width: 100%;
	}
	.games-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 2rem;
		width: 100%;
		max-width: 1100px;
	}
	.game-card {
		background: var(--panel);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 2rem;
		text-decoration: none;
		color: var(--text);
		transition:
			transform 0.3s ease,
			box-shadow 0.3s ease;
		cursor: pointer;
		position: relative;
		overflow: hidden;
	}
	.game-card:hover {
		transform: translateY(-5px);
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
	}
	.game-card.coming-soon {
		background: var(--bg);
		border-color: var(--border);
		cursor: not-allowed;
		opacity: 0.6;
	}
	.game-card.coming-soon:hover {
		transform: none;
	}
	.game-content h2 {
		font-size: 1.8rem;
		margin-bottom: 0.5rem;
	}

	.game-content p {
		font-size: 1rem;
		margin-bottom: 1rem;
		opacity: 0.9;
	}
	.play-button {
		display: inline-block;
		background: rgba(62, 207, 192, 0.15);
		padding: 0.5rem 1rem;
		border-radius: var(--radius);
		font-family: var(--font-mono);
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 500;
		color: var(--mint);
	}

	.game-card:hover .play-button {
		background: rgba(0, 212, 170, 0.25);
	}

	.coming-soon {
		display: inline-block;
		background: rgba(255, 255, 255, 0.1);
		padding: 0.5rem 1rem;
		border-radius: var(--radius);
		font-family: var(--font-mono);
		font-size: 0.75rem;
		letter-spacing: 0.5px;
		color: var(--text-muted);
	}

	@media (max-width: 768px) {
		h1 {
			font-size: 2rem;
		}

		.games-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
