<script lang="ts">
    import { onMount } from "svelte";
    import { SpaceInvaders } from "./game";
    import { Play, X, RotateCcw } from "lucide-svelte"

    let canvas: HTMLCanvasElement;
    let game: SpaceInvaders;
    let score = $state(0);
    let gameStarted = $state(false);
    let imagesLoaded = $state(false);
    let gameOver = $state(false);
    let finalScore = $state(0);

    let consoleMessages = $state<string[]>([]);
    const MAX_MESSAGES = 3;


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
        game.restart();
    }
    function addMessage(msg: string) {
        if (consoleMessages.length >= MAX_MESSAGES) {
            consoleMessages.shift();
        }
        consoleMessages.push(msg);
        consoleMessages = consoleMessages;
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

</style>