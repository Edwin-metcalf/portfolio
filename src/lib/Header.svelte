<script lang="ts">
    import {onMount} from 'svelte';
    import { fly } from 'svelte/transition';
    import { quintOut } from 'svelte/easing';
   let mounted = false;

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
            <a href="#skills" in:fly={{ y: -50, duration: 800, delay: 400, easing: quintOut}}>Skills</a>
            <a href="#portfolio" in:fly={{ y: -50, duration: 800, delay: 600, easing: quintOut}}>Portfolio</a>
            <a href="#contact" in:fly={{ y: -50, duration: 800, delay: 800, easing: quintOut}}>Contact</a>
        </div>
    {/if}

</div>

<style>
    .header {
        overflow: hidden;
        background: linear-gradient(90deg, #0d1117, #161b22, #21262d);
        padding: 20px 10px;
        color: #f0f6fc;
        transform: translateY(-2px);
    }
    .header a {
        float: left;
        color: #f0f6fc;
        text-align: center;
        padding: 12px;
        text-decoration: none;
        font-size: 18px;
        line-height: 25px;
        border-radius: 4px;
    }
    .header a:hover {
        background-color: #00d4aa;
        color: #0d1117;
        transform: translateY(-2px);
    }
    .header-right {
        float: right;
    }
    .logo {
        width: 50px;
        height: auto;
        display: block;
        transition: filter 0.3s ease;
    }
    .logo:hover {
        filter: brightness(0.5);
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