<script lang="ts">
    import { onMount } from "svelte";
    import { SpaceInvaders } from "./game";

    let canvas: HTMLCanvasElement;
    let game: SpaceInvaders;

    onMount(() => {
        game = new SpaceInvaders(canvas);
        game.start()
    });

</script>

<h1>this is the space invaders</h1>
<canvas bind:this={canvas}></canvas>

<style>

</style>