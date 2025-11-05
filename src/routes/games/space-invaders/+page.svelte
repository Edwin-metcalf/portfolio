<script lang="ts">
    import { onMount } from "svelte";
    import { SpaceInvaders } from "./game";

    let canvas: HTMLCanvasElement;
    let game: SpaceInvaders;
    let score = $state(0);

    onMount(() => {
        game = new SpaceInvaders(canvas, (newScore) => {
            score = newScore
        });
        
        game.start()

    });

</script>

<div style="position: relative">
    <p class="score"><span>Score:</span> <span>{score}</span></p>
    <canvas bind:this={canvas}></canvas>
</div>

<style>
    canvas {
        image-rendering: pixelated;
        image-rendering: crisp-edges;
    }
    .score {
        position: absolute;
        z-index: 10;
        color: white;
        margin-left: 10px;
        margin: 0;
        font-family: sans-serif;
        font-size: 20px;
    }

</style>