<script lang="ts">
	import { fly, scale } from 'svelte/transition';
    import { quintOut } from 'svelte/easing';
    import{ onMount } from 'svelte';

    import Header from '../lib/Header.svelte';
    import Aboutme from '../lib/Aboutme.svelte';
    import Languages from '$lib/Languages.svelte';
    import Portfolio from '$lib/Portfolio.svelte';
    import Contact from '$lib/Contact.svelte';
    import Sidebar from '$lib/Sidebar.svelte';

    let mounted: boolean = false;

    onMount(() => {
        mounted = true;
    });

</script>
<Header></Header>
<Sidebar></Sidebar>
<div class="main-page">
    <div class="hero-section">
        <div class="hero-content">
            <div class="name-header">
                {#if mounted}
                    <p class="name-subtitle" in:fly={{ y: -20, duration: 880, delay: 200, easing: quintOut}}>
                        Hey! I'm
                    </p>
                    <h1 class="name-title" in:fly={{ x: -50, duration: 1000, delay: 400, easing: quintOut}}>
                        Edwin Metcalf
                    </h1>
                    <a href="/ResumeEdwin2025oct.pdf" target="_blank" rel="noopener noreferrer" class="resume-btn" in:fly={{ y: 20, duration: 800, delay: 600, easing: quintOut}}>
                        Resume
                    </a>
                    <!--<p class="name-subtitle">Check out what I am working on</p>-->
                    <!-- prolly want to replace with cool buttons -->
                {/if}

            </div>
            {#if mounted}
                <div class="profile-wrapper" in:scale= {{ duration: 1000, delay: 300, start: 0.8, easing: quintOut}}>
                    <img src="/images/ProfilePicV1.jpg" alt="Me hiking up the highest peak in Wales" class="profile-pic" />
                </div>
            {/if}
        </div>
    </div>
    <Aboutme></Aboutme>
    <Languages></Languages>
    <Portfolio></Portfolio>
    <Contact></Contact>

</div>

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
        min-height: 80vh;
        display: flex;
        align-items: center;
        padding: 0 5%;
    }
    .profile-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 50px;
    }
    .profile-pic {
        width: 60vh;
        height: 60vh;
        border-radius: 50%;
        border: 5px solid #0d1117;
        object-fit: cover;
        box-shadow: 0 0 15px rgba(0, 212, 170, 0.4);
    }
    @keyframes glow {
        from {
            filter: drop-shadow(0 0 5px rgba(0, 212, 170, 0.3));
        }
        to {
            filter: drop-shadow(0 0 20px rgba(0, 212, 170, 0.6));
        }
        html {
            scroll-behavior: smooth;
        }
    }

    .resume-btn {
        background: transparent;
        color: #00d4aa;
        border: 2px solid #00d4aa;
        padding: 10px 26px;
        font-size: 1rem;
        font-weight: 600;
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;
        display: inline-flex;
        align-items: center;
        gap: 8px;
        text-decoration: none;
    }

    .resume-btn:hover {
        background: rgba(0, 212, 170, 0.15);
        transform: translateY(-2px);
        box-shadow: 0 5px 15px rgba(0, 212, 170, 0.2);
    }

    .resume-btn:active {
        transform: translateY(0);
    }

    .resume-btn {
        background: transparent;
        color: #00d4aa;
        border: 2px solid #00d4aa;
        padding: 10px 26px;
        font-size: 1rem;
        font-weight: 600;
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;
        display: inline-flex;
        align-items: center;
        gap: 8px;
        text-decoration: none;
        margin-top: 20px;  /* ADD THIS for spacing */
    }

    @media (max-width: 768px) {
        .hero-section {
            min-height: auto;
            padding: 2rem 5%;
        }
        
        .hero-content {
            flex-direction: column-reverse;
            gap: 2rem;
        }
        
        .name-header {
            padding: 20px;
            text-align: center;
        }
        
        .name-header p {
            text-align: center;
        }
        
        .name-title {
            font-size: 2.5rem;
        }
        
        .name-subtitle {
            font-size: 1.2rem;
        }
        
        .profile-wrapper {
            margin-right: 0;
            margin-top: 2rem;
        }
        
        .profile-pic {
            width: 250px;
            height: 250px;
            border: 3px solid #0d1117;
        }
    }

</style>