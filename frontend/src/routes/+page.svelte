<script lang="ts">
	import { fly, scale } from 'svelte/transition';
    import { quintOut } from 'svelte/easing';
    import{ onMount } from 'svelte';
    import { ArrowDown } from 'lucide-svelte';

    import Header from '../lib/Header.svelte';
    import Aboutme from '../lib/Aboutme.svelte';
    import Languages from '$lib/Languages.svelte';
    import Portfolio from '$lib/Portfolio.svelte';
    import Contact from '$lib/Contact.svelte';
    import Sidebar from '$lib/Sidebar.svelte';
    import Courses from '$lib/Courses.svelte';

    

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
                    <div class="button-holder">
                        <a href="/Resume.pdf" target="_blank" rel="noopener noreferrer" class="resume-btn" in:fly={{ y: 20, duration: 800, delay: 600, easing: quintOut}}>
                            Resume
                        </a>
                    </div>
                    <a href="#about-me" class="scroll-down" aria-label="Scroll down for more" in:fly={{y: 20, duration: 800, delay: 600, easing: quintOut}}>
                        <ArrowDown size={20}/>
                    </a>
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
    <Portfolio></Portfolio>
    <Courses></Courses>
    <Languages></Languages>
    <!-- section here for going to my projects I am currently working on -->
    <section class="explore-section">
        <div class="explore-container">
            <h2 class="explore-title">Want to check out what I am working on?</h2>
            <div class="explore-buttons">
                <a href="/games" class="explore-button">
                    <span class="explore-button-title">Current Projects</span>
                    <span class="explore-button-desc">See what else I've been building</span>
                </a>
                <a href="/games/game-stats/clash-royale/matchup-generator" class="explore-button">
                    <span class="explore-button-title">Clash Royale Matchup</span>
                    <span class="explore-button-desc">Check your head-to-head stats against a friend!</span>
                </a>
                <a href="/games/space-invaders" class="explore-button">
                    <span class="explore-button-title">Space Invaders</span>
                    <span class="explore-button-desc">A vampire survivors twist on the classic, with a global leaderboard</span>
                </a>
            </div>
        </div>
    </section>
    <Contact></Contact>

</div>

<style>
 .explore-section {
    padding: 80px 5%;
    text-align: center;

  }

  .explore-container {
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
  }

  .explore-title {
    font-size: 2.5rem;
    font-weight: 600;
    color: var(--text);
    font-family: var(--font-display);
    margin: 0;
  }

  .explore-buttons {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 1rem;
        margin-top: 1.5rem;
    }

  .explore-button {
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: 0.4rem;
        font-family: var(--font-mono);
        background: var(--panel);
        border: 1px solid var(--border);
        border-radius: var(--radius);
        padding: 1.2rem 1rem;
        text-decoration: none;
        text-align: left;
        transition: transform 0.15s ease, border-color 0.15s ease;
    }

  .explore-button:hover {
        transform: translateY(-2px);
        border-color: var(--mint);
    }

    .explore-button-title {
        color: var(--mint);
        font-weight: bold;
        font-size: 1rem;
    }

    .explore-button-desc {
        color: var(--text-muted);
        font-size: 0.85rem;
        line-height: 1.3;
    }
    .main-page {
        background: var(--bg);
        min-height: 100vh;
    }
    .name-header {
        text-align: center;
        padding: 100px 20px 60px;
        background: transparent;
    }
    .name-header p {
        text-align: left;
    }
    .name-title {
        font-size: 4rem;
        font-weight: 700;
        color: var(--text);
        margin: 0;
    }
    .name-subtitle {
        font-size: 2rem;
        color: var(--text-muted);
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
        width: 50vh;
        height: 50vh;
        border-radius: 50%;
        border: 1px solid var(--border);
        object-fit: cover;
    }
    .resume-btn:hover {
        background: rgba(62, 207, 192, 0.15);
        transform: translateY(-1px);
    }

    .resume-btn:active {
        transform: translateY(0);
    }

    .resume-btn {
        background: transparent;
        color: var(--mint);
        border: 1px solid var(--mint);
        padding: 10px 26px;
        font-size: 1rem;
        font-weight: 500;
        border-radius: var(--radius);
        cursor: pointer;
        transition: all 0.15s ease;
        display: inline-flex;
        align-items: center;
        gap: 8px;
        text-decoration: none;
        margin-top: 20px; 
    }
    .button-holder {
        display: flex;
        justify-content: center;
        gap: 20px;
    }

    .scroll-down {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        margin-top: 28px;
        padding: 8px;
        color: var(--text-muted);
        border-radius: var(--radius);
        transition: color 0.15s ease, transform 0.15s ease;
    }

    .scroll-down:hover {
        color: var(--mint);
        transform: translateY(2px);
    }

    :global(body) {
        background: var(--bg); 
        margin: 0;
        padding: 0;
        overflow-x: hidden;
    }

    :global(html) {
        background: var(--bg);
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
            border: 1px solid var(--border);
        }
    }

</style>