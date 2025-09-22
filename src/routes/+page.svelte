<script lang="ts">
	//import { fly } from 'svelte/transition';
    import Header from '../lib/Header.svelte';
    import Aboutme from '../lib/Aboutme.svelte';
</script>
<Header></Header>
<div class="main-page">
    <div class="hero-section">
        <div class="hero-content">
            <div class="name-header">
                <p class="name-subtitle">Hey! I'm</p>
                <h1 class="name-title">Edwin Metcalf</h1>
                <p class="name-subtitle">Check out what I am working on</p>
                <!-- prolly want to replace with cool buttons -->

            </div>
            <!-- add a fuckin photo brother -->
        </div>
    </div>
    <Aboutme></Aboutme>

</div>
<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>

<style>
    .main-page {
        background: linear-gradient(-45deg, #0d1117, #161b22, #21262d, #00d4aa);
        background-size: 400% 400%;
        animation: gradient 30s linear infinite;
        min-height: 100vh;
    }
    @keyframes gradient {
        0% { background-position: 0% 50%; }
        50% { background-position: 100% 50%; }
        100% { background-position: 0% 50%; }
    }
    .name-header {
        text-align: center;
        padding: 100px, 20px, 60px;
        background: transparent;
    }
    .name-header p {
        text-align: left;
    }
    .name-title {
        font-size: 4rem;
        font-weight: 700;
        background: linear-gradient(135deg, #f0f6fc, #00d4aa, #58a6ff);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin: 0;
    }
    .name-subtitle {
        font-size: 2rem;
        color: #8b949e;
        margin: auto;
        font-weight: 300;
    }
    .hero-content {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        margin: 0 auto;
    }
    .hero-section {
        min-height: 90vh;
        display: flex;
        align-items: center;
        padding: 0 5%;
    }
    @keyframes glow {
    from {
        filter: drop-shadow(0 0 5px rgba(0, 212, 170, 0.3));
    }
    to {
        filter: drop-shadow(0 0 20px rgba(0, 212, 170, 0.6));
    }
}
</style>