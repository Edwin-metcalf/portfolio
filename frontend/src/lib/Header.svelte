<script lang="ts">
    import {onMount} from 'svelte';
    import { fly } from 'svelte/transition';
    import { quintOut } from 'svelte/easing';
   let mounted: boolean = false;

    onMount(() => {
        mounted = true;
    });


</script>

<div class="header">
    {#if mounted}
        <a href="/" class="logo-wrapper" in:fly={{ y: -50, duration: 1000, easing: quintOut}}>
            <img src="/images/EMlogo.png" alt="Edwin Metcalf Logo" class="logo" />
        </a>
        <div class="header-right">
            <a href="#about-me" in:fly={{ y: -50, duration: 800, delay: 200, easing: quintOut}}>About me</a>
            <a href="#portfolio" in:fly={{ y: -50, duration: 800, delay: 400, easing: quintOut}}>Portfolio</a>
            <a href="#skills" in:fly={{ y: -50, duration: 800, delay: 600, easing: quintOut}}>Skills</a>
            <a href="#contact" in:fly={{ y: -50, duration: 800, delay: 800, easing: quintOut}}>Contact</a>
            <a href="/games" in:fly={{ y: -50, duration: 800, delay: 1000, easing: quintOut}}>Current Projects</a>
        </div>
    {/if}

</div>

<style>
    .header {
        overflow: hidden;
        background: var(--bg);
        border-bottom: 1px solid var(--border);
        padding: 20px 10px;
        color: var(--text);
        transform: none;
    }
    .header a {
        float: left;
        color: var(--text);
        text-align: center;
        padding: 12px;
        text-decoration: none;
        font-size: 30px;
        line-height: 25px;
        border-radius: var(--radius);
        transition: all 0.15s ease;
    }
    .header a:hover {
        background-color: rgba(62, 207, 192, 0.15);
        color: var(--mint);
        transform: translateY(-1px);
    }
    .header-right {
        float: right;
    }
    .logo {
        width: 65px;
        height: auto;
        display: block;
        transition: filter 0.3s ease;
    }
    .logo:hover {
        filter: brightness(0.7);
    }


    @media (max-width: 768px) {
        .header {
            padding: 15px 5%;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 15px;
        }
        
        .logo-wrapper {
            float: none;
            margin-bottom: 10px;
        }
        
        .logo {
            width: 40px;
        }
        
        .header-right {
            float: none;
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            gap: 5px;
            width: 100%;
        }
        
        .header a {
            float: none;
            font-size: 16px;
            padding: 10px 15px;
            flex: 0 1 auto;
        }
    }
</style>