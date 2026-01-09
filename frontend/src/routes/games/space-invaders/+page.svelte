<script lang="ts">
    import { fetchAPI } from "$lib/api";
    import { onMount } from "svelte";
    import { SpaceInvaders } from "./game";
    import { Play, X, RotateCcw } from "lucide-svelte"

    interface leaderboardEntry {
        name: string,
        score: number
    }
    

    let canvas: HTMLCanvasElement;
    let game: SpaceInvaders;
    let score = $state(0);
    let gameStarted = $state(false);
    let imagesLoaded = $state(false);
    let gameOver = $state(false);
    let finalScore = $state(0);

    let consoleMessages = $state<string[]>([]);
    const MAX_MESSAGES = 3;
    let scoreSubmitted: boolean = $state(false)
    let entry = $state<leaderboardEntry>({name: "", score: 0})
    let leaderboard: Array<leaderboardEntry>  | undefined = $state()

    onMount(() => {
        game = new SpaceInvaders(canvas, 
        (message) => {
            addMessage(message)
        }, 
        (newScore) => {
            score = newScore;
        }, () => {
            gameOver = true;
            finalScore = score;
        });

        game.waitForImages().then(() => {
            imagesLoaded = true;
        });

        getLeaderboard()

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.code === 'Space' && gameStarted && !gameOver) {
                e.preventDefault();
            }
        };
        window.addEventListener('keydown', handleKeyDown)

        return () => {
            window.removeEventListener('keydown', handleKeyDown)
        }

    });

    function startGame() {
        if (imagesLoaded && !gameStarted) {
            gameStarted = true;
            gameOver = false;
            game.start();

            consoleMessages = ['SYSTEM ONLINE'];

        }
    }
    function restartGame() {
        gameOver = false;
        score = 0;
        scoreSubmitted = false
        game.restart();
    }
    function addMessage(msg: string) {
        if (consoleMessages.length >= MAX_MESSAGES) {
            consoleMessages.shift();
        }
        consoleMessages.push(msg);
        consoleMessages = consoleMessages;
    }

    //database stuff here


    async function SubmitScore(name: string, score: number) {
        entry.score = score
        entry.name = name
        try {
            const result = await fetchAPI('/api/space-invaders/entry', {
                method: 'POST',
                body: JSON.stringify(entry)

            });

        } catch (err) {
            console.error('Error: ', err);
        }
    }

    async function getLeaderboard() {
        try {
            const result = await fetchAPI('/api/space-invaders/leaderboard',{
                method: 'GET',
            });
            leaderboard = result
        } catch (err) {
            console.error('Error ', err);
        }
    }
    async function getAndSubmitScore(name: string, score: number) {
        try {
            await SubmitScore(name, score)
            await getLeaderboard()
            scoreSubmitted = true
        } catch (err) {
            console.error('failed to submit score:', err)
        }
    }

</script>
<div class="game-container">
    <a href="/games" class="exit-button" aria-label="Exit game">
        <X size={24} />
    </a>
    {#if gameStarted}
        <p class="score"><span>Score:</span> <span>{score}</span></p>
    {/if}
    {#if !gameStarted}
        <div class="play-overlay">
            <button class="play-button" onclick={startGame} disabled={!imagesLoaded}>
                {#if !imagesLoaded}
                    <span>loading...</span>
                {:else}
                    <Play size={32}/>
                    <span>Play Game</span>
                {/if}
            </button>
        </div>
    {/if}

    {#if gameOver}
        <div class="game-over-overlay">
            <div class="game-over-modal">
                <h2 class="game-over-title">Game Over</h2>
                <p class="final-score">Final Score: {finalScore}</p>



                <div class="score-submission">
                    <input type="text" bind:value={entry.name} placeholder="enter name" class="name-input"/>
                    {#if !scoreSubmitted}
                        <button class="submit-button" onclick={() => getAndSubmitScore(entry.name, finalScore)} disabled={!entry.name.trim()}> Submit</button>
                    {/if}
                </div>
                <button class="restart-button" onclick={restartGame}>
                    <RotateCcw size={24}/>
                    <span>Play Again</span>
                
                </button>
            </div>

        </div>
    {/if}
    <canvas bind:this={canvas}></canvas>

    <!-- I should do more comments this is the console area-->
    {#if gameStarted && !gameOver}
        <div class="terminal-console">
            {#each consoleMessages as message}
                <div class="terminal-message">{message}</div>
            {/each}
        </div>
    {/if}

</div>
<!--this is the leaderboard area for real type beat-->
    <div class="leaderboard-section">
        <h2 class="leaderboard-title">Leaderboard</h2>
        {#if leaderboard && leaderboard.length > 0}
            <table class="leaderboard-table">
                <thead>
                    <tr>
                        <th>Rank</th>
                        <th>Name</th>
                        <th>Score</th>
                    </tr>
                </thead>
                <tbody>
                    {#each leaderboard as entryLeaderboard, index}
                        <tr class="rank-{index + 1}">
                            <td class="rank-cell">#{index + 1}</td>
                            <td class="name-cell">{entryLeaderboard.name}</td>
                            <td class="score-cell">{entryLeaderboard.score}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        {:else}
            <p class="no-entries">No entries yet... be the first!</p>

        {/if}
    </div>

<style>
    .game-container {
        position: relative;
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        padding-bottom: 10px;
    }
    canvas {
        image-rendering: pixelated;
        image-rendering: crisp-edges;
        border: 2px solid rgba(0, 212, 170, 0.3);
        border-radius: 8px;
    }
    .score {
        position: absolute;
        top: 20px;
        left: 20px;
        z-index: 10;
        color: white;
        margin: 0;
        font-family: sans-serif;
        font-size: 20px;
        font-weight: 600;
        background: rgba(0, 0, 0, 0.7);
        padding: 10px 20px;

        border-radius: 8px;
        border: 2px solid rgba(0, 212, 170, 0.3);

    }

    .play-overlay {
        position: absolute;
        z-index: 20;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 15px;
    }

    .play-button {
        display: flex;
        align-items: center;
        gap: 10px;
        background: rgba(0, 212, 170, 0.9);
        color: white;
        border: none;
        border-radius: 8px;
        padding: 15px 30px;
        font-size: 1.1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .play-button:hover:not(:disabled) {
        background: rgba(0, 212, 170, 1);
        transform: scale(1.05);
    }

    .play-button:disabled {
        opacity: 0.5;
        cursor: wait;
        background: rgba(100, 100, 100, 0.6);
    }

    .exit-button {
        position: absolute;
        top: 20px;
        right: 20px;
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        color: #f0f6fc;
        transition: all 0.3s ease;
        z-index: 1;
    }
    
    .exit-button:hover {
        background: rgba(0, 212, 170, 0.2);
        border-color: #00d4aa;
        transform: rotate(90deg);

    }
    .game-over-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 30;
        animation: fadeIn 0.3s ease;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    .game-over-modal {
        background: rgba(22, 27, 34, 0.95);
        border: 1px solid rgba(0, 212, 170, 0.3);
        border-radius: 12px;
        padding: 50px 70px;
        text-align: center;
        animation: slideUp 0.4s ease;
    }

    @keyframes slideUp {
        from {
            transform: translateY(30px);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
    }

    .game-over-title {
        font-size: 2.5rem;
        font-weight: 600;
        color: #f0f6fc;
        margin: 0 0 15px 0;
    }

    .final-score {
        font-size: 1.3rem;
        color: #8b949e;
        margin: 0 0 35px 0;
        font-weight: 500;
    }

    .final-score::before {
        content: '';
        display: block;
        width: 40px;
        height: 2px;
        background: rgba(0, 212, 170, 0.4);
        margin: 0 auto 15px;
    }

    .score-submission {
        display: flex;
        flex-direction: column;
        gap: 15px;
        margin-bottom: 25px;
        width: 100%;
    }

    .name-input {
        padding: 12px 16px;
        background: rgba(0, 0, 0, 0.3);
        border: 1px solid rgba(0, 212, 170, 0.3);
        border-radius: 6px;
        color: #f0f6fc;
        font-size: 1rem;
        font-family: inherit;
        transition: all 0.2s ease;
        outline: none;
    }

    .name-input::placeholder {
        color: #8b949e;
    }

    .name-input:focus {
        border-color: #00d4aa;
        background: rgba(0, 0, 0, 0.5);
        box-shadow: 0 0 0 2px rgba(0, 212, 170, 0.1);
    }
    .submit-button {
        padding: 12px 24px;
        color: white;
        background: rgba(0, 0, 0, 0.5);
        border: none;
        border-radius: 6px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .submit-button:hover:not(:disabled) {
        background: #00e6bb;
        transform: translateY(-2px);
    }

    .submit-button:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .restart-button {
        display: flex;
        align-items: center;
        gap: 10px;
        background: transparent;
        color: #00d4aa;
        border: 1px solid #00d4aa;
        border-radius: 8px;
        padding: 12px 28px;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
        margin: 0 auto;
    }

    .restart-button:hover {
        background: rgba(0, 212, 170, 0.1);
        transform: translateY(-2px);
    }

    .terminal-console {
        position: absolute;
        bottom: 20px;
        left: 50%;
        transform: translateX(-50%);
        width: 1024px; 
        max-width: calc(100vw - 40px);
        max-height: 80px;
        background: rgba(0, 0, 0, 0.85);
        border: 1px solid rgba(0, 212, 170, 0.3);
        border-radius: 4px;
        padding: 8px 12px;
        font-family: 'Courier New', monospace;
        font-size: 13px;
        color: #00d4aa;
        overflow: hidden;
        z-index: 15;
    }

    .terminal-message {
        line-height: 1.4;
        margin: 2px 0;
        opacity: 0;
        animation: fadeIn 0.3s ease forwards;
    }

    @keyframes fadeIn {
        to {
            opacity: 1;
        }
    }
    .leaderboard-section {
        width: 100%;
        max-width: 600px;
        margin: 60px auto 40px;
        padding: 30px;
        background: rgba(0, 0, 0, 0.6);
        border: 5px solid rgba(0, 212, 170, 0.3);
        border-radius: 8px;
    }

    .leaderboard-title {
        text-align: center;
        color: #f0f6fc;
        font-size: 1.8rem;
        font-weight: 600;
        margin: 0 0 25px 0;
        text-transform: uppercase;
        letter-spacing: 1px;
    }

    .leaderboard-table {
        width: 100%;
        border-collapse: collapse;
        font-family: 'Courier New', monospace;
    }

    .leaderboard-table thead {
        border-bottom: 2px solid rgba(0, 212, 170, 0.3);
    }

    .leaderboard-table th {
        padding: 12px 8px;
        text-align: left;
        color: #00d4aa;
        font-weight: 600;
        font-size: 0.9rem;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .leaderboard-table td {
        padding: 12px 8px;
        color: #9299a5;
        border-bottom: 1px solid rgba(139, 148, 158, 0.1);
    }

    .leaderboard-table tbody tr {
        transition: background 0.2s ease;
    }

    .leaderboard-table tbody tr:hover {
        background: rgba(0, 212, 170, 0.05);
    }

    /* Top 3 special styling */
    .leaderboard-table .rank-1 .rank-cell {
        color: #ffd700;
        font-weight: 700;
        font-size: 1.1rem;
    }

    .leaderboard-table .rank-1 .name-cell,
    .leaderboard-table .rank-1 .score-cell {
        color: #f0f6fc;
        font-weight: 700;
    }

    .leaderboard-table .rank-2 .rank-cell {
        color: #c0c0c0;
        font-weight: 600;
    }

    .leaderboard-table .rank-2 .name-cell,
    .leaderboard-table .rank-2 .score-cell {
        color: #c9d1d9;
        font-weight: 500;
    }

    .leaderboard-table .rank-3 .rank-cell {
        color: #cd7f32;
        font-weight: 600;
    }

    .leaderboard-table .rank-3 .name-cell,
    .leaderboard-table .rank-3 .score-cell {
        color: #b1bac4;
        font-weight: 500;
    }

    .rank-cell {
        width: 60px;
        text-align: center;
    }

    .score-cell {
        text-align: right;
        font-weight: 500;
    }

    .no-entries {
        text-align: center;
        color: #8b949e;
        font-style: italic;
        padding: 20px;
        margin: 0;
    }

</style>